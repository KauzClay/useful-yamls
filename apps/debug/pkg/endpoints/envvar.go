package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type EnvVarEndpoint struct {
	EndpointData
}

func NewEnvVarEndpoint() *EnvVarEndpoint {
	return &EnvVarEndpoint{
		EndpointData: EndpointData{
			Name: "EnvVar",
			Paths: map[string]string{
				"/env": "responds with the value of the environment variable MY_VAR",
			},
		},
	}
}

func (c *EnvVarEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint (%s) received a request.", c.Name)

	target := os.Getenv("MY_VAR")
	if target == "" {
		target = "no color provided"
	}
	fmt.Fprintf(w, "%s\n", target)
}
