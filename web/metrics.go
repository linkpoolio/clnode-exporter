package web

import (
	"github.com/linkpoolio/clnode-exporter/client"
	"github.com/linkpoolio/clnode-exporter/model"
	"time"
	"sync"
	log "github.com/sirupsen/logrus"
	"github.com/linkpoolio/clnode-exporter/service"
)

var (
	clientMetrics    map[*client.NodeClient]*model.Metrics
	supportedMetrics []func(nodeClient *client.NodeClient, metrics *model.Metrics)
)

func init() {
	supportedMetrics = append(supportedMetrics,
		service.CollectNodeStats,
		service.CollectBalanceMetrics,
		service.CollectBridgeMetrics,
	)
}

// StartMonitor creates the ticker to gather metrics based on user input
func StartMonitor() {
	log.Info("collecting metrics for the first time")
	collectMetrics()

	ticker := time.NewTicker(Config.TickerInterval)
	go func() {
		for range ticker.C {
			collectMetrics()
			log.Debug("metrics collected")
		}
	}()
}

func collectMetrics() {
	SetNodeConfig()

	var wg sync.WaitGroup
	wg.Add(len(clientMetrics))

	for c, m := range clientMetrics {
		go func(c *client.NodeClient, m *model.Metrics) {
			defer wg.Done()

			var mwg sync.WaitGroup
			mwg.Add(len(supportedMetrics))

			tm := &model.Metrics{}
			if m != nil {
				tm = clientMetrics[c]
			}

			for _, f := range supportedMetrics {
				go func(f func(nodeClient *client.NodeClient, metrics *model.Metrics), c *client.NodeClient, m *model.Metrics) {
					defer mwg.Done()
					f(c, m)
				}(f, c, tm)
			}

			if m == nil {
				clientMetrics[c] = tm
			}

			mwg.Wait()
		}(c, m)
	}
	wg.Wait()

	UpdateWsClients()
	UpdatePromMetrics()
}
