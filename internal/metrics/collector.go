package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Collector struct {
	// Ad metrics
	adsCreatedTotal     prometheus.Counter
	adsDeactivatedTotal prometheus.Counter
	adsActiveCurrent    prometheus.Gauge
	adsInactiveCurrent  prometheus.Gauge
	adsExpiredCurrent   prometheus.Gauge

	// HTTP metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec

	// System metrics
	uptime prometheus.Counter
}

var (
	collector *Collector
	startTime = time.Now()
)

// Init initializes the metrics collector
func Init() {
	collector = &Collector{
		// Ad counters
		adsCreatedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "admoai_ads_created_total",
			Help: "Total number of ads created",
		}),
		adsDeactivatedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "admoai_ads_deactivated_total",
			Help: "Total number of ads deactivated",
		}),

		// Ad gauges
		adsActiveCurrent: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "admoai_ads_active_current",
			Help: "Current number of active ads",
		}),
		adsInactiveCurrent: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "admoai_ads_inactive_current",
			Help: "Current number of inactive ads",
		}),
		adsExpiredCurrent: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "admoai_ads_expired_current",
			Help: "Current number of expired ads",
		}),

		// HTTP metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "admoai_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "admoai_http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),

		// System metrics
		uptime: promauto.NewCounter(prometheus.CounterOpts{
			Name: "admoai_uptime_seconds",
			Help: "Total uptime in seconds",
		}),
	}

	// Start uptime counter
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			collector.uptime.Add(1)
		}
	}()
}

// GetCollector returns the singleton collector instance
func GetCollector() *Collector {
	return collector
}

// IncrementAdCreated increments the total ads created counter
func (c *Collector) IncrementAdCreated() {
	c.adsCreatedTotal.Inc()
}

// IncrementAdDeactivated increments the total ads deactivated counter
func (c *Collector) IncrementAdDeactivated() {
	c.adsDeactivatedTotal.Inc()
}

// UpdateAdCounts updates the current ad counts
func (c *Collector) UpdateAdCounts(active, inactive, expired int64) {
	c.adsActiveCurrent.Set(float64(active))
	c.adsInactiveCurrent.Set(float64(inactive))
	c.adsExpiredCurrent.Set(float64(expired))

	// Check for alerts
	c.checkAlerts(active)
}

// RecordHTTPRequest records an HTTP request
func (c *Collector) RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	c.httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	c.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// checkAlerts checks for alert conditions
func (c *Collector) checkAlerts(activeCount int64) {
	if activeCount > 10 {
		fmt.Printf("⚠️  ALERT: High number of active ads detected: %d (threshold: 10)\n", activeCount)
	}
} 