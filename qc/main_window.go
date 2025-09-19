package qc

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

func Show_window() {

	log.Println("Info: Process started")
	// log.Println(time.Now().UTC().UnixNano())

	//
	// build window
	windigo.DefaultFont = windigo.NewFont("MS Shell Dlg 2", GUI.BASE_FONT_SIZE, windigo.FontNormal)
	mainWindow := NewQCWindow(nil)
	mainWindow.set_font_size()

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop() // Must call to start event loop.
}

func wndOnClose(arg *windigo.Event) {
	GUI.OKPen.Dispose()
	GUI.ErroredPen.Dispose()
	windigo.Exit()
}
