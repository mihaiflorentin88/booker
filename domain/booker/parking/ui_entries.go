package parking

import (
	"booking/domain/renderer/component"
	"booking/port/config"
	"booking/port/contract"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type UIEntries struct {
	Config                      *config.Config
	parkingClient               contract.ParkingClientInterface
	App                         fyne.App
	Window                      fyne.Window
	UsernameEntry               *widget.Entry
	PasswordEntry               *widget.Entry
	GoogleAPIKeyEntry           *widget.Entry
	StandByCheckbox             *widget.Check
	StandByStartTimePicker      *fyne.Container
	StandByEndTimePicker        *fyne.Container
	BookStartTimePicker         *fyne.Container
	OffsetEntry                 *widget.Entry
	ParkingSpotContainer        *fyne.Container
	ParkingSpotScrollContainer  *container.Scroll
	AddParkingSlotBtn           *widget.Button
	ExecuteBtn                  *widget.Button
	Tab                         *container.TabItem
	RunnerOutputWindow          fyne.Window
	RunnerOutputLabel           *widget.Label
	RunnerOutputScrollContainer *container.Scroll
}

func NewUIEntries(config *config.Config, parkingClient contract.ParkingClientInterface, app fyne.App, window fyne.Window) *UIEntries {
	ui := &UIEntries{
		Config:               config,
		parkingClient:        parkingClient,
		App:                  app,
		Window:               window,
		UsernameEntry:        widget.NewEntry(),
		PasswordEntry:        widget.NewPasswordEntry(),
		GoogleAPIKeyEntry:    widget.NewEntry(),
		StandByCheckbox:      widget.NewCheck("", nil),
		OffsetEntry:          widget.NewEntry(),
		ParkingSpotContainer: container.NewVBox(),
		RunnerOutputWindow:   app.NewWindow("Output Window"),
		RunnerOutputLabel:    widget.NewLabel("Output will appear here."),
	}
	ui.UsernameEntry.SetPlaceHolder("Username")
	ui.PasswordEntry.SetPlaceHolder("Password")
	ui.OffsetEntry.SetPlaceHolder("Offset")
	ui.GoogleAPIKeyEntry.SetPlaceHolder("Google API Key")
	ui.StandByStartTimePicker = component.CreateCustomDatePicker()
	ui.StandByEndTimePicker = component.CreateCustomDatePicker()
	ui.BookStartTimePicker = component.CreateCustomDatePicker()
	ui.ParkingSpotScrollContainer = container.NewVScroll(ui.ParkingSpotContainer)
	ui.ParkingSpotScrollContainer.SetMinSize(fyne.NewSize(600, 500))
	ui.AddParkingSlotBtn = ui.getParkingSpotBtn()
	ui.RunnerOutputLabel.Wrapping = fyne.TextWrapWord
	ui.RunnerOutputScrollContainer = container.NewVScroll(ui.RunnerOutputLabel)
	ui.RunnerOutputScrollContainer.SetMinSize(fyne.NewSize(700, 300))
	ui.ExecuteBtn = ui.getExecuteBtn()
	ui.Tab = container.NewTabItem("Parking", ui.parkingTab())
	return ui
}

func (ui *UIEntries) UpdateConfig() {
	ui.Config.Parking.Credentials.Username = ui.UsernameEntry.Text
	ui.Config.Parking.Credentials.Password = ui.PasswordEntry.Text
	ui.Config.Parking.Credentials.GoogleAPIKey = ui.GoogleAPIKeyEntry.Text
	ui.Config.Parking.StandBy.CanStandby = ui.StandByCheckbox.Checked

	ui.Config.Parking.StandBy.StartTime = component.AssembleDateFromPicker(ui.StandByStartTimePicker, ui.Window)
	ui.Config.Parking.StandBy.EndTime = component.AssembleDateFromPicker(ui.StandByEndTimePicker, ui.Window)
	ui.Config.Parking.Schedule.Date = component.AssembleDateFromPicker(ui.BookStartTimePicker, ui.Window)
	ui.Config.Parking.Schedule.Offset = 0
	if offset, err := strconv.Atoi(ui.OffsetEntry.Text); err == nil {
		ui.Config.Parking.Schedule.Offset = offset
	}
	ui.updateParkingSpotsConfig()
}

func (ui *UIEntries) UpdateUI() {
	if ui.Config == nil {
		fmt.Println("Error: Configuration data is missing.")
		return
	}

	ui.UsernameEntry.SetText(ui.Config.Parking.Credentials.Username)
	ui.PasswordEntry.SetText(ui.Config.Parking.Credentials.Password)
	ui.GoogleAPIKeyEntry.SetText(ui.Config.Parking.Credentials.GoogleAPIKey)

	ui.StandByCheckbox.SetChecked(ui.Config.Parking.StandBy.CanStandby)

	ui.OffsetEntry.SetText(strconv.Itoa(ui.Config.Parking.Schedule.Offset))
	component.UpdateDatePicker(ui.StandByStartTimePicker, ui.Config.Parking.StandBy.StartTime)
	component.UpdateDatePicker(ui.StandByEndTimePicker, ui.Config.Parking.StandBy.EndTime)
	component.UpdateDatePicker(ui.BookStartTimePicker, ui.Config.Parking.Schedule.Date)
	ui.ParkingSpotContainer.Objects = nil
	for i := range ui.Config.Parking.ParkingSpot {
		ui.ParkingSpotContainer.Add(ui.createParkingSpotEntry(&ui.Config.Parking.ParkingSpot[i]))
		ui.ParkingSpotContainer.Add(widget.NewSeparator())
	}
	ui.ParkingSpotContainer.Refresh()
}

func (ui *UIEntries) updateParkingSpotsConfig() {
	if ui.ParkingSpotContainer == nil {
		fmt.Println("Error: Parking spots container is nil.")
		return
	}
	spotIndex := 0
	for _, object := range ui.ParkingSpotContainer.Objects {
		if vbox, ok := object.(*fyne.Container); ok {
			if len(vbox.Objects) == 2 {
				parkingIDContainer := vbox.Objects[0].(*fyne.Container)
				spotIDContainer := vbox.Objects[1].(*fyne.Container)

				if parkingIDEntry, ok := parkingIDContainer.Objects[1].(*widget.Entry); ok {
					ui.Config.Parking.ParkingSpot[spotIndex].ParkingID = parkingIDEntry.Text
				}

				if spotIDEntry, ok := spotIDContainer.Objects[1].(*widget.Entry); ok {
					ui.Config.Parking.ParkingSpot[spotIndex].SpotID = spotIDEntry.Text
				}
				spotIndex++
			}
		}
	}
}

func (ui *UIEntries) createParkingSpotEntry(spot *config.ParkingSpot) *fyne.Container {
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
func (ui *UIEntries) credentialsContent() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Credentials"),
		widget.NewForm(
			widget.NewFormItem("Username", ui.UsernameEntry),
			widget.NewFormItem("Password", ui.PasswordEntry),
			widget.NewFormItem("Google API Key", ui.GoogleAPIKeyEntry),
		),
	)
}

