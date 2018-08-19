package web

import (
	"os"
	"io/ioutil"
	"github.com/linkpoolio/clnode-exporter/model"
	"encoding/json"
	"github.com/linkpoolio/clnode-exporter/client"
	log "github.com/sirupsen/logrus"
)

// SetNodeConfig refreshes the client nodes from file
func SetNodeConfig() {
	if Config.ClientApiEnabled { return }
	jsonFile, err := os.Open(Config.NodeConfigLocation)
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var configs []*model.NodeConfig
	json.Unmarshal(bytes, &configs)
	UpdateConfigs(configs)
}

// Initialises the config by parsing flags, then refreshing config from file (if enabled)
func InitialiseConfig() {
	Config = model.Config{}

	parseFlags()
	SetNodeConfig()
}

// UpdateConfig refreshes the client metrics map used by all controllers based on the config
func UpdateConfigs(cs []*model.NodeConfig) {
	if clientMetrics == nil {
		clientMetrics = make(map[*client.NodeClient]*model.Metrics)
	}
	for _, c := range cs {
		if !configExists(c) {
			nc, err := client.NewNodeClient(c)
			if err != nil {
				log.Error(err)
			} else {
				clientMetrics[nc] = nil
			}
		}
	}
	for nc, m := range clientMetrics {
		exists := false
		for _, c := range cs {
			if nc.Config.URL == c.URL {
				exists = true
				break
			}
		}
		if !exists {
			delete(clientMetrics, nc)
			DeletePromMetricPerNode(m)
		}
	}
}

func configExists(nc *model.NodeConfig) bool {
	for c := range clientMetrics {
		if c.Config.URL == nc.URL {
			return true
		}
	}
	return false
}