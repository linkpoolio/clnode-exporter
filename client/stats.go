package client

import (
	"github.com/linkpoolio/clnode-exporter/model"
	"fmt"
)

// GetStats calls the /v2/stats endpoint on the node with a given page & size
func (nc *NodeClient) GetStats(page int, size int) (*model.NodeStats, error) {
	var nodeStats model.NodeStats
	err := nc.HttpGet(fmt.Sprintf("/v2/stats?page=%d&size=%d", page, size), &nodeStats)
	if err != nil {
		return nil, err
	}
	return &nodeStats, nil
}