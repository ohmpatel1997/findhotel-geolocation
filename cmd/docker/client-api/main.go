package main

import (
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/router"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/controller"
	model_manager "github.com/ohmpatel1997/findhotel-geolocation/internal/model-manager"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/service"
	"net/http"
	"os"
)

func main() {
	l := log.NewLogger()

	l.Info("### Starting up client api ###")

	sslModeCoreDB := os.Getenv("DB_SSL_MODE")
	if sslModeCoreDB == "" {
		sslModeCoreDB = repository.SSLModeRequire
	}

	connStr := "dbname=ohmpatel user=ohmpatel password=ohmpatel host=localhost sslmode=disable\n\n"
	//os.Getenv("DATABASE_URL")

	l.Info("connection string-->", connStr, "\n\n")
	if len(connStr) == 0 {
		l.Panic("no conn string found")
	}

	rdb, err := repository.NewPGConnection(nil, &connStr)
	if err != nil {
		l.PanicD("Error getting read connection", log.Fields{"err": err.Error()})
	}

	f := repository.NewFinder(rdb)
	c := repository.NewCuder(rdb)
	manager := model_manager.NewGeoLocationManager(l, c, f)

	srv := service.NewGeolocationService(l, manager)
	cntrl := controller.NewController(l, srv)
	router := registerRoutes(cntrl, l)

	err = router.ListenAndServeTLS(os.Getenv("PORT"), nil)
	if err != nil {
		l.Panic(err.Error())
	}
}

func registerRoutes(clientCntrl controller.ClientController, l log.Logger) router.Router {
	r := router.NewBasicRouter()

	r.Route(clientCntrl.GetAPIVersionPath("/ip-info"), func(r router.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			clientCntrl.GetGeolocationData(w, r)
		})
	})

	return r
}
