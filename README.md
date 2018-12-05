### Midriff

[![GoDoc](https://godoc.org/github.com/yawboakye/midriff?status.svg)](https://godoc.org/github.com/yawboakye/midriff)

Because we have to inspect requests before they reach our
inner parts, and our response before they are sent out into
the open, unforgiving world.

Many a times you write a simple web application, in Go,
without a framework. Which usually means you have the
honor and pleasure of composing your own parts to make
a whole that uniquely fits you, looks like you even.

In those situations you'd probably need something like a
middleware, and just like me, you've probably created one
for yourself already. In that case you don't need me and
what I have to offer. Unless you want to entertain some
competition for your little servant, keep it on its toes.

For the rest of us who haven't yet abstracted our middleware
functions into a neat package, this is how your routes/mux
setup probably looks like:

```go
package main

import "net/http"

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/route1", one(two(three(four(handler1)))))
  mux.HandleFunc("/route2", one(two(three(four(handler2)))))
}
```

A little Lisp in your Go. Or, you wrote a little combinator
(aka chaining function) and that helped you bring them
parentheses under control. Does your code look like this?

```go
package main

import (
  "net/http"
)

func chain(handlers ...[]http.HandlerFunc) http.HandlerFunc {
  // Combine all the handlers and return a single
  // handler that knows how to do what the constituent
  // handlers did.
}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/route1", chain(one, two, three, four, handler1))
  mux.HandleFunc("/route2", chain(one, two, three, four, handler2))
}
```

Still gnarly, and it doesn't even depend on how you look at
it or who you ask. Flat out gnarly. I know because like you
I've written loads of this and variants. Many years later
how this appeared simple and clean beats me. It goes without
saying that I don't write that anymore. I _midriff_ my middleware
like so:

```go
package main

import (
  "net/http"

  "github.com/yawboakye/midriff"
)

func main() {
  // A basic group of middleware functions.
  basicGroup := midriff.NewGroup("json-api")
  basicGroup.Append(one, two, three, four)

  mux.HandleFunc("/route1", group.And(handler1))
  mux.HandleFunc("/route2", group.And(handler2))

  // A new group of middleware that extends
  // the basic middleware. It ensures that
  // the request has enough information to be
  // granted admin privileges.
  adminGroup := midriff.NewGroup("must-be-admin")
  adminGroup.Extend(basicGroup)
  adminGroup.Prepend(mustBeAdmin)

  mux.HandleFunc("/admin/route1", adminGroup.And(handler1))
  mux.HandleFunc("/admin/route2", adminGroup.And(handler2))
}
```

In the interest of keeping things simple we'll keep it at
that. Something I've thought of is to move the `HandleFunc`
function of `mux` onto the middleware group. But I resist.
