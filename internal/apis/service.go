package apis

import (
	"context"
	"fmt"
	"github.com/jijiechen/dami-ultra/internal/ai"
	kong_api "github.com/jijiechen/dami-ultra/internal/apis/kong-api"
	"github.com/jijiechen/dami-ultra/internal/business"
)

var OpenAIKey = "5VoC5nUaLiGuVpVmHuoBcoLGyQ6ezgsbyqRqXtSqw3yPDKkk7R7OJQQJ99BCACi0881XJ3w3AAAAACOGKfkY"
var KongAdminURL = "http://ec2-54-166-250-69.compute-1.amazonaws.com:8001/routes"

type DamiService struct {
	OpenAISDK           *ai.OpenAI
	KongGatewayAdminUrl string
}

func NewService() *DamiService {
	sdk := &ai.OpenAI{APIKey: OpenAIKey}
	return &DamiService{OpenAISDK: sdk, KongGatewayAdminUrl: KongAdminURL}
}

func (s *DamiService) PostMessage(ctx context.Context, list business.MessageList) (string, error) {
	lastMessage := list.Messages[len(list.Messages)-1]

	aiResp, err := s.OpenAISDK.ValidateKongConfiguration(lastMessage.Content)
	if err != nil {
		return "", err
	}

	return s.applyValidatedKongConfig(aiResp)
}

func (s *DamiService) applyValidatedKongConfig(aiResp ai.ValidateOpenAIResponse) (string, error) {
	if aiResp.Valid {
		fmt.Println("applying configuration", aiResp.RawConfiguration)
		err := kong_api.ApplyKongConfig(s.KongGatewayAdminUrl, aiResp.RawConfiguration)
		if err != nil {
			return "", fmt.Errorf("your configuration look perfect, but unfortunately it's not possible to be applied at the moment: %w", err)
		}

		return "Your configuration has been applied successfully", nil
	} else {
		return aiResp.ErrorMessages, nil
	}
}
