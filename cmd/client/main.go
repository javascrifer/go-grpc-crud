package main

import (
	"context"
	"fmt"
	"log"

	"github.com/javascrifer/go-grpc-crud/internal/pkg/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	targetAddress = "localhost:50051"
)

func main() {
	fmt.Println("Initializing client")

	cc, err := grpc.Dial(targetAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	createBlog(c)
}

func createBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "niko",
			Title:    "How are you?",
			Content:  "Just asking.",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		logError(err)
		return
	}
	log.Printf("newly created blog: %v\n", res.GetBlog())
}

func logError(err error) {
	s, ok := status.FromError(err)
	if ok {
		log.Fatalf("[%v] error while creating blog: %v\n", s.Code(), s.Message())
	} else {
		log.Fatalf("error while creating blog: %v\n", err)
	}
}
