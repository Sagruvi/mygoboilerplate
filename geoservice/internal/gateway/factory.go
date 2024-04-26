package gateway

import (
	"errors"
	"main/internal/service"
)

type Factory interface {
	GetGateway() GatewayFactory
}
type RPCGW struct {
	factory GatewayFactory
}

func (r RPCGW) GetFactory(protocol string, geoService service.Servicer) (GatewayFactory, error) {

	switch protocol {
	case "RPC":
		return &RPCGatewayFactory{service: geoService}, nil
	case "JSON-RPC":
		return &JsonRPCGatewayFactory{service: geoService}, nil
	case "GRPC":
		return &GRPCGatewayFactory{service: geoService}, nil
	default:
		return nil, errors.New("unknown protocol")
	}
}
