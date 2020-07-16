package starfish

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	// A general request for all tests
	t.Run("an null router", func(t *testing.T) {
		router := new(Router)
		t.Run("has no routes", func(t *testing.T) {
			routes := router.Routes()
			if len(routes) != 0 {
				t.Errorf("Expected 0 routes, got %d", len(routes))
			}
		})
		t.Run("always serves 404", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
			assertStatus(t, router, r, http.StatusNotFound)
		})
	})

	t.Run("a configured router", func(t *testing.T) {
		router := new(Router)
		router.Match(Never).Handle(writeString("never hit"))
		router.Match(Always).Handle(writeString("Hello!"))
		t.Run("has the configured routes", func(t *testing.T) {
			routes := router.Routes()
			if len(routes) != 2 {
				t.Errorf("Expected 2 routes, got %d", len(routes))
			}
		})
	})
}

func TestServeHTTP(t *testing.T) {
	Match(Never).Handle(writeString("never hit"))
	Match(Always).Handle(writeString("Hello!"))
	t.Run("serves the first matching route", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
		w := httptest.NewRecorder()
		ServeHTTP(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		assertEqual(t, resp.StatusCode, http.StatusOK)
		assertEqual(t, string(body), "Hello!")
	})
}

func TestPush(t *testing.T) {
	r := new(Router)
	r.Push(&testRoute{})
	assertEqual(t, 1, len(r.Routes()))
	r.Push(&testRoute{})
	r.Push(&testRoute{})
	assertEqual(t, 3, len(r.Routes()))
}

func TestPop(t *testing.T) {
	r := Router([]Route{&testRoute{}, &testRoute{}})
	assertEqual(t, 2, len(r.Routes()))
	r.Pop()
	assertEqual(t, 1, len(r.Routes()))
	r.Pop()
	assertEqual(t, 0, len(r.Routes()))
	r.Pop()
	assertEqual(t, 0, len(r.Routes()))
}

func TestAppend(t *testing.T) {
	r := new(Router)
	r.Append(&testRoute{}, &testRoute{})
	assertEqual(t, 2, len(r.Routes()))
	r.Append(&testRoute{})
	assertEqual(t, 3, len(r.Routes()))
	r.Append()
	assertEqual(t, 3, len(r.Routes()))
}

func TestReplace(t *testing.T) {
	o1 := new(testRoute)
	o2 := new(testRoute)
	n1 := new(testRoute)
	n2 := new(testRoute)
	n3 := new(testRoute)
	r := Router([]Route{o1, o2})
	assertEqual(t, o1, r.Routes()[0])
	assertEqual(t, o2, r.Routes()[1])
	r.Replace([]Route{n1, n2, n3})
	assertEqual(t, n1, r.Routes()[0])
	assertEqual(t, n2, r.Routes()[1])
	assertEqual(t, n3, r.Routes()[2])
}

func TestClear(t *testing.T) {
	r := Router([]Route{&testRoute{}, &testRoute{}})
	assertEqual(t, 2, len(r.Routes()))
	r.Clear()
	assertEqual(t, 0, len(r.Routes()))
}

func ExampleRouter() {
	host := func(host string) Matcher {
		return MatcherFunc(func(r *http.Request) bool { return r.Host == host })
	}
	stringer := func(s string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, s)
		})
	}
	r := new(Router)
	test := func(url string) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, url, nil))
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	r.Match(host("foo.com")).Handle(stringer("Hello from foo!"))
	test("http://foo.com")
	test("http://bar.com")

	r.Match(host("bar.com")).Handle(stringer("Hello from bar!"))
	test("http://foo.com")
	test("http://bar.com")

	r.Pop()
	test("http://foo.com")
	test("http://bar.com")

	r.Clear()
	test("http://foo.com")
	test("http://bar.com")

	r.MatchFunc(func(r *http.Request) bool { return true }).HandleFunc(http.NotFound)

	// Output:
	// Hello from foo!
	// 404 page not found
	//
	// Hello from foo!
	// Hello from bar!
	// Hello from foo!
	// 404 page not found
	//
	// 404 page not found
	//
	// 404 page not found
	//
}
