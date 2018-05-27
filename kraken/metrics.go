package kraken

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
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
	prometheus.MustRegister(krakenDurations)
	prometheus.MustRegister(krakenErrors)
	prometheus.MustRegister(krakenInFlight)
}

type krakenObserver struct {
}

func (o krakenObserver) StartRequest(api string) requestObserver {
	krakenInFlight.Inc()
	return requestObserver{
		api:   api,
		begin: time.Now(),
	}
}

type requestObserver struct {
	api   string
	begin time.Time
}

func (o requestObserver) Finish() {
	krakenDurations.With(prometheus.Labels{"api": o.api}).
		Observe(time.Since(o.begin).Seconds())
	krakenInFlight.Dec()
}

func (o krakenObserver) OnError(api string, err error) {
	krakenErrors.With(prometheus.Labels{"api": api}).
		Inc()
}
