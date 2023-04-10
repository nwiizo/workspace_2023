# HTTP Ping Tool

⚠ 注意 現在ではこちらのツールを利用することはできません

HTTP Pingツールは、Go言語で作成されたWebAssemblyアプリケーションです。このツールは、指定されたURLに対してHTTP GETリクエストを行い、レイテンシとHTTPステータスコードを取得します。また、エラーが発生した場合は、エラーメッセージも表示します。このツールは、ブラウザで実行されるシンプルなWeb UIを提供します。

## Requirements

- Go 1.11 以降
- モダンなウェブブラウザ（WebAssemblyに対応している必要があります）

## Usage

### 1. WebAssemblyバイナリをビルドする

次のコマンドを実行して、GoのコードをWebAssemblyバイナリにコンパイルします。


```
GOOS=js GOARCH=wasm go build -o http_ping.wasm http_ping.go
```


### 2. 簡易Webサーバーを立ち上げる

`main.go` を使って、簡易Webサーバーを立ち上げます。

```
go run main.go
```


### 3. ブラウザでアプリケーションにアクセスする

ブラウザで `http://localhost:8181` にアクセスして、HTTP Pingツールを使用します。

## License

[MIT License](LICENSE)

