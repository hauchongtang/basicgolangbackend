package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var busStopIds = []string{
	"378204", "383050", "378202", "383049", "382998", "378237", "378233", "378230",
	"378229", "378228", "378227", "382995", "378224", "378226", "383010", "383009",
	"383006", "383004", "378234", "383003", "378222", "383048", "378203", "382999",
	"378225", "383014", "383013", "383011", "377906", "383018", "383015", "378207",
}

var busLineIds = []string{"44478", "44479", "44480", "44481"}

type BusStopGet struct { // map the entire get response from the API
	ExternalID string `json:"external_id"`
	Forecast   []struct {
		ForecastSeconds float64 `json:"forecast_seconds"`
		Route           struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			ShortName string `json:"short_name"`
		} `json:"route"`
		RvID      int     `json:"rv_id"`
		TotalPass float64 `json:"total_pass"`
		Vehicle   string  `json:"vehicle"`
		VehicleID int     `json:"vehicle_id"`
	} `json:"forecast"`
	Geometry []struct {
		ExternalID interface{} `json:"external_id"`
		Lat        string      `json:"lat"`
		Lon        string      `json:"lon"`
		Seq        int         `json:"seq"`
	} `json:"geometry"`
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	NameEn      interface{} `json:"name_en"`
	NameRu      interface{} `json:"name_ru"`
	Nameslug    string      `json:"nameslug"`
	ResourceURI string      `json:"resource_uri"`
}

type BusLineGet struct { // map the entire GET response from the external API
	ExternalID  interface{} `json:"external_id"`
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	NameEn      interface{} `json:"name_en"`
	NameRu      interface{} `json:"name_ru"`
	Nameslug    interface{} `json:"nameslug"`
	ResourceURI string      `json:"resource_uri"`
	Routename   string      `json:"routename"`
	Vehicles    []struct {
		Bearing    int    `json:"bearing"`
		DeviceTs   string `json:"device_ts"`
		Enterprise struct {
			EnterpriseID   int    `json:"enterprise_id"`
			EnterpriseName string `json:"enterprise_name"`
		} `json:"enterprise"`
		Lat  string `json:"lat"`
		Lon  string `json:"lon"`
		Park struct {
			ParkID   int    `json:"park_id"`
			ParkName string `json:"park_name"`
		} `json:"park"`
		Position struct {
			Bearing  int    `json:"bearing"`
			DeviceTs int    `json:"device_ts"`
			Lat      string `json:"lat"`
			Lon      string `json:"lon"`
			Speed    int    `json:"speed"`
			Ts       int    `json:"ts"`
		} `json:"position"`
		Projection struct {
			EdgeDistance    string `json:"edge_distance"`
			EdgeID          int    `json:"edge_id"`
			EdgeProjection  string `json:"edge_projection"`
			EdgeStartNodeID int    `json:"edge_start_node_id"`
			EdgeStopNodeID  int    `json:"edge_stop_node_id"`
			Lat             string `json:"lat"`
			Lon             string `json:"lon"`
			OrigLat         string `json:"orig_lat"`
			OrigLon         string `json:"orig_lon"`
			RoutevariantID  int    `json:"routevariant_id"`
			Ts              int    `json:"ts"`
		} `json:"projection"`
		RegistrationCode string `json:"registration_code"`
		RoutevariantID   int    `json:"routevariant_id"`
		Speed            string `json:"speed"`
		Stats            struct {
			AvgSpeed    string `json:"avg_speed"`
			Bearing     int    `json:"bearing"`
			CummSpeed10 string `json:"cumm_speed_10"`
			CummSpeed2  string `json:"cumm_speed_2"`
			DeviceTs    int    `json:"device_ts"`
			Lat         string `json:"lat"`
			Lon         string `json:"lon"`
			Speed       int    `json:"speed"`
			Ts          int    `json:"ts"`
		} `json:"stats"`
		Ts        string `json:"ts"`
		VehicleID int    `json:"vehicle_id"`
	} `json:"vehicles"`
	Via interface{} `json:"via"`
}

type BusStop struct {
	STOPNAME    string        `json:"stopname"`
	DATA        []BusStopData `json:"data"`
	id          int           `json:"id"`
	COORDINATES [2]string     `json:"coordinates"` // coordinate of the busstop.
}

