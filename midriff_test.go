package midriff

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewGroup(t *testing.T) {
	name := "test-group"

	group := NewGroup(name)
	if group.name != name {
		t.Fatalf("midriff_test: expected %s; got=%s instead", group.name, name)
	}

	if len(group.units) != 0 {
		t.Fatalf("midriff_test: expected 0 units; got=%d instead", len(group.units))
	}
}
func TestAppend(t *testing.T) {
	group := NewGroup("test-append")
	f := func(http.ResponseWriter, *http.Request) {}
	g := func(http.ResponseWriter, *http.Request) {}

	group.Append(f)
	group.Append(g)

	// Expect two units in the group
	if len(group.units) != 2 {
		t.Fatalf("expected 2 units; got=%d instead", len(group.units))
	}

	// Order of units
	// This isn't a great test, mainly because while functions can
	// be passed around as values in Go, they cannot be compared to
	// each other. In the future I hope to investigate why this is
	// so.
	expected := []http.HandlerFunc{f, g}
	for i, unit := range group.units {
		if !funcEqual(unit, expected[i]) {
			t.Fatal("order of units is different than expected.")
		}
	}

}
func TestPrepend(t *testing.T) {
	group := NewGroup("test-prepend")
	f := func(http.ResponseWriter, *http.Request) {}
	g := func(http.ResponseWriter, *http.Request) {}

	group.Append(f)
	group.Prepend(g)

	// Expect two units in the group
	if len(group.units) != 2 {
		t.Fatalf("expected 2 units; got=%d instead", len(group.units))
	}

	// Test that g comes before f
	expected := []http.HandlerFunc{g, f}
	for i, unit := range group.units {
		if !funcEqual(unit, expected[i]) {
			t.Fatal("order of units is different than expected.")
		}
	}
}
func TestExtend(t *testing.T) {
	basic := NewGroup("basic-group")
	basicF := func(http.ResponseWriter, *http.Request) {}
	basicG := func(http.ResponseWriter, *http.Request) {}
	basic.Append(basicF, basicG)

	authed := NewGroup("auth-group")
	authedF := func(http.ResponseWriter, *http.Request) {}
	authed.Append(authedF)
	authed.Extend(basic)

	// Expect 3 units in the auth-group group
	if len(authed.units) != 3 {
		t.Fatalf("expected 3 units; got=%d instead", len(authed.units))
	}

	// Expected order of units in the group
	expected := []http.HandlerFunc{authedF, basicF, basicG}
	for i, unit := range authed.units {
		if !funcEqual(unit, expected[i]) {
			t.Fatal("wrong order of units in group")
		}
	}

}
func TestAnd(t *testing.T) {
	respBody := "hello, world!"
	group := NewGroup("test-and-group")
	f := func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Request-ID", "request-id")
	}
	g := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(respBody))
	}

	group.Append(f)
	handler := group.And(g)

	// Test that `And` returns an http.HandlerFunc which
	// calls all of the units when it is invoked.
	r := httptest.NewRequest("GET", "https://example.org", nil)
	w := httptest.NewRecorder()
	handler(w, r)

	if r.Header.Get("Request-ID") != "request-id" {
		t.Fatal("expected unit to set Request-ID but wasn't ran")
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != respBody {
		t.Fatalf("expected handler to write %s as body; got=%s instead", respBody, string(body))
	}
}

func funcEqual(f, g http.HandlerFunc) bool {
	ptrToF := reflect.ValueOf(f).Pointer()
	ptrToG := reflect.ValueOf(g).Pointer()

	return ptrToF == ptrToG
}
