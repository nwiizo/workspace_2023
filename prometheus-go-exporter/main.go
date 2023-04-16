package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheusのメトリクスを定義しています。
// これらのメトリクスは、HTTPリクエストの情報や3-shake.comへのアクセス情報を収集するために使用されます。
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests processed",
		},
		[]string{"method", "path"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests",
			Buckets: prometheus.ExponentialBuckets(128, 2, 10),
		},
		[]string{"method", "path"},
	)

	httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses",
			Buckets: prometheus.ExponentialBuckets(128, 2, 10),
		},
		[]string{"method", "path"},
	)

	httpResponseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_response_time_seconds",
			Help: "Time of the last HTTP response",
		},
		[]string{"method", "path"},
	)

	externalAccessDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "external_access_duration_seconds",
			Help:    "Duration of external access to 3-shake.com",
			Buckets: prometheus.DefBuckets,
		},
	)

	lastExternalAccessStatusCode = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "last_external_access_status_code",
			Help: "Last status code of external access to 3-shake.com",
		},
	)
)

// init関数内で、メトリクスをPrometheusに登録しています。
func init() {
	registerMetrics()
}

// registerMetrics関数では、Prometheusにメトリクスを登録しています。
// これにより、Prometheusがメトリクスを収集できるようになります。
func registerMetrics() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestSize)
	prometheus.MustRegister(httpResponseSize)
	prometheus.MustRegister(httpResponseTime)
	prometheus.MustRegister(externalAccessDuration)
	prometheus.MustRegister(lastExternalAccessStatusCode)
}

// updateMetrics関数では、受信したHTTPリクエストのメトリクスを更新しています。
// これにより、各リクエストに関する情報が収集されます。
func updateMetrics(method, path string, requestSize, responseSize int, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, path).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	httpRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	httpResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
	httpResponseTime.WithLabelValues(method, path).Set(float64(time.Now().Unix()))
}

// prometheusMiddleware関数では、Echoのミドルウェアとして、受信したHTTPリクエストに関するメトリクスを更新する機能を追加しています。
func prometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()
		err := next(c)
		duration := time.Since(startTime)

		requestSize := c.Request().ContentLength
		responseSize := c.Response().Size

		updateMetrics(c.Request().Method, c.Path(), int(requestSize), int(responseSize), duration)

		return err
	}
}

// measureExternalAccess関数では、3-shake.comへの外部アクセスを定期的に計測し、そのアクセス時間とステータスコードをメトリクスに格納しています。
// この関数はメイン関数内で呼び出され、別のゴルーチンで実行されます。
func measureExternalAccess() {
	client := &http.Client{Timeout: 10 * time.Second}

	go func() {
		for {
			startTime := time.Now()
			resp, err := client.Get("https://3-shake.com")
			duration := time.Since(startTime)

			if err == nil {
				externalAccessDuration.Observe(duration.Seconds())
				lastExternalAccessStatusCode.Set(float64(resp.StatusCode))
				resp.Body.Close()
			}

			time.Sleep(1 * time.Minute)
		}
	}()
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware for Prometheus Exporter
	e.Use(prometheusMiddleware)

	// Enable request logger
	e.Use(middleware.Logger())

	e.GET("/3-shake-status", func(c echo.Context) error {
		status := lastExternalAccessStatusCode.Desc().String()
		return c.String(http.StatusOK, fmt.Sprintf("Last 3-shake.com access status: %s", status))
	})

	// Prometheus Exporter endpoint
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Measure external access to 3-shake.com
	measureExternalAccess()

	// Start the server
	e.Start(":2121")
}
