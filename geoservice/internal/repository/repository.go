package repository

import (
	"context"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"main/internal/entity"
	"main/metrics"
	"time"
)

type RepositoryCacher interface {
	CacheSearchHistory(response entity.SearchResponse) error
	CacheAddress(address *entity.Address) error
}
type Repositorer interface {
	CacheSearchHistory(request string) error
	CacheAddress(geocodeResponse entity.GeocodeResponse) error
	GetSearchHistory(response entity.SearchResponse) (entity.SearchRequest, error)
	GetCache(request string) (entity.SearchResponse, error)
}
type Repository struct {
	client *redis.Client
}

func NewRepository() Repositorer {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Repository{client: client}
}

func (r *Repository) CacheSearchHistory(request string) error {
	start := time.Now()
	redisErr := r.client.Set(context.Background(), "search_history"+request, request, 0)

	if redisErr.Err() != nil {
		return redisErr.Err()
	}
	metrics.CacheAccessDuration.WithLabelValues("redis", "set").Observe(time.Since(start).Seconds())
	return nil
}

func (r *Repository) CacheAddress(geocodeResponse entity.GeocodeResponse) error {
	start := time.Now()
	redisErr := r.client.Set(context.Background(), "address"+geocodeResponse.Addresses[0].Lat+" "+geocodeResponse.Addresses[0].Lng, geocodeResponse, 0)
	if redisErr.Err() != nil {
		return redisErr.Err()
	}
	metrics.CacheAccessDuration.WithLabelValues("redis", "cache").Observe(time.Since(start).Seconds())
	return nil
}
func (r *Repository) GetSearchHistory(response entity.SearchResponse) (entity.SearchRequest, error) {
	start := time.Now()
	var searchRequest entity.SearchRequest
	res := r.client.Get(context.Background(), "search_history"+response.Addresses[0].Lat+" "+response.Addresses[0].Lng)
	if res.Err() != nil {
		return entity.SearchRequest{}, res.Err()
	}
	if res.Val() != "" {
		metrics.CacheAccessDuration.WithLabelValues("redis", "getSearchHistory").Observe(time.Since(start).Seconds())
		err := json.Unmarshal([]byte(res.Val()), &searchRequest)
		if err != nil {
			return entity.SearchRequest{}, err
		}
	}
	metrics.CacheAccessDuration.WithLabelValues("redis", "getSearchHistory").Observe(time.Since(start).Seconds())

	return searchRequest, nil
}

func (r *Repository) GetCache(request string) (entity.SearchResponse, error) {

	start := time.Now()
	var searchResponse entity.SearchResponse

	res := r.client.Get(context.Background(), "address"+request)
	metrics.CacheAccessDuration.WithLabelValues("redis", "getCache").Observe(time.Since(start).Seconds())
	if res.Err() != nil {
		return entity.SearchResponse{}, res.Err()
	}
	if res.Val() != "" {
		err := json.Unmarshal([]byte(res.Val()), &searchResponse)
		if err != nil {
			return entity.SearchResponse{}, err
		}
	}

	metrics.CacheAccessDuration.WithLabelValues("db", "getCache").Observe(time.Since(start).Seconds())
	return searchResponse, nil
}
