.PHONY: clean install image

build: starfish

starfish: main.go
	go build -d -v ./...

install:
	go install -v ./...

test:
	go test ./...

image: Dockerfile
	docker build -t bjjb/starfish .

clean:
	rm -rf starfish
