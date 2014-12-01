package global

import (
	"github.com/motemen/go-loghttp"
	"net/http"
)

func init() {
	http.DefaultClient.Transport = loghttp.DefaultTransport
}
