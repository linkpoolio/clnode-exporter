package web

import (
	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
)

// ClientApi sets the route and starts the server for client management
func ClientApi() {
	log.Infof("starting client management api on port %d", Config.ClientApiPort)
	api := rest.NewApi()
	api.Use(rest.DefaultProdStack...)
	router, err := rest.MakeRouter(
		rest.Post("/clients", SetClients),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Config.ClientApiPort), api.MakeHandler()))
}
