package starfish

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type API string

func (api API) Match(r *http.Request) bool {
	return string(api) == r.URL.Host
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"foo": "bar"}
	json.NewEncoder(w).Encode(data)
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
}

func TestAPI(t *testing.T) {
	api := API(*new(Router))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	api.ServeHTTP(w, r)
	resp := w.Result()
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		t.Fatal(err)
	}
}
