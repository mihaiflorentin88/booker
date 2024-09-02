package component

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

func CreateCustomDatePicker() *fyne.Container {
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

func AssembleDateFromPicker(picker *fyne.Container, window fyne.Window) string {
	if err := ValidateDatePicker(picker); err != nil {
		dialog.ShowError(err, window)
		return ""
	}
	day := picker.Objects[1].(*fyne.Container).Objects[0].(*widget.Select).Selected
	month := picker.Objects[1].(*fyne.Container).Objects[1].(*widget.Select).Selected[:2]
	year := picker.Objects[1].(*fyne.Container).Objects[2].(*widget.Select).Selected

	hour := picker.Objects[3].(*fyne.Container).Objects[0].(*widget.Select).Selected
	minute := picker.Objects[3].(*fyne.Container).Objects[1].(*widget.Select).Selected
	second := picker.Objects[3].(*fyne.Container).Objects[2].(*widget.Select).Selected

	return fmt.Sprintf("%s-%s-%s %s:%s:%s", day, month, year, hour, minute, second)
}

func ValidateDatePicker(picker *fyne.Container) error {
	day := picker.Objects[1].(*fyne.Container).Objects[0].(*widget.Select).Selected
	month := picker.Objects[1].(*fyne.Container).Objects[1].(*widget.Select).Selected
	year := picker.Objects[1].(*fyne.Container).Objects[2].(*widget.Select).Selected

	if day == "" || month == "" || year == "" {
		return fmt.Errorf("date is not fully selected")
	}

	hour := picker.Objects[3].(*fyne.Container).Objects[0].(*widget.Select).Selected
	minute := picker.Objects[3].(*fyne.Container).Objects[1].(*widget.Select).Selected
	second := picker.Objects[3].(*fyne.Container).Objects[2].(*widget.Select).Selected

	if hour == "" || minute == "" || second == "" {
		return fmt.Errorf("time is not fully selected")
	}

	return nil
}

func UpdateDatePicker(picker *fyne.Container, dateStr string) {
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
