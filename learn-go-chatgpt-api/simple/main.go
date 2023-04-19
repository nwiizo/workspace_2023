package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("OPENAI_API_KEY")

	// 新しいOpenAIクライアントを作成
	client := openai.NewClient(apiKey)

	// メッセージ配列を作成
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Hello world",
		},
	}

	// GPT-3.5-turboモデルを使ってチャットの会話を生成するリクエストを作成
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    messages,
			Temperature: 0.9,
			MaxTokens:   200,
		},
	)
	// エラーチェック
	if err != nil {
		// エラーがあれば表示して終了
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	// 正常に会話が生成された場合、最初の選択肢の内容を出力
	fmt.Println(resp.Choices[0].Message.Content)
}
