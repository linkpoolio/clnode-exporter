package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"github.com/linkpoolio/clnode-exporter/model"
	"encoding/json"
	"fmt"
)

// SetClients receives POST data and updates the nodes to scrape metrics
func SetClients(w rest.ResponseWriter, r *rest.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	var configs []*model.NodeConfig
	err := json.Unmarshal(bytes, &configs)
	if err != nil {
		sendError(500, err, w)
	}
	UpdateConfigs(configs)
	UpdateWsClients()
}

func sendError(sc int, err error, w rest.ResponseWriter) {
	m := map[string]string {
		"statusCode": fmt.Sprintf("%d", sc),
		"error": err.Error(),
	}
	w.WriteJson(m)
}
