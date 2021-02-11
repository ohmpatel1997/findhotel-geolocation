package model

import (
	"github.com/google/uuid"
	"time"
)

type Geolocation struct {
	ID           uuid.UUID
	IP           string
	Country      string
	CountryCode  string
	City         string
	Latitude     string
	Longitude    string
	MysteryValue string
	CreatedAt    time.Time `sql:"DEFAULT:current_timestamp"`
	ModifiedAt   time.Time `sql:"DEFAULT:current_timestamp"`
}

func (p *Geolocation) TableName() string {
	return "geolocation"
}
