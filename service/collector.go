package service

import (
	"github.com/linkpoolio/clnode-exporter/model"
	"github.com/linkpoolio/clnode-exporter/client"
	log "github.com/sirupsen/logrus"
)

// CollectNodeStats fetches all the stats for each job spec and then aggregates it all
func CollectNodeStats(c *client.NodeClient, m *model.Metrics) {
	p := 1
	s := 100

	jss := model.JobSpecStats{}
	jss.ID = "*" // All Job Spec stats aggregated
	jss.AdaptorCount = make(map[string]int)
	jss.StatusCount = make(map[string]int)
	jss.ParamCount = make(map[string][]model.ParamCount)

	for {
		ns, err := c.GetStats(p, s)
		if err != nil {
			metricError(c,"stats", err)
			return
		}

		m.TotalSpecs = ns.Meta.Count

		for _, js := range ns.Data.Attributes.JobSpecStats {
			jss.RunCount += js.RunCount

			for a, ac := range js.AdaptorCount {
				jss.AdaptorCount[a] += ac
			}

			for s, sc := range js.StatusCount {
				jss.StatusCount[s] += sc
			}

			for p, pc := range js.ParamCount {
				for _, vc := range pc {
					f := false
					fi := 0
					for i, cvc := range jss.ParamCount[p] {
						if cvc.Value == vc.Value {
							f = true
							fi = i
							break
						}
					}
					if !f {
						jss.ParamCount[p] = append(jss.ParamCount[p], model.ParamCount{
							Value: vc.Value,
							Count: vc.Count,
						})
					} else {
						jss.ParamCount[p][fi].Count += vc.Count
					}
				}
			}
		}

		if ns.Links.Next == "" {
			break
		}

		p++
	}

	m.JobSpecStats = jss
}

// CollectBalanceMetrics sets the balance of both ETH/LINK from the node
func CollectBalanceMetrics(c *client.NodeClient, m *model.Metrics) {
	balance, err := c.GetBalance()
	if err != nil {
		metricError(c,"balance", err)
		return
	}
	m.Address = balance.Data.Id
	m.EthBalance = balance.Data.Attributes.EthBalance
	m.LinkBalance = balance.Data.Attributes.LinkBalance
}

// CollectBridgeMetrics sets the total amount of bridges created on the node
func CollectBridgeMetrics(c *client.NodeClient, m *model.Metrics) {
	bridges, err := c.GetBridges()
	if err != nil {
		metricError(c,"bridge", err)
		return
	}
	m.TotalBridges = bridges.Meta.Count
}

func metricError(c *client.NodeClient, metric string, err error) {
	log.WithFields(log.Fields{
		"hostname": c.Config.URL,
		"metric": metric,
	}).Error(err)
}
