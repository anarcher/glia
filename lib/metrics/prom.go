package metrics

import (
	"github.com/go-kit/kit/metrics"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"time"
)

var (
	Fetching = kitprometheus.NewGauge(stdprometheus.GaugeOpts{
		Namespace: "glia",
		Subsystem: "fetcher",
		Name:      "fetching",
		Help:      "Number of fetching operations waiting to be processed.",
	}, []string{})
	Sending = kitprometheus.NewGauge(stdprometheus.GaugeOpts{
		Namespace: "glia",
		Subsystem: "sender",
		Name:      "sending",
		Help:      "Number of sending operations waiting to be processed.",
	}, []string{})

	FetchLatency = metrics.NewTimeHistogram(time.Microsecond, kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "glia",
		Subsystem: "fetcher",
		Name:      "fetch_latency_microseconds",
		Help:      "Total duration of fetching in microseconds.",
	}, []string{}))

	SendLatency = metrics.NewTimeHistogram(time.Microsecond, kitprometheus.NewSummary(stdprometheus.SummaryOpts{
		Namespace: "glia",
		Subsystem: "sender",
		Name:      "send_latency_microseconds",
		Help:      "Total duration of sending in microseconds.",
	}, []string{}))
)
