package server

import (
	"context"
	"log"

	"github.com/javascrifer/go-grpc-crud/internal/pkg/blogpb"
	"go.mongodb.org/mongo-driver/bson"
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
		log.Printf("failed to create blog: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create blog")
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("failed to parse object id: %v", res.InsertedID)
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

// GetBlog returns one blog in the database
func (s *GRPCServer) GetBlog(
	ctx context.Context,
	req *blogpb.GetBlogRequest,
) (*blogpb.GetBlogResponse, error) {
	id := req.GetId()
	mongoBlog, err := getBlogByID(s, id)
	if err != nil {
		return nil, err
	}
	return &blogpb.GetBlogResponse{Blog: mongoBlogToGrpc(mongoBlog)}, nil
}

// UpdateBlog updates value of the blog in the database and returns updated value
func (s *GRPCServer) UpdateBlog(
	ctx context.Context,
	req *blogpb.UpdateBlogRequest,
) (*blogpb.UpdateBlogResponse, error) {
	grpcBlog := req.GetBlog()
	mongoBlog, err := getBlogByID(s, grpcBlog.GetId())
	if err != nil {
		return nil, err
	}
	mongoBlog.AuthorID = grpcBlog.GetAuthorId()
	mongoBlog.Title = grpcBlog.GetTitle()
	mongoBlog.Content = grpcBlog.GetContent()

	filter := getBlogFilter(mongoBlog.ID)
	_, err = s.blogCollection.ReplaceOne(context.Background(), filter, mongoBlog)
	if err != nil {
		log.Printf("failed to update blog: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update blog")
	}
	return &blogpb.UpdateBlogResponse{Blog: mongoBlogToGrpc(mongoBlog)}, nil
}

func getBlogByID(s *GRPCServer, id string) (*blog, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("failed to parse object id: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "incorrect id")
	}

	filter := getBlogFilter(oid)
	res := s.blogCollection.FindOne(context.Background(), filter)
	if err := res.Err(); err != nil {
		log.Printf("failed to find blog: %v", err)
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "blog not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find blog")
	}

	mongoBlog := &blog{}
	if err := res.Decode(mongoBlog); err != nil {
		log.Printf("failed to decode blog: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to decode blog")
	}

	return mongoBlog, nil
}

func getBlogFilter(oid primitive.ObjectID) primitive.M {
	return bson.M{"_id": oid}
}

func mongoBlogToGrpc(b *blog) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       b.ID.Hex(),
		AuthorId: b.AuthorID,
		Title:    b.Title,
		Content:  b.Content,
	}
}
