.PHONY:
pb-generate:
	protoc internal/pkg/blogpb/blog.proto --go_out=plugins=grpc:.

.PHONY:
run-server:
	go run ./cmd/server/main.go

.PHONY:
run-client:
	go run ./cmd/client/main.go

.PHONY:
docker-build-server:
	rm -f ./bin/grpc-server/server
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/grpc-server/server ./cmd/server/main.go
	docker image build -t nikasdocker/go-grpc-crud-server .

.PHONY:
db-run:
	docker container run -d -p 27017-27019:27017-27019 --name mongodb mongo:4.2.3

.PHONY:
db-attach:
	docker container exec -it mongodb bin/bash

.PHONY:
db-clean:
	docker container rm -f mongodb
