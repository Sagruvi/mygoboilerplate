package gateway

import (
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	"main/internal/entity"
	"net/http"
)

type GatewayFactory interface {
	GetGateway(service, port string) Gateway
}
type Factory struct{}

func (f Factory) GetGateway(service, port string) Gateway {
	switch service {
	case "geoservice":
		return NewGeoGateway(service + ":" + port)
	case "authservice":
		return NewAuthGateway(service + ":" + port)
	case "userservice":
		return NewUserGateway(service + ":" + port)
	default:
		return nil
	}
}

type Gateway interface {
	Proxy1(w http.ResponseWriter, r *http.Request)
	Proxy2(w http.ResponseWriter, r *http.Request)
}
type GeoGateway struct {
	client *grpc.ClientConn
}

func NewGeoGateway(url string) Gateway {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &GeoGateway{
		client: conn,
	}
}
func (g *GeoGateway) Proxy1(w http.ResponseWriter, r *http.Request) {
	var request entity.GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
	}

}
func (g *GeoGateway) Proxy2(w http.ResponseWriter, r *http.Request) {
	var request entity.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
	}
}

type AuthGateway struct {
	client *grpc.ClientConn
}

func NewAuthGateway(url string) Gateway {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &AuthGateway{
		client: conn,
	}
}
func (g *AuthGateway) Proxy1(w http.ResponseWriter, r *http.Request) {
	var res entity.User
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		log.Println(err)
	}
}
func (g *AuthGateway) Proxy2(w http.ResponseWriter, r *http.Request) {
	var res entity.User
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		log.Println(err)
	}
}

type UserGateway struct {
	client *grpc.ClientConn
}

func NewUserGateway(url string) Gateway {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &UserGateway{
		client: conn,
	}
}
func (g *UserGateway) Proxy1(w http.ResponseWriter, r *http.Request) {
	var res entity.User
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		log.Println(err)
	}
}
func (g *UserGateway) Proxy2(w http.ResponseWriter, r *http.Request) {

}
