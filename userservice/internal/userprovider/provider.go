package userprovider

import (
	"context"
	"google.golang.org/grpc"
	"log"
	mygoboilerplate "main/proto"
	"net"

	"main/internal/entity"
	"main/internal/service"
)

type Provider interface {
	GetUser(ctx context.Context, req *mygoboilerplate.AuthOrLogin) *mygoboilerplate.User
	Register(ctx context.Context, req *mygoboilerplate.User) *mygoboilerplate.User
}
type userProvider struct {
	service service.Servicer
	server  mygoboilerplate.UnimplementedUserServer
}

func NewUserProvider(service service.Servicer) Provider {
	return &userProvider{
		service,
		mygoboilerplate.UnimplementedUserServer{},
	}
}
func (u *userProvider) GetUser(ctx context.Context, req *mygoboilerplate.AuthOrLogin) *mygoboilerplate.User {
	user, err := u.service.GetUser(req.Email, req.Password)
	if err != nil {
		log.Fatal(err)
	}
	pbUser := mygoboilerplate.User{
		Id:       int64(user.Id),
		Name:     user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
	return &pbUser
}
func (u *userProvider) Register(ctx context.Context, req *mygoboilerplate.User) *mygoboilerplate.User {
	user := entity.User{
		Username: req.Name,
		Password: req.Password,
	}
	err := u.service.CreateUser(user)
	if err != nil {
		log.Fatal(err)
	}
	pbUser := mygoboilerplate.User{
		Id:       int64(user.Id),
		Name:     user.Username,
		Password: user.Password,
	}
	return &pbUser
}
func (u *userProvider) Run(port string) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	mygoboilerplate.RegisterUserServer(server, u.server)
	return server.Serve(listen)
}
