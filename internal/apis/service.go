package apis

import (
	"context"
	"encoding/json"
	"fmt"
	kong_api "github.com/jijiechen/dami-ultra/internal/apis/kong-api"
	"github.com/jijiechen/dami-ultra/internal/business"
	"strings"
)

var OpenAIKey = "5VoC5nUaLiGuVpVmHuoBcoLGyQ6ezgsbyqRqXtSqw3yPDKkk7R7OJQQJ99BCACi0881XJ3w3AAAAACOGKfkY"
var KongAdminURL = "http://ec2-54-166-250-69.compute-1.amazonaws.com:8001/routes"

func WrapQuestion(question string) string {
	return fmt.Sprintf(PromptTemplate, LuaValidator, question)
}

type DamiService struct {
	OpenAISDK           *OpenAI
	KongGatewayAdminUrl string
}

func NewService() *DamiService {
	sdk := &OpenAI{APIKey: OpenAIKey}
	return &DamiService{OpenAISDK: sdk, KongGatewayAdminUrl: KongAdminURL}
}

func (s *DamiService) PostMessage(ctx context.Context, list business.MessageList) (string, error) {
	lastMessage := list.Messages[len(list.Messages)-1]

	prompt := WrapQuestion(lastMessage.Content)

	aiResp, err := s.OpenAISDK.CallAI(prompt)
	fmt.Println(aiResp)

	aiResp = strings.ReplaceAll(aiResp, "```json", "")
	aiResp = strings.ReplaceAll(aiResp, "```", "")

	var respObj AIResponse
	err = json.Unmarshal([]byte(aiResp), &respObj)
	if err != nil {
		return "", err
	}

	if respObj.Valid {
		fmt.Println("applying configuration", respObj.RawConfiguration)
		err = kong_api.ApplyKongConfig(s.KongGatewayAdminUrl, respObj.RawConfiguration)
		if err != nil {
			return "", err
		}

		return "Your configuration has been applied successfully", nil
	} else {
		return "", fmt.Errorf(respObj.ErrorMessages)
	}
}
