package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	// プルリクエストのファイルの変更点を取得
	opts := &github.ListOptions{
		PerPage: 10, // 一度に取得するファイルの数を制限する（必要に応じて調整）
	}
	files, _, err := client.PullRequests.ListFiles(ctx, *owner, *repo, *prNumber, opts)
	if err != nil {
		log.Fatalf("Error listing files for pull request: %v", err)
	}

	// 変更されたファイルのリストをマークダウン形式で整形
	var fileListBuilder strings.Builder
	fileListBuilder.WriteString("変更されたファイル:\n\n") // マークダウンにおける改行は、行の終わりに2つのスペースを入れるか、空の行を挟む

	for _, file := range files {
		// ファイル名をバッククォートで囲むことで、マークダウン内でコードとして表示されるようにします。
		fileListBuilder.WriteString(fmt.Sprintf("- `%s`\n", *file.Filename)) // マークダウンのリスト形式で出力
	}

	commentBody := fileListBuilder.String()
	comment := &github.IssueComment{Body: &commentBody}

	// コメントを投稿
	_, _, err = client.Issues.CreateComment(ctx, *owner, *repo, *prNumber, comment)
	if err != nil {
		log.Fatalf("Error creating comment: %v", err)
	}

	fmt.Printf("Comment created successfully on PR %d\n", *prNumber)
}
