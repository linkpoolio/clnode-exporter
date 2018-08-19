package web

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/linkpoolio/clnode-exporter/model"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var totalSpecs = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_total_specs",
		Help: "Total number of specs on the node.",
	},
	[]string{"address"},
)
var totalSpecRuns = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_total_spec_runs",
		Help: "Total number of specs on the node.",
	},
	[]string{"address"},
)
var totalBridges = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_total_bridges",
		Help: "Total number of bridges on the node.",
	},
	[]string{"address"},
)
var ethBalance = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_eth_balance",
		Help: "The Ethereum balance of the nodes wallet.",
	},
	[]string{"address"},
)
var linkBalance = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_link_balance",
		Help: "The LINK balance of the nodes wallet.",
	},
	[]string{"address"},
)
var specTaskCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_spec_task_count",
		Help: "The total of spec task runs of a specific adaptor.",
	},
	[]string{"address", "adaptor"},
)
var specRunStatusCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_spec_run_status_count",
		Help: "The count of spec run statuses.",
	},
	[]string{"address", "status"},
)
var specRunParamCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "node_spec_param_count",
		Help: "The count of spec run urls used.",
	},
	[]string{"address", "param", "value"},
)


func init() {
	prometheus.MustRegister(
		totalSpecs,
		totalSpecRuns,
		totalBridges,
		ethBalance,
		linkBalance,
		specTaskCount,
		specRunStatusCount,
		specRunParamCount,
	)
}

// ServeProm transforms the go-json-rest reader/writer to the in-built for the prom client
func ServeProm(w rest.ResponseWriter, r *rest.Request) {
	prom := promhttp.Handler()
	prom.ServeHTTP(w.(http.ResponseWriter), r.Request)
}

// UpdatePromMetrics refreshes all the prometheus client gauges for output
func UpdatePromMetrics() {
	for _, m := range clientMetrics {
		if m != nil {
			updatePromMetricPerNode(m)
		}
	}
}

// DeletePromMetricsPerNode deletes the values stored in the prometheus gauges
func DeletePromMetricPerNode(m *model.Metrics) {
	totalSpecs.DeleteLabelValues(m.Address)
	totalSpecRuns.DeleteLabelValues(m.Address)
	totalBridges.DeleteLabelValues(m.Address)
	ethBalance.DeleteLabelValues(m.Address)
	linkBalance.DeleteLabelValues(m.Address)

	for a := range m.JobSpecStats.AdaptorCount {
		specTaskCount.DeleteLabelValues(m.Address, a)
	}

	for s := range m.JobSpecStats.StatusCount {
		specRunStatusCount.DeleteLabelValues(m.Address, s)
	}

	for p, pc := range m.JobSpecStats.ParamCount {
		for _, vc := range pc {
			specRunParamCount.DeleteLabelValues(m.Address, p, vc.Value)
		}
	}
}

func updatePromMetricPerNode(m *model.Metrics) {
	totalSpecs.WithLabelValues(m.Address).Set(float64(m.TotalSpecs))
	totalSpecRuns.WithLabelValues(m.Address).Set(float64(m.JobSpecStats.RunCount))
	totalBridges.WithLabelValues(m.Address).Set(float64(m.TotalBridges))
	ethBalance.WithLabelValues(m.Address).Set(ToFloat64(m.EthBalance))
	linkBalance.WithLabelValues(m.Address).Set(ToFloat64(m.LinkBalance))

	for a, c := range m.JobSpecStats.AdaptorCount {
		specTaskCount.WithLabelValues(m.Address, a).Set(float64(c))
	}

	for s, c := range m.JobSpecStats.StatusCount {
		specRunStatusCount.WithLabelValues(m.Address, s).Set(float64(c))
	}

	for p, pc := range m.JobSpecStats.ParamCount {
		for _, vc := range pc {
			specRunParamCount.WithLabelValues(m.Address, p, vc.Value).Set(float64(vc.Count))
		}
	}
}