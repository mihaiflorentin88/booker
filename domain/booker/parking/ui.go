package parking

import (
	"booking/port/config"
	"booking/port/contract"
	"fyne.io/fyne/v2"
)

type Renderer struct {
	config        *config.Config
	UI            *UIEntries
	parkingClient contract.ParkingClientInterface
}

func NewParkingRenderer(config *config.Config, app fyne.App, parkingClient contract.ParkingClientInterface, window fyne.Window) *Renderer {
	return &Renderer{
		config:        config,
		UI:            NewUIEntries(config, parkingClient, app, window),
		parkingClient: parkingClient,
	}
}

func (r *Renderer) UpdateConfig() {
	r.UI.UpdateConfig()
}

func (r *Renderer) UpdateUI() {
	r.UI.UpdateUI()
}
