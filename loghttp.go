package loghttp

import (
	"log"
	"net/http"
)

// Transport implements http.RoundTripper. It executes HTTP requests with logging.
type Transport struct {
	Transport   http.RoundTripper
	LogRequest  func(req *http.Request)
	LogResponse func(req *http.Request, resp *http.Response)
}

// DefaultTransport is a default logging transport that wraps http.DefaultTransport.
var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

// Used if transport.LogRequest is not set.
var DefaultLogRequest = func(req *http.Request) {
	log.Printf("---> %s %s", req.Method, req.URL)
}

// Used if transport.LogResponse is not set.
var DefaultLogResponse = func(req *http.Request, resp *http.Response) {
	log.Printf("<--- %d %s", resp.StatusCode, req.URL)
}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.logRequest(req)

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(req, resp)

	return resp, err
}

func (t *Transport) logRequest(req *http.Request) {
	if t.LogRequest != nil {
		t.LogRequest(req)
	} else {
		DefaultLogRequest(req)
	}
}

func (t *Transport) logResponse(req *http.Request, resp *http.Response) {
	if t.LogResponse != nil {
		t.LogResponse(req, resp)
	} else {
		DefaultLogResponse(req, resp)
	}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
