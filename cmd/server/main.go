package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/javascrifer/go-grpc-crud/internal/pkg/blogpb"
	"github.com/javascrifer/go-grpc-crud/internal/pkg/infrastructure"
	"github.com/javascrifer/go-grpc-crud/internal/pkg/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server/dev.toml", "path to server config file")
}

func main() {
	log.Println("parsing flags")
	flag.Parse()

	log.Println("initializing config")
	config := infrastructure.NewServerConfig()
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	log.Println("connecting to database")
	dbClient, err := newMongoClient(config.Mongo.URL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	collection := dbClient.Database(config.Mongo.DB).Collection(config.Mongo.Collection)

	log.Println("initializing listener")
	listener, err := net.Listen(config.Listener.Network, config.Listener.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := server.NewGRPCServer(collection)
	blogpb.RegisterBlogServiceServer(s, srv)

	go func() {
		log.Println("starting a server")
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Println("stopping server, database and listener")
	s.Stop()
	dbClient.Disconnect(context.Background())
	listener.Close()
}

func newMongoClient(dbURI string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(dbURI)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConnect()
	if err := client.Connect(connectCtx); err != nil {
		return nil, err
	}

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
