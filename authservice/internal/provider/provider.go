package provider

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"main/internal/entity"
	pb "main/internal/proto"
	"main/internal/service"
	"net"
)

type Provider interface {
	Get(ctx context.Context, req *pb.AuthOrLogin) (*pb.User, error)
	Register(ctx context.Context, req *pb.User) (*pb.User, error)
}

type AuthProvider struct {
	service.Servicer
	pb.UnimplementedUserServer
	serverPort string
}

func NewProvider(clientPort, serverPort string) *AuthProvider {
	return &AuthProvider{service.NewService(clientPort), pb.UnimplementedUserServer{}, serverPort}
}

func (a *AuthProvider) Get(ctx context.Context, req *pb.AuthOrLogin) (*pb.User, error) {
	user, err := a.Servicer.CheckUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.User{Email: user.Email, Id: int64(user.Id), Name: user.Username, Password: user.Password}, nil
}
func (a *AuthProvider) Register(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := a.Servicer.SaveUser(entity.User{
		Id:       int(req.Id),
		Email:    req.Email,
		Password: req.Password,
		Username: req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.User{Email: user.Email, Id: int64(user.Id), Name: user.Username, Password: user.Password}, nil
}

func (a *AuthProvider) Run() {
	g := grpc.NewServer()
	pb.RegisterUserServer(g, a)

	lis, err := net.Listen("tcp", ":"+a.serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := g.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
