package service

import (
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/router"
	model_manager "github.com/ohmpatel1997/findhotel-geolocation/internal/model-manager"
)

type GeoLocationService interface {
	GetIPData(*GetRequest) (GeoLocationResponse, *router.HttpError)
}

type geolocation struct {
	l       log.Logger
	manager model_manager.GeoLocationManager
}

func NewGeolocationService(l log.Logger, mn model_manager.GeoLocationManager) GeoLocationService {
	return &geolocation{
		l, mn,
	}
}

func (g *geolocation) GetIPData(request *GetRequest) (GeoLocationResponse, *router.HttpError) {
	var resp GeoLocationResponse

	if len(request.IP) == 0 {
		return resp, router.NewHttpError("invalid ip", 400)
	}

	data, found, err := g.manager.FindDataByIP(request.IP)
	if err != nil {
		return resp, router.NewHttpError(err.Error(), 500)
	}

	if !found {
		return resp, router.NewHttpError("not found", 404)
	}

	return GeoLocationResponse{
		IP:           data.IP,
		CountryCode:  data.CountryCode,
		Country:      data.Country,
		City:         data.City,
		Latitude:     data.Latitude,
		Longitude:    data.Longitude,
		MysteryValue: data.MysteryValue,
	}, nil
}
