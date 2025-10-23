package blender_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
	"github.com/samuel-jimenez/qc_data_entry/config"
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

	newMonikerMenu_OnClick(*windigo.Event)
}

/*
 * BlenderWindow
 *
 */
type BlenderWindow struct {
	*windigo.Form

	Blend_product_view *BlendStrappingProductView
}

func BlenderWindow_from_new(parent windigo.Controller) *BlenderWindow {
	window_title := "Blender"

	// build window
	view := new(BlenderWindow)
	view.Form = windigo.NewForm(parent)
	// mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	view.SetText(window_title)

	menu := view.NewMenu()
	// TODO settings
	fileMenu := menu.AddSubMenu("File")
	newMenu := fileMenu.AddSubMenu("New")
	newMonikerMenu := newMenu.AddItem("Moniker", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyN,
	})
	newMonikerMenu.OnClick().Bind(view.newMonikerMenu_OnClick)

	dock := windigo.NewSimpleDock(view)
	// BlendVessel

	view.Blend_product_view = BlendStrappingProductView_from_new(view)

	// threads.Status_bar = windigo.NewStatusBar(mainWindow)
	// mainWindow.SetStatusBar(threads.Status_bar)
	// dock.Dock(threads.Status_bar, windigo.Bottom)

	// Dock
	dock.Dock(view.Blend_product_view, windigo.Fill)

	return view
}

// functionality

func (view *BlenderWindow) SetFont(font *windigo.Font) {

	view.Blend_product_view.SetFont(windigo.DefaultFont)
	// threads.Status_bar.SetFont(windigo.DefaultFont)

}

func (view *BlenderWindow) RefreshSize() {
	Refresh_globals(config.BASE_FONT_SIZE)
	view.Blend_product_view.RefreshSize()
}

func (view *BlenderWindow) AddShortcuts() {
	toplevel_ui.AddShortcuts(view)
}

func (view *BlenderWindow) Set_font_size() {
	toplevel_ui.Set_font_size()
	view.SetFont(windigo.DefaultFont)
	view.RefreshSize()
}
func (view *BlenderWindow) Increase_font_size() bool {
	config.BASE_FONT_SIZE += 1
	view.Set_font_size()
	return true
}
func (view *BlenderWindow) Decrease_font_size() bool {
	config.BASE_FONT_SIZE -= 1
	view.Set_font_size()
	return true
}

func (view *BlenderWindow) newMonikerMenu_OnClick(*windigo.Event) {
	MonikerView := views.MonikerView_from_new(view)
	MonikerView.SetModal(false)
	MonikerView.Show()
	MonikerView.RefreshSize()
}
