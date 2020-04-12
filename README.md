# Course for gRPC usage in Go

## Motivation

Learn about gRPC, protocol buffers and improve Go skills.

## About the project

This project is a part of [course](https://www.udemy.com/course/grpc-golang/). Project contains blog service which exposes 4 different RPC calls:

- create blog - stores blog in the database.
- read blog - retrieves blog by id from the database.
- update - replaces blog in the database by given blog id.
- delete - deletes blog from the database by given id.
- list - streams all the blogs to the client one by one.

## Local development

### gRPC Server

- Start MongoDB docker container.

```bash
  make run-db
  make db-attach # attaches to container if you need to execute commands in it
```

- Start Go gRPC server.

```bash
  make run-server
```

### gRPC Client

- Start server and db using commands mentioned in **gRPC Server** section or run them using docker compose.

```bash
  make run-server-compose
```

- Start Go gRPC client.

```bash
  make run-client
```
