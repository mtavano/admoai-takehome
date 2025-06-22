# Testing Metrics Endpoint

## Overview

The application now includes a `/metrics` endpoint that provides Prometheus-formatted metrics using the official Prometheus client library for Go.

## Endpoint

**GET** `/metrics`

Returns metrics in Prometheus format using the official `promhttp.Handler()`.

## Example Response

```bash
curl -X GET http://localhost:9001/metrics
```

**Response:**
```
# HELP admoai_ads_created_total Total number of ads created
# TYPE admoai_ads_created_total counter
admoai_ads_created_total 5

# HELP admoai_ads_deactivated_total Total number of ads deactivated
# TYPE admoai_ads_deactivated_total counter
admoai_ads_deactivated_total 2

# HELP admoai_ads_active_current Current number of active ads
# TYPE admoai_ads_active_current gauge
admoai_ads_active_current 3

# HELP admoai_ads_inactive_current Current number of inactive ads
# TYPE admoai_ads_inactive_current gauge
admoai_ads_inactive_current 2

# HELP admoai_ads_expired_current Current number of expired ads
# TYPE admoai_ads_expired_current gauge
admoai_ads_expired_current 0

# HELP admoai_http_requests_total Total number of HTTP requests
# TYPE admoai_http_requests_total counter
admoai_http_requests_total{method="POST",endpoint="/v1/ads",status="201"} 5
admoai_http_requests_total{method="GET",endpoint="/v1/ads",status="200"} 3
admoai_http_requests_total{method="GET",endpoint="/metrics",status="200"} 2

# HELP admoai_http_request_duration_seconds Duration of HTTP requests
# TYPE admoai_http_request_duration_seconds histogram
admoai_http_request_duration_seconds_bucket{method="POST",endpoint="/v1/ads",le="0.005"} 2
admoai_http_request_duration_seconds_bucket{method="POST",endpoint="/v1/ads",le="0.01"} 3
admoai_http_request_duration_seconds_bucket{method="POST",endpoint="/v1/ads",le="0.025"} 5
admoai_http_request_duration_seconds_bucket{method="POST",endpoint="/v1/ads",le="+Inf"} 5
admoai_http_request_duration_seconds_sum{method="POST",endpoint="/v1/ads"} 0.075
admoai_http_request_duration_seconds_count{method="POST",endpoint="/v1/ads"} 5

# HELP admoai_uptime_seconds Total uptime in seconds
# TYPE admoai_uptime_seconds counter
admoai_uptime_seconds 3600
```

## Available Metrics

### Counters
- **admoai_ads_created_total**: Total number of ads created since startup
- **admoai_ads_deactivated_total**: Total number of ads deactivated since startup
- **admoai_http_requests_total**: Total HTTP requests per method, endpoint, and status
- **admoai_uptime_seconds**: Application uptime in seconds

### Gauges
- **admoai_ads_active_current**: Current number of active ads
- **admoai_ads_inactive_current**: Current number of inactive ads
- **admoai_ads_expired_current**: Current number of expired ads

### Histograms
- **admoai_http_request_duration_seconds**: Request duration distribution with buckets

## Alert System

### Active Ads Alert
The system includes a simulated alert that triggers when more than 10 active ads are detected.

**Alert Condition:**
- When `admoai_ads_active_current > 10`

**Alert Output:**
```
⚠️  ALERT: High number of active ads detected: 12 (threshold: 10)
```

### Testing the Alert

1. Create more than 10 ads:
```bash
for i in {1..12}; do
  curl -X POST http://localhost:9001/v1/ads \
    -H "Content-Type: application/json" \
    -d "{
      \"title\": \"Ad $i\",
      \"image_url\": \"https://example.com/image$i.jpg\",
      \"placement\": \"homepage\",
      \"ttl\": 30
    }"
done
```

2. Check metrics to see the alert:
```bash
curl -X GET http://localhost:9001/metrics
```

3. Look for the alert message in the application logs.

## Implementation Details

### Using Official Prometheus Client
The application now uses the official Prometheus client library:
- `github.com/prometheus/client_golang/prometheus`
- `github.com/prometheus/client_golang/prometheus/promhttp`

### Benefits of Official Library
- ✅ Standard Prometheus format
- ✅ Thread-safe metrics collection
- ✅ Built-in histogram buckets
- ✅ Automatic metric registration
- ✅ Production-ready and well-tested
- ✅ Rich ecosystem and documentation

### Metrics Collection
- **Counters**: For cumulative values (ads created, requests)
- **Gauges**: For current values (active ads count)
- **Histograms**: For request duration distribution
- **Labels**: Method, endpoint, and status for HTTP metrics

## Prometheus Integration

### Configuration
Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'admoai-api'
    static_configs:
      - targets: ['localhost:9001']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboard
You can create a Grafana dashboard using these metrics:

- **Ads Overview**: Show current active/inactive/expired counts
- **Request Rate**: Monitor HTTP requests per second
- **Response Time**: Track request duration percentiles
- **Uptime**: Monitor application availability

## Monitoring Queries

### Active Ads Alert Rule
```promql
admoai_ads_active_current > 10
```

### Request Rate
```promql
rate(admoai_http_requests_total[5m])
```

### 95th Percentile Response Time
```promql
histogram_quantile(0.95, rate(admoai_http_request_duration_seconds_bucket[5m]))
```

### Error Rate
```promql
rate(admoai_http_requests_total{status=~"4..|5.."}[5m])
```

### Uptime
```promql
admoai_uptime_seconds
```

## Advantages Over Custom Implementation

1. **Standard Compliance**: Uses official Prometheus format
2. **Performance**: Optimized for high-frequency metrics collection
3. **Reliability**: Well-tested in production environments
4. **Features**: Built-in histogram buckets, labels, and metric types
5. **Maintenance**: No custom code to maintain
6. **Ecosystem**: Compatible with all Prometheus tools and dashboards 