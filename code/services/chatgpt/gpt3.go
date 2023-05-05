package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"start-feishubot/initialization"
	customOpenai "start-feishubot/services/openai"
)

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPT struct {
	config *initialization.Config
}

type Gpt3 interface {
	StreamChat() error
	StreamChatWithHistory() error
}

func NewGpt3(config *initialization.Config) *ChatGPT {
	return &ChatGPT{config: config}
}

func (c *ChatGPT) StreamChat(ctx context.Context,
	msg []customOpenai.Messages,
	responseStream chan string) error {
	//change msg type from Messages to openai.ChatCompletionMessage
	chatMsgs := make([]openai.ChatCompletionMessage, len(msg))
	for i, m := range msg {
		chatMsgs[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return c.StreamChatWithHistory(ctx, chatMsgs, 2000,
		responseStream)
}

func (c *ChatGPT) StreamChatWithHistory(ctx context.Context, msg []openai.ChatCompletionMessage, maxTokens int,
	responseStream chan string,
) error {
	config := openai.DefaultConfig(c.config.OpenaiApiKeys[0])
	config.BaseURL = c.config.OpenaiApiUrl + "/v1"

	proxyClient, parseProxyError := customOpenai.GetProxyClient(c.config.HttpProxy)
	if parseProxyError != nil {
		return parseProxyError
	}
	config.HTTPClient = proxyClient

	client := openai.NewClientWithConfig(config)
	//pp.Printf("client: %v", client)
	req := openai.ChatCompletionRequest{
		Model:       c.config.OpenaiModel,
		Messages:    msg,
		N:           1,
		Temperature: 0.7,
		MaxTokens:   maxTokens,
		TopP:        1,
		//Moderation:     true,
		//ModerationStop: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Errorf("CreateCompletionStream returned error: %v", err)
	}

	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			//fmt.Println("Stream finished")
			return nil
		}
		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return err
		}
		responseStream <- response.Choices[0].Delta.Content
	}
	return nil

}
