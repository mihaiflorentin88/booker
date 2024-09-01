package main

import (
	"booking/domain/booker/parking"
	parking2 "booking/infrastructure/client/parking"
	"booking/port/config"
	"context"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"
	"strconv"
	"time"
)

func main() {
	a := app.NewWithID("com.booking.app")
	w := a.NewWindow("Booking")

	cfg := config.Config{
		Parking: &config.Parking{},
	}

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Enter Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter Password")

	googleAPIKeyEntry := widget.NewEntry()
	googleAPIKeyEntry.SetPlaceHolder("Enter Google API key")

	canStandbyEntry := widget.NewCheck("Enable Standby", nil)

	startTimePicker := createCustomDatePicker()
	endTimePicker := createCustomDatePicker()
	datePicker := createCustomDatePicker()

	offsetEntry := widget.NewEntry()
	offsetEntry.SetPlaceHolder("Enter Offset")

	parkingSpotContainer := container.NewVBox()
	for _, spot := range cfg.Parking.ParkingSpot {
		parkingSpotContainer.Add(createParkingSpotEntry(&spot))
		parkingSpotContainer.Add(widget.NewSeparator())
	}

	scrollableParkingSpots := container.NewVScroll(parkingSpotContainer)
	scrollableParkingSpots.SetMinSize(fyne.NewSize(600, 200))

	addSpotButton := widget.NewButton("Add Parking Spot", func() {
		newSpot := config.ParkingSpot{}
		cfg.Parking.ParkingSpot = append(cfg.Parking.ParkingSpot, newSpot)
		parkingSpotContainer.Add(createParkingSpotEntry(&newSpot))
		parkingSpotContainer.Add(widget.NewSeparator())
		parkingSpotContainer.Refresh()
	})

	parkingTab := container.NewVBox(
		widget.NewLabel("Credentials"),
		widget.NewForm(
			widget.NewFormItem("Username", usernameEntry),
			widget.NewFormItem("Password", passwordEntry),
			widget.NewFormItem("Google API Key", googleAPIKeyEntry),
		),
		widget.NewLabel("Standby Prevention"),
		widget.NewForm(
			widget.NewFormItem("Can Standby", canStandbyEntry),
			widget.NewFormItem("Start Time", startTimePicker),
			widget.NewFormItem("End Time", endTimePicker),
		),
		widget.NewLabel("Schedule"),
		widget.NewForm(
			widget.NewFormItem("Date", datePicker),
			widget.NewFormItem("Offset", offsetEntry),
		),
		widget.NewLabel("Parking Spots"),
		scrollableParkingSpots,
		addSpotButton,
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Parking", parkingTab),
	)

	importButton := widget.NewButton("Import Config", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			newConfig, err := importConfig(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			if parkingTab != nil && parkingSpotContainer != nil {
				cfg = newConfig
				updateUI(usernameEntry, passwordEntry, canStandbyEntry, startTimePicker, endTimePicker, datePicker, offsetEntry, googleAPIKeyEntry, parkingSpotContainer, &cfg)
			} else {
				dialog.ShowError(fmt.Errorf("UI components not initialized"), w)
			}
		}, w).Show()
	})

	exportButton := widget.NewButton("Export Config", func() {
		cfg.Parking.Credentials.Username = usernameEntry.Text
		cfg.Parking.Credentials.Password = passwordEntry.Text
		cfg.Parking.Credentials.GoogleAPIKey = googleAPIKeyEntry.Text
		cfg.Parking.StandBy.CanStandby = canStandbyEntry.Checked

		cfg.Parking.StandBy.StartTime = assembleDateFromPicker(startTimePicker)
		cfg.Parking.StandBy.EndTime = assembleDateFromPicker(endTimePicker)
		cfg.Parking.Schedule.Date = assembleDateFromPicker(datePicker)

		if offset, err := strconv.Atoi(offsetEntry.Text); err == nil {
			cfg.Parking.Schedule.Offset = offset
		} else {
			dialog.ShowError(fmt.Errorf("invalid offset value"), w)
			return
		}

		updateParkingSpots(cfg.Parking.ParkingSpot, parkingSpotContainer)

		dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()

			if err := exportConfig(cfg, writer.URI().Path()); err != nil {
				dialog.ShowError(err, w)
			} else {
				fmt.Println("Configuration exported to", writer.URI().Path())
			}
		}, w).Show()
	})

	runButton := widget.NewButton("Run", func() {
		cfg.Parking.Credentials.Username = usernameEntry.Text
		cfg.Parking.Credentials.Password = passwordEntry.Text
		cfg.Parking.Credentials.GoogleAPIKey = googleAPIKeyEntry.Text
		cfg.Parking.StandBy.CanStandby = canStandbyEntry.Checked

		cfg.Parking.StandBy.StartTime = assembleDateFromPicker(startTimePicker)
		cfg.Parking.StandBy.EndTime = assembleDateFromPicker(endTimePicker)
		cfg.Parking.Schedule.Date = assembleDateFromPicker(datePicker)

		if offset, err := strconv.Atoi(offsetEntry.Text); err == nil {
			cfg.Parking.Schedule.Offset = offset
		} else {
			dialog.ShowError(fmt.Errorf("invalid offset value"), w)
			return
		}
		updateParkingSpots(cfg.Parking.ParkingSpot, parkingSpotContainer)

		outputWindow := a.NewWindow("Output Window")
		outputLabel := widget.NewLabel("Output will appear here.")
		outputLabel.Wrapping = fyne.TextWrapWord
		scrollContainer := container.NewVScroll(outputLabel)
		scrollContainer.SetMinSize(fyne.NewSize(700, 300))

		var ctx context.Context
		ctx, cancelFunc := context.WithCancel(context.Background())

		stopButton := widget.NewButton("Stop (F6)", func() {
			if cancelFunc != nil {
				cancelFunc()
			}
			outputWindow.Close()
		})
		w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
			if ev.Name == fyne.KeyF6 {
				if cancelFunc != nil {
					cancelFunc()
				}
				outputWindow.Close()
			}
		})
		outputWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
			if ev.Name == fyne.KeyF6 {
				if cancelFunc != nil {
					cancelFunc()
				}
				outputWindow.Close()
			}
		})

		outputWindow.SetContent(container.NewVBox(
			stopButton,
			scrollContainer,
		))
		outputWindow.Resize(fyne.NewSize(600, 400))
		outputWindow.Show()
		parkingClient := parking2.NewParkingClient(&cfg.Parking.Credentials)
		parkingBooker := parking.NewParkingBooker(&cfg, parkingClient, ctx, outputLabel, updateOutput)
		go func() {
			err := parkingBooker.Start()
			if err != nil {
				updateOutput(outputLabel, fmt.Sprintf("Error: %v", err))
			}
		}()
	})

	content := container.NewVBox(
		tabs,
	)

	scrollableContent := container.NewVScroll(content)
	scrollableContent.SetMinSize(fyne.NewSize(600, 700))

	w.SetContent(container.NewVBox(scrollableContent, container.NewHBox(importButton, exportButton, runButton)))

	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}

