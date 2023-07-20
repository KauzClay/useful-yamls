package main

import (
	"examples/debug/pkg/endpoints"
	"examples/debug/pkg/server"
	"log"
)

func main() {
	log.Print("Debug app started.")

	endpoints := []endpoints.Endpoint{
		endpoints.NewRequestInfoEndpoint(),
		endpoints.NewEnvVarEndpoint(),
		endpoints.NewMakeRequestEndpoint(),
	}

	s := server.NewAppServer(endpoints)

	log.Fatal(s.ListenAndServe())
}
