package model

// Metrics is the type returned for the rest and ws clients
type Metrics struct {
	TotalSpecs          int 		   `json:"totalSpecs"`
	TotalBridges        int            `json:"totalBridges"`
	Address             string         `json:"address"`
	EthBalance          string         `json:"ethBalance"`
	LinkBalance         string         `json:"linkBalance"`
	JobSpecStats        JobSpecStats   `json:"job_spec_stats"`
}