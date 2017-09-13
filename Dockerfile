FROM alpine
ENV GOPATH /usr/local
RUN apk add --no-cache musl-dev git go 
RUN mkdir -p /usr/local/src/github.com/bjjb/starfish
COPY main.go /usr/local/src/github.com/bjjb/starfish/
RUN cd /usr/local/src/github.com/bjjb/starfish && go install
EXPOSE 80
CMD ["starfish"]
