package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const Interval = 2 * time.Second

type MakeRequestEndpoint struct {
	EndpointData
	URL   string
	done  chan bool
	count int
	mu    sync.Mutex
}

func NewMakeRequestEndpoint() *MakeRequestEndpoint {

	requestURLVar := os.Getenv("REQUEST_URL")
	requestURL := requestURLVar

	if requestURL == "" {
		requestURL = "REQUEST_URL not set"
	}
	paths := map[string]string{
		"/request/repeat/start": fmt.Sprintf("start making repeated calls to (%s) every %s", requestURL, Interval),
		"/request/repeat/stop":  fmt.Sprintf("stop making repeated calls to (%s) every %s", requestURL, Interval),
		"/request/single":       fmt.Sprintf("make a single request to (%s)", requestURL),
	}

	return &MakeRequestEndpoint{
		EndpointData: EndpointData{
			Name:  "MakeRequest",
			Paths: paths,
		},
		URL:  requestURLVar,
		done: make(chan bool),
	}
}

func (ri *MakeRequestEndpoint) Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint (%s) received a request.", ri.Name)

	if ri.URL == "" {
		w.Write([]byte("Env var REQUEST_URL is unset, no request to be made.\n"))
	} else {
		switch {
		case strings.HasPrefix(r.URL.Path, "/request/repeat/start"):
			ri.startRepeatCalls(w, r)
			msg := fmt.Sprintf("(%v) running goroutines running", ri.count)
			w.Write([]byte(msg))
		case strings.HasPrefix(r.URL.Path, "/request/repeat/stop"):
			ri.stopRepeatCalls(w, r)
			msg := fmt.Sprintf("(%v) running goroutines running", ri.count)
			w.Write([]byte(msg))
		case strings.HasPrefix(r.URL.Path, "/request/single"):
			result, err := ri.singleRequest()
			if err != nil {
				http.Error(w, result, http.StatusInternalServerError)
			}
			w.Write([]byte(result))
		default:
			data := map[string]string{
				"/request/repeat/start": fmt.Sprintf("start making repeated calls to %s every %s", ri.URL, Interval),
				"/request/repeat/stop":  fmt.Sprintf("stop making repeated calls to %s every %s", ri.URL, Interval),
				"/request/single":       fmt.Sprintf("make a single request to %s", ri.URL),
			}

			resp, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				log.Fatalf("failed to load index page: %s", err)
			}

			fmt.Fprint(w, string(resp))
		}
	}

}

func (ri *MakeRequestEndpoint) startRepeatCalls(w http.ResponseWriter, r *http.Request) {
	ticker := time.NewTicker(Interval)

	ri.mu.Lock()
	ri.count += 1
	ri.mu.Unlock()

	go func() {
		for {
			select {
			case <-ri.done:
				log.Printf("from repeated call goroutine: received stop event")
				ri.mu.Lock()
				ri.count -= 1
				ri.mu.Unlock()
				return
			case <-ticker.C:
				log.Printf("from repeated call goroutine: making single request")
				ri.singleRequest()
			}
		}
	}()

	log.Printf("Started repeat calls to %s (goroutine count: %d)", ri.URL, ri.count)
}

func (ri *MakeRequestEndpoint) stopRepeatCalls(w http.ResponseWriter, r *http.Request) {
	log.Printf("Stopping repeat calls to %s", ri.URL)
	if ri.count != 0 {
		ri.done <- true
	}
}

func (ri *MakeRequestEndpoint) singleRequest() (string, error) {
	resp, err := http.Get(ri.URL)
	if err != nil {
		msg := fmt.Sprintf("error: failed to reach endpoint %s: %s", ri.URL, err)
		log.Print(msg)
		return msg, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("error: failed to read response body from %s : %s", ri.URL, err)
		log.Print(msg)
		return msg, err
	}
	return fmt.Sprintf("Response from %s: \n\n%s", ri.URL, body), nil
}
