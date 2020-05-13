package mux

import "net/http"

// Context type is an empty interface that represents all the things you wish to
// expose to your Views (handler functions). Feel free to declare your own
// context types and use them with the router.
//
// TIP: make sure to pass your custom context type as a referense; it's more
// memory efficient that way since all routers would share the same context.
// For example:
//
//     package main
//
//     import (
//         "fmt"
//         "net/http"
//
//         "github.com/sharpvik/mux"
//     )
//
//     type Context struct {
//         // some fields
//     }
//
//     func NewContext(...) *Context {...}
//
// The NewContext function here returns a pointer to Context, now you can just
// pass it to the root Router.
//
//     func main() {
//         ctx := NewContext(...)
//         root := mux.New(ctx)
//         root.View = homeView
//     }
//
// When trying to use context within your handler function (view), you must use
// a type assertion to convert mux.Context (which is just an alias for an empty
// interface) to your custom Context type.
//
//     func homeView(w http.ResponseWriter, r *http.Request, ctx mux.Context) {
//         cont := ctx.(*Context) // using our local Context type here!
//         articles := cont.DB.Articles()
//         fmt.Fprint(w, articles)
//     }
//
type Context interface{}

// View is a special function type that represents a handler function. The last
// parameter it expects, represents context you wish to expose to the View. It
// is an empty interface in order to allow you to create your own context type.
type View func(http.ResponseWriter, *http.Request, Context)

// pathVarType is an alias for int that we use as a custom type for path vars.
type pathVarType int
