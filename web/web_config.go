package web

import (
	"time"
	"flag"
	"github.com/linkpoolio/clnode-exporter/model"
)

// Global config used throughout the service, set on flag parsing
var Config model.Config

func parseFlags() {
	flag.IntVar(&Config.Port, "port", 8080, "Port number to serve for REST/WS API")
	flag.IntVar(
		&Config.ClientApiPort,
		"clientApiPort", 8082,
		"Port number to serve the client api")
	flag.DurationVar(
		&Config.TickerInterval,
		"ticker",
		time.Second * 15,
		"Ticker interval for the monitor to refresh the nodes metrics: s, m, h")
	flag.StringVar(
		&Config.NodeConfigLocation,
		"configFile",
		"nodes.json",
		"Your nodes configuration file, containing a JSON list of each node and its credentials.")
	flag.BoolVar(
		&Config.PromEnabled,
		"prom",
		true,
		"Boolean flag to whether to export metrics for Prometheus at /metrics.")
	flag.BoolVar(
		&Config.ClientApiEnabled,
		"clientApi",
		false,
		"Boolean flag to whether to enable the client management API. This disables node configuration." +
		"Default port 8181. Do not publicly expose this port!")
	flag.BoolVar(
		&Config.StatsEnabled,
		"stats",
		true,
		"Boolean flag for whether to enable the /stats endpoint. Default is true.")
	flag.BoolVar(
		&Config.WsEnabled,
		"ws",
		true,
		"Boolean flag for whether to enable the WebSocket /ws endpoint. Default is true.")
	flag.BoolVar(
		&Config.Debug,
		"debug",
		false,
		"true/false to whether to enable debug mode")
	flag.Parse()
}