.PHONY:
pb-generate:
	protoc internal/pkg/blogpb/blog.proto --go_out=plugins=grpc:.

.PHONY:
run-db:
	docker container rm -f mongodb
	docker container run -d -p 27017-27019:27017-27019 --name mongodb mongo:4.2.3

.PHONY:
run-server:
	go run ./cmd/server/main.go

.PHONY:
run-server-compose:
	docker-compose up

.PHONY:
run-client:
	go run ./cmd/client/main.go

.PHONY:
build-server:
	rm -f ./bin/grpc-server/server
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/grpc-server/server ./cmd/server/main.go

.PHONY:
build-server-image:
	rm -f ./bin/grpc-server/server
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/grpc-server/server ./cmd/server/main.go
	docker image build -t nikasdocker/go-grpc-crud-server .

.PHONY:
db-attach:
	docker container exec -it mongodb bin/bash

.PHONY:
db-clean:
	docker container rm -f mongodb
