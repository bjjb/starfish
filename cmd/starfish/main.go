// Package main defines a CLI which can run or use a starfish server.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bjjb/starfish"
	"github.com/bjjb/starfish/cmd"
	"github.com/spf13/cobra"
)

// Name is the name used as the programme name
var Name = "starfish"

// Version of the programme
var Version = "0.0.1"

func main() {
	cli().Execute()
}

func cli() *cobra.Command {
	return cmd.New(
		cmd.Name(Name),
		cmd.Version(Version),
		cmd.Summary("A small, simple HTTP frontend server/client"),
		cmd.Description(`
A client/server application which can run or configure a HTTP frontend
server.`),
		cmd.Commands(server(), completion()),
	)
}

func server() *cobra.Command {
	handler := starfish.New()
	server := &http.Server{
		Addr:    getenv("STARFISH_ADDR", ""),
		Handler: handler,
	}
	tls := struct{ cert, key string }{
		cert: getenv("STARFISH_TLS_CERT", ""),
		key:  getenv("STARFISH_TLS_KEY", ""),
	}
	return cmd.New(
		cmd.Name("server"),
		cmd.Aliases("serve", "daemon"),
		cmd.Summary("a simple, dynamic HTTP proxy and file-server"),
		cmd.Description(`
Starfish listens for HTTP requests and decides what to do based on dynamically
configurable rules. Rules consist of a matcher and a handler; the first rule
which matches will handle the request; no further rules are checked.

Matchers can check any HTTP headers (defaulting to Host), the request method,
or the request URL.

The types of handlers available are: 'serve' (which serves static files from a
directory location, which must be readable by Starfish) and 'proxy' (which
proxies requests to a target server). There's another handler type: 'api',
which serves a REST API that allows the rules to be configured on the fly;
this handler also serves its API in OpenAPI v3 format. It can be protected
either by using matcher rules, or by configuring the handler itself.

Configuration may be initialized from JSON, YAML or TOML file, or in a concise
custom format specific to Starfish. These files must be readable by Starfish
when it starts, at a file:// or a http(s):// URL.

If Starfish is running in Docker, it can be told to read configuration from
service labels or container labels, thus reconfiguring itself dynamically when
services or containers are started or stopped. These labels are prefixed with
"starfish.", and Starfish makes every effort to use Docker defaults such as
discovering the exposed port and honouring health-checks.

You can tell Starfish to expose Prometheus metrics, making it a handy edge
server for balanced apps. You can also give it an email address to use to
automatically generate Let's Encrypt certificates for domains; if you're
running multiple instances of Starfish in front of one cluster, be sure to
specify a shared storage location.`),
		cmd.Args(cobra.NoArgs),
		cmd.String(
			&server.Addr,
			"address",
			"a",
			server.Addr,
			"the address on which to listen",
		),
		cmd.String(
			&tls.cert,
			"tls.cert",
			"",
			tls.cert,
			"the location of a TLS certificate file",
		),
		cmd.String(
			&tls.key,
			"tls.key",
			"",
			tls.key,
			"the location of a TLS server key file",
		),
		cmd.Run(func(c *cobra.Command, args []string) {
			fmt.Fprintf(c.OutOrStdout(), "%s v%s listening on %s...\n", Name, Version, server.Addr)
			server.ListenAndServe()
		}),
	)
}

func completion() *cobra.Command {
	bash := cmd.New(
		cmd.Name("bash"),
		cmd.Summary("command-line completion for Bash"),
		cmd.Description("Generate a command-line completion bash function."),
		cmd.Args(cobra.NoArgs),
		cmd.Run(func(c *cobra.Command, args []string) {
			if err := c.Root().GenBashCompletion(c.OutOrStdout()); err != nil {
				fmt.Fprintln(c.ErrOrStderr(), err)
			}
		}),
	)
	fish := cmd.New(
		cmd.Name("bash"),
		cmd.Summary("command-line completion for Fish"),
		cmd.Description("Generate a command-line completion Fish function."),
		cmd.Args(cobra.NoArgs),
		cmd.Run(func(c *cobra.Command, args []string) {
			if err := c.Root().GenFishCompletion(c.OutOrStdout(), true); err != nil {
				fmt.Fprintln(c.ErrOrStderr(), err)
			}
		}),
	)
	return cmd.New(
		cmd.Name("completion"),
		cmd.Summary("command-line completion functions"),
		cmd.Description("Generate command-line completion functions for shells"),
		cmd.Hidden(),
		cmd.Commands(bash, fish),
	)
}

func getenv(k, fallback string) string {
	if v, found := os.LookupEnv(k); found {
		return v
	}
	return fallback
}
