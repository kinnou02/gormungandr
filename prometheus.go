package gormungandr

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
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
)

func init() {
	prometheus.MustRegister(httpDurations)
	prometheus.MustRegister(krakenDurations)
	prometheus.MustRegister(krakenErrors)
	prometheus.MustRegister(krakenInFlight)
	prometheus.MustRegister(httpInFlight)
}

func InstrumentGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		begin := time.Now()
		httpInFlight.Inc()
		c.Next()
		httpInFlight.Dec()
		observer := httpDurations.With(prometheus.Labels{"handler": c.HandlerName(), "code": strconv.Itoa(c.Writer.Status())})
		observer.Observe(time.Since(begin).Seconds())
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
