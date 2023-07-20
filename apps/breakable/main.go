package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil {
		log.Print("helloworld: request has tls context")
		log.Printf("%#v", r.TLS)
	} else {
		log.Print("helloworld: request doesn't have tls context")
	}
	log.Print("helloworld: received a request")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", target)
}

func breakit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this app is not long for this world!\n")

	log.Fatalln("going down, byeeee")
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)
	http.HandleFunc("/break", breakit)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
