package gateway

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"main/internal/entity"
	"main/internal/service"
	pb "main/proto"
	"net"
)

type GRPCGateway struct {
	service service.Servicer
	pb.UnimplementedRpcServer
}

func NewGRPCGateway(service service.Servicer) *GRPCGateway {
	return &GRPCGateway{service: service}
}

func (r GRPCGateway) GeoCode(lat, lng float64, reply *pb.Input) error {
	rep, err := r.service.DadataGeocodeApi(entity.GeocodeRequest{Lat: lat, Lng: lng})
	if err != nil {
		return err
	}
	reply.Input = rep
	return nil
}

func (r GRPCGateway) AddressSearch(ctx context.Context, input *pb.Input) (*pb.Addresses, error) {
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

type GRPCGatewayFactory struct {
	service service.Servicer
}

func (f *GRPCGatewayFactory) CreateGateway() Gateway {
	return NewGRPCGateway(f.service)
}
func (r GRPCGateway) Run(port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterRpcServer(server, r)

	log.Println("Запуск gRPC сервера...")
	if err := server.Serve(listen); err != nil {
		return err
	}
	return nil
}