func updateOutput(outputLabel *widget.Label, newLog string) {
	outputLabel.SetText(fmt.Sprintf("%s\n%s", outputLabel.Text, newLog))
}

func createParkingSpotEntry(spot *config.ParkingSpot) *fyne.Container {
	parkingIDEntry := widget.NewEntry()
	spotIDEntry := widget.NewEntry()

	parkingIDEntry.SetPlaceHolder("Enter Parking ID")
	spotIDEntry.SetPlaceHolder("Enter Spot ID")

	parkingIDEntry.SetText(spot.ParkingID)
	spotIDEntry.SetText(spot.SpotID)

	parkingIDContainer := container.NewGridWithColumns(2,
		widget.NewLabel("Parking ID"),
		parkingIDEntry,
	)
	spotIDContainer := container.NewGridWithColumns(2,
		widget.NewLabel("Spot ID"),
		spotIDEntry,
	)

	return container.NewVBox(
		parkingIDContainer,
		spotIDContainer,
	)
}

func createCustomDatePicker() *fyne.Container {
	days := make([]string, 31)
	for i := 1; i <= 31; i++ {
		days[i-1] = fmt.Sprintf("%02d", i)
	}
	daySelect := widget.NewSelect(days, func(value string) {})

	months := []string{
		"01 - January", "02 - February", "03 - March",
		"04 - April", "05 - May", "06 - June",
		"07 - July", "08 - August", "09 - September",
		"10 - October", "11 - November", "12 - December",
	}
	monthSelect := widget.NewSelect(months, func(value string) {})

	years := make([]string, 51)
	currentYear := time.Now().Year()
	for i := 0; i < 51; i++ {
		years[i] = strconv.Itoa(currentYear - 25 + i)
	}
	yearSelect := widget.NewSelect(years, func(value string) {})

	hours := make([]string, 24)
	for i := 0; i < 24; i++ {
		hours[i] = fmt.Sprintf("%02d", i)
	}
	hourSelect := widget.NewSelect(hours, func(value string) {})

	minutes := make([]string, 60)
	seconds := make([]string, 60)
	for i := 0; i < 60; i++ {
		minutes[i] = fmt.Sprintf("%02d", i)
		seconds[i] = fmt.Sprintf("%02d", i)
	}
	minuteSelect := widget.NewSelect(minutes, func(value string) {})
	secondSelect := widget.NewSelect(seconds, func(value string) {})

	dateTimeContainer := container.NewVBox(
		widget.NewLabel("Day/Month/Year"),
		container.NewHBox(
			daySelect,
			monthSelect,
			yearSelect,
		),
		widget.NewLabel("Hour/Minute/Second"),
		container.NewHBox(
			hourSelect,
			minuteSelect,
			secondSelect,
		),
	)

	return dateTimeContainer
}

