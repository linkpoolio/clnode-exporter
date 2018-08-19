package web

import (
	"github.com/gorilla/websocket"
	"github.com/linkpoolio/clnode-exporter/model"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
)

var wsClients = make(map[*websocket.Conn]bool)
var clientSubscriptions = make(map[*websocket.Conn][]string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewConnection upgrades a HTTP connection to ws, leaving it open until error or disconnect
func NewConnection(w rest.ResponseWriter, r *rest.Request) {
	ws, err := upgrader.Upgrade(w.(http.ResponseWriter), r.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	wsClients[ws] = true
	for {
		var msg model.Subscription
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Error(err)
			delete(wsClients, ws)
			delete(clientSubscriptions, ws)
			break
		}
		clientSubscriptions[ws] = msg.Addresses
		updateClient(ws)
	}
}

// UpdateWsClients refreshes the output given to every connected client
func UpdateWsClients() {
	for c := range wsClients {
		updateClient(c)
	}
}

func updateClient(c *websocket.Conn) {
	msg := make(map[string]*model.Metrics)
	for _, m := range clientMetrics {
		if isSubscribed(clientSubscriptions[c], m.Address) {
			msg[m.Address] = m
		}
	}
	err := c.WriteJSON(msg)
	if err != nil {
		log.Error(err)
		c.Close()
		delete(wsClients, c)
	}
}

func isSubscribed(subs []string, addr string) bool {
	for _, s := range subs {
		if s == addr {
			return true
		}
	}
	return false
}
