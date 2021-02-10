package service

import (
	"fmt"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/model"
)

type GeoLocationService interface {
	GetIPData(*GetRequest) (GeoLocationResponse, error)
}
type geolocation struct {
	l      log.Logger
	cuder  repository.Cuder
	finder repository.Finder
}

func NewGeolocationService(l log.Logger, c repository.Cuder, f repository.Finder) GeoLocationService {
	return &geolocation{
		l, c, f,
	}
}

func (g *geolocation) GetIPData(request *GetRequest) (GeoLocationResponse, error) {
	var resp GeoLocationResponse

	if len(request.IP) == 0 {
		return resp, fmt.Errorf("empty id found")
	}

	data, found, err := g.finder.FindManaged(&model.Geolocation{IP: request.IP})
	if err != nil {
		return resp, err
	}

	if !found {
		return resp, fmt.Errorf("no data found")
	}

	v, ok := data.(*model.Geolocation)
	if !ok {
		return resp, fmt.Errorf("can not able to type assert the response: %v", data)
	}

	return GeoLocationResponse{
		IP:           v.IP,
		CountryCode:  v.CountryCode,
		City:         v.City,
		Latitude:     v.Latitude,
		Longitude:    v.Longitude,
		MysteryValue: v.MysteryValue,
	}, nil
}
