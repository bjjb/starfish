.PHONY: clean image

default: image

main: main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

starfish: main.go
	go build .

image: main
	docker build -t bjjb/starfish .

clean:
	rm -f main starfish
