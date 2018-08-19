package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
)

// Api sets the routes for websocket and stats, prom if enabled
// also starts the ticker for gathering metrics
func Api() *rest.Api{
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	var routes []*rest.Route
	if Config.StatsEnabled {
		routes = append(routes, rest.Get("/stats", GetMetrics))
	}
	if Config.WsEnabled {
		routes = append(routes, rest.Get("/ws", NewConnection))
	}
	if Config.PromEnabled {
		routes = append(routes, rest.Get("/metrics", ServeProm))
	}
	router, err := rest.MakeRouter(routes...)
	if err != nil {
		log.Fatal(err)
	}

	StartMonitor()
	api.SetApp(router)
	log.Print("apis started")
	return api
}