func assembleDateFromPicker(picker *fyne.Container) string {
	day := picker.Objects[1].(*fyne.Container).Objects[0].(*widget.Select).Selected
	month := picker.Objects[1].(*fyne.Container).Objects[1].(*widget.Select).Selected[:2]
	year := picker.Objects[1].(*fyne.Container).Objects[2].(*widget.Select).Selected

	hour := picker.Objects[3].(*fyne.Container).Objects[0].(*widget.Select).Selected
	minute := picker.Objects[3].(*fyne.Container).Objects[1].(*widget.Select).Selected
	second := picker.Objects[3].(*fyne.Container).Objects[2].(*widget.Select).Selected

	return fmt.Sprintf("%s-%s-%s %s:%s:%s", day, month, year, hour, minute, second)
}

func updateUI(usernameEntry, passwordEntry *widget.Entry, canStandbyEntry *widget.Check, startTimePicker, endTimePicker, datePicker *fyne.Container, offsetEntry *widget.Entry, googleAPIKeyEntry *widget.Entry, parkingSpotContainer *fyne.Container, config *config.Config) {
	if config == nil {
		fmt.Println("Error: Configuration data is missing.")
		return
	}

	usernameEntry.SetText(config.Parking.Credentials.Username)
	passwordEntry.SetText(config.Parking.Credentials.Password)
	googleAPIKeyEntry.SetText(config.Parking.Credentials.GoogleAPIKey)

	canStandbyEntry.SetChecked(config.Parking.StandBy.CanStandby)

	offsetEntry.SetText(strconv.Itoa(config.Parking.Schedule.Offset))

	updateDatePicker(startTimePicker, config.Parking.StandBy.StartTime)
	updateDatePicker(endTimePicker, config.Parking.StandBy.EndTime)
	updateDatePicker(datePicker, config.Parking.Schedule.Date)
	parkingSpotContainer.Objects = nil
	for i := range config.Parking.ParkingSpot {
		parkingSpotContainer.Add(createParkingSpotEntry(&config.Parking.ParkingSpot[i]))
		parkingSpotContainer.Add(widget.NewSeparator())
	}
	parkingSpotContainer.Refresh()
}

func updateDatePicker(picker *fyne.Container, dateStr string) {
	date, err := time.Parse("02-01-2006 15:04:05", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	day := fmt.Sprintf("%02d", date.Day())
	month := fmt.Sprintf("%02d - %s", int(date.Month()), date.Month().String())
	year := strconv.Itoa(date.Year())
	hour := fmt.Sprintf("%02d", date.Hour())
	minute := fmt.Sprintf("%02d", date.Minute())
	second := fmt.Sprintf("%02d", date.Second())

	if len(picker.Objects) >= 4 {
		dateContainer, ok := picker.Objects[1].(*fyne.Container)
		if ok && len(dateContainer.Objects) >= 3 {
			if daySelect, ok := dateContainer.Objects[0].(*widget.Select); ok {
				daySelect.SetSelected(day)
			}
			if monthSelect, ok := dateContainer.Objects[1].(*widget.Select); ok {
				monthSelect.SetSelected(month)
			}
			if yearSelect, ok := dateContainer.Objects[2].(*widget.Select); ok {
				yearSelect.SetSelected(year)
			}
		} else {
			fmt.Println("Date container does not have the expected structure")
		}

		timeContainer, ok := picker.Objects[3].(*fyne.Container)
		if ok && len(timeContainer.Objects) >= 3 {
			if hourSelect, ok := timeContainer.Objects[0].(*widget.Select); ok {
				hourSelect.SetSelected(hour)
			}
			if minuteSelect, ok := timeContainer.Objects[1].(*widget.Select); ok {
				minuteSelect.SetSelected(minute)
			}
			if secondSelect, ok := timeContainer.Objects[2].(*widget.Select); ok {
				secondSelect.SetSelected(second)
			}
		} else {
			fmt.Println("Time container does not have the expected structure")
		}
	} else {
		fmt.Println("Picker container does not have the expected structure")
	}
}

func updateParkingSpots(spots []config.ParkingSpot, container *fyne.Container) {
	if container == nil {
		fmt.Println("Error: Parking spots container is nil.")
		return
	}

	spotIndex := 0

	for _, object := range container.Objects {
		if vbox, ok := object.(*fyne.Container); ok {
			if len(vbox.Objects) == 2 {
				parkingIDContainer := vbox.Objects[0].(*fyne.Container)
				spotIDContainer := vbox.Objects[1].(*fyne.Container)

				if parkingIDEntry, ok := parkingIDContainer.Objects[1].(*widget.Entry); ok {
					spots[spotIndex].ParkingID = parkingIDEntry.Text
				}

				if spotIDEntry, ok := spotIDContainer.Objects[1].(*widget.Entry); ok {
					spots[spotIndex].SpotID = spotIDEntry.Text
				}

				spotIndex++
			}
		}
	}
}

func exportConfig(config config.Config, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func importConfig(filename string) (config.Config, error) {
	var cfg config.Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, nil
}
