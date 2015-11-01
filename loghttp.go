// Package loghttp provides automatic logging functionalities to http.Client.
package loghttp

import (
	"log"
	"net/http"
	"time"
)

type RoundTripper func(*http.Request) (*http.Response, error)

// Transport implements http.RoundTripper. When set as Transport of http.Client, it executes HTTP requests with logging.
// No field is mandatory.
type Transport struct {
	Transport http.RoundTripper
	DoAround  func(req *http.Request, roundtrip RoundTripper) (*http.Response, error)
}

// THe default logging transport that wraps http.DefaultTransport.
var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

// Used if transport.LogRequest is not set.
var DefaultLogRequest = func(req *http.Request) {
}

// Used if transport.LogResponse is not set.
var DefaultLogResponse = func(resp *http.Response) {
}

var DefaultDoAround = func(req *http.Request, roundtrip RoundTripper) (*http.Response, error) {
	start := time.Now()
	log.Printf("---> %s %s", req.Method, req.URL)
	resp, err := roundtrip(req)
	if err != nil {
		return resp, err
	}
	log.Printf("<--- %d %s (%s)", resp.StatusCode, resp.Request.URL, time.Since(start))
	return resp, err
}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.doAround(req, func(r *http.Request) (*http.Response, error) {
		return t.transport().RoundTrip(r)
	})
}

func (t *Transport) doAround(req *http.Request, roundtrip RoundTripper) (*http.Response, error) {
	if t.DoAround != nil {
		return t.DoAround(req, roundtrip)
	} else {
		return DefaultDoAround(req, roundtrip)
	}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
