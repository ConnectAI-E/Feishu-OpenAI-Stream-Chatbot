package handlers

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
	"math/rand"
	"start-feishubot/services/openai"
	"time"
)

type MessageAction struct { /*消息*/
}

func (m *MessageAction) Execute(a *ActionInfo) bool {
	cardId, err2 := sendOnProcess(a)
	if err2 != nil {
		return false
	}
	pp.Println("cardId", cardId)

	updateMsg(a, cardId)
	return false

	msg := a.handler.sessionCache.GetMsg(*a.info.sessionId)
	msg = append(msg, openai.Messages{
		Role: "user", Content: a.info.qParsed,
	})
	completions, err := a.handler.gpt.Completions(msg)
	if err != nil {
		replyMsg(*a.ctx, fmt.Sprintf(
			"🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.info.msgId)
		return false
	}
	msg = append(msg, completions)
	a.handler.sessionCache.SetMsg(*a.info.sessionId, msg)
	//if new topic
	if len(msg) == 2 {
		//fmt.Println("new topic", msg[1].Content)
		sendNewTopicCard(*a.ctx, a.info.sessionId, a.info.msgId,
			completions.Content)
		return false
	}
	err = replyMsg(*a.ctx, completions.Content, a.info.msgId)
	if err != nil {
		replyMsg(*a.ctx, fmt.Sprintf(
			"🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.info.msgId)
		return false
	}
	return true
}

func sendOnProcess(a *ActionInfo) (*string, error) {
	// send 正在处理中
	cardId, err := sendOnProcessCard(*a.ctx, a.info.sessionId, a.info.msgId)
	if err != nil {
		return nil, err
	}
	return cardId, nil

}

func updateMsg(a *ActionInfo, cardId *string) bool {
	// update 正在处理中
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop() // 注意在函数结束时停止 ticker
	context := *a.ctx
	done := context.Done() // 获取 context 的取消信号

	count := 0     // 计数器
	maxCount := 15 // 最大循环次数

	msgStr := "demo"
	for {
		select {
		case <-done:
			// context 被取消，或者执行次数达到最大值，退出循环
			return false
		case <-ticker.C:
			msgStr = msgStr + randomWord() + " "
			updateTextCard(*a.ctx, msgStr, cardId)
			count++
			if count == maxCount {
				// 执行次数达到最大值，退出循环
				updateFinalCard(*a.ctx, msgStr, cardId)
				return false
			}
		}
	}
}
func randomWord() string {
	words := []string{"apple", "banana", "cherry", "orange", "pear"}
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子
	index := rand.Intn(len(words))   // 生成 0 到 len(words)-1 之间的随机数
	return words[index]
}
