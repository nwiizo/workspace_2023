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
		PerPage: 10,
	}
	files, _, err := client.PullRequests.ListFiles(ctx, *owner, *repo, *prNumber, opts)
	if err != nil {
		log.Fatalf("Error listing files for pull request: %v", err)
	}

	// 変更されたファイルとその差分をコメントとして整理
	var changesBuilder strings.Builder
	changesBuilder.WriteString("変更の詳細:\n\n")

	for _, file := range files {
		changesBuilder.WriteString(fmt.Sprintf("### %s\n", *file.Filename))
		changesBuilder.WriteString("\n```diff\n")             // diffをコードブロックとしてマークダウンに表示
		changesBuilder.WriteString(truncateDiff(*file.Patch)) // 大きな差分を適切に取り扱うため、差分内容をある程度まで切り詰める
		changesBuilder.WriteString("\n```\n")
	}

	commentBody := changesBuilder.String()
	comment := &github.IssueComment{Body: &commentBody}

	// コメントを投稿
	_, _, err = client.Issues.CreateComment(ctx, *owner, *repo, *prNumber, comment)
	if err != nil {
		log.Fatalf("Error creating comment: %v", err)
	}

	fmt.Printf("Comment created successfully on PR %d\n", *prNumber)
}

// truncateDiff は、表示する差分のサズが大きすぎないように、差分のテキストを切り詰めます。
// これは、コメントに表示する情報量を制御するための単純な方法です。
func truncateDiff(diff string) string {
	const maxDiffSize = 4000 // GitHubのコメントはある程度の文字数に制限があるため、それを超えないようにする
	if len(diff) > maxDiffSize {
		return diff[:maxDiffSize] + "\n...（差分は省略されました）"
	}
	return diff
}
