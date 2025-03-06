package apis

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

type OpenAI struct {
	APIKey string
}

const azureOpenAIEndpoint = "https://hackathonchina2025.services.ai.azure.com"

// The latest API versions, including previews, can be found here:
// https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versioning
const azureOpenAIAPIVersion = "2024-06-01"

func (o *OpenAI) CallAI(msg string) (string, error) {
	client := openai.NewClient(
		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
		// Choose between authenticating using a TokenCredential or an API Key
		azure.WithAPIKey(o.APIKey),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(msg),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})

	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
