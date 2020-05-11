// Use of this source code is governed by the Mozilla Public License Version 2.0
// that can be found in the LICENSE file.

/*
Package mux seeks to provide customizable and convenient internal routing for
Go's http.Server.

When building a server-side app or a REST API, you may wish to implement
complex routing mechanisms that allow you to effectively process incoming
requests by checking their signatures against different routes with preset
filters.

There already exist a few packages that provide this functionality. However,
there is a catch: those packages are written as though the only things your
handler functions are going to be working with are http.ResponseWriter and
http.Request sent by the user. In reality, you may wish to query a database
and/or log some data through log.Logger. Of course, there are ways to allow a
function to acces those interfaces via a global variable or something like that,
but those ways are a pain to write tests for.

In this package, you can find the Router sturct that is written in a way that
resolves those problems by supporting the router-embedded context.
*/
package mux
