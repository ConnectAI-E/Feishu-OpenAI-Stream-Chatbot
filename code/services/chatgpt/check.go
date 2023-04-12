package chatgpt

import (
	"errors"
	"github.com/sashabaranov/go-openai"
)

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

func CheckChatCompletionMessages(messages []openai.ChatCompletionMessage) error {
	hasSystemMsg := false
	for _, msg := range messages {
		if msg.Role != ChatMessageRoleSystem && msg.Role != ChatMessageRoleUser && msg.Role != ChatMessageRoleAssistant {
			return errors.New("invalid message role")
		}
		if msg.Role == ChatMessageRoleSystem {
			if hasSystemMsg {
				return errors.New("more than one system message")
			}
			hasSystemMsg = true
		} else {
			// 对于非 system 角色的消息，Content 不能为空
			if msg.Content == "" {
				return errors.New("empty content in non-system message")
			}
		}
	}
	return nil
}
