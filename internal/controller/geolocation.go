package controller

import (
	"fmt"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/router"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/service"
	"net/http"
)

const (
	ParamIP = "ip"
)

func (cntrl *clientController) GetGeolocationData(w http.ResponseWriter, r *http.Request) {
	l := router.GetLogger(r)
	req := new(service.GetRequest)
	req.IP = r.URL.Query().Get(ParamIP)
	if len(req.IP) == 0 {
		router.RenderJSON(router.Response{
			Writer: w,
			Data:   fmt.Errorf("path param could not be found"),
			Logger: l,
			Status: 400,
		})
		return
	}

	response, err := cntrl.geolocationSrv.GetIPData(req)
	if err != nil {
		router.RenderJSON(router.Response{
			Writer: w,
			Data:   err,
			Logger: l,
			Status: 500,
		})

		return
	}

	router.RenderJSON(router.Response{
		Writer: w,
		Data:   response,
		Logger: l,
		Status: 200,
	})

}
