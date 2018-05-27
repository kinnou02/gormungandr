package gormungandr

import (
	"context"
	"database/sql"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/gchaincl/sqlhooks"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
)

type contextKey string

var (
	sqlBeginKey contextKey = "begin"

	httpDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gormungandr",
		Subsystem: "http",
		Name:      "durations_seconds",
		Help:      "http request latency distributions.",
		Buckets:   prometheus.ExponentialBuckets(0.001, 1.5, 25),
	},
		[]string{"handler", "code"},
	)

	httpInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gormungandr",
		Subsystem: "http",
		Name:      "in_flight",
		Help:      "current number of http request being served",
	},
	)

	krakenDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gormungandr",
		Subsystem: "kraken",
		Name:      "durations_seconds",
		Help:      "kraken request latency distributions.",
		Buckets:   prometheus.ExponentialBuckets(0.001, 1.5, 25),
	},
		[]string{"api"},
	)

	krakenErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gormungandr",
		Subsystem: "kraken",
		Name:      "errors_count",
		Help:      "kraken request errors count",
	},
		[]string{"api"},
	)
	krakenInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gormungandr",
		Subsystem: "kraken",
		Name:      "in_flight",
		Help:      "current number of request being called",
	},
	)
	sqlDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "gormungandr",
		Subsystem: "sql",
		Name:      "durations_seconds",
		Help:      "sql request latency distributions.",
		Buckets:   prometheus.ExponentialBuckets(0.001, 1.5, 15),
	},
	)
)

func init() {
	prometheus.MustRegister(httpDurations)
	prometheus.MustRegister(krakenDurations)
	prometheus.MustRegister(krakenErrors)
	prometheus.MustRegister(krakenInFlight)
	prometheus.MustRegister(httpInFlight)
	prometheus.MustRegister(sqlDurations)
	sql.Register("postgresInstrumented", sqlhooks.Wrap(&pq.Driver{}, &instrumentHook{}))
}

func Instrument(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		begin := time.Now()
		httpInFlight.Inc()
		if err := next(c); err != nil {
			c.Error(err)
		}
		httpInFlight.Dec()
		observer := httpDurations.With(prometheus.Labels{"handler": runtime.FuncForPC(reflect.ValueOf(c.Handler()).Pointer()).Name(), "code": strconv.Itoa(c.Response().Status)})
		observer.Observe(time.Since(begin).Seconds())
		return nil
	}
}

type RequestObserver struct {
	api   string
	begin time.Time
}

func (o RequestObserver) Finish() {
	krakenDurations.With(prometheus.Labels{"api": o.api}).
		Observe(time.Since(o.begin).Seconds())
	krakenInFlight.Dec()
}

func (o KrakenObserver) OnError(api string, err error) {
	krakenErrors.With(prometheus.Labels{"api": api}).
		Inc()
}

type KrakenObserver struct {
}

func (o KrakenObserver) StartRequest(kraken *Kraken, api string) RequestObserver {
	krakenInFlight.Inc()
	return RequestObserver{
		api:   api,
		begin: time.Now(),
	}
}

// Hooks satisfies the sqlhook.Hooks interface
type instrumentHook struct{}

// register the timestamp in the context
func (h *instrumentHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, sqlBeginKey, time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and feed prometheus
func (h *instrumentHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value(sqlBeginKey).(time.Time)
	sqlDurations.Observe(time.Since(begin).Seconds())
	return ctx, nil
}
