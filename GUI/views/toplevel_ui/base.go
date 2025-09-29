package toplevel_ui

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/windigo"
)

/*
 * TopLevelWindow
 *
 */
type TopLevelWindow interface {
	// windigo.Window
	Center()
	Show()
	OnClose() *windigo.EventManager
	RunMainLoop() int
	AddShortcut(shortcut windigo.Shortcut, action func() bool)

	// SetFont(font *windigo.Font)
	// RefreshSize()
	AddShortcuts()
	Set_font_size()
	Increase_font_size() bool
	Decrease_font_size() bool
}

func Show_window(mainWindow TopLevelWindow) {

	log.Println("Info: Process started")

	// build window
	windigo.DefaultFont = windigo.NewFont("MS Shell Dlg 2", GUI.BASE_FONT_SIZE, windigo.FontNormal)
	mainWindow.AddShortcuts()
	mainWindow.Set_font_size()

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)
	mainWindow.RunMainLoop()
}

func wndOnClose(arg *windigo.Event) {
	GUI.OKPen.Dispose()
	GUI.ErroredPen.Dispose()
	windigo.Exit()
}

func Set_font_size() {

	config.Main_config.Set("font_size", GUI.BASE_FONT_SIZE)
	config.Write_config(config.Main_config)

	old_font := windigo.DefaultFont
	windigo.DefaultFont = windigo.NewFont(old_font.Family(), GUI.BASE_FONT_SIZE, 0)
	old_font.Dispose()
}

func AddShortcuts(view TopLevelWindow) {

	// Resize handling
	view.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMPlus}, // +
		view.Increase_font_size,
	)

	view.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModControl, Key: windigo.KeyOEMMinus}, // -
		view.Decrease_font_size,
	)

}
