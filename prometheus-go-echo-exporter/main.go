package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	prom "github.com/prometheus/client_golang/prometheus"
)

// Prometheus のメトリクスを定義しています。
// これらのメトリクスは、3-shake.com への外部アクセスの情報を収集するために使用されます。
var (
	externalAccessDuration = prom.NewHistogram(
		prom.HistogramOpts{
			Name:    "external_access_duration_seconds",
			Help:    "Duration of external access to 3-shake.com",
			Buckets: prom.DefBuckets,
		},
	)

	lastExternalAccessStatusCode = prom.NewGauge(
		prom.GaugeOpts{
			Name: "last_external_access_status_code",
			Help: "Last status code of external access to 3-shake.com",
		},
	)
)

// init 関数内で、メトリクスを Prometheus に登録しています。
func init() {
	prom.MustRegister(externalAccessDuration)
	prom.MustRegister(lastExternalAccessStatusCode)
}

// 3-shake.com の外部アクセスを計測するミドルウェアを作成します。
func measureExternalAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// HTTP クライアントを作成し、タイムアウトを 10 秒に設定します。
		client := &http.Client{Timeout: 10 * time.Second}
		// 現在の時刻を取得し、アクセス開始時間として保持します。
		startTime := time.Now()
		// 3-shake.com に対して HTTP GET リクエストを送信します。
		resp, err := client.Get("https://3-shake.com")
		// アクセス開始時間から現在の時刻までの経過時間を計算し、duration に格納します。
		duration := time.Since(startTime)
		// エラーが発生しない場合（リクエストが成功した場合）
		if err == nil {
			// アクセス時間（duration）をヒストグラムメトリクスに追加します。
			externalAccessDuration.Observe(duration.Seconds())
			// ステータスコードをゲージメトリクスに設定します。
			lastExternalAccessStatusCode.Set(float64(resp.StatusCode))
			// レスポンスのボディを閉じます。
			resp.Body.Close()
		}
		// 次のミドルウェアまたはハンドラ関数に処理を移します。
		return next(c)
	}
}

func unstableEndpoint(c echo.Context) error {
	// 0 から 4 までのランダムな整数を生成します。
	randomNumber := rand.Intn(5)

	// 生成された整数が 4 の場合、HTTP ステータスコード 500 を返します。
	if randomNumber == 4 {
		return c.String(http.StatusInternalServerError, "Something went wrong!")
	}

	// それ以外の場合、HTTP ステータスコード 200 を返します。
	return c.String(http.StatusOK, "Success!")
}

func main() {
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Prometheus ミドルウェアを有効にします。
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// 3-shake.com への外部アクセスを計測するミドルウェアを追加します。
	e.Use(measureExternalAccess)

	// ルートのエンドポイントを設定します。
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// /unstable エンドポイントを設定します。
	// 20% の確率で HTTP ステータスコード 500 を返します。
	e.GET("/unstable", unstableEndpoint)

	// サーバーを開始します。
	e.Start(":2121")
}
