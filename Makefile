.PHONY: clean image

USER=bjjb
REPO=starfish

default: image

image: main.go
	docker build -t bjjb/starfish .

push: image
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	docker push $(USER)/$(REPO)

clean:
	rm -f main starfish
