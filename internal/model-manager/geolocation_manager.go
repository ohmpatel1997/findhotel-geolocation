package model_manager

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/model"
	"strings"
)

type GeoLocationManager interface {
	FindDataByIP(ip string) (model.Geolocation, bool, error)
	UpsertGeolocation(geolocation *model.Geolocation) (model.Geolocation, error)
	BulkInsert(geolocation []model.Geolocation) error
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
	if !ok {
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

	geolocation.ID, err = uuid.NewUUID()
	if err != nil {
		return resp, err
	}

	err = m.cuder.Insert(geolocation)
	if err != nil {
		return resp, nil
	}

	return *geolocation, nil
}

func (m *manager) BulkInsert(geolocation []model.Geolocation) error {
	return m.cuder.Transact(func(db *gorm.DB) error {
		valueStrings := []string{}
		valueArgs := []interface{}{}

		for _, geo := range geolocation {

			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")

			id, _ := uuid.NewUUID()

			valueArgs = append(valueArgs, id)
			valueArgs = append(valueArgs, geo.IP)
			valueArgs = append(valueArgs, geo.Country)
			valueArgs = append(valueArgs, geo.Longitude)
			valueArgs = append(valueArgs, geo.Latitude)
			valueArgs = append(valueArgs, geo.MysteryValue)
			valueArgs = append(valueArgs, geo.City)
			valueArgs = append(valueArgs, geo.CountryCode)

		}

		stmt := fmt.Sprintf("INSERT INTO geolocation (id, ip, country, longitude, latitude, mystery_value, city, country_code) VALUES %s", strings.Join(valueStrings, ","))
		return db.Exec(stmt, valueArgs...).Error
	})
}
