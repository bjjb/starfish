FROM golang:alpine AS builder
# Install dependencies
RUN apk --no-cache add git tzdata ca-certificates
# Build the executable
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
WORKDIR src/github.com/bjjb/starfish
ADD . .
WORKDIR cmd/starfish
RUN go mod download
RUN go mod verify
RUN go build -ldflags '-w -s' -a -installsuffix cgo -o /go/bin/starfish

# Make the final image
FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/starfish /bin/starfish
ENTRYPOINT ["/bin/starfish"]
EXPOSE 80
