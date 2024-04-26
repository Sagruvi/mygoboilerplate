package provider

import (
	"google.golang.org/grpc"
	"log"
	"main/internal/service"
	pb "main/proto"
	"net"
)

type AuthProvider interface {
	Login(reply *pb.User)
	Register(reply *pb.User)
}

type Provider struct {
	service service.Servicer
}

func Run(port string) {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	server := grpc.NewServer()
	pb.RegisterRpcServer(server, &pb.UnimplementedRpcServer{})

	log.Println("Запуск gRPC сервера...")
	if err := server.Serve(listen); err != nil {
		log.Println(err)
		panic(err)
	}

}
func (p *Provider) Login(loginRequest *pb.AuthOrLogin) (*pb.User, error) {
	rep, err := p.service.SaveUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Println(err)
		return &pb.User{}, err
	}
	res := pb.User{
		Id:       int64(rep.Id),
		Name:     rep.Username,
		Password: "",
		Email:    rep.Email,
	}
	return &res, nil
}
func (p *Provider) Register(loginRequest *pb.AuthOrLogin) (*pb.User, error) {
	err := p.service.CheckUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Println(err)
		return &pb.User{}, err
	}
	rep, err := p.service.SaveUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		log.Println(err)
		return &pb.User{}, err
	}
	res := pb.User{
		Id:       int64(rep.Id),
		Name:     rep.Username,
		Password: rep.Password,
		Email:    rep.Email,
	}
	return &res, nil
}
