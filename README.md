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

- Start MongoDB docker container or install MongoDB locally to your device.

```bash
  make db-run # starts docker container
  make db-attach # attaches to the docker container using interactive mode
  make db-stop # stops docker container
  make db-delete # deletes docker container
```

- Start Go gRPC server.

```bash
  make run-server
```

- Start Go gRPC client.

```bash
  make run-client
```

P.S. All the configs for the client and server are hard-coded.
