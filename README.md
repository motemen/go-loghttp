go-loghttp
==========

Log http.Client's requests and responses automatically.

[GoDoc](http://godoc.org/github.com/motemen/go-loghttp)

## Synopsis

To log all the HTTP requests/responses, import `github.com/motemen/go-loghttp/global`.

```go
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/motemen/go-loghttp/global" // Just this line!
)

func main() {
	resp, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, resp.Body)
}
```

Or set `loghttp.Transport` to `http.Client`'s `Transport` field.

```
import "github.com/motemen/go-loghttp"

client := &http.Client{
	Transport: &loghttp.Transport{},
}
```

You can modify [loghttp.Transport](http://godoc.org/github.com/motemen/go-loghttp#Transport)'s `LogRequest` and `LogResponse` to customize logging function.

## Author

motemen <motemen@gmail.com>
