# Starfish

A silly little HTTP reverse-proxy and file server.

Starfish listens for TCP connections, and then either serves a directory or
proxies to an upstream. It's configured with a simple file (starfish.cfg, by
default), and reloaded with a SIGUSR1.

## Installation

    go install github.com/bjjb/starfish

## Usage

    starfish

... will start the server on the default port, i.e. 80 if you have
permission. This can be overridden with the `-b` flag, for example:

    starfish -b :8080

For each entry in the config, it will create a HTTP handler which is either a
file server or a proxy. For example, suppose you have a directory at
/var/www/example containing HTML files, and a service running on port
tcp/9000. Given the following configuration:

```
example.com serve   /var/www/example
test.com    forward localhost:9000
```

all requests to `example.com` will serve the website at `/var/www/example,`
and all requests to `test.com` will be proxied to `http://localhost:9000`.

## Future features

- Ability to load config rules from JSON
- Ability to load config rules from any io.Reader
- Specific endpoint for reconfiguring with a HTTP request
- Ability to auto-reconfigure by reading its docker environment

## Credits

This is a shameless ripoff of [this project][webfront], by [Andrew
Gerrand][nf]. The only additions were the new config file format and the
reload on signal.

[webfront]: https://github.com/nf/webfront
[nf]: https://github.com/nf
