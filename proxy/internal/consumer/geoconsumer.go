package consumer

import (
	"context"
	"google.golang.org/grpc"
	"main/internal/entity"
	pb "main/internal/proto"
)

type GeoConsumer struct {
	pb.RpcClient
}

func NewGeoConsumer(port string) GeoConsumer {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return GeoConsumer{
		pb.NewRpcClient(conn),
	}
}
func (c GeoConsumer) Geocode(address entity.GeocodeRequest) (string, error) {
	req := pb.Request{Lat: float32(address.Lat), Lon: float32(address.Lng)}
	res, err := c.RpcClient.Geocode(context.Background(), &req, nil)
	if err != nil {
		return "", err
	}
	return res.Input, nil
}
func (c GeoConsumer) AddressSearch(address string) ([]entity.Address, error) {
	req := pb.Input{Input: address}
	res, err := c.RpcClient.AddressSearch(context.Background(), &req, nil)
	if err != nil {
		return nil, err
	}
	addresses := make([]entity.Address, 0)
	for _, v := range res.Addresses {
		addresses = append(addresses, entity.Address{Lat: v.Lat, Lng: v.Lon})
	}
	return addresses, nil
}
