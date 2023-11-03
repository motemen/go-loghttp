package loghttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("200 OK"))
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	var (
		reqs  = []*http.Request{}
		resps = []*http.Response{}
	)

	client := &http.Client{
		Transport: &Transport{
			LogRequest: func(req *http.Request) {
				reqs = append(reqs, req)
			},
			LogResponse: func(resp *http.Response) {
				resps = append(resps, resp)
			},
		},
	}

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	if len(reqs) != 1 {
		t.Errorf("expected 1 request, got %d", len(reqs))
	}

	if len(resps) != 1 {
		t.Errorf("expected 1 response, got %d", len(resps))
	}

	if reqs[0].URL.String() != ts.URL {
		t.Errorf("expected request URL %q, got %q", ts.URL, reqs[0].URL.String())
	}

	if reqs[0].Method != "GET" {
		t.Errorf("expected request method GET, got %q", reqs[0].Method)
	}

	if resps[0].StatusCode != 200 {
		t.Errorf("expected response status code 200, got %d", resps[0].StatusCode)
	}

	if resps[0].Request.URL.String() != ts.URL {
		t.Errorf("expected response request URL %q, got %q", ts.URL, resps[0].Request.URL.String())
	}
}
