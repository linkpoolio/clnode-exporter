package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/linkpoolio/clnode-exporter/web"
)

func main() {
	log.Print("chainlink node monitor")

	web.InitialiseConfig()

	if web.Config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if web.Config.ClientApiEnabled {
		go web.ClientApi()
	}

	log.Printf("starting to serve on port %d", web.Config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", web.Config.Port), web.Api().MakeHandler()))
}
