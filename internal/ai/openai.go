package ai

import (
	"context"
	"github.com/jijiechen/dami-ultra/internal/business"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"strings"
)

type OpenAI struct {
	APIKey string
}

const azureOpenAIEndpoint = "https://hackathonchina2025.services.ai.azure.com"

// The latest API versions, including previews, can be found here:
// https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versioning
const azureOpenAIAPIVersion = "2024-06-01"

func (o *OpenAI) CallAI(systemMessage string, messages []business.Message) (string, error) {
	client := openai.NewClient(
		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
		// Choose between authenticating using a TokenCredential or an API Key
		azure.WithAPIKey(o.APIKey),
	)

	var llmMessages []openai.ChatCompletionMessageParamUnion
	if systemMessage != "" {
		llmMessages = append(llmMessages, openai.SystemMessage(systemMessage))
	}

	for _, msg := range messages {
		switch msg.Author {
		case "system":
			llmMessages = append(llmMessages, openai.AssistantMessage(msg.Content))
		case "user":
			llmMessages = append(llmMessages, openai.UserMessage(msg.Content))
		}
	}

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F(llmMessages),
		Model:    openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func (o *OpenAI) CallAISingle(message string) (string, error) {
	aiResp, err := o.CallAI(
		"",
		[]business.Message{
			{Author: "user", Content: message},
		})

	if err != nil {
		return "", err
	}

	aiResp = strings.TrimFunc(aiResp, func(r rune) bool {
		return r == ' ' || r == '\n' || r == '"' || r == '\''
	})
	return aiResp, nil
}
