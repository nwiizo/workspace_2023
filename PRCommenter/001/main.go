package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	// コマンドライン引数を解析する
	var (
		repo        = flag.String("repo", "", "Repository name")
		owner       = flag.String("owner", "", "Repository owner")
		prNumber    = flag.Int("pr", 0, "Pull request number")
		commentBody = flag.String("comment", "", "Comment body")
		token       = flag.String("token", "", "GitHub token")
	)
	flag.Parse()

	if *repo == "" || *owner == "" || *prNumber == 0 || *commentBody == "" || *token == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// OAuth2トークンを使用してGitHubクライアントを作成
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// PRにコメントを投稿
	comment := &github.IssueComment{Body: github.String(*commentBody)}
	_, _, err := client.Issues.CreateComment(ctx, *owner, *repo, *prNumber, comment)
	if err != nil {
		log.Fatalf("Error creating comment: %s", err)
	}

	log.Printf("Comment created successfully on PR %d", *prNumber)
}
