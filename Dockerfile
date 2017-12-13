FROM scratch
COPY main /starfish
ENTRYPOINT ["/starfish"]
EXPOSE 80
