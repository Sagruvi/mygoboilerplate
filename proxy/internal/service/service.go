package service

import pb "main/internal/proto"

type Servicer interface {
	Geocode(address pb.Address)
	Search(input pb.Input)
	Profile(user pb.User)
	List()
	Register(login pb.AuthOrLogin)
	Login(login pb.AuthOrLogin)
}
type Service struct {
}
