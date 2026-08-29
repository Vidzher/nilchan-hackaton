package resource

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"nilchan-hackaton/internal/quiz"
	quizgen "nilchan-hackaton/internal/quiz/gen"
)

const maxConcurrentGeneration = 5

type quizGenerator interface {
	Generate(ctx context.Context, request quizgen.GenerationRequest) (*quizgen.GeneratedQuiz, error)
}

type processingRepository interface {
	CompleteGeneration(ctx context.Context, resourceID int64, title string, questions []quiz.Question) error
	FailGeneration(ctx context.Context, resourceID int64) error
}

type Processor struct {
	repo      processingRepository
	generator quizGenerator
	ctx       context.Context
	cancel    context.CancelFunc
	semaphore chan struct{}
	workers   sync.WaitGroup
}

func NewProcessor(repo processingRepository, generator quizGenerator) *Processor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Processor{
		repo:      repo,
		generator: generator,
		ctx:       ctx,
		cancel:    cancel,
		semaphore: make(chan struct{}, maxConcurrentGeneration),
	}
}

func (p *Processor) Submit(resource *Resource) {
	p.workers.Go(func() {

		select {
		case p.semaphore <- struct{}{}:
			defer func() { <-p.semaphore }()
		case <-p.ctx.Done():
			return
		}

		if err := p.process(resource); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("resource quiz generation failed resource_id=%d error=%v", resource.ID, err)
			p.fail(resource.ID)
		}
	})
}

func (p *Processor) Close() {
	p.cancel()
	p.workers.Wait()
}

func (p *Processor) process(resource *Resource) error {
	generated, err := p.generator.Generate(p.ctx, quizgen.GenerationRequest{
		SourceTitle: resource.Title,
		SourceText:  resource.Content,
	})
	if err != nil {
		return err
	}

	questions := make([]quiz.Question, len(generated.Questions))
	for index, generatedQuestion := range generated.Questions {
		questions[index] = quiz.Question{
			Text:         generatedQuestion.Text,
			CorrectIndex: generatedQuestion.CorrectIndex,
			Explanation:  generatedQuestion.Explanation,
			Evidence:     generatedQuestion.Evidence,
		}
		copy(questions[index].Options[:], generatedQuestion.Options)
	}

	return p.repo.CompleteGeneration(p.ctx, resource.ID, resource.Title, questions)
}

func (p *Processor) fail(resourceID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := p.repo.FailGeneration(ctx, resourceID); err != nil && !errors.Is(err, ErrStateConflict) {
		log.Printf("resource failure transition failed resource_id=%d error=%v", resourceID, err)
	}
}
