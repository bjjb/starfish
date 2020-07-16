package starfish

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func assertStatus(t *testing.T, h http.Handler, r *http.Request, expected int) {
	if h == nil {
		h = defaultRouter
	}
	t.Helper()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()
	actual := resp.StatusCode
	if actual != expected {
		t.Errorf("Expected status code to be %d, got %d", expected, resp.StatusCode)
	}
}

func assertBody(t *testing.T, h http.Handler, r *http.Request, expected string) {
	if h == nil {
		h = defaultRouter
	}
	t.Helper()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	resp := w.Result()
	actual, _ := ioutil.ReadAll(resp.Body)
	if string(actual) != expected {
		t.Errorf("Expected body to be %q, got %q", expected, string(actual))
	}
}

type testRoute struct {
	match      bool
	statusCode int
	body       string
}

func (tr *testRoute) Match(r *http.Request) bool {
	return tr.match
}

func (tr *testRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		statusCode int
		body       string
	)
	if statusCode = tr.statusCode; statusCode == 0 {
		statusCode = http.StatusTeapot
	}
	if body = tr.body; body == "" {
		body = "Hello!"
	}
	w.WriteHeader(statusCode)
	io.Copy(w, strings.NewReader(body))
}

func writeString(text string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(w, strings.NewReader(text)); err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	})
}
