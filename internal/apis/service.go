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

func (s *DamiService) PostMessages(ctx context.Context, list business.MessageList) (string, error) {
	lastMessage := list.Messages[len(list.Messages)-1]

	aiResp, err := s.OpenAISDK.ValidateKongConfiguration(lastMessage.Content)
	if err != nil {
		return "", err
	}

	if aiResp.Valid {
		return s.applyKongConfig(aiResp.RawConfiguration)
	} else {
		return aiResp.ErrorMessages, nil
	}
}

func (s *DamiService) PostOperationMessage(ctx context.Context, list business.MessageList) (string, error) {
	operationName, err := s.OpenAISDK.GetOperation(list.Messages)
	if err != nil {
		return "", err
	}

	fmt.Println(fmt.Sprintf("user operation: %s", operationName))
	if operationName == ai.OperationNone {
		return s.OpenAISDK.OperationNotUnderstood()
	}

	switch operationName {
	case ai.OperationCheckValidity:
		lastMessage := list.Messages[len(list.Messages)-1]
		aiResp, err := s.OpenAISDK.ValidateKongConfiguration(lastMessage.Content)
		if err != nil {
			return "", err
		}
		if aiResp.Valid {
			return s.OpenAISDK.AskIfApplyConfig()
		} else {
			return aiResp.ErrorMessages, nil
		}
	case ai.OperationApplyConfigYes:
		validatedConfig, err := s.OpenAISDK.ExtractValidatedConfig(list.Messages)
		if err != nil {
			return "", err
		}
		return s.applyKongConfig(validatedConfig)
	case ai.OperationApplyConfigNo:
		return s.OpenAISDK.ShowConfigDiscarded()
	case ai.OperationShowHelp:
		return s.OpenAISDK.ShowHelpMessage("https://docs.konghq.com/gateway/")
	}
	return "", fmt.Errorf("operation %s is not supported", operationName)
}

func (s *DamiService) applyKongConfig(kongConfig string) (string, error) {
	fmt.Println("applying configuration", kongConfig)
	err := kong_api.ApplyKongConfig(s.KongGatewayAdminUrl, kongConfig)
	if err != nil {
		return "", fmt.Errorf("your configuration look perfect, but unfortunately it's not possible to be applied at the moment: %w", err)
	}

	return "Your configuration has been applied successfully", nil
}
