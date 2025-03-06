package apis

import (
	"context"
	"fmt"
	"github.com/jijiechen/dami-ultra/internal/business"
)

var OpenAIKey = "5VoC5nUaLiGuVpVmHuoBcoLGyQ6ezgsbyqRqXtSqw3yPDKkk7R7OJQQJ99BCACi0881XJ3w3AAAAACOGKfkY"

func WrapQuestion(question string) string {
	return fmt.Sprintf(PromptTemplate, LuaValidator, question)
}

type DamiService struct {
	OpenAISDK            *OpenAI
	KongGatewayAdminUrl  string
	KongGatewayPublicUrl string
}

func NewService() *DamiService {
	sdk := &OpenAI{APIKey: OpenAIKey}
	return &DamiService{OpenAISDK: sdk}
}

func (s *DamiService) PostMessage(ctx context.Context, list business.MessageList) (string, error) {
	lastMessage := list.Messages[len(list.Messages)-1]

	prompt := WrapQuestion(lastMessage.Content)

	aiResp, err := s.OpenAISDK.CallAI(prompt)
	return aiResp, err
}
