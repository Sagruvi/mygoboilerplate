package consumer

import (
	"context"
	"google.golang.org/grpc"
	"main/internal/entity"
	"main/internal/proto"
)

type UserConsumer interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	CheckUser(ctx context.Context, email, password string) (entity.User, error)
}

type UserConsumerImpl struct {
	rpc mygoboilerplate.UserClient
}

func NewUserConsumer(serverurl string) UserConsumer {
	return &UserConsumerImpl{}
}
func Run(serverurl string) mygoboilerplate.RpcClient {
	conn, err := grpc.Dial(serverurl, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := mygoboilerplate.NewRpcClient(conn)
	return c
}

func (u *UserConsumerImpl) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	us := mygoboilerplate.User{
		Email:    user.Email,
		Password: user.Password,
		Id:       int64(user.Id),
	}
	resp, err := u.rpc.Register(ctx, &us, nil)
	if err != nil {
		return entity.User{}, err
	}

	res := entity.User{
		Id:       int(resp.Id),
		Email:    resp.Email,
		Password: resp.Password,
		Username: resp.Name,
	}
	return res, nil
}

func (u *UserConsumerImpl) CheckUser(ctx context.Context, email, password string) (entity.User, error) {
	req := mygoboilerplate.AuthOrLogin{
		Email:    email,
		Password: password,
	}
	resp, err := u.rpc.Get(ctx, &req, nil)
	if err != nil {
		return entity.User{}, err
	}
	res := entity.User{
		Id:       int(resp.Id),
		Email:    resp.Email,
		Password: resp.Password,
		Username: resp.Name,
	}
	return res, nil
}
