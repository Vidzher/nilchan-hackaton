package llm

import (
	"context"
	"fmt"
	"strings"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/optionalnullable"
)

type OpenRouterClient struct {
	client    *openrouter.OpenRouter
	modelName string
}

type Request struct {
	SystemPrompt       string
	UserPrompt         string
	ResponseSchemaName string
	ResponseSchema     map[string]any
	Temperature        float64
}

func NewOpenRouterClient(apiKey, modelName string) (*OpenRouterClient, error) {
	apiKey = strings.TrimSpace(apiKey)
	modelName = strings.TrimSpace(modelName)

	if apiKey == "" {
		return nil, ErrInvalidApiKey
	}
	if modelName == "" {
		return nil, ErrInvalidModelName
	}

	return &OpenRouterClient{
		client:    openrouter.New(openrouter.WithSecurity(apiKey)),
		modelName: modelName,
	}, nil
}

func (c *OpenRouterClient) Complete(ctx context.Context, request Request) (string, error) {
	responseFormat := components.CreateResponseFormatJSONSchema(
		components.ChatFormatJSONSchemaConfig{
			JSONSchema: components.ChatJSONSchemaConfig{
				Name:   request.ResponseSchemaName,
				Schema: request.ResponseSchema,
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
			openrouter.Pointer(request.Temperature),
		),
	}, nil)
	if err != nil {
		return "", fmt.Errorf("openrouter request failed: %w", err)
	}

	if result == nil || result.ChatResult == nil || len(result.ChatResult.Choices) == 0 {
		return "", ErrEmptyResponse
	}

	content, exists := result.ChatResult.Choices[0].Message.Content.Get()
	if !exists || content == nil || content.Str == nil || strings.TrimSpace(*content.Str) == "" {
		return "", ErrEmptyResponse
	}

	return *content.Str, nil
}
