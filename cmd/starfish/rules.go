package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	// ConfigOK indicates that all routes were created OK
	ConfigOK = 0
	// ConfigInvalid indicates that one or more routes couldn't be created
	ConfigInvalid = 1
)

type rule struct {
	Host   string `json:"host"`
	Action string `json:"action"`
	Target string `json:"target"`
}

var commentPattern = `\A\s*#`
var commentRegexp = regexp.MustCompile(commentPattern)
var hostPattern = `\S+`
var actionPattern = `(?:f|forward|s|serve|a|api)`
var actionRegexp = regexp.MustCompile(actionPattern)
var targetPattern = `\S+`
var rulePattern = fmt.Sprintf(`\A\s*(%s)\s+(%s)\s+(%s)\s*\z`, actionPattern, hostPattern, targetPattern)
var ruleRegexp = regexp.MustCompile(rulePattern)

func parseConfigFile(filename string) ([]*rule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return parseConfig(file)
}

func parseConfig(r io.Reader) ([]*rule, error) {
	scanner := bufio.NewScanner(r)

	rules := []*rule{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || commentRegexp.MatchString(line) {
			continue
		}
		rule, err := parseRule(strings.NewReader(line))
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func parseRule(r io.Reader) (*rule, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)

	if len(data) == 0 || commentRegexp.Match(data) {
		return nil, nil
	}

	matches := ruleRegexp.FindStringSubmatch(string(data))
	if matches == nil || len(matches) != 4 {
		return nil, fmt.Errorf("malformed rule: %s (%v)", data, matches)
	}

	action, host, target := matches[1], matches[2], matches[3]

	return &rule{
		Host:   host,
		Action: action,
		Target: target,
	}, nil
}
