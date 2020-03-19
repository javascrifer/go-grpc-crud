package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/javascrifer/go-grpc-crud/internal/pkg/blogpb"
	"github.com/javascrifer/go-grpc-crud/internal/pkg/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	dbURI            = "mongodb://localhost:27017"
	dbName           = "bloggrpc"
	dbCollectionName = "blog"
	listenerNetwork  = "tcp"
	listenerAddress  = "0.0.0.0:50051"
)

func newMongoClient() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(dbURI)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	fmt.Println("initializing server")

	dbClient, err := newMongoClient()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
		return
	}
	collection := dbClient.Database(dbName).Collection(dbCollectionName)

	listener, err := net.Listen(listenerNetwork, listenerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	srv := server.NewGRPCServer(collection)
	blogpb.RegisterBlogServiceServer(s, srv)

	go func() {
		fmt.Println("starting a server")
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	fmt.Println("stopping server, database and listener")
	s.Stop()
	dbClient.Disconnect(context.Background())
	listener.Close()
}
