package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/linkpoolio/clnode-exporter/model"
)

// GetMetrics returns all the NodeStats gathered in an array keyed by their address
func GetMetrics(w rest.ResponseWriter, _ *rest.Request) {
	nm := make(map[string]*model.Metrics)
	for _, m := range clientMetrics {
		if m != nil {
			nm[m.Address] = m
		}
	}
	w.WriteJson(nm)
}