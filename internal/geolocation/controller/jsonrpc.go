package controller

import (
	"log"
	"mygoboilerplate/internal/geolocation/entity"
	"mygoboilerplate/internal/geolocation/service"
	pb "mygoboilerplate/proto"
	"net"
	"net/http"
	"net/rpc"
)

type Gateway interface {
	GeoCode(lat, lng float64, reply *pb.Input) error
	AddressSearch(input string, reply *pb.Addresses) error
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

func (r JsonRPCGateway) AddressSearch(input string, reply *pb.Addresses) error {
	rep, err := r.service.DadataSearchApi(input)
	if err != nil {
		return err
	}
	for _, address := range rep.Addresses {
		reply.Addresses = append(reply.Addresses, &pb.Address{Lat: address.Lat, Lon: address.Lng})
	}
	return nil
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
