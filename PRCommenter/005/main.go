package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2"
)

func main() {
	// コマンドライン引数の解析
	var (
		repo      = flag.String("repo", "", "Repository name")
		owner     = flag.String("owner", "", "Repository owner")
		prNumber  = flag.Int("pr", 0, "Pull request number")
		token     = flag.String("token", "", "GitHub token")
		openaiKey = flag.String("openai-key", "", "OpenAI API Key")
		model     = flag.String("model", "gpt-3.5-turbo", "OpenAI model ID")
	)
	flag.Parse()

	if *repo == "" || *owner == "" || *prNumber == 0 || *token == "" || *openaiKey == "" {
		fmt.Println("All flags are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// GitHubクライアントの設定
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	// プルリクエストのフルの変更点を取得
	opts := &github.ListOptions{
		PerPage: 10,
	}
	files, _, err := githubClient.PullRequests.ListFiles(ctx, *owner, *repo, *prNumber, opts)
	if err != nil {
		log.Fatalf("Error listing files for pull request: %v", err)
	}

	// OpenAIのクライアントを作成
	openaiClient := openai.NewClient(*openaiKey)

	for _, file := range files {
		if file.Patch == nil {
			continue
		}
		fmt.Printf("Reviewing changes in file: %s\n", *file.Filename)

		// 変更内容を取得
		fullOutput := fmt.Sprintf("```diff\n%s\n```", *file.Patch)

		// OpenAIでチャットコンプリーションを作成
		resp, err := openaiClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: *model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: fullOutput,
					},
				},
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}

		// GPT-3からのレスポンスをGitHubのコメントとして整形
		comment := fmt.Sprintf(
			"GPT-3 code review for file %s:\n\n%s",
			*file.Filename,
			resp.Choices[0].Message.Content,
		)
		commentBody := &github.IssueComment{Body: &comment}

		// コメントをプルリクエストに投稿
		_, _, err = githubClient.Issues.CreateComment(ctx, *owner, *repo, *prNumber, commentBody)
		if err != nil {
			log.Printf("Error posting comment to GitHub: %v", err)
		} else {
			fmt.Printf("Posted GPT-3 review comment for file %s\n", *file.Filename)
		}
	}
}
