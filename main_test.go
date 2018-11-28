package main

import (
	"net/http"
	"testing"
)

func responseAsserter(t *testing.T, name string) func(*http.Response, *http.Response) {
	return func(expected, actual *http.Response) {
		t.Run(name, func(t *testing.T) {
			if expected.Status != actual.Status {
				t.Errorf("Response status incorrect (expected %s, got %s)", expected.Status, actual.Status)
			}
			if expected.StatusCode != actual.StatusCode {
				t.Errorf("Response status code incorrect (expected %d, got %d)", expected.StatusCode, actual.StatusCode)
			}
			if expected.Proto != actual.Proto {
				t.Errorf("Response protocol incorrect (expected %s, got %s)", expected.Proto, actual.Proto)
			}
			for h, _ := range expected.Header {
				if expected.Header.Get(h) != actual.Header.Get(h) {
					t.Errorf("Response header '%s' incorrect (expected %s, got %s)", h, expected.Header.Get(h), actual.Header.Get(h))
				}
			}
		})
	}
}

func Test_loadRules(t *testing.T) {
}
