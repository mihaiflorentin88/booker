package renderer

import (
	"booking/domain/booker/parking"
	"booking/port/config"
	"booking/port/contract"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"os"
)

type Renderer struct {
	Config          *config.Config
	App             fyne.App
	Window          fyne.Window
	ParkingBooker   *parking.Renderer
	Tabs            *container.AppTabs
	ButtonContainer *fyne.Container
}

func NewRenderer(config *config.Config, parkingClient contract.ParkingClientInterface) *Renderer {
	renderer := &Renderer{
		Config: config,
		App:    app.NewWithID("com.booking.app"),
	}
	renderer.Window = renderer.App.NewWindow("Booking")
	parkingBooker := parking.NewParkingRenderer(config, renderer.App, parkingClient, renderer.Window)
	renderer.ParkingBooker = parkingBooker
	return renderer
}

func (r *Renderer) Render() {
	r.Tabs = container.NewAppTabs(
		r.ParkingBooker.UI.Tab,
	)
	layout := container.NewBorder(container.NewHBox(r.fileMenu()), nil, nil, nil, r.Tabs)
	r.Window.SetContent(layout)
	r.Window.Resize(fyne.NewSize(600, 700))
	r.Window.ShowAndRun()
}

func (r *Renderer) exportSettings(filename string) {
	data, err := json.MarshalIndent(r.Config, "", "  ")
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to marshal config: %w", err), r.Window)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		dialog.ShowError(fmt.Errorf("failed to write config file: %w", err), r.Window)
		return
	}
}

func (r *Renderer) importConfig(filename string) {

	data, err := os.ReadFile(filename)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read config file: %w", err), r.Window)
	}

	if err := json.Unmarshal(data, r.Config); err != nil {
		dialog.ShowError(fmt.Errorf("failed to unmarshal config: %w", err), r.Window)
	}
	r.ParkingBooker.UpdateUI()
}
