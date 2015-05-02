package main

import (
	"dumbdumb"
	"dumbdumb/handler"
	"dumbdumb/listener"
	"os"
)

func main() {
	server := &dumbdumb.Server{}

	server.AddListener(listener.SMTPListener{})

	server.AddHandler("weather.*", handler.WeatherHandler{
		WuApiKey: os.Getenv("DUMBDUMB_WEATHERUNDERGROUND_API_KEY"),
	})

	server.AddHandler("transit mbta.*", handler.TransitMBTAHandler{
		ApiKey: os.Getenv("DUMBDUMB_TRANSIT_MBTA_API_KEY"),
	})

	server.ListenAndServe()
}
