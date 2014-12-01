// Package global automatically sets http.DefaultTransport to loghttp.DefaultTransport when loaded.
package global

import (
	"github.com/motemen/go-loghttp"
	"net/http"
)

func init() {
	http.DefaultTransport = loghttp.DefaultTransport
}
