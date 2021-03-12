# Package sharpvik/mux

This package seeks to provide customizable and convenient internal routing for
Go's `http.Server`.

When building a server-side app or a REST API, you may wish to implement
complex routing mechanisms that allow you to effectively process incoming
requests by checking its signature against different routes with preset filters.

There already exist a few packages that provide this functionality. However,
there is a catch: those packages are written as though the only things your
handler functions are going to be working with are `http.ResponseWriter` and
`http.Request` sent by the user. In reality, you may wish to query a database
and/or log some data through `log.Logger`. Of course, there are ways to allow
function to acces those interfaces via a global variable or something like that,
but those are a pain to write tests for.

In this package, you can find the Router sturct that is written in a way that
resolves those problems by supporting the router-embedded context.

## Theory

Mux allows you to build flexible and fairly complex routing schemas by giving
you the `Router`. To create new `Router` you should use the `New` function.

```go
func New() *Router
```

Having a single `Router` may be exactly what you want, but it is unlikely that
you decided to use this package for such basic use case. The _cool_ things come
when you make _more_ `Router`s! Or, to be more precise, more `Subrouter`s.

```go
func (rtr *Router) Subrouter() *Router
```

## License

Use of this source code is governed by the _Mozilla Public License Version 2.0_
that can be found in the [LICENSE](LICENSE) file.
