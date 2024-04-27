package consumer

import (
	"context"
	"google.golang.org/grpc"
	"main/internal/entity"
	pb "main/internal/proto"
)

type AuthConsumer struct {
	pb.UserClient
}

func NewAuthConsumer(port string) AuthConsumer {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return AuthConsumer{
		pb.NewUserClient(conn),
	}
}

func (a AuthConsumer) Register(user entity.User) (entity.User, error) {
	us, err := a.UserClient.Register(context.Background(), &pb.User{
		Name:     user.Username,
		Email:    user.Email,
		Password: user.Password,
		Id:       int64(user.Id),
	})
	if err != nil {
		return entity.User{}, err
	}
	return entity.User{
		Id:       int(us.Id),
		Username: us.Name,
		Email:    us.Email,
		Password: us.Password,
	}, nil
}

func (a AuthConsumer) Login(email, password string) (entity.User, error) {
	us, err := a.UserClient.Get(context.Background(), &pb.AuthOrLogin{
		Email:    email,
		Password: password,
	}, nil)
	if err != nil {
		return entity.User{}, err
	}
	return entity.User{
		Id:       int(us.Id),
		Username: us.Name,
		Email:    us.Email,
		Password: us.Password,
	}, nil
}
