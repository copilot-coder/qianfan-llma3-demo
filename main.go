package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/copilot-coder/qianfan-llma3-demo/chatengine"
)

func main() {
	cfg := chatengine.Config{
		AccessKey: os.Getenv("QIANFAN_ACCESS_KEY"),
		SecretKey: os.Getenv("QIANFAN_SECRET_KEY"),
	}
	engine := chatengine.NewEngine(cfg)

	req := chatengine.ChatReq{
		Model: "llama_3_8b",
		Messages: []chatengine.Message{
			{Role: "user", Content: "北京有什么好玩的景点？用中文回答"},
		},
		MaxOutputTokens: 1000,
	}

	// stream request
	fmt.Println(">>> test for stream request")
	ch, err := engine.StreamRequest(context.TODO(), req)
	if err != nil {
		log.Panic(err)
	}
	for resp := range ch {
		if resp.Err != nil {
			log.Panic("recv fail.", resp.Err)
		}
		fmt.Print(resp.Result)
	}

	// non-stream request
	fmt.Println("\n>>> test for non-stream request")
	req.Stream = false

	resp, err := engine.ChatRequest(context.TODO(), req)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(resp.Result)
}
