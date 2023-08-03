// package prometheus is used to define and register Prometheus metrics for an HTTP server.
package prometheus

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/os"
	tstrings "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

// ServerMetrics is a struct that holds all server metrics for prometheus.
// It defines Prometheus counters and histograms to measure server and RPC attributes, like a total number of server startups and number of requests completed.
type ServerMetrics struct {
	serverStartedCounter   *prometheus.CounterVec
	serverHandledCounter   *prometheus.CounterVec
	serverMsgReceived      *prometheus.CounterVec
	serverMsgSent          *prometheus.CounterVec
	serverHandledHistogram *prometheus.HistogramVec
}

// NewServerMetrics is a function that initializes a new ServerMetrics type with Prometheus counter and histogram objects.
// It returns a pointer to the new ServerMetrics object.
// This function uses the fx.Lifecycle package to register the metrics with Prometheus.
// This function takes in two arguments: a version.Version object and an fx.Lifecycle object.
// Inside the function, this function initializes Prometheus counters and histograms to measure server and RPC attributes, like the total number of server startups and number of requests completed.
// These metrics are initialized with Prometheus labels that include the executable name and version number of the server.
// Finally, the function adds the metrics to the Prometheus registry.
func NewServerMetrics(lc fx.Lifecycle, version version.Version) *ServerMetrics {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ServerMetrics{
		serverStartedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_started_total",
				Help:        "Total number of RPCs started on the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		serverHandledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_handled_total",
				Help:        "Total number of RPCs completed on the server, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method", "http_code"}),
		serverMsgReceived: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_msg_received_total",
				Help:        "Total number of RPC messages received on the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		serverMsgSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "http_server_msg_sent_total",
				Help:        "Total number of RPC messages sent by the server.",
				ConstLabels: labels,
			}, []string{"http_service", "http_method"}),
		serverHandledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_server_handling_seconds",
				Help:        "Histogram of response latency (seconds) of HTTP that had been application-level handled by the server.",
				Buckets:     prometheus.DefBuckets,
				ConstLabels: labels,
			}, []string{"http_service", "http_method"},
		),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(metrics)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(metrics)

			return nil
		},
	})

	return metrics
}

// Describe is a method for the ServerMetrics struct that sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once the last descriptor has been sent.
func (m *ServerMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.serverStartedCounter.Describe(ch)
	m.serverHandledCounter.Describe(ch)
	m.serverMsgReceived.Describe(ch)
	m.serverMsgSent.Describe(ch)
	m.serverHandledHistogram.Describe(ch)
}

// Collect is a method for the ServerMetrics struct that is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the provided channel and returns once the last metric has been sent.
func (m *ServerMetrics) Collect(ch chan<- prometheus.Metric) {
	m.serverStartedCounter.Collect(ch)
	m.serverHandledCounter.Collect(ch)
	m.serverMsgReceived.Collect(ch)
	m.serverMsgSent.Collect(ch)
	m.serverHandledHistogram.Collect(ch)
}

// Handler is a method for the ServerMetrics struct that returns a new http.Handler that wraps the provided
// http.Handler and collects Prometheus metrics.
func (m *ServerMetrics) Handler(h http.Handler) http.Handler {
	return &handler{metrics: m, Handler: h}
}

type handler struct {
	metrics *ServerMetrics
	http.Handler
}

// ServeHTTP is a method for the handler struct that is called when an HTTP request is received.
// The function maps the HTTP request path and method to metrics labels and initializes a new serverReporter to report metrics for this request.
// It then initializes a new responseWriter to collect HTTP response status codes.
// The function uses the handler's ServeHTTP function to handle the request.
// It then collects the final response status code from the responseWriter, reports metrics for the request, and sends the response to the client.
func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Use the HTTP request path and method to generate Prometheus metrics
	service, method := req.URL.Path, strings.ToLower(req.Method)
	// Ignore health check endpoints
	if tstrings.IsHealth(service) {
		h.Handler.ServeHTTP(resp, req)

		return
	}

	monitor := newServerReporter(h.metrics, service, method)
	monitor.ReceivedMessage()

	res := &responseWriter{ResponseWriter: resp, Status: http.StatusOK}
	h.Handler.ServeHTTP(res, req)

	monitor.Handled(res.Status)

	// Send message metrics only for successful requests (HTTP status codes 2xx)
	if res.Status >= 200 && res.Status <= 299 {
		monitor.SentMessage()
	}
}


type responseWriter struct {
	http.ResponseWriter
	Status int
}

// WriteHeader is a method for the responseWriter struct that is called when an HTTP response is sent.
// The function sets response status codes and updates the responseWriter status code.
func (r *responseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type serverReporter struct {
	metrics     *ServerMetrics
	serviceName string
	methodName  string
	startTime   time.Time
}

// newServerReporter is a method for the serverReporter struct that initializes a new serverReporter object with Prometheus labels for a specific HTTP request.
func newServerReporter(m *ServerMetrics, service, method string) *serverReporter {
	r :=   &serverReporter{metrics: m, startTime: time.Now(), serviceName: service, methodName: method}
	r.metrics.serverStartedCounter.WithLabelValues(r.serviceName, r.methodName).Inc()

	return r
}

// ReceivedMessage is a method for the serverReporter struct that reports metrics for an HTTP request message received.
func (r *serverReporter) ReceivedMessage() {
	r.metrics.serverMsgReceived.WithLabelValues(r.serviceName, r.methodName).Inc()
}

// SentMessage is a method for the serverReporter struct that reports metrics for an HTTP request message sent.
func (r *serverReporter) SentMessage() {
	r.metrics.serverMsgSent.WithLabelValues(r.serviceName, r.methodName).Inc()
}

// Handled is a method for the serverReporter struct that reports metrics for an HTTP request when it has been handled.
func (r *serverReporter) Handled(code int) {
	r.metrics.serverHandledCounter.WithLabelValues(r.serviceName, r.methodName, strconv.Itoa(code)).Inc()
	r.metrics.serverHandledHistogram.WithLabelValues(r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
}