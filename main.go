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

func main() {
	flag.Parse()
	s := &server{cfg: *cfg}
	go func() {
		signal.Notify(s.signal(), syscall.SIGUSR1)
	}()
	s.reload()
	log.Printf("Listening on %s", *bind)
	log.Fatal(http.ListenAndServe(*bind, s))
}

type server struct {
	mutex sync.RWMutex
	cfg   string
	rules map[string]http.Handler
}

func (s *server) reload() {
	b, err := ioutil.ReadFile(s.cfg)
	if err != nil {
		log.Printf("Error loading rules: %s", err)
		return
	}

	rules, err := loadRules(bytes.NewReader(b))
	if err != nil {
		log.Printf("Error loading rules: %s", err)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.rules = rules
	return
}

func (s *server) signal() chan os.Signal {
	c := make(chan os.Signal, 1)
	go func() {
		for {
			<-c
			s.reload()
		}
	}()
	return c
}

// ServeHTTP proxies the request to a Handler in s's rules, or serves a 404.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	say("%s %s", r.Method, r.URL.String())

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	host := strings.Split(r.Host, ":")[0]

	if s.rules != nil {
		if handler, found := s.rules[host]; found {
			say("=> %s", host)
			handler.ServeHTTP(w, r)
			return
		}
	}

	say("Not found")
	http.NotFound(w, r)
}

const (
	rulePattern = `^(\S+)\s+(serve|forward)\s+(\S+)$`
)

var ruleRE *regexp.Regexp

func init() {
	ruleRE = regexp.MustCompile(rulePattern)
}

func loadRules(r io.Reader) (map[string]http.Handler, error) {
	rules := make(map[string]http.Handler)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	for n, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		host, handler, err := parseRule(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %s", n, err)
		}
		rules[host] = handler
	}
	return rules, nil
}

func parseRule(b []byte) (string, http.Handler, error) {
	result := ruleRE.FindStringSubmatch(string(b))

	if result == nil {
		return "", nil, fmt.Errorf("`%s` - invalid syntax", b)
	}

	if len(result) != 4 {
		return "", nil, fmt.Errorf("`%s` - invalid syntax", b)
	}

	say("%s %s %s", result[2], result[1], result[3])

	switch result[2] {
	case "serve":
		return result[1], makeServer(result[3]), nil
	case "forward":
		return result[1], makeProxy(result[3]), nil
	}
	return "", nil, fmt.Errorf("`%s` - unknown rule", result[2], result)
}

type fileServer struct {
	dir     string
	handler http.Handler
}

func (s *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s]: %s %s", s.dir, r.Method, r.URL.String())
	s.handler.ServeHTTP(w, r)
}

func makeServer(dir string) http.Handler {
	return &fileServer{
		dir:     dir,
		handler: http.FileServer(http.Dir(dir)),
	}
}

func makeProxy(upstream string) http.Handler {
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
