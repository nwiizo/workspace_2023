package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	// コマンドライン引数の解析
	var (
		repo     = flag.String("repo", "", "Repository name")
		owner    = flag.String("owner", "", "Repository owner")
		prNumber = flag.Int("pr", 0, "Pull request number")
		token    = flag.String("token", "", "GitHub token")
	)
	flag.Parse()

	if *repo == "" || *owner == "" || *prNumber == 0 || *token == "" {
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
	client := github.NewClient(tc)

	// PRの内容を取得
	pr, _, err := client.PullRequests.Get(ctx, *owner, *repo, *prNumber)
	if err != nil {
		log.Fatalf("Error getting pull request: %v", err)
	}

	// nilチェック
	if pr == nil || pr.Body == nil {
		log.Fatalf("Pull request not found or has no body.")
	}

	// PRの本文からコメントを作成
	commentBody := fmt.Sprintf("PR内容:\n%s", *pr.Body)
	comment := &github.IssueComment{Body: &commentBody}

	// コメントを投稿
	_, _, err = client.Issues.CreateComment(ctx, *owner, *repo, *prNumber, comment)
	if err != nil {
		log.Fatalf("Error creating comment: %v", err)
	}

	fmt.Printf("Comment created successfully on PR %d\n", *prNumber)
}
