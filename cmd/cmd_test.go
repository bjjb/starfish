package cmd

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	cmd := New(Use("foo"))
	assertEquals(t, cmd.Use, "foo")
	assertEquals(t, os.Stdin, cmd.InOrStdin())
	assertEquals(t, os.Stdout, cmd.OutOrStdout())
	assertEquals(t, os.Stderr, cmd.ErrOrStderr())
}

func assertEquals(t *testing.T, a ...interface{}) {
	t.Helper()
	if len(a) < 2 {
		return
	}
	if a[0] != a[1] {
		t.Errorf("expected %v to equal %v", a[0], a[1])
		return
	}
	assertEquals(t, a[1:]...)
}
