package provider

import (
	"gopkg.in/webdeskltd/dadata.v2"
	"mygoboilerplate/internal/geolocation/entity"
	"mygoboilerplate/internal/metrics"
	"strconv"
	"time"
)

type GeoProvider interface {
	AddressSearch(input string) ([]*entity.Address, error)
	GeoCode(lat, lng string) ([]*entity.Address, error)
}

func Geocode(lat, lng string) ([]*entity.Address, error) {
	start := time.Now()
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")
	lt, err := strconv.ParseFloat(lat, 32)
	if err != nil {
		return nil, err
	}
	lg, err := strconv.ParseFloat(lng, 32)
	if err != nil {
		return nil, err
	}
	req := dadata.GeolocateRequest{
		Lat:          float32(lt),
		Lon:          float32(lg),
		Count:        5,
		RadiusMeters: 100,
	}
	addresses, err := api.GeolocateAddress(req)
	if err != nil {
		return nil, err
	}
	metrics.ExternalAPIAccessDuration.WithLabelValues("dadata", "geocode").Observe(time.Since(start).Seconds())
	Addresses := []*entity.Address{{Lat: addresses[0].Data.City, Lng: addresses[0].Data.Street + " " + addresses[0].Data.House}}
	return Addresses, nil
}
func AddressSearch(input string) ([]*entity.Address, error) {
	start := time.Now()
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	addresses, err := api.CleanAddresses(input)
	if err != nil {
		return nil, err
	}
	res := []*entity.Address{{Lat: addresses[0].GeoLat, Lng: addresses[0].GeoLon}}
	metrics.CacheAccessDuration.WithLabelValues("dadata", "search").Observe(time.Since(start).Seconds())
	return res, nil
}
