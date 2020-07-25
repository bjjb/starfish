package starfish

import (
	"net/http"
	"testing"
)

func TestStatic(t *testing.T) {
	t.Run("returns a new file server", func(t *testing.T) {
		x := Static(".")
		if typ, ok := x.(http.Handler); !ok {
			t.Errorf("expected a http.FileServer, got a %t", typ)
		}
	})
}
