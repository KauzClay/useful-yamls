package main

import (
	// "fmt"
	// "io"
	"fmt"
	"log"
	"net/http"
)

//based on https://gist.github.com/denji/12b3a568f092ab951456

func HelloServer(w http.ResponseWriter, req *http.Request) {
	if req.TLS != nil {
		fmt.Println("TLS config present")
		fmt.Printf("%#v\n", req.TLS)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func HelloServerTLS(w http.ResponseWriter, req *http.Request) {
	if req.TLS != nil {
		fmt.Println("TLS config present")
		fmt.Printf("%#v\n", req.TLS)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example tls server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	//http.ListenAndServe(":8080", nil)
	err := http.ListenAndServeTLS(":7890", "tls.crt", "tls.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
