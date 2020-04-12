FROM alpine:latest

EXPOSE 50051

WORKDIR /grpc-server

ADD /bin/grpc-server/server ./server
ADD /configs/server/docker.toml ./config.toml

CMD ["./server", "-config-path=config.toml"]
