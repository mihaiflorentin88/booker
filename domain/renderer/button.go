package renderer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (r *Renderer) fileMenu() *widget.Button {
	var fileMenuButton *widget.Button
	fileMenuButton = widget.NewButton("File", func() {
		fileMenu := fyne.NewMenu("File",
			r.importSettingsMenuItem(),
			r.exportSettingsMenuItem(),
		)
		widget.ShowPopUpMenuAtPosition(fileMenu, r.Window.Canvas(), fileMenuButton.Position().Add(fyne.NewPos(0, fileMenuButton.Size().Height)))
	})
	fileMenuButton.Importance = widget.LowImportance
	return fileMenuButton
}

func (r *Renderer) importSettingsMenuItem() *fyne.MenuItem {
	return fyne.NewMenuItem("Import settings", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			r.importConfig(reader.URI().Path())
			r.ParkingBooker.UpdateConfig()
		}, r.Window).Show()
	})
}

func (r *Renderer) exportSettingsMenuItem() *fyne.MenuItem {
	return fyne.NewMenuItem("Export settings", func() {
		r.ParkingBooker.UpdateConfig()
		dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			r.exportSettings(writer.URI().Path())
		}, r.Window).Show()
	})
}
