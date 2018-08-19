package model

// Subscription is the type required to be sent from a ws client to subscribe to any addresses
type Subscription struct {
	Addresses []string `json:"addresses"`
}
