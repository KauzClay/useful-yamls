package server

//inspired by https://github.com/Heathland/healthProbes/blob/master/main.go

import (
	"encoding/json"
	"examples/debug/pkg/endpoints"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type AppServer struct {
	Port              int
	T                 time.Time
	JobStart          time.Time
	WaitStartupTime   time.Duration
	WaitLivenessTime  time.Duration
	WaitReadinessTime time.Duration
	JobDuration       time.Duration
	Endpoints         []endpoints.Endpoint
	http.Server
}

const (
	startupPath   = "/startup"
	livenessPath  = "/liveness"
	readinessPath = "/readiness"
	startJobPath  = "/startJob"
)

var AppServerEndpointData = map[string]string{
	startupPath:   "",
	livenessPath:  "",
	readinessPath: "",
	startJobPath:  "",
}

func NewAppServer(endpoints []endpoints.Endpoint) *AppServer {
	var s AppServer

	err := s.getEnvValues()
	if err != nil {
		log.Fatal(err)
	}

	s.Endpoints = endpoints
	s.registerEndpoints()

	// Start time
	s.T = time.Now()

	return &s
}

func (s *AppServer) registerEndpoints() {
	sm := http.NewServeMux()

	sm.Handle(startupPath, http.HandlerFunc(s.startupProbe))
	sm.Handle(livenessPath, http.HandlerFunc(s.livenessProbe))
	sm.Handle(readinessPath, http.HandlerFunc(s.readinessProbe))
	sm.Handle(startJobPath, http.HandlerFunc(s.startJob))

	sm.Handle("/", http.HandlerFunc(s.index))

	for _, e := range s.Endpoints {
		for k := range e.GetPaths() {
			sm.Handle(k, http.HandlerFunc(e.Handle))
		}
	}

	s.Handler = sm
}

func (s *AppServer) index(w http.ResponseWriter, r *http.Request) {
	log.Print("Index Page received a request.")

	data := map[string]string{}

	for _, e := range s.Endpoints {
		for k, v := range e.GetPaths() {
			data[k] = v
		}
	}

	for k, v := range AppServerEndpointData {
		data[k] = v
	}

	resp, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("failed to load index page: %s", err)
	}

	fmt.Fprint(w, string(resp))
}

func getEnvToDuration(e string) (time.Duration, error) {
	envValue := os.Getenv(e)
	if envValue != "" {
		envValueInt, err := strconv.Atoi(os.Getenv(e))
		return time.Duration(envValueInt) * time.Second, err
	} else {
		log.Printf("(%s) not set, defaulting to 0\n", e)
		return 0 * time.Second, nil
	}
}

func (s *AppServer) getEnvValues() (err error) {
	s.WaitStartupTime, err = getEnvToDuration("WAIT_STARTUP_TIME")
	if err != nil {
		log.Fatalf("error getting WAIT_STARTUP_TIME: %s\n", err)
	}
	s.WaitLivenessTime, err = getEnvToDuration("WAIT_LIVENESS_TIME")
	if err != nil {
		log.Fatalf("error getting WAIT_LIVENESS_TIME: %s\n", err)
	}
	s.WaitReadinessTime, err = getEnvToDuration("WAIT_READINESS_TIME")
	if err != nil {
		log.Fatalf("error getting WAIT_READINESS_TIME: %s\n", err)
	}
	s.JobDuration, err = getEnvToDuration("JOB_DURATION_TIME")
	if err != nil {
		log.Fatalf("error getting JOB_DURATION_TIME: %s\n", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Print("PORT not set, defaulting to 8080\n")
		s.Addr = ":8080"
	} else {
		s.Addr = ":" + port
	}

	return
}

func (s *AppServer) startupProbe(w http.ResponseWriter, r *http.Request) {
	if time.Since(s.T) > s.WaitStartupTime {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(503)
	}
}

func (s *AppServer) livenessProbe(w http.ResponseWriter, r *http.Request) {
	if time.Since(s.T) > s.WaitLivenessTime {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(503)
	}
}

func (s *AppServer) readinessProbe(w http.ResponseWriter, r *http.Request) {
	if time.Since(s.T) > s.WaitReadinessTime && time.Since(s.JobStart) > s.JobDuration {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(503)
	}
}

func (s *AppServer) startJob(w http.ResponseWriter, r *http.Request) {
	if time.Since(s.JobStart) > s.JobDuration {
		s.JobStart = time.Now()
		fmt.Fprintf(w, "Pod (%s)\nStarting job. Unavailable till: %s", os.Getenv("HOSTNAME"), s.JobStart.Add(s.JobDuration).Format("Mon Jan _2 15:04:05 2006"))
	} else {
		fmt.Fprintf(w, "Still running job. Unavailable till: %v", s.JobStart.Add(s.JobDuration))
	}
}
