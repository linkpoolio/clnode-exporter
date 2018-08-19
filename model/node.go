package model

// ParamCount reflects the `param_count` obj from the stats endpoint
type ParamCount struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// JobSpecStats reflects the `job_spec_stats` obj from the stats endpoint
type JobSpecStats struct {
	ID           string                  `json:"id"`
	RunCount     int                     `json:"run_count"`
	AdaptorCount map[string]int          `json:"adaptor_count"`
	StatusCount  map[string]int          `json:"status_count"`
	ParamCount   map[string][]ParamCount `json:"param_count"`
}

// NodeStatsAttributes reflects the `attributes` obj from the stats endpoint
type NodeStatsAttributes struct {
	JobSpecStats []JobSpecStats `json:"job_spec_stats"`
}

// NodeStatsData reflects the `data` obj from the stats endpoint
type NodeStatsData struct {
	ID string `json:"id"`
	Attributes NodeStatsAttributes `json:"attributes"`
}

// NodeStats is the parent type returned from the stats endpoint
type NodeStats struct {
	Data  NodeStatsData `json:"data"`
	Links Links         `json:"links"`
	Meta  Meta          `json:"meta"`
}

// NodeConfig is the type needed for the client management api and the config file
type NodeConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Session is the type needed for creating/deleting sessions on the node
type Session struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AccountBalanceAttributes reflects `attributes` returned from the account balance endpoint
type AccountBalanceAttributes struct {
	Address     string `json:"address"`
	EthBalance  string `json:"eth_balance"`
	LinkBalance string `json:"link_balance"`
}

// AccountBalanceData reflects `data` returned from the account balance endpoint
type AccountBalanceData struct {
	Id 	       string                   `json:"id"`
	Attributes AccountBalanceAttributes `json:"attributes"`
}

// AccountBalance is the parent type for the account balance endpoint
type AccountBalance struct {
	Data AccountBalanceData `json:"data"`
}

// BridgeTypeAttributes reflects the `attributes` obj returned from the bridge endpoint
type BridgeTypeAttributes struct {
	URL string `json:"url"`
}

// BridgeTypeData reflects the `data` obj returned from the bridge endpoint
type BridgeTypeData struct {
	Type       string               `json:"type"`
	Id         string               `json:"id"`
	Attributes BridgeTypeAttributes `json:"attributes"`
}

// BridgeTypes is the parent type for the bridge type api
type BridgeTypes struct {
	Data []*BridgeTypeData `json:"data"`
	Meta Meta              `json:"meta"`
}

// Links is part of the paginated response from a node
type Links struct {
	Next string `json:"next"`
}

// Meta is part of the paginated response from a node
type Meta struct {
	Count int `json:"count"`
}