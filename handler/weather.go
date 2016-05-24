package handler

import (
	"fmt"
	"github.com/davemt/dumbdumb"
	"github.com/jmcvetta/napping"
	"github.com/jmoiron/jsonq"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type WeatherHandler struct {
	// Weather Underground API key
	WuApiKey string
}

func (h WeatherHandler) HandleRequest(request dumbdumb.Request) error {
	log.Printf("Weather handler got request: %v", request.GetPayload())
	// TODO "weather" without query causes panic

	// ["weather", "<location query>"]
	parts := strings.SplitAfterN(request.GetPayload(), " ", 2)
	_, locQuery := parts[0], parts[1]

	params := url.Values{}
	params.Add("query", locQuery)
	params.Add("c", "US")

	var data map[string]interface{}

	// make request to get location id based on location query

	// http://www.wunderground.com/weather/api/d/docs?d=autocomplete-api
	resp, err := napping.Get("http://autocomplete.wunderground.com/aq", &params, &data, nil)
	if resp.Status() != 200 {
		return err
	}

	jq := jsonq.NewQuery(data)

	// assume the first location match is the best
	locId, err := jq.String("RESULTS", "0", "l")

	// use location id to get daily forecast data

	// http://www.wunderground.com/weather/api/d/docs?d=data/forecast
	resp, err = napping.Get(
		fmt.Sprintf("http://api.wunderground.com/api/%v/forecast/%v.json", h.WuApiKey, locId),
		nil,
		&data,
		nil,
	)
	if resp.Status() != 200 {
		return err
	}

	jq = jsonq.NewQuery(data)

	// gather forecasts for today, tonight, tomorrow; output contains all three
	for i := 0; i <= 2; i++ {
		forecast, err := jq.Object("forecast", "txt_forecast", "forecastday", strconv.Itoa(i))
		if err != nil {
			return err
		}
		dayName := forecast["title"].(string)
		forecastStr := forecast["fcttext"].(string)
		err = request.SendOutput(fmt.Sprintf("Weather %v: %v", dayName, forecastStr))
		if err != nil {
			return err
		}
	}
	return err
}
