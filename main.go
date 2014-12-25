package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/jingweno/negroni-gorelic"
	"github.com/justinas/nosurf"
	"github.com/tristanoneil/badger/routes"
)

func main() {
	n := negroni.Classic()
	n.UseHandler(nosurf.New(routes.Router()))

	if os.Getenv("NEWRELIC_KEY") != "" {
		n.Use(negronigorelic.New(os.Getenv("NEWRELIC_KEY"), "badger", true))
	}

	log.Println(fmt.Sprintf("Listening on port %s", os.Getenv("PORT")))
	n.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
