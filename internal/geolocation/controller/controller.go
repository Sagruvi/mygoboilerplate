package controller

import (
	"encoding/json"
	"log"
	"mygoboilerplate/internal/geolocation/entity"
	"mygoboilerplate/internal/geolocation/service"
	"mygoboilerplate/internal/metrics"
	"net/http"
	"time"
)

type Controllerer interface {
	Geocode(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}
type Controller struct {
	service service.Servicer
}
type GeocodeRequest struct {
	entity.GeocodeRequest
}
type SearchRequest struct {
	entity.SearchRequest
}

func NewController(service2 service.Service) Controller {
	return Controller{service: &service2}
}

// SearchAddress godoc
//
//	@Summary		Search for address suggestions
//	@Description	Search for address suggestions by latitude and longitude
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			lat				body		repository.GeocodeRequest	true	"Lat and Lon"
//	@Param			Authorization	header		string			true	"Authorization token"
//	@Param			X-Secret		header		string			true	"API Private token"
//
// @Param        Authorization header string true "Bearer token"
//
//	@Success		200					"Successful operation"
//	@Failure		400				"Bad request"
//	@Failure		401				"Unauthorized"
//	@Failure		404				"Not found"
//	@Failure		500				"Internal server error"
//	@Router			/geocode [post]
func (c *Controller) Geocode(w http.ResponseWriter, r *http.Request) {
	var geocodeRequest GeocodeRequest
	start := time.Now()
	err := json.NewDecoder(r.Body).Decode(&geocodeRequest.GeocodeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	geocodeResponse, err := c.service.DadataGeocodeApi(geocodeRequest.GeocodeRequest)
	metrics.RequestDuration.WithLabelValues("POST", "/geocode").Observe(time.Since(start).Seconds())
	metrics.RequestCount.WithLabelValues("POST", "/geocode").Inc()
	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// SearchAddress godoc
//
//	@Summary		Search for address
//	@Description	Search for latitude and longitude by address
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			lat				body		repository.SearchRequest	true	"Address"
//	@Param			Authorization	header		string			true	"Authorization token"
//	@Param			X-Secret		header		string			true	"API Private token"
//
// @Param        Authorization header string true "Bearer token"
//
//	@Success		200					"Successful operation"
//	@Failure		400				"Bad request"
//	@Failure		401				"Unauthorized"
//	@Failure		404				"Not found"
//	@Failure		500				"Internal server error"
//	@Router			/search [post]
func (c *Controller) Search(w http.ResponseWriter, r *http.Request) {
	var searchRequest SearchRequest
	start := time.Now()
	err := json.NewDecoder(r.Body).Decode(&searchRequest.SearchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	cachedResponse, err := c.service.GetCache(searchRequest.SearchRequest)
	if err == nil {
		err = json.NewEncoder(w).Encode(&cachedResponse)
		if err != nil {
			log.Println("data not found")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	resp, err := c.service.DadataSearchApi(searchRequest.SearchRequest.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	metrics.RequestDuration.WithLabelValues("POST", "/search").Observe(time.Since(start).Seconds())
	metrics.RequestCount.WithLabelValues("POST", "/search").Inc()
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Println("Status 500, dadata.ru is not responding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
