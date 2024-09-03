package parking

import (
	"booking/port/config"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (ui *UIEntries) getParkingSpotBtn() *widget.Button {
	return widget.NewButton("Add Parking Spot", func() {
		newSpot := config.ParkingSpot{}
		ui.Config.Parking.ParkingSpot = append(ui.Config.Parking.ParkingSpot, newSpot)
		ui.ParkingSpotContainer.Add(ui.createParkingSpotEntry(&newSpot))
		ui.ParkingSpotContainer.Add(widget.NewSeparator())
		ui.ParkingSpotContainer.Refresh()
	})
}

func (ui *UIEntries) getExecuteBtn() *widget.Button {
	return widget.NewButton("Execute", func() {
		ui.UpdateConfig()
		isValid, missingFields := ui.Config.IsValid()
		if !isValid {
			dialog.ShowError(fmt.Errorf("Cannot execute. Reason: missing the following fields  %s", missingFields), ui.Window)
			return
		}
		ui.RunnerOutputWindow = ui.App.NewWindow("Output Window")
		ctx, stopButton := ui.getRunnerCancelBtn()
		ui.RunnerOutputWindow.SetContent(container.NewVBox(
			stopButton,
			ui.RunnerOutputScrollContainer,
		))
		ui.RunnerOutputWindow.Resize(fyne.NewSize(600, 400))
		ui.RunnerOutputWindow.Show()
		parkingBooker := NewParkingBooker(ui.Config, ui.parkingClient, ctx, ui.RunnerOutputLabel, updateRunnerOutput)
		go func() { parkingBooker.Start() }()
	})
}

func (ui *UIEntries) getRunnerCancelBtn() (context.Context, *widget.Button) {
	var ctx context.Context
	ctx, cancelFunc := context.WithCancel(context.Background())
	stopButton := widget.NewButton("Stop (F6)", func() {
		if cancelFunc != nil {
			cancelFunc()
		}
		ui.RunnerOutputWindow.Close()
	})
	ui.Window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyF6 {
			if cancelFunc != nil {
				cancelFunc()
			}
			ui.RunnerOutputWindow.Close()
		}
	})
	ui.RunnerOutputWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyF6 {
			if cancelFunc != nil {
				cancelFunc()
			}
			ui.RunnerOutputWindow.Close()
		}
	})
	return ctx, stopButton
}

func updateRunnerOutput(runnerOutputLabel *widget.Label, newLog string) {
	runnerOutputLabel.SetText(fmt.Sprintf("%s\n%s", runnerOutputLabel.Text, newLog))
}
