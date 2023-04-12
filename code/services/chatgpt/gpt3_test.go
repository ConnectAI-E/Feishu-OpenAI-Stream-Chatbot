package chatgpt

import (
	"context"
	"fmt"
	"start-feishubot/initialization"
	"start-feishubot/services/openai"
	"testing"
	"time"
)

func TestChatGPT_streamChat(t *testing.T) {
	// 初始化配置
	config := initialization.LoadConfig("../../config.yaml")

	// 准备测试用例
	testCases := []struct {
		msg        []openai.Messages
		wantOutput string
		wantErr    bool
	}{
		{
			msg: []openai.Messages{
				{
					Role:    "system",
					Content: "从现在起你要化身职场语言大师，你需要用婉转的方式回复老板想你提出的问题，或像领导提出请求。",
				},
				{
					Role:    "user",
					Content: "领导，我想请假一天",
				},
			},
			wantOutput: "",
			wantErr:    false,
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		// 准备输入和输出
		responseStream := make(chan string)
		ctx := context.Background()
		c := &ChatGPT{config: config}

		// 启动一个协程来模拟流式聊天
		go func() {
			err := c.StreamChat(ctx, tc.msg, responseStream)
			if err != nil {
				t.Errorf("streamChat() error = %v, wantErr %v", err, tc.wantErr)
			}
		}()

		// 等待输出并检查是否符合预期
		select {
		case gotOutput := <-responseStream:
			fmt.Printf("gotOutput: %v\n", gotOutput)

		case <-time.After(5 * time.Second):
			t.Errorf("streamChat() timeout, expected output not received")
		}
	}
}
