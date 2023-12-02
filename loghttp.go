// Package loghttp provides automatic logging functionalities to http.Client.
package loghttp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/motemen/go-nuts/roundtime"
)

// Transport implements http.RoundTripper. When set as Transport of http.Client, it executes HTTP requests with logging.
// No field is mandatory.
type Transport struct {
	Transport   http.RoundTripper
	LogRequest  func(req *http.Request)
	LogResponse func(resp *http.Response)
}

// THe default logging transport that wraps http.DefaultTransport.
var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

// Used if transport.LogRequest is not set.
var DefaultLogRequest = func(req *http.Request) {
	log.Printf("--> %s %s", req.Method, req.URL)
}

// Used if transport.LogResponse is not set.
var DefaultLogResponse = func(resp *http.Response) {
	ctx := resp.Request.Context()
	if start, ok := ctx.Value(ContextKeyRequestStart).(time.Time); ok {
		log.Printf("<-- %d %s (%s)", resp.StatusCode, resp.Request.URL, roundtime.Duration(time.Now().Sub(start), 2))
	} else {
		log.Printf("<-- %d %s", resp.StatusCode, resp.Request.URL)
	}
}

type contextKey struct {
	name string
}

var ContextKeyRequestStart = &contextKey{"RequestStart"}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), ContextKeyRequestStart, time.Now())
	req = req.WithContext(ctx)

	t.logRequest(req)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	return resp, err
}

func (t *Transport) logRequest(req *http.Request) {
	if t.LogRequest != nil {
		t.LogRequest(req)
	} else {
		DefaultLogRequest(req)
	}
}

func (t *Transport) logResponse(resp *http.Response) {
	if t.LogResponse != nil {
		t.LogResponse(resp)
	} else {
		DefaultLogResponse(resp)
	}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}

func (t *Transport) CloseIdleConnections() {
	if tr, ok := t.transport().(interface {
		CloseIdleConnections()
	}); ok {
		tr.CloseIdleConnections()
	}
}

func (t *Transport) CancelRequest(req *http.Request) {
	if tr, ok := t.transport().(interface {
		CancelRequest(*http.Request)
	}); ok {
		tr.CancelRequest(req)
	}
}