func (ui *UIEntries) standbyContent() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Standby Prevention"),
		widget.NewForm(
			widget.NewFormItem("Prevent OS Standby", ui.StandByCheckbox),
			widget.NewFormItem("Start Time", ui.StandByStartTimePicker),
			widget.NewFormItem("End Time", ui.StandByEndTimePicker),
		),
	)
}

func (ui *UIEntries) scheduleContent() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Schedule"),
		widget.NewForm(
			widget.NewFormItem("Date", ui.BookStartTimePicker),
			widget.NewFormItem("Offset", ui.OffsetEntry),
		),
	)
}

func (ui *UIEntries) parkingSpotsContent() *fyne.Container {
	content := container.NewVBox(
		widget.NewLabel("Parking Spots"),
		ui.ParkingSpotScrollContainer,
		//	ui.AddParkingSlotBtn,
	)
	return container.NewBorder(nil, ui.AddParkingSlotBtn, nil, nil, content)
}

func (ui *UIEntries) parkingTab() *fyne.Container {
	content := container.NewMax()
	content.Add(ui.credentialsContent())
	navbar := container.NewVBox(
		widget.NewButton("Credentials", func() {
			content.Objects = []fyne.CanvasObject{ui.credentialsContent()}
			content.Refresh()
		}),
		widget.NewButton("Standby Prevention", func() {
			content.Objects = []fyne.CanvasObject{ui.standbyContent()}
			content.Refresh()
		}),
		widget.NewButton("Schedule", func() {
			content.Objects = []fyne.CanvasObject{ui.scheduleContent()}
			content.Refresh()
		}),
		widget.NewButton("Parking Spots", func() {
			content.Objects = []fyne.CanvasObject{ui.parkingSpotsContent()}
			content.Refresh()
		}),
	)
	//buttons := container.NewHBox(ui.ExecuteBtn)
	//layout := container.NewBorder(nil, buttons, navbar, nil, content)
	buttons := container.NewBorder(nil, nil, nil, nil, ui.ExecuteBtn)
	layout := container.NewBorder(nil, buttons, navbar, nil, content)
	return layout
}
