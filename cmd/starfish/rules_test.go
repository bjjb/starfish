package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_parseRule(t *testing.T) {
	if !ruleRegexp.MatchString("s foo.com foo") {
		t.Fatalf("rules expression is wrong %s", ruleRegexp)
	}
	for _, tc := range []struct {
		in, err string
		rule    *rule
	}{
		{"s foo.com foo", "", &rule{"foo.com", "serve", "foo"}},
		{"f foo.com http://foo", "", &rule{"foo.com", "forward", "http://foo"}},
	} {
		t.Run(tc.in, func(t *testing.T) {
			actual, err := parseRule(strings.NewReader(tc.in))
			if err != nil {
				if tc.err == "" {
					t.Fatalf("expected no error, got %q", err)
				}
				if tc.err != err.Error() {
					t.Fatalf("expected an error of %q. got %q", tc.err, err.Error())
				}
			}
			if actual == nil && tc.rule != nil {
				t.Fatal("expected a result, got nil")
			}
			if actual != nil && tc.rule == nil {
				t.Fatal("got result, when we didn't expect one")
			}
			if actual.Host != tc.rule.Host {
				t.Fatalf("expected host to be %q, got %q", tc.rule.Host, actual.Host)
			}
		})
	}
}

func Test_parseConfigFile(t *testing.T) {
	t.Skip()
	f, _ := ioutil.TempFile("", "cfg")
	defer os.Remove(f.Name())
	ioutil.WriteFile(f.Name(), []byte(`
		foo.com serve hello
		bar.com forward blah
	`), 0644)
	rules, err := parseConfigFile(f.Name())
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}
