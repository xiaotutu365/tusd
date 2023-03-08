package cli

import (
	"net/http"

	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/prometheuscollector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MetricsOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "tusd_connections_open",
	Help: "Current number of open connections.",
})

var MetricsHookErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tusd_hook_errors_total",
		Help: "Total number of execution errors per hook type.",
	},
	[]string{"hooktype"},
)

var MetricsHookInvocationsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "tusd_hook_invocations_total",
		Help: "Total number of invocations per hook type.",
	},
	[]string{"hooktype"},
)

func SetupMetrics(mux *http.ServeMux, handler *handler.Handler) {
	prometheus.MustRegister(MetricsOpenConnections)
	prometheus.MustRegister(MetricsHookErrorsTotal)
	prometheus.MustRegister(MetricsHookInvocationsTotal)
	prometheus.MustRegister(prometheuscollector.New(handler.Metrics))

	stdout.Printf("Using %s as the metrics path.\n", Flags.MetricsPath)
	mux.Handle(Flags.MetricsPath, promhttp.Handler())
}
