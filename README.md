# Package sharpvik/mux

This package seeks to provide customizable and convenient internal routing for
Go's `http.Server`.

When building a server-side app or a REST API, you may wish to implement
complex routing mechanisms that allow you to effectively process incoming
requests by checking its signature against different routes with preset filters.

There already exist a few packages that provide this functionality. However,
there is a catch: those packages are written as though the only thing your
handler function are going to be working with are `http.ResponseWriter` and
`http.Request` sent by the user. In reality, you may wish to query a database
and/or log some data through `log.Logger`. Of course, there are ways to allow
function to acces those interfaces via a global variable or something like that,
but those are a pain to write tests for.

In this package, you can find the Router sturct that is written in a way that
resolves those problems by supporting the router-embedded context.



## Install

```bash
go get github.com/sharpvik/router
```



## Example

```go
package main

import (
	"net/http"
	"fmt"
	"log"
	"os"

	"github.com/sharpvik/mux"
)

// Define your custom Context type.
type Context struct {
	logger  *log.Logger
	message string
}

func main() {
	// Initialize new Router.
	rtr := mux.RootRouter(Context{
		log.New(os.Stdout, "", log.Ltime),
		"Cheer up, life's beautiful :)",
	})

	// Set router's View.
	rtr.View = func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		// Type assertion is required here to retrieve stuff from ctx.
		context := ctx.(Context)

		// Fetch logger and message from context. Remember, depending on your
		// definition of the Context type, you may have other things there.
		lgr := context.logger
		msg := context.message

		// Log some things.
		lgr.Printf("Request: %s", r.URL.String())
		lgr.Printf("Response: %s", msg)

		// Write response to the client.
		fmt.Fprintf(w, msg)
	}

	http.ListenAndServe(":5050", rtr)
}
```



## License

Use of this source code is governed by the *Mozilla Public License Version 2.0*
that can be found in the [LICENSE](LICENSE) file.