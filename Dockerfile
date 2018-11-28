FROM golang:alpine
WORKDIR /go/src/github.com/bjjb/starfish
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o starfish .
FROM scratch
COPY --from=0 /go/src/github.com/bjjb/starfish/starfish .
EXPOSE 80
CMD ["./starfish"]
