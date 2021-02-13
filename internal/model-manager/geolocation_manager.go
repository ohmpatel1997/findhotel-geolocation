package model_manager

import (
	"fmt"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/model"
)

type GeoLocationManager interface {
	FindDataByIP(ip string) (model.Geolocation, bool, error)
	UpsertGeolocation(geolocation *model.Geolocation) (model.Geolocation, error)
}

type manager struct {
	l      log.Logger
	cuder  repository.Cuder
	finder repository.Finder
}

func NewGeoLocationManager(l log.Logger, c repository.Cuder, f repository.Finder) GeoLocationManager {
	return &manager{
		l:      l,
		cuder:  c,
		finder: f,
	}
}

func (m *manager) FindDataByIP(ip string) (model.Geolocation, bool, error) {
	var resp model.Geolocation

	if len(ip) == 0 {
		return resp, false, fmt.Errorf("ip can not be empty")
	}

	response, found, err := m.finder.FindManaged(&model.Geolocation{IP: ip})

	if err != nil {
		return resp, false, err
	}

	if !found {
		return resp, false, fmt.Errorf("data not found with given ip")
	}

	v, ok := response.(*model.Geolocation)
	if ok {
		return resp, false, fmt.Errorf("can not able to type assert")
	}

	return *v, true, nil
}

func (m *manager) UpsertGeolocation(geolocation *model.Geolocation) (model.Geolocation, error) {
	var resp model.Geolocation

	response, found, err := m.finder.FindManaged(geolocation)
	if err != nil {
		return resp, err
	}

	if found {
		v, ok := response.(*model.Geolocation)
		if !ok {
			return resp, fmt.Errorf("can not able to type assert")
		}

		return *v, nil
	}

	err = m.cuder.Insert(geolocation)
	if err != nil {
		return resp, nil
	}

	return *geolocation, nil
}
