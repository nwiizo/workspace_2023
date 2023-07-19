package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "An app that executes a shell command and sends the result to OpenAI GPT",
}

var language string

var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute a shell command and send the result to OpenAI GPT",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if language != "en" && language != "ja" {
			fmt.Println(
				"Invalid language. Please select either 'en' for English or 'ja' for Japanese.",
			)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// コマンドの実行
		shellCmd := exec.Command(args[0], args[1:]...)
		var out bytes.Buffer
		shellCmd.Stdout = &out
		err := shellCmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return
		}

		// 実行したコマンドとその結果を結合
		fullOutput := ""
		if language == "en" {
			fullOutput = fmt.Sprintf(
				"Command executed: %v\nResult:\n%v\nCan you explain this result?",
				shellCmd.String(),
				out.String(),
			)
		} else if language == "ja" {
			fullOutput = fmt.Sprintf("実行したコマンド: %v\n結果:\n%v\nこの結果について説明していただけますか？", shellCmd.String(), out.String())
		}

		// APIキーの取得
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("Error: OPENAI_API_KEY is not set")
			return
		}

		// 実行したコマンドとその結果を表示
		color.New(color.FgCyan).Printf("Command executed: %v\n", shellCmd.String())
		color.New(color.FgGreen).Printf("Result:\n%v\n\n", out.String())

		// スピナーの作成
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		// OpenAIのクライアントの作成
		client := openai.NewClient(apiKey)

		// ChatGPTへのリクエストの作成
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: fullOutput,
					},
				},
			},
		)

		s.Stop()

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}

		// レスポンスの表示
		fmt.Println(resp.Choices[0].Message.Content)
	},
}

func main() {
	executeCmd.Flags().
		StringVarP(&language, "language", "l", "en", "Language for the command execution (en/ja)")
	rootCmd.AddCommand(executeCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
