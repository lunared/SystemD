package bot

import (
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"

	"../helpers"
)

//WeatherData ...
//Holds some of the weather data we get from the server
type WeatherData struct {
	Temp      float32
	Humidity  string
	Rain      string
	Windchill string
	FullName  string
}

type jsonRoot struct {
	CurrentObservation jsonObservation `json:"current_observation"`
}

type jsonObservation struct {
	DisplayLocation   jsonDisplay `json:"display_location"`
	TempF             float32     `json:"temp_f"`
	RelativeHumidity  string      `json:"relative_humidity"`
	WindchillF        string      `json:"windchill_f"`
	PrecipTodayString string      `json:"precip_today_string"`
}

type jsonDisplay struct {
	Full string `json:"full"`
}

func printWeatherData(conn net.Conn, channel string, msg string, apiKey string) {
	args := strings.Fields(msg)
	state := "PA"
	location := args[1]
	if len(args) == 3 {
		state = args[1]
		location = args[2]
	}

	url := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/%s/%s.json", apiKey, state, location)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR getting data from : " + url)
		return
	}

	jsonText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR getting data from : " + url)
		return
	}

	var data WeatherData
	var decodedJSON jsonRoot
	err = json.Unmarshal(jsonText, &decodedJSON)
	if err != nil {
		fmt.Println("ERROR getting data from : " + url + "\n" + err.Error())
		return
	}

	data.FullName = decodedJSON.CurrentObservation.DisplayLocation.Full
	data.Rain = decodedJSON.CurrentObservation.PrecipTodayString
	data.Humidity = decodedJSON.CurrentObservation.RelativeHumidity
	data.Temp = decodedJSON.CurrentObservation.TempF
	data.Windchill = decodedJSON.CurrentObservation.WindchillF

	printData(conn, channel, data)
}

func printData(conn net.Conn, channel string, data WeatherData) {
	helpers.SendPrivMsg(conn, channel, fmt.Sprintf("Data for %s", data.FullName))
	helpers.SendPrivMsg(conn, channel,
		fmt.Sprintf("Temp(F): %5.2f\t\tWindchill: %s", data.Temp, data.Windchill))
	helpers.SendPrivMsg(conn, channel,
		fmt.Sprintf("Precipitation: %s\t\tHumidity: %s", data.Rain, data.Humidity))
}

func checkWeatherCommand(msg string) bool {
	reg := regexp.MustCompile(`!weather (([A-Za-z]+) ([A-Za-z]?)|([A-Za-z]+))`)
	return reg.MatchString(msg)
}
