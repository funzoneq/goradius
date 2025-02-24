package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	authRejectedRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goradius_auth_rejected_requests",
			Help: "The total number of rejected RADIUS authentication requests",
		},
		[]string{"identifier", "site"},
	)
	authAcceptedRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goradius_auth_accepted_requests",
			Help: "The total number of accepted RADIUS authentication requests",
		},
		[]string{"identifier", "site"},
	)
	acctRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goradius_acct_total_requests",
			Help: "The total number of accepted RADIUS authentication requests",
		},
		[]string{"identifier", "site"},
	)
	authDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "goradius_auth_duration_seconds",
			Help:    "Histogram of response time of authentication handler in seconds",
			Buckets: []float64{0, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 30, 60, 120, 300},
		},
		[]string{"auth_status"},
	)
	acctDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "goradius_acct_duration_seconds",
			Help:    "Histogram of response time of accounting handler in seconds",
			Buckets: []float64{0, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 30, 60, 120, 300},
		},
	)
)

func MetricsEndpoint() {
	log.Printf("Starting metrics server on %s", Config.MetricsListenAddress)

	// Expose the registered metrics via HTTP.
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(Config.MetricsListenAddress, nil))
}
