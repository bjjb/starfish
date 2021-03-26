module github.com/bjjb/starfish/cmd/starfish

go 1.14

require (
	github.com/bjjb/starfish v0.0.0-20171213154043-77bcd2ec9d97
	github.com/bjjb/starfish/cmd v0.0.0-00010101000000-000000000000
	github.com/caddyserver/certmagic v0.12.0 // indirect
	github.com/spf13/cobra v1.0.0
)

replace github.com/bjjb/starfish => ../..

replace github.com/bjjb/starfish/cmd => ../
