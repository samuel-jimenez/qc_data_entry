package blender_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlenderWinder
 *
 */
type BlenderWinder interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	AddShortcuts()
	Set_font_size()
	Increase_font_size() bool
	Decrease_font_size() bool
}

/*
 * BlenderWindow
 *
 */
type BlenderWindow struct {
	*windigo.Form

	// Blend_product_view *views.BlendProductView
	Blend_product_view *BlendProductView
}

func NewBlenderWindow(parent windigo.Controller) *BlenderWindow {
	window_title := "QC Data Blender"

	// build window
	mainWindow := new(BlenderWindow)
	mainWindow.Form = windigo.NewForm(parent)
	// mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	mainWindow.SetText(window_title)

	dock := windigo.NewSimpleDock(mainWindow)

	// mainWindow.Blend_product_view = views.NewBlendProductView(mainWindow)
	mainWindow.Blend_product_view = NewBlendProductView(mainWindow)

	// threads.Status_bar = windigo.NewStatusBar(mainWindow)
	// mainWindow.SetStatusBar(threads.Status_bar)
	// dock.Dock(threads.Status_bar, windigo.Bottom)

	// Dock
	dock.Dock(mainWindow.Blend_product_view, windigo.Fill)

	return mainWindow
}

// functionality

func (view *BlenderWindow) SetFont(font *windigo.Font) {

	view.Blend_product_view.SetFont(windigo.DefaultFont)
	// threads.Status_bar.SetFont(windigo.DefaultFont)

}

func (view *BlenderWindow) RefreshSize() {
	Refresh_globals(GUI.BASE_FONT_SIZE)
	view.Blend_product_view.RefreshSize()
}

func (view *BlenderWindow) AddShortcuts() {
	toplevel_ui.AddShortcuts(view)
}

func (mainWindow *BlenderWindow) Set_font_size() {
	toplevel_ui.Set_font_size()
	mainWindow.SetFont(windigo.DefaultFont)
	mainWindow.RefreshSize()
}
func (view *BlenderWindow) Increase_font_size() bool {
	GUI.BASE_FONT_SIZE += 1
	view.Set_font_size()
	return true
}
func (view *BlenderWindow) Decrease_font_size() bool {
	GUI.BASE_FONT_SIZE -= 1
	view.Set_font_size()
	return true
}
