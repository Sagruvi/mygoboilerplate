package consumer

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"main/internal/entity"
	pb "main/internal/proto"
)

type GeoConsumer interface {
	Search(input *pb.Input) (entity.Address, error)
	Geocode(address *pb.Request) (entity.Address, error)
}

type geoConsumer struct {
	pb.RpcClient
}

func NewGeoConsumer(url string) GeoConsumer {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	defer conn.Close()
	return &geoConsumer{
		RpcClient: pb.NewRpcClient(conn),
	}
}
func (g *geoConsumer) Search(input *pb.Input) (entity.Address, error) {
	res, err := g.RpcClient.AddressSearch(context.Background(), input)
	if err != nil {
		return entity.Address{}, err
	}
	return entity.Address{
		Lat: res.String(),
		Lng: res.String(),
	}, nil
}
func (g *geoConsumer) Geocode(address *pb.Request) (entity.Address, error) {
	res, err := g.RpcClient.Geocode(context.Background(), address)
	if err != nil {
		return entity.Address{}, err
	}
	return entity.Address{
		Lat: res.String(),
		Lng: res.String(),
	}, nil
}
