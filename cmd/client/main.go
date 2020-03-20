package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/javascrifer/go-grpc-crud/internal/pkg/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	targetAddress = "localhost:50051"
)

func main() {
	log.Println("initializing client")

	cc, err := grpc.Dial(targetAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	// createBlog(c)
	// getBlog(c)
	updateBlog(c)
}

func createBlog(c blogpb.BlogServiceClient) {
	log.Println("creating blog")
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
	}
	log.Printf("newly created blog: %v\n", res.GetBlog())
}

func getBlog(c blogpb.BlogServiceClient) {
	log.Println("receiving blog")
	id := "5e73bd378a15391517ed0cfc"
	req := &blogpb.GetBlogRequest{Id: id}

	res, err := c.GetBlog(context.Background(), req)
	if err != nil {
		logError(err)
		return
	}
	log.Printf("blog %s: %v\n", id, res.GetBlog())
}

func updateBlog(c blogpb.BlogServiceClient) {
	log.Println("updating blog")
	id := "5e73bcd78a15391517ed0cfa"
	now := time.Now().Unix()
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       id,
			AuthorId: fmt.Sprintf("author %v", now),
			Title:    fmt.Sprintf("title %v", now),
			Content:  fmt.Sprintf("content %v", now),
		},
	}

	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		logError(err)
		return
	}
	log.Printf("blog %s: %v\n", id, res.GetBlog())
}

func logError(err error) {
	s, ok := status.FromError(err)
	if ok {
		log.Fatalf("[%v] blog grpc call error: %v\n", s.Code(), s.Message())
	} else {
		log.Fatalf("blog grpc call error: %v\n", err)
	}
}
