package gateway

import (
	"context"
	"log"
	"main/internal/entity"
	"main/internal/service"
	pb "main/proto"
	"net"
	"net/http"
	"net/rpc"
)

type Gateway interface {
	GeoCode(lat, lng float64, reply *pb.Input) error
	AddressSearch(ctx context.Context, input *pb.Input) (*pb.Addresses, error)
	Run(port string) error
}

type GatewayFactory interface {
	CreateGateway() Gateway
}

type JsonRPCGateway struct {
	service service.Servicer
}

type JsonRPCGatewayFactory struct {
	service service.Servicer
}

func (f JsonRPCGatewayFactory) CreateGateway() Gateway {
	return &JsonRPCGateway{service: f.service}
}

func (r JsonRPCGateway) GeoCode(lat, lng float64, reply *pb.Input) error {
	rep, err := r.service.DadataGeocodeApi(entity.GeocodeRequest{Lat: lat, Lng: lng})
	if err != nil {
		return err
	}
	reply.Input = rep
	return nil
}

func (r JsonRPCGateway) AddressSearch(ctx context.Context, input *pb.Input) (*pb.Addresses, error) {
	rep, err := r.service.DadataSearchApi(input.Input)
	if err != nil {
		return nil, err
	}
	reply := &pb.Addresses{}
	for _, address := range rep.Addresses {
		reply.Addresses = append(reply.Addresses, &pb.Address{Lat: address.Lat, Lon: address.Lng})
	}
	return reply, nil
}
func (r JsonRPCGateway) Run(port string) error {
	rpc.Register(r)

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Println("Сервер запущен на порту 8080")

	http.Serve(l, nil)
	return nil
}
