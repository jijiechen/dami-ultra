package business

import "context"

type Message struct {
	Author  string `json:"author"` // system/user
	Content string `json:"content"`
}

type MessageList struct {
	Messages []Message `json:"messages"`
}

type IDamiUltraService interface {
	PostMessage(ctx context.Context, list MessageList) (string, error)
}
