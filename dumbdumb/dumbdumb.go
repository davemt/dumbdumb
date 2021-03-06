package main

import (
	"github.com/davemt/dumbdumb"
	"github.com/davemt/dumbdumb/handler"
	"github.com/davemt/dumbdumb/listener"
	"log"
	"os"
	"strings"
)

func InitializeServer() *dumbdumb.Server {
	server := dumbdumb.NewServer()

	// Whitespace-separated list of email address domains
	domainWhitelist := os.Getenv("DUMBDUMB_SMTP_SENDER_DOMAIN_WHITELIST")
	server.AddListener(listener.SMTPListener{
		DomainWhitelist: strings.Fields(domainWhitelist),
	})

	server.AddHandler("[Ww]eather.*", handler.WeatherHandler{
		WuApiKey: os.Getenv("DUMBDUMB_WEATHERUNDERGROUND_API_KEY"),
	})

	server.AddHandler("[Ff]ind.*", handler.PlaceDirectoryHandler{
		GoogleAPIKey: os.Getenv("DUMBDUMB_GOOGLE_API_KEY"),
	})

	server.AddHandler("[Tt]ranslate.*", handler.TranslateHandler{
		GoogleAPIKey: os.Getenv("DUMBDUMB_GOOGLE_API_KEY"),
	})

	return server
}

func main() {
	server := InitializeServer()

	log.Printf("Starting dumbdumb server...")
	server.ListenAndServe()
}
