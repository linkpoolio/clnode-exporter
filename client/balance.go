package client

import (
	"github.com/linkpoolio/clnode-exporter/model"
)

// GetBalance calls the /v2/account_balance endpoint on the node
func (nc *NodeClient) GetBalance() (*model.AccountBalance, error) {
	var accountBalance model.AccountBalance
	err := nc.HttpGet("/v2/account_balance", &accountBalance)
	if err != nil {
		return nil, err
	}
	return &accountBalance, nil
}

