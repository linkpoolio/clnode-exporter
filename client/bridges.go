package client

import "github.com/linkpoolio/clnode-exporter/model"

// GetBridges calls the /v2/bridges endpoint on the node
func (nc *NodeClient) GetBridges() (*model.BridgeTypes, error) {
	var bt model.BridgeTypes
	err := nc.HttpGet("/v2/bridge_types", &bt)
	if err != nil {
		return nil, err
	}
	return &bt, nil
}