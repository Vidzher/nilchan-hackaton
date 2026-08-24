package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/optionalnullable"
)

type OpenRouterClient struct {
	client    *openrouter.OpenRouter
	modelName string
	schema    map[string]any
}

func NewOpenRouterClient(apiKey, modelName string) (*OpenRouterClient, error) {
	if apiKey == "" {
		return nil, ErrInvalidApiKey
	}
	if modelName == "" {
		return nil, ErrInvalidModelName
	}
	var schema map[string]any

	if err := json.Unmarshal([]byte(quizSchemaJSON), &schema); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidSchema, err)
	}
	sdkClient := openrouter.New(
		openrouter.WithSecurity(apiKey),
	)
	return &OpenRouterClient{
		client:    sdkClient,
		modelName: modelName,
		schema:    schema,
	}, nil
}

func (c *OpenRouterClient) Generate(
	ctx context.Context,
	request CompletionRequest,
) (string, error) {
	responseFormat := components.CreateResponseFormatJSONSchema(
		components.ChatFormatJSONSchemaConfig{
			JSONSchema: components.ChatJSONSchemaConfig{
				Name:   "learning_backlog_quiz",
				Schema: c.schema,
				Strict: optionalnullable.From(openrouter.Pointer(true)),
			},
		},
	)

	result, err := c.client.Chat.Send(ctx, components.ChatRequest{
		Model: openrouter.Pointer(c.modelName),
		Messages: []components.ChatMessages{
			components.CreateChatMessagesSystem(
				components.ChatSystemMessage{
					Role: components.ChatSystemMessageRoleSystem,
					Content: components.CreateChatSystemMessageContentStr(
						request.SystemPrompt,
					),
				},
			),
			components.CreateChatMessagesUser(
				components.ChatUserMessage{
					Role: components.ChatUserMessageRoleUser,
					Content: components.CreateChatUserMessageContentStr(
						request.UserPrompt,
					),
				},
			),
		},
		ResponseFormat: &responseFormat,
		Temperature: optionalnullable.From(
			openrouter.Pointer(0.2),
		),
	}, nil)
	if err != nil {
		return "", fmt.Errorf("openrouter request failed: %w", err)
	}

	if result == nil ||
		result.ChatResult == nil ||
		len(result.ChatResult.Choices) == 0 {
		return "", ErrEmptyLLMResponse
	}

	content, exists := result.ChatResult.Choices[0].Message.Content.Get()
	if !exists ||
		content == nil ||
		content.Str == nil ||
		strings.TrimSpace(*content.Str) == "" {
		return "", ErrEmptyLLMResponse
	}

	return *content.Str, nil
}
