.PHONY:
pb-generate:
	protoc internal/pkg/blogpb/blog.proto --go_out=plugins=grpc:.

.PHONY:
run-server:
	go run ./cmd/server/main.go

.PHONY:
db-run:
	docker run -d -p 27017-27019:27017-27019 --name mongodb mongo:4.2.3

.PHONY:
db-attach:
	docker exec -it mongodb bin/bash

.PHONY:
db-stop:
	docker stop mongodb

.PHONY:
db-delete:
	docker rm mongodb