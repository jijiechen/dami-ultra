package business

import "context"

type MessageList struct {
	Messages []string `json:"items"`
}

type IDamiUltraService interface {
	PostMessage(ctx context.Context, list MessageList) (string, error)
}
