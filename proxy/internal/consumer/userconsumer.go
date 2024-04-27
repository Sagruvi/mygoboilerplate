package consumer

import (
	"context"
	"google.golang.org/grpc"
	"main/internal/entity"
	pb "main/internal/proto"
)

type UserConsumer struct {
	pb.UserClient
}

func NewUserConsumer(port string) UserConsumer {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return UserConsumer{
		pb.NewUserClient(conn),
	}
}

func (u UserConsumer) Get(email, password string) (entity.User, error) {
	us, err := u.UserClient.Get(context.Background(), &pb.AuthOrLogin{
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
func (u UserConsumer) List() ([]entity.User, error) {
	us, err := u.UserClient.List(context.Background(), &pb.Empty{}, nil)
	if err != nil {
		return []entity.User{}, err
	}
	users := make([]entity.User, 0)
	for _, v := range us.Users {
		users = append(users, entity.User{
			Id:       int(v.Id),
			Username: v.Name,
			Email:    v.Email,
			Password: v.Password})
	}
	return users, nil
}
