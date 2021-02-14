package main

import (
	"fmt"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	model_manager "github.com/ohmpatel1997/findhotel-geolocation/internal/model-manager"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/service"
	"os"
	"strconv"
)

func main() {
	args := os.Args

	l := log.NewLogger()

	if len(args) != 2 {
		l.Panic("please only specify the data file location")
	}

	file, err := os.Open(args[1])
	if err != nil {
		l.Panic(err.Error())
	}

	rport, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		l.PanicD("Unable to read PORT var", log.Fields{"err": err.Error()})
	}

	sslModeCoreDB := os.Getenv("DB_SSL_MODE")
	if sslModeCoreDB == "" {
		sslModeCoreDB = repository.SSLModeRequire
	}

	rpgc := repository.PGConfig{
		Host:     os.Getenv("HOST"),
		Port:     rport,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  sslModeCoreDB,
	}

	rdb, err := repository.NewPGConnection(&rpgc, nil)
	if err != nil {
		l.PanicD("Error getting read connection", log.Fields{"err": err.Error()})
	}

	fmt.Println("successfully connected===>", rdb)
	c := repository.NewCuder(rdb)
	f := repository.NewFinder(rdb)

	manager := model_manager.NewGeoLocationManager(l, c, f)
	parserService := service.NewParser(l, file, manager)

	timeTaken, invalid, validData, err := parserService.ParseAndStore()
	if err != nil {
		l.Error(err.Error())
		return
	}

	l.Info("Successfully Parsed And Stored")

	l.InfoD("Metrics: ", log.Fields{"Time Taken": timeTaken, "Valid Data": validData, "Invalid Data": invalid})
}
