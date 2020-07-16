package starfish

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func ExampleProxy() {
	fooHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from upstream!")
	})
	fooServer := httptest.NewServer(fooHandler)
	fooURL, _ := url.Parse(fooServer.URL)
	mainHandler := Proxy(fooURL)
	r := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
	w := httptest.NewRecorder()
	mainHandler.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// Hello from upstream!
}
