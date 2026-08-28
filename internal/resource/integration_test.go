package resource_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"nilchan-hackaton/internal/llm"
	"nilchan-hackaton/internal/parser"
	"nilchan-hackaton/internal/quiz"
	quizgen "nilchan-hackaton/internal/quiz/gen"
	"nilchan-hackaton/internal/resource"
	resourcerepo "nilchan-hackaton/internal/resource/repository"
	"nilchan-hackaton/internal/storage"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type llmCompleter interface {
	Complete(context.Context, llm.Request) (string, error)
}

type recordingCompleter struct {
	inner     llmCompleter
	mu        sync.Mutex
	responses []string
	errors    []error
}

func (r *recordingCompleter) Complete(ctx context.Context, request llm.Request) (string, error) {
	response, err := r.inner.Complete(ctx, request)
	r.mu.Lock()
	defer r.mu.Unlock()
	if err != nil {
		r.errors = append(r.errors, err)
	} else {
		r.responses = append(r.responses, response)
	}
	return response, err
}

func (r *recordingCompleter) logAttempts(t *testing.T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for index, response := range r.responses {
		t.Logf("      raw LLM response %d:\n%s", index+1, response)
	}
	for index, err := range r.errors {
		t.Logf("      LLM request error %d: %v", index+1, err)
	}
}

