package starfish

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleProxy() {
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from upstream!")
	})
	fooServer := httptest.NewServer(fooHandler)
	mainHandler := Proxy(fooServer.URL)
	r := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
	w := httptest.NewRecorder()
	mainHandler.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// Hello from upstream!
}

func TestProxy(t *testing.T) {
	t.Run("a bad url gives a bad gateway", func(t *testing.T) {
		x := Proxy(":") // invalid URL
		r := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
		w := httptest.NewRecorder()
		x.ServeHTTP(w, r)
		statusCode := w.Result().StatusCode
		if statusCode != http.StatusBadGateway {
			t.Errorf("expected %d, got %d", http.StatusBadGateway, statusCode)
		}
	})
}
