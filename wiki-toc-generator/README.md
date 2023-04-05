# GitLab Wiki Table of Contents Generator

GitLabのWikiリポジトリをクローンし、Markdown形式の目次を生成するGo製のコマンドラインツールです。

## インストール方法

リポジトリをクローンし、ビルドしてください。

```sh
git clone https://github.com/nwiizo/wiki-toc-generator.git
cd wiki-toc-generator
go build -o wiki-toc-generator main.go
```

## 使い方

以下のコマンドで目次を生成できます。
```sh
./wiki-toc-generator <wikiRepoURL>
```

<wikiRepoURL>には、GitLab WikiリポジトリのURLを指定してください。

### 例

```sh
./wiki-toc-generator https://gitlab.com/yourusername/yourproject.wiki.git
```

目次が生成されたら、必要に応じてWikiページにコピーしてください。
