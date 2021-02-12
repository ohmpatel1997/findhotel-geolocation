package service

import (
	"fmt"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/router"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/model"
)

type GeoLocationService interface {
	GetIPData(*GetRequest) (GeoLocationResponse, *router.HttpError)
}
type geolocation struct {
	l      log.Logger
	finder repository.Finder
}

func NewGeolocationService(l log.Logger, f repository.Finder) GeoLocationService {
	return &geolocation{
		l, f,
	}
}

func (g *geolocation) GetIPData(request *GetRequest) (GeoLocationResponse, *router.HttpError) {
	var resp GeoLocationResponse

	if len(request.IP) == 0 {
		return resp, router.NewHttpError("invalid ip", 400)
	}

	data, found, err := g.finder.FindManaged(&model.Geolocation{IP: request.IP})
	if err != nil {
		return resp, router.NewHttpError(err.Error(), 500)
	}

	if !found {
		return resp, router.NewHttpError("not found", 404)
	}

	v, ok := data.(*model.Geolocation)
	if !ok {
		return resp, router.NewHttpError(fmt.Sprintf("can not able to type assert the response: %v", data), 500)
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
