package server

import (
	"context"

	"github.com/javascrifer/go-grpc-crud/internal/pkg/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blog struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

// GRPCServer is a gRPC server which handles blog requests
type GRPCServer struct {
	blogCollection *mongo.Collection
}

// NewGRPCServer is constructor function for server
func NewGRPCServer(collection *mongo.Collection) *GRPCServer {
	return &GRPCServer{blogCollection: collection}
}

// CreateBlog stores blog in the database
func (s *GRPCServer) CreateBlog(
	ctx context.Context,
	req *blogpb.CreateBlogRequest,
) (*blogpb.CreateBlogResponse, error) {
	grpcBlog := req.GetBlog()
	mongoBlog := &blog{
		AuthorID: grpcBlog.GetAuthorId(),
		Title:    grpcBlog.GetTitle(),
		Content:  grpcBlog.GetContent(),
	}

	res, err := s.blogCollection.InsertOne(context.Background(), mongoBlog)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create blog")
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to parse id")
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: grpcBlog.GetAuthorId(),
			Title:    grpcBlog.GetTitle(),
			Content:  grpcBlog.GetContent(),
		},
	}, nil
}
