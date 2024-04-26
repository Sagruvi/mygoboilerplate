package consumer

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"main/internal/entity"
	pb "main/internal/proto"
)

type UserConsumer interface {
	CheckUser(login *pb.AuthOrLogin) error
	CreateUser(user *pb.User) (entity.User, error)
}

type userConsumer struct {
	pb.UserClient
}

func NewUserConsumer(url string) UserConsumer {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserClient(conn)
	return &userConsumer{
		client,
	}
}

func (u *userConsumer) CheckUser(login *pb.AuthOrLogin) error {
	user, err := u.UserClient.Get(context.Background(), login)
	if err != nil {
		return err
	}

	if user.GetPassword() != login.GetPassword() || user.GetEmail() != login.GetEmail() {
		return err
	}
	return nil
}

func (u *userConsumer) CreateUser(user *pb.User) (entity.User, error) {
	us, err := u.UserClient.Register(context.Background(), user)
	if err != nil {
		return entity.User{}, nil
	}
	res := entity.User{
		Password: us.GetPassword(),
		Email:    us.GetName(),
		Id:       int(us.GetId()),
	}
	return res, nil
}
