package userprovider

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"main/internal/entity"
	"main/internal/service"
	pb "main/proto"
	"net"
)

type UserService struct {
	service.Servicer
	pb.UnimplementedUserServer
}

func (u UserService) mustEmbedUnimplementedUserServer() {}

func (u UserService) Get(ctx context.Context, login *pb.AuthOrLogin) (*pb.User, error) {
	user, err := u.Servicer.GetUser(login.Email, login.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "method Geocode not implemented")
	}
	res := &pb.User{
		Id:       int64(user.Id),
		Name:     user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	return res, nil
}
func (u UserService) Register(ctx context.Context, login *pb.User) (*pb.User, error) {
	user := entity.User{
		Username: login.Name,
		Email:    login.Email,
		Password: login.Password,
		Id:       int(login.Id),
	}
	err := u.Servicer.CreateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "method Geocode not implemented")
	}
	return login, nil
}
func (u UserService) List(ctx context.Context, none *pb.Empty) (*pb.ListOfUsers, error) {

	users, err := u.Servicer.ListUsers()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "method Geocode not implemented")
	}
	res := &pb.ListOfUsers{
		Users: []*pb.User{},
	}
	for _, user := range users {
		res.Users = append(res.Users, &pb.User{
			Id:       int64(user.Id),
			Name:     user.Username,
			Email:    user.Email,
			Password: user.Password,
		})
	}
	return res, nil
}
func NewUserProvider() *UserService {
	return &UserService{
		Servicer: service.NewService(),
	}
}
func (u UserService) Run(port string) error {

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServer(s, u)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
