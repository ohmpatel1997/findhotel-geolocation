package controller

import (
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/service"
	"net/http"
)

const (
	clientApiVersion = "v1"
)

type ClientController interface {

	//metadata
	GetAPIVersion() string
	GetAPIVersionPath(string) string

	GetGeolocationData(http.ResponseWriter, *http.Request)
}

type clientController struct {
	l              log.Logger
	geolocationSrv service.GeoLocationService
}

func NewController(l log.Logger, geolocation service.GeoLocationService) ClientController {
	return &clientController{
		l:              l,
		geolocationSrv: geolocation,
	}
}

func (cntrl *clientController) GetAPIVersion() string {
	return clientApiVersion
}

func (cntrl *clientController) GetAPIVersionPath(p string) string {
	return "/" + clientApiVersion + p
}
