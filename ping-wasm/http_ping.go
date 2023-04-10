package main

import (
	"encoding/json"
	"net/http"
	"syscall/js"
	"time"
)

type PingResult struct {
	URL        string `json:"url"`
	Latency    int64  `json:"latency"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func ping(this js.Value, args []js.Value) interface{} {
	url := args[0].String()
	callback := args[1]

	start := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(start).Milliseconds()

	result := PingResult{
		URL:     url,
		Latency: latency,
	}

	if err != nil {
		result.Error = err.Error()
	} else {
		result.StatusCode = resp.StatusCode
		resp.Body.Close()
	}

	resultJSON, _ := json.Marshal(result)
	callback.Invoke(string(resultJSON))
	return nil
}

func main() {
	js.Global().Set("ping", js.FuncOf(ping))
	select {} // 無限ループでGoプログラムを終了させない
}
