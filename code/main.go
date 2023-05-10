package main

import (
	"context"
	"encoding/json"
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"start-feishubot/handlers"
	"start-feishubot/initialization"
	"start-feishubot/services/openai"
	"start-feishubot/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"

	sdkginext "github.com/larksuite/oapi-sdk-gin"

	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func main() {
	initialization.InitRoleList()
	pflag.Parse()
	globalConfig := initialization.GetConfig()

	// 打印一下实际读取到的配置
	globalConfigPrettyString, _ := json.MarshalIndent(globalConfig, "", "    ")
	log.Println(string(globalConfigPrettyString))

	initialization.LoadLarkClient(*globalConfig)
	gpt := openai.NewChatGPT(*globalConfig)
	handlers.InitHandlers(gpt, *globalConfig)

	if globalConfig.EnableLog {
		logger := enableLog()
		defer utils.CloseLogger(logger)
	}

	eventHandler := dispatcher.NewEventDispatcher(
		globalConfig.FeishuAppVerificationToken, globalConfig.FeishuAppEncryptKey).
		OnP2MessageReceiveV1(handlers.Handler).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			return handlers.ReadHandler(ctx, event)
		})

	cardHandler := larkcard.NewCardActionHandler(
		globalConfig.FeishuAppVerificationToken, globalConfig.FeishuAppEncryptKey,
		handlers.CardHandler())

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/webhook/event",
		sdkginext.NewEventHandlerFunc(eventHandler))
	r.POST("/webhook/card",
		sdkginext.NewCardActionHandlerFunc(
			cardHandler))

	err := initialization.StartServer(*globalConfig, r)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

func enableLog() *lumberjack.Logger {
	// Set up the logger
	var logger *lumberjack.Logger

	logger = &lumberjack.Logger{
		Filename: "logs/app.log",
		MaxSize:  100,      // megabytes
		MaxAge:   365 * 10, // days
	}

	fmt.Printf("logger %T\n", logger)

	// Set up the logger to write to both file and console
	log.SetOutput(io.MultiWriter(logger, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	// Write some log messages
	log.Println("Starting application...")

	return logger
}