func TestAddResourceHappyPathRealProviders(t *testing.T) {
	t.Log("[1/8] Loading provider configuration")
	_ = godotenv.Load("../../.env.local")

	var providerConfig struct {
		OpenRouter struct {
			ModelName string `yaml:"model_name" env:"OPENROUTER_MODEL"`
		} `yaml:"openrouter"`
	}
	if err := cleanenv.ReadConfig("../config/config.yml", &providerConfig); err != nil {
		t.Fatalf("read provider config: %v", err)
	}

	firecrawlKey := os.Getenv("FIRECRAWL_KEY")
	openRouterKey := os.Getenv("OPENROUTER_API_KEY")
	openRouterModel := providerConfig.OpenRouter.ModelName
	if firecrawlKey == "" || openRouterKey == "" {
		t.Skip("real-provider test requires FIRECRAWL_KEY and OPENROUTER_API_KEY")
	}
	if openRouterModel == "" {
		t.Fatal("OpenRouter model is not configured")
	}

	resourceURL := os.Getenv("RESOURCE_INTEGRATION_URL")
	if resourceURL == "" {
		resourceURL = "https://go.dev/blog/pipelines"
	}
	t.Logf("      model=%s url=%s", openRouterModel, resourceURL)

	t.Log("[2/8] Creating isolated SQLite database and test user")
	store, err := storage.NewStorage(filepath.Join(t.TempDir(), "integration.db"))
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	defer store.Close()

	result, err := store.DB.Exec(`
		INSERT INTO users(email, username, password_hash)
		VALUES('learner@example.com', 'learner', 'hash')
	`)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get user id: %v", err)
	}
	if _, err := store.DB.Exec("INSERT INTO user_progress(user_id) VALUES(?)", userID); err != nil {
		t.Fatalf("create user progress: %v", err)
	}
	t.Logf("      user_id=%d", userID)

	t.Log("[3/8] Creating real Firecrawl and OpenRouter clients")
	const firecrawlTimeout = 30 * time.Second
	firecrawlClient, err := parser.NewFirecrawlClient(
		firecrawlKey,
		os.Getenv("FIRECRAWL_BASE_URL"),
		&http.Client{Timeout: firecrawlTimeout},
	)
	if err != nil {
		t.Fatalf("create Firecrawl client: %v", err)
	}
	openRouterClient, err := llm.NewOpenRouterClient(openRouterKey, openRouterModel)
	if err != nil {
		t.Fatalf("create OpenRouter client: %v", err)
	}
	recorder := &recordingCompleter{inner: openRouterClient}
	generator, err := quizgen.NewGenerator(recorder)
	if err != nil {
		t.Fatalf("create quiz generator: %v", err)
	}

	repo := resourcerepo.New(store)
	processor := resource.NewProcessor(repo, generator)
	defer processor.Close()
	service := resource.NewService(repo, firecrawlClient, processor, firecrawlTimeout)

	t.Log("[4/8] Normalizing URL and fetching content from Firecrawl")
	addStarted := time.Now()
	created, err := service.Add(context.Background(), userID, resource.CreateResourceRequest{URL: resourceURL})
	if err != nil {
		t.Fatalf("add resource: %v", err)
	}
	t.Logf("      Firecrawl completed in %s", time.Since(addStarted).Round(time.Millisecond))
	t.Logf("      page title: %s", created.Title)
	t.Logf("      page content:\n%s", created.Content)
	t.Logf("      tags=%v content_chars=%d", created.Tags, len([]rune(created.Content)))
	if created.Status != resource.StatusProcessing {
		t.Fatalf("initial status = %q, want %q", created.Status, resource.StatusProcessing)
	}
	t.Logf("[5/8] Resource persisted: id=%d status=%s normalized_url=%s", created.ID, created.Status, created.URL)

	t.Log("[6/8] Waiting for background OpenRouter quiz generation")
	generationStarted := time.Now()
	deadline := time.Now().Add(4*time.Minute + 30*time.Second)
	lastStatus := resource.Status("")
	for {
		var status resource.Status
		if err := store.DB.QueryRow("SELECT status FROM resources WHERE id = ?", created.ID).Scan(&status); err != nil {
			t.Fatalf("read resource status: %v", err)
		}
		if status != lastStatus {
			t.Logf("      status=%s elapsed=%s", status, time.Since(generationStarted).Round(time.Millisecond))
			lastStatus = status
		}
		if status == resource.StatusNotCompleted {
			break
		}
		if status == resource.StatusFailed {
			t.Log("      generation failed; dumping rejected provider responses")
			recorder.logAttempts(t)
			t.Fatal("quiz generation failed and resource transitioned to FAILED")
		}
		if time.Now().After(deadline) {
			recorder.logAttempts(t)
			t.Fatalf("resource status = %q, want %q before timeout", status, resource.StatusNotCompleted)
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Log("[7/8] Reading atomically persisted quiz and checking verification data")
	var quizID int64
	var quizTitle string
	var questionsJSON string
	if err := store.DB.QueryRow(`
		SELECT id, title, questions_json FROM quizzes WHERE resource_id = ?
	`, created.ID).Scan(&quizID, &quizTitle, &questionsJSON); err != nil {
		t.Fatalf("read generated quiz: %v", err)
	}

	var questions []quiz.Question
	if err := json.Unmarshal([]byte(questionsJSON), &questions); err != nil {
		t.Fatalf("decode quiz questions: %v", err)
	}
	if len(questions) < 5 || len(questions) > 10 {
		t.Fatalf("question count = %d, want 5..10", len(questions))
	}
	for index, question := range questions {
		if question.VerificationSalt == "" || question.CorrectAnswerHash == "" {
			t.Fatalf("question %d has missing verification data", index)
		}
	}
	t.Logf("      quiz_id=%d questions=%d hashes=present", quizID, len(questions))

	t.Log("[8/8] Rendering final stored resource and quiz")
	output := struct {
		Resource struct {
			ID           int64           `json:"id"`
			URL          string          `json:"url"`
			Title        string          `json:"title"`
			Content      string          `json:"content"`
			Tags         []string        `json:"tags"`
			Status       resource.Status `json:"status"`
			ContentChars int             `json:"contentChars"`
		} `json:"resource"`
		Quiz quiz.Quiz `json:"quiz"`
	}{}
	output.Resource.ID = created.ID
	output.Resource.URL = created.URL
	output.Resource.Title = created.Title
	output.Resource.Content = created.Content
	output.Resource.Tags = created.Tags
	output.Resource.Status = resource.StatusNotCompleted
	output.Resource.ContentChars = len([]rune(created.Content))
	output.Quiz = quiz.Quiz{
		ID:         quizID,
		ResourceID: created.ID,
		Title:      quizTitle,
		Questions:  questions,
	}

	pretty, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		t.Fatalf("encode inspection output: %v", err)
	}
	t.Logf("      completed in %s\n%s", time.Since(addStarted).Round(time.Millisecond), strings.TrimSpace(string(pretty)))
}
