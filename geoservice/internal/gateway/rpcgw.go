package gateway

import (
	"context"
	"log"
	"main/internal/entity"
	"main/internal/service"
	pb "main/proto"
	"net"
	"net/rpc"
)

type RPCGateway struct {
	service service.Servicer
}

type RPCGatewayFactory struct {
	service service.Servicer
}

func (f RPCGatewayFactory) CreateGateway() Gateway {
	return &RPCGateway{service: f.service}
}
func (r RPCGateway) GeoCode(lat, lng float64, reply *pb.Input) error {
	rep, err := r.service.DadataGeocodeApi(entity.GeocodeRequest{Lat: lat, Lng: lng})
	if err != nil {
		return err
	}
	reply.Input = rep
	return nil
}

func (r RPCGateway) AddressSearch(context context.Context, input *pb.Input) (*pb.Addresses, error) {
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
func (r RPCGateway) Run(port string) error {
	rpc.Register(r)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Println("Сервер запущен на порту 1234")
	rpc.Accept(l)

	return nil
}
