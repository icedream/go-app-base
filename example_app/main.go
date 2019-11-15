package main

import (
	"log"
	"net"
	"net/url"

	appbase "github.com/icedream/go-app-base"
)

var app = appbase.New(appbase.ApplicationDescription{
	ID:                                   "example",
	VendorID:                             "icedream",
	Name:                                 "Example application",
	AllowConfigurationInCurrentDirectory: true,
})

var (
	setting1 = app.Settings().Create("message", "The message to echo.").String("Hello world")
	setting2 = app.Settings().Create("ip", "An IP").IP(net.IP{127, 0, 0, 1})
	setting3 = app.Settings().Create("url", "A URL").URL(&url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/",
	}, true)
)

func main() {
	app.Settings().ReadIn()

	app.Configuration().Debug()

	log.Printf("%+v", *setting1)
	log.Printf("%+v", *setting2)
	log.Printf("%+v", *setting3)
}
