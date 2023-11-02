// Package global automatically sets http.DefaultTransport to loghttp.DefaultTransport when loaded.
package global

import (
	"net/http"

	"github.com/motemen/go-loghttp"
)

func init() {
	http.DefaultTransport = loghttp.DefaultTransport
}
