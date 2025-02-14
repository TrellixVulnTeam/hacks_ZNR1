package main

import (
	"log"
	"net/http"
	"time"
)

const (
	numPollers     = 100
	pollInterval   = 2 * time.Second
	statusInterval = 5 * time.Second
	errTimeout     = 10 * time.Second
)

var urls = []string{
	"http://www.google.com",
	"http://golang.org/",
	"http://blog.golang.org/",
	"https://mattjmcnaughton.com",
	"https://github.com/",
}

type State struct {
	url    string
	status string
}

// StateMonitor maintains a map that stores the state of the URLs being polled
// and prints the current state every `updateInterval` nanoseconds.
// It returns a chan State to which resource state should be sent.
func StateMonitor(updateInterval time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string)
	ticker := time.NewTicker(updateInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()

	return updates
}

// logState prints a state map.
func logState(s map[string]string) {
	log.Println("Current state:")
	for k, v := range s {
		log.Printf("%s %s", k, v)
	}
}

// Resource represents an HTTP URL to be polled by this program.
type Resource struct {
	url      string
	errCount int
}

// Poll executes an HTTP HEAD request for url
// and returns the HTTP status string or an error string.
func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}
	r.errCount = 0
	return resp.Status
}

// Sleep sleeps for an appropriate interval before sending the Resource to done.
func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(pollInterval + errTimeout*time.Duration(r.errCount))
	done <- r
}

func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	for r := range in {
		s := r.Poll()
		status <- State{r.url, s}
		out <- r
	}
}

func main() {
	// Create our input and output channels.
	pending, complete := make(chan *Resource), make(chan *Resource)

	// Launch the StateMonitor. It returns the channel to which we should
	// send its events.
	status := StateMonitor(statusInterval)

	// Channels provide the main means of communication between different go
	// routines.

	for i := 0; i < numPollers; i++ {
		go Poller(pending, complete, status)
	}

	go func() {
		for _, url := range urls {
			pending <- &Resource{url: url}
		}
	}()

	for r := range complete {
		// Send back to the pending queue after the sleep.
		go r.Sleep(pending)
	}
}
