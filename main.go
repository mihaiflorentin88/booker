package main

import (
	"booking/domain/renderer"
	"booking/infrastructure/client/parking"
	"booking/port/config"
)

func main() {
	//cfg := config.NewEmptyConfig()
	cfg := config.Config{
		Parking: config.Parking{},
	}
	parkingClient := parking.NewParkingClient(&cfg.Parking.Credentials)
	r := renderer.NewRenderer(&cfg, parkingClient)
	r.Render()
}
