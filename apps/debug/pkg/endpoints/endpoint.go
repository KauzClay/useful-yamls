package endpoints

import (
	"net/http"
)

type Endpoint interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetName() string
	GetPaths() map[string]string
}

type EndpointData struct {
	Name  string
	Paths map[string]string
}

func (ed *EndpointData) GetName() string {
	return ed.Name
}

func (ed *EndpointData) GetPaths() map[string]string {
	return ed.Paths
}
