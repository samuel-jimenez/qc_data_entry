package blender_ui

import (
	"database/sql"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
	"github.com/samuel-jimenez/windigo"
)

/*
 * FormulatorWinder
 *
 */
type FormulatorWinder interface {
	SetFont(font *windigo.Font)
	RefreshSize()
	AddShortcuts()
	Set_font_size()
	Increase_font_size() bool
	Decrease_font_size() bool
}

/*
 * FormulatorWindow
 *
 */
type FormulatorWindow struct {
	*windigo.Form

	Recipe_product_view *views.RecipeProductView
}

func NewFormulatorWindow(parent windigo.Controller) *FormulatorWindow {

	window_title := "QC Data Formulator"

	product_data := make(map[string]int64)

	var (
		product_list []string
	)

	// build window
	mainWindow := new(FormulatorWindow)
	mainWindow.Form = windigo.NewForm(parent)
	// mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	mainWindow.SetText(window_title)

	dock := windigo.NewSimpleDock(mainWindow)

	Recipe_product_view := views.NewRecipeProductView(mainWindow)

	// threads.Status_bar = windigo.NewStatusBar(mainWindow)
	// mainWindow.SetStatusBar(threads.Status_bar)
	// dock.Dock(threads.Status_bar, windigo.Bottom)

	// Dock
	dock.Dock(Recipe_product_view, windigo.Fill)

	// combobox
	GUI.Fill_combobox_from_query_rows(Recipe_product_view.Product_Field, func(row *sql.Rows) error {
		var (
			id                   int64
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		product_data[name] = id

		product_list = append(product_list, name)

		Recipe_product_view.Product_Field.AddItem(name)
		return nil
	},

		DB.DB_Select_product_info_all)
	Recipe_product_view.Update_products(product_list, product_data)

	// build object
	mainWindow.Recipe_product_view = Recipe_product_view

	return mainWindow
}

// functionality

func (view *FormulatorWindow) SetFont(font *windigo.Font) {
	view.Recipe_product_view.SetFont(windigo.DefaultFont)
	// threads.Status_bar.SetFont(windigo.DefaultFont)
}

func (view *FormulatorWindow) RefreshSize() {
	Refresh_globals(GUI.BASE_FONT_SIZE)
	view.Recipe_product_view.RefreshSize()
}

func (view *FormulatorWindow) AddShortcuts() {
	toplevel_ui.AddShortcuts(view)
}

func (mainWindow *FormulatorWindow) Set_font_size() {
	toplevel_ui.Set_font_size()
	mainWindow.SetFont(windigo.DefaultFont)
	mainWindow.RefreshSize()
}
func (view *FormulatorWindow) Increase_font_size() bool {
	GUI.BASE_FONT_SIZE += 1
	view.Set_font_size()
	return true
}
func (view *FormulatorWindow) Decrease_font_size() bool {
	GUI.BASE_FONT_SIZE -= 1
	view.Set_font_size()
	return true
}