type BusType struct {
	TYPE string `json:"TYPE"`
	ID   int    `json:"id"` // vehicle ID
}

type BusStopData struct { // coord is coord of busstop not bus live loc
	BUS         BusType   `json:"bus"`
	COORDINATES [2]string `json:"coordinates"`
	ARRIVE_IN   float64   `json:"arrive_in"`
}

type BusLineData struct { // coord here is bus live loc
	BUS         BusType   `json:"bus"`
	COORDINATES [2]string `json:"live_coordinates"`
}

type Bus struct {
	DATA []BusLineData `json:"data"`
}

func readBusLineAPI(id string) []BusLineData {
	response, err := http.Get(fmt.Sprintf("https://baseride.com/routes/apigeo/routevariantvehicle/%s/?format=json", id))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var busLineObject BusLineGet
	json.Unmarshal(responseData, &busLineObject)

	vehiclesArr := busLineObject.Vehicles
	var busesOnline []BusLineData
	line_id, err := strconv.Atoi(id)
	for _, vehicle := range vehiclesArr {
		toAdd := BusLineData{
			BUS:         createBusType(line_id, vehicle.VehicleID),
			COORDINATES: [2]string{vehicle.Position.Lat, vehicle.Position.Lon},
		}

		busesOnline = append(busesOnline, toAdd)
	}
	return busesOnline
}

func readBusStopAPI(id string, busLineArr []Bus) ([]BusStopData, [2]string, string) {
	response, err := http.Get(fmt.Sprintf("https://baseride.com/routes/api/platformbusarrival/%s/?format=json", id))
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var busStopObject BusStopGet
	json.Unmarshal(responseData, &busStopObject)

	busesArr := busStopObject.Forecast
	location := [2]string{busStopObject.Geometry[0].Lat, busStopObject.Geometry[0].Lon}
	stopName := busStopObject.Name
	var data []BusStopData

	for _, bussvc := range busesArr {
		toAdd := BusStopData{
			BUS: createBusType(bussvc.RvID, bussvc.VehicleID),
			// change to bus live location
			COORDINATES: getLiveLocation(busLineArr, bussvc.VehicleID, bussvc.RvID),
			ARRIVE_IN:   bussvc.ForecastSeconds,
		}
		data = append(data, toAdd)
	}
	return data, location, stopName
}

func getLiveLocation(busLineArr []Bus, veh_id int, route int) [2]string {
	bus := busLineArr[0]
	switch route {
	case 44478:
		bus = busLineArr[0]
	case 44479:
		bus = busLineArr[1]
	case 44480:
		bus = busLineArr[2]
	case 44481:
		bus = busLineArr[3]
	}

	for _, b := range bus.DATA {
		if b.BUS.ID == veh_id {
			return b.COORDINATES
		}
	}
	return [2]string{"", ""}
}

func createBusType(s int, id int) BusType {
	bustype := BusType{TYPE: "NONE", ID: id}
	switch s {
	case 44478:
		bustype.TYPE = "RED"
	case 44479:
		bustype.TYPE = "BLUE"
	case 44480:
		bustype.TYPE = "GREEN"
	case 44481:
		bustype.TYPE = "BROWN"
	}
	return bustype
}

func constructAPI() []BusStop {
	var busStops []BusStop
	var busLines []Bus
	for _, line_id := range busLineIds {
		event := Bus{
			DATA: readBusLineAPI(line_id),
		}
		busLines = append(busLines, event)
	}

	for _, stop_id := range busStopIds {
		busStopData, coordinate, stopName := readBusStopAPI(stop_id, busLines)
		id, err := strconv.Atoi(stop_id)

		if err != nil {
			log.Fatal(err)
		}
		event := BusStop{
			STOPNAME:    stopName,
			DATA:        busStopData,
			id:          id,
			COORDINATES: coordinate,
		}

		busStops = append(busStops, event)
	}
	return busStops
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("End point hit: Homepage")
}

func returnBusEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("End point hit: GET all")
	json.NewEncoder(w).Encode(constructAPI())
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/events", returnBusEvents)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	//fmt.Println(constructAPI())
	//fmt.Println(readBusLineAPI("44478"))
	handleRequests()
}
