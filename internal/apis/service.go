package apis

import (
	"context"

	"github.com/jijiechen/dami-ultra/internal/business"
)

type DamiService struct {
	OpenAIKey string

	KongGatewayPublicUrl string
}

func NewService() *DamiService {
	return &DamiService{}
}

func (s *DamiService) PostMessage(ctx context.Context, list business.MessageList) (string, error) {
	return "hello world", nil
}
