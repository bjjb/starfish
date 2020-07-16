.PHONY: test build image clean

test:
	go test

build:
	go build

image:
	docker build -t bjjb/starfish .

clean:
	git clean -Xfd
