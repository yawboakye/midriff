package midriff

import (
	"net/http"
)

// Group is a group of middleware functions
// that should be called sequentially, according
// to their order of insertion.
type Group struct {
	name  string
	units []http.HandlerFunc
}

// NewGroup returns a new middleware group.
func NewGroup(name string) *Group { return &Group{name: name} }

// Append appends a new middleware function to the end
// of the group.
func (g *Group) Append(fns ...http.HandlerFunc) { g.units = append(g.units, fns...) }

// Prepend prepends a new middleware function to the
// start of the group.
func (g *Group) Prepend(fns ...http.HandlerFunc) { g.units = append(fns, g.units...) }

// Extend extends a group of middleware functions with
// that of another. Extended functions are appended to
// the end of the list but this isn't always ensured as
// further calls to Append will put new functions at the
// end of the list.
func (g *Group) Extend(o *Group) {
	g.units = append(g.units, o.units...)
}

// And defines the actual handler for the request.
// It allows you to reuse a group for different kinds
// of handlers that require the same pre-processing
// of the request. The function is not appended to
// the list.
func (g *Group) And(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, mw := range g.units {
			mw.ServeHTTP(w, r)
		}
		f.ServeHTTP(w, r)
	}
}
