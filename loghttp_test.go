package loghttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			DoAround: func(req *http.Request, roundtrip RoundTripper) (*http.Response, error) {
				reqs = append(reqs, req)
				resp, err := roundtrip(req)
				if err != nil {
					return resp, err
				}
				resps = append(resps, resp)
				return resp, err
			},
		},
	}

	resp, err := client.Get(ts.URL)
	require.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)

	require.Equal(t, len(reqs), 1)
	require.Equal(t, len(resps), 1)

	assert.Equal(t, reqs[0].URL.String(), ts.URL)
	assert.Equal(t, reqs[0].Method, "GET")

	assert.Equal(t, resps[0].StatusCode, 200)
	assert.Equal(t, resps[0].Request.URL.String(), ts.URL)
}
