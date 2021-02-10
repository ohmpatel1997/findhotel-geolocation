package main

import (
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/router"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/controller"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/service"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	l := log.NewLogger()

	l.Info("### Starting up client api ###")

	rport, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		l.PanicD("Unable to read CORE_DB_READ_PORT var", log.Fields{"err": err.Error()})
	}

	sslModeCoreDB := os.Getenv("DB_SSL_MODE")
	if sslModeCoreDB == "" {
		sslModeCoreDB = repository.SSLModeRequire
	}

	connMaxLife := 4 * time.Minute
	maxIdleConn := 2
	maxOpenConn := 3

	rpgc := repository.PGConfig{
		Host:               os.Getenv("DATABASE_URL"),
		Port:               rport,
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASSWD"),
		DBName:             os.Getenv("DB_NAME"),
		SSLMode:            sslModeCoreDB,
		SetConnMaxLifetime: &connMaxLife,
		SetMaxOpenConns:    &maxOpenConn,
		SetMaxIdleConns:    &maxIdleConn,
	}

	rdb, err := repository.NewPGConnection(rpgc)
	if err != nil {
		l.PanicD("Error getting read connection", log.Fields{"err": err.Error()})
	}

	c := repository.NewCuder(rdb)
	f := repository.NewFinder(rdb)

	srv := service.NewGeolocationService(l, c, f)
	cntrl := controller.NewController(l, srv)
	router := registerRoutes(cntrl, l)

	err = router.ListenAndServeTLS(os.Getenv("CONTAINER_LISTEN_PORT"), nil)
	if err != nil {
		l.Panic(err.Error())
	}
}

func registerRoutes(clientCntrl controller.ClientController, l log.Logger) router.Router {
	r := router.NewBasicRouter()

	r.Route(clientCntrl.GetAPIVersionPath("/ip"), func(r router.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			clientCntrl.GetGeolocationData(w, r)
		})
	})

	return r
}
