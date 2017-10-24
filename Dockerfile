FROM golang:alpine
WORKDIR /go/src/starfish
COPY main.go ./
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["starfish"]
EXPOSE 80
