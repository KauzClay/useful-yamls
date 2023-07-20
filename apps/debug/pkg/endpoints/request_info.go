package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type RequestInfoEndpoint struct {
	EndpointData
}

func NewRequestInfoEndpoint() *RequestInfoEndpoint {
	return &RequestInfoEndpoint{
		EndpointData: EndpointData{
			Name: "RequestInfo",
			Paths: map[string]string{
				"/requestinfo": "responds with the httputil.DumpRequest data",
			},
		},
	}
}

func (ri *RequestInfoEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint (%s) received a request.", ri.Name)

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Request Info: \n\n%s\n", string(requestDump))

	fmt.Fprintf(w, "Request Info: \n\n%s\n", string(requestDump))

}
