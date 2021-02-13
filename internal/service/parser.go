package service

import (
	"bufio"
	"fmt"
	model_manager "github.com/ohmpatel1997/findhotel-geolocation/internal/model-manager"

	"github.com/ohmpatel1997/findhotel-geolocation/integration/log"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/common"
	"github.com/ohmpatel1997/findhotel-geolocation/internal/model"
	"io"
	"math"
	"os"

	"strconv"
	"strings"
	"sync"
	"time"
)

type ParserService interface {
	ParseAndStore() (float64, int64, int64, error)
}

type parser struct {
	l       log.Logger
	f       *os.File
	manager model_manager.GeoLocationManager
}

func NewParser(l log.Logger, f *os.File, mn model_manager.GeoLocationManager) ParserService {
	return &parser{
		l:       l,
		f:       f,
		manager: mn,
	}
}

func (p *parser) ParseAndStore() (float64, int64, int64, error) {

	timeThen := time.Now()

	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(p.f)

	firstLine, _, err := r.ReadLine()
	if err != nil {
		return 0, 0, 0, err
	}
	firstLineSlice := strings.Split(string(firstLine), ",")
	positions := make(map[int]string)

	//map the positions
	for i, header := range firstLineSlice {
		positions[i] = header
	}

	var invalidDataCountFromFirstPass int64 = 0
	var invalidDataCountFromSecondPass int64 = 0
	var validDataCount int64 = 0

	outPutChan := make(chan model.Geolocation, 10000)
	savToDbChan := make(chan model.Geolocation, 10000)
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	go p.ExtractAndLoad(outPutChan, &invalidDataCountFromSecondPass, &validDataCount, &wg, savToDbChan, &wg2)
	go p.SaveToDB(savToDbChan, &wg2)

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			break
		}

		nextUntillNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		invalidDataCountFromFirstPass += ProcessChunk(buf, &linesPool, &stringPool, positions, outPutChan, &wg)
	}

	close(outPutChan)

	wg.Wait()  //wait ExtractAndLoad
	wg2.Wait() //wait until saving data to db

	return time.Since(timeThen).Seconds(), invalidDataCountFromFirstPass + invalidDataCountFromSecondPass, validDataCount, nil
}

func ProcessChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool, positions map[int]string, outPutChan chan<- model.Geolocation, wg *sync.WaitGroup) int64 {

	var wg2 sync.WaitGroup
	var invalid int64 = 0
	logs := stringPool.Get().(string)
	logs = string(chunk)

	//put back the old chunk
	linesPool.Put(chunk)

	logsSlice := strings.Split(logs, "\n")

	//put back the slice
	stringPool.Put(logs)

	chunkSize := 500
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < (noOfThread); i++ {

		wg2.Add(1) //span out locally

		go func(textSlice []string) {
			for _, text := range textSlice { //first stage of cleaning

				if len(text) == 0 { //in case there is line gap
					continue
				}
				logSlice := strings.Split(text, ",")

				if len(logSlice) != 7 { //if not valid number of fields
					invalid++
					continue
				}

				geoloc := model.Geolocation{}
				invalidData := false
				for i, value := range logSlice {

					if len(value) == 0 { //if empty value
						invalid++
						invalidData = true
						break
					}
					col := positions[i]
					switch col {
					case common.IP:
						geoloc.IP = value
					case common.CountryCode:
						geoloc.CountryCode = value
					case common.Country:
						geoloc.Country = value
					case common.Longitude:
						geoloc.Longitude = value
					case common.Latitude:
						geoloc.Latitude = value
					case common.MysteryValue:
						geoloc.MysteryValue = value
					case common.City:
						geoloc.City = value
					default: //if some other columns come in
						invalidData = true
						invalid++
						break
					}
				}

				if !invalidData {
					wg.Add(1)            //increment counter for data processing
					outPutChan <- geoloc //send to output chan
				}

			}
			wg2.Done() //done processing a chunk
		}(logsSlice[i*chunkSize : int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice))))]) //prevent overflow
	}

	wg2.Wait()
	logsSlice = nil //free up the log slice

	return invalid //return the invalid data count
}

//will extract the data, checks the validity and load it into database
func (p *parser) ExtractAndLoad(outPutChan <-chan model.Geolocation, invalidCount *int64, validCount *int64, wg *sync.WaitGroup, saveToDbChan chan<- model.Geolocation, wg2 *sync.WaitGroup) {

	visitedIP := make(map[string]bool)          // will keep track of already visited ip address
	visitedCoordinates := make(map[string]bool) // will keep track of already visited coordinates

	for data := range outPutChan { // second stage of cleaning
		local_data := data
		IPValid := common.IsIpv4Regex(data.IP)
		if !IPValid {
			*invalidCount++
			wg.Done()
			continue
		}

		latitude, err := strconv.ParseFloat(data.Latitude, 64)
		if err != nil {
			*invalidCount++
			wg.Done()
			continue
		}

		longitude, err := strconv.ParseFloat(data.Longitude, 64)
		if err != nil {
			*invalidCount++
			wg.Done()
			continue
		}

		if latitude > 90 || latitude < -90 { //invalid latitude coordinates
			*invalidCount++
			wg.Done()
			continue
		}

		if longitude > 180 || longitude < -180 { //invalid longitude coordinates
			*invalidCount++
			wg.Done()
			continue
		}

		if ok := visitedIP[data.IP]; ok {
			*invalidCount++
			wg.Done()
			continue
		}

		coordinates := fmt.Sprintf("%s+%s", data.Latitude, data.Longitude)
		if ok := visitedCoordinates[coordinates]; ok {
			*invalidCount++
			wg.Done()
			continue
		}

		visitedIP[data.IP] = true
		visitedCoordinates[coordinates] = true
		*validCount++

		wg2.Add(1)                 //add count for saving data to db
		saveToDbChan <- local_data //push to save data to db

		wg.Done() //decrement for process done
	}

	close(saveToDbChan) //once all data have been processed, close save to db chan
}

func (p *parser) SaveToDB(savChan <-chan model.Geolocation, wg2 *sync.WaitGroup) {
	for data := range savChan {

		local_data := data

		go func() {
			defer wg2.Done() //decrement after data is saved

			_, err := p.manager.UpsertGeolocation(&local_data)
			if err != nil {
				p.l.ErrorD("failed to check data already exists", log.Fields{"data": local_data, "Error": err.Error()})
				return
			}
		}()
	}
}
