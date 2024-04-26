package service

import (
	"main/internal/entity"
	"main/internal/provider"
	"main/internal/repository"
	"main/metrics"
	"strconv"
	"time"
)

type Servicer interface {
	CacheSearchHistory(request entity.SearchRequest) error
	CacheAddress(request entity.GeocodeResponse) error
	GetSearchHistory(response entity.SearchResponse) (entity.SearchRequest, error)
	GetCache(request entity.SearchRequest) (entity.SearchResponse, error)
	DadataGeocodeApi(geocodeRequest entity.GeocodeRequest) (string, error)
	DadataSearchApi(query string) (entity.SearchResponse, error)
}

type Service struct {
	Repository repository.Repositorer
}

func NewService() Servicer {
	return &Service{repository.NewRepository()}
}

func (s *Service) CacheSearchHistory(request entity.SearchRequest) error {
	return s.Repository.CacheSearchHistory(request.Query)
}
func (s *Service) CacheAddress(request entity.GeocodeResponse) error {
	return s.Repository.CacheAddress(request)
}
func (s *Service) GetSearchHistory(response entity.SearchResponse) (entity.SearchRequest, error) {
	return s.Repository.GetSearchHistory(response)
}
func (s *Service) GetCache(request entity.SearchRequest) (entity.SearchResponse, error) {
	return s.Repository.GetCache(request.Query)
}

func (s *Service) DadataGeocodeApi(geocodeRequest entity.GeocodeRequest) (string, error) {
	address := entity.Address{
		Lat: strconv.FormatFloat(geocodeRequest.Lat, 'f', -1, 64),
		Lng: strconv.FormatFloat(geocodeRequest.Lng, 'f', -1, 64),
	}
	request := entity.SearchResponse{Addresses: []*entity.Address{&address}}
	cachedResponse, err := s.Repository.GetSearchHistory(request)
	if err == nil {
		return cachedResponse.Query, nil
	}
	addresses, err := provider.Geocode(address.Lat, address.Lng)
	if err != nil {
		return "", err
	}
	var geocodeResponse entity.GeocodeResponse
	geocodeResponse.Addresses = addresses
	res := geocodeResponse.Addresses[0].Lat + " " + geocodeResponse.Addresses[0].Lng
	err = s.CacheAddress(geocodeResponse)
	if err != nil {
		return res, err
	}
	return res, nil

}

func (s *Service) DadataSearchApi(query string) (entity.SearchResponse, error) {
	cachedResponse, err := s.GetCache(entity.SearchRequest{Query: query})
	if err == nil {
		return cachedResponse, nil
	}
	addresses, err := provider.AddressSearch(query)
	if err != nil {
		return entity.SearchResponse{}, err
	}
	var searchResponse entity.SearchResponse
	searchResponse.Addresses = addresses
	metrics.CacheAccessDuration.WithLabelValues("dadata", "search").Observe(time.Since(time.Now()).Seconds())
	err = s.CacheSearchHistory(entity.SearchRequest{Query: query})
	return searchResponse, nil
}
