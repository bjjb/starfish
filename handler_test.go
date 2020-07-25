package starfish

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadGateway(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()
	BadGateway.ServeHTTP(w, r)
	actual := w.Result().StatusCode
	if actual != http.StatusBadGateway {
		t.Errorf("expected %d, got %d", http.StatusBadGateway, actual)
	}
}
