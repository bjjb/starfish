package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

var (
	bind    = flag.String("b", ":http", "bind address")
	cfg     = flag.String("c", "starfish.cfg", "configuration file")
	verbose = flag.Bool("v", false, "print more info while running")
)

// main starts the application; parsing command-line flags, adding the signal
// handler, loading the rules, and starting the server.
func main() {
	flag.Parse()
	s := &Router{cfg: *cfg}
	go func() {
		signal.Notify(s.Signal(), syscall.SIGUSR1)
	}()
	s.Reload()
	log.Printf("Listening on %s", *bind)
	log.Fatal(http.ListenAndServe(*bind, s))
}

// A Router is the http.Handler which proxies to the handlers in its rules
// table.
type Router struct {
	mutex sync.RWMutex
	cfg   string
	rules map[string]http.Handler
}

// Reload reloads the router's configuration.
func (router *Router) Reload() {
	b, err := ioutil.ReadFile(router.cfg)
	if err != nil {
		log.Printf("Error loading rules: %s", err)
		return
	}

	rules, err := LoadRules(bytes.NewReader(b))
	if err != nil {
		log.Printf("Error loading rules: %s", err)
		return
	}

	router.mutex.Lock()
	defer router.mutex.Unlock()

	router.rules = rules
	return
}

// Signal reloads the router's config in a goroutine, returning a channel upon
// which the caller can wait.
func (router *Router) Signal() chan os.Signal {
	c := make(chan os.Signal, 1)
	go func() {
		for {
			<-c
			router.Reload()
		}
	}()
	return c
}

// ServeHTTP proxies the request to a Handler in s's rules, or serves a 404.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	say("%s %s", r.Method, r.URL.String())

	router.mutex.RLock()
	defer router.mutex.RUnlock()

	host := strings.Split(r.Host, ":")[0]

	if router.rules != nil {
		if handler, found := router.rules[host]; found {
			say("=> %s", host)
			handler.ServeHTTP(w, r)
			return
		}
	}

	say("Not found")
	http.NotFound(w, r)
}

const (
	// RulePattern defines how a line in a config file should be parsed.
	RulePattern = `^(\S+)\s+(serve|forward)\s+(\S+)$`
)

// RuleRE matches a valid line in a configuration file.
var RuleRE *regexp.Regexp

func init() {
	RuleRE = regexp.MustCompile(RulePattern)
}

// LoadRules loads a set of rules from the given reader.
func LoadRules(r io.Reader) (map[string]http.Handler, error) {
	rules := make(map[string]http.Handler)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	for n, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		host, handler, err := ParseRule(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %s", n, err)
		}
		rules[host] = handler
	}
	return rules, nil
}

// ParseRule parses an individual rule and gets an appropriate http.Handler.
func ParseRule(b []byte) (string, http.Handler, error) {
	result := RuleRE.FindStringSubmatch(string(b))

	if result == nil {
		return "", nil, fmt.Errorf("`%s` - invalid syntax", b)
	}

	if len(result) != 4 {
		return "", nil, fmt.Errorf("`%s` - invalid syntax", b)
	}

	say("%s %s %s", result[2], result[1], result[3])

	switch result[2] {
	case "serve":
		return result[1], NewFileServer(result[3]), nil
	case "forward":
		return result[1], NewProxy(result[3]), nil
	}
	return "", nil, fmt.Errorf("`%s` - unknown rule (%s)", result[2], result)
}

// A FileServer is a http.FileServer with a webroot.
type FileServer struct {
	dir     string
	handler http.Handler
}

// ServeHTTP implements http.Handler for a FileServer
func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s]: %s %s", s.dir, r.Method, r.URL.String())
	s.handler.ServeHTTP(w, r)
}

// NewFileServer makes a new file server with the given webroot.
func NewFileServer(dir string) http.Handler {
	return &FileServer{
		dir:     dir,
		handler: http.FileServer(http.Dir(dir)),
	}
}

// NewProxy makes a new reverse proxy to the given upstream.
func NewProxy(upstream string) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Host = upstream
			r.URL.Scheme = "http"
		},
	}
}

func say(fmt string, a ...interface{}) {
	if verbose != nil && *verbose {
		log.Printf(fmt, a...)
	}
}
