package model

import "time"

// Config stores all the configurable options parsed by flags
type Config struct {
	Port 		   	   int
	TickerInterval 	   time.Duration
	NodeConfigLocation string
	NodeConfigs 	   []*NodeConfig
	PromPort           int
	PromEnabled        bool
	ClientApiEnabled   bool
	ClientApiPort      int
	Debug              bool
	WsEnabled          bool
	StatsEnabled       bool
}
