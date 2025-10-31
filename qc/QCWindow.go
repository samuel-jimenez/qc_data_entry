package qc

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

/*
 * QCWinder
 *
 */
type QCWinder interface {
	clear_lot(product_id int)
	SetFont()
	RefreshSize()
	Set_font_size()
	Increase_font_size() bool
	Decrease_font_size() bool
	keygrab_start() bool
	keygrab_end() bool

	AddShortcuts()
	ChangeContainer(*product.QCProduct)
	SetCurrentTab(i int)
	UpdateProduct(*product.QCProduct)
	SetBlend(*product.QCProduct)
}

/*
 * QCWindow
 *
 */
type QCWindow struct {
	*windigo.Form

	product_panel *TopPanelView

	tabs              *windigo.TabView
	panel_water_based *WaterBasedPanelView
	panel_oil_based   *OilBasedPanelView
	panel_fr          *FrictionReducerPanelView
	component_panel   *views.QCBlendView

	keygrab *windigo.Edit
}

func QCWindow_from_new(parent windigo.Controller) *QCWindow {
	window_title := "QC Data Entry"

	view := new(QCWindow)
	view.Form = windigo.NewForm(parent)
	view.SetText(window_title)

	// menu := view.NewMenu()
	//
	// popupMenu := windigo.NewContextMenu()
	// copyMenu := popupMenu.AddItem("Copy", windigo.Shortcut{
	// 	Modifiers: windigo.ModControl,
	// 	Key:       windigo.KeyC,
	// })

	view.keygrab = windigo.NewEdit(view)
	view.keygrab.Hide()

	dock := windigo.NewSimpleDock(view)

	product_panel := NewTopPanelView(view)

	tabs := windigo.NewTabView(view)
	tab_wb := tabs.AddAutoPanel(formats.BLEND_WB)
	tab_oil := tabs.AddAutoPanel(formats.BLEND_OIL)
	tab_fr := tabs.AddAutoPanel(formats.BLEND_FR)

	view.component_panel = views.QCBlendView_from_new(view)

	threads.Status_bar = windigo.NewStatusBar(view)
	view.SetStatusBar(threads.Status_bar)

	//
	//
	// Dock
	dock.Dock(threads.Status_bar, windigo.Bottom)
	dock.Dock(product_panel, windigo.Top)
	dock.Dock(view.component_panel, windigo.Bottom)
	dock.Dock(tabs, windigo.Top)           // tabs should prefer docking at the top
	dock.Dock(tabs.Panels(), windigo.Fill) // tab panels dock just below tabs and fill area

	//
	//
	//
	//
	// functionality

	panel_water_based := Show_water_based(tab_wb, product_panel)
	panel_oil_based := Show_oil_based(tab_oil, product_panel)
	panel_fr := Show_fr(tab_fr, product_panel)

	view.AddShortcuts()

	// build object
	view.product_panel = product_panel
	view.tabs = tabs
	view.panel_water_based = panel_water_based
	view.panel_oil_based = panel_oil_based
	view.panel_fr = panel_fr

	return view
}

func (view *QCWindow) SetFont(font *windigo.Font) {
	view.Form.SetFont(font)

	view.product_panel.SetFont(font)
	view.tabs.SetFont(font)
	view.panel_water_based.SetFont(font)
	view.panel_oil_based.SetFont(font)
	view.panel_fr.SetFont(font)

	view.component_panel.SetFont(font)

	threads.Status_bar.SetFont(font)
}

func (view *QCWindow) RefreshSize() {
	Refresh_globals(config.BASE_FONT_SIZE)

	view.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)

	// view.product_panel.RefreshSize()
	view.product_panel.RefreshSize(config.BASE_FONT_SIZE)

	view.tabs.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.panel_water_based.RefreshSize()
	view.panel_oil_based.RefreshSize()
	view.panel_fr.RefreshSize()

	view.component_panel.RefreshSize()
}

func (view *QCWindow) Set_font_size() {
	toplevel_ui.Set_font_size()
	view.SetFont(windigo.DefaultFont)
	view.RefreshSize()
}

func (view *QCWindow) Increase_font_size() bool {
	config.BASE_FONT_SIZE += 1
	view.Set_font_size()
	return true
}

func (view *QCWindow) Decrease_font_size() bool {
	config.BASE_FONT_SIZE -= 1
	view.Set_font_size()
	return true
}

func (view *QCWindow) keygrab_start() bool {
	view.keygrab.SetText("")
	view.keygrab.SetFocus()
	return true
}

func (view *QCWindow) keygrab_end() bool {
	var product QR.QRJson

	qr_json := view.keygrab.Text() + "}"
	log.Println("debug: ReadFromScanner: ", qr_json)
	err := json.Unmarshal([]byte(qr_json), &product)
	if err == nil {
		view.product_panel.PopQRData(product)
	} else {
		log.Printf("error: [%s]: %q\n", "qr_json_mainWindow.keygrab", err)
	}
	view.keygrab.SetText("")
	view.SetFocus()
	return true
}

func (view *QCWindow) AddShortcuts() {
	// QR keyboard handling

	view.keygrab.OnSetFocus().Bind(func(e *windigo.Event) {
		view.keygrab.SetText("{")
		view.keygrab.SelectText(1, 1)
	})

	view.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM4}, // {
		view.keygrab_start,
	)

	view.AddShortcut(windigo.Shortcut{Modifiers: windigo.ModShift, Key: windigo.KeyOEM6}, // }
		view.keygrab_end,
	)

	// Resize handling
	toplevel_ui.AddShortcuts(view)
}

func (view *QCWindow) ChangeContainer(qc_product *product.QCProduct) {
	view.panel_fr.ChangeContainer(qc_product)
}

func (view *QCWindow) SetCurrentTab(i int) {
	view.tabs.SetCurrent(i)
}

func (view *QCWindow) UpdateProduct(QC_Product *product.QCProduct) {
	view.product_panel.UpdateProduct(QC_Product)
	view.panel_water_based.Update(QC_Product)
	view.panel_oil_based.Update(QC_Product)
	view.panel_fr.Update(QC_Product)
	view.ChangeContainer(QC_Product) // TODO recip00
	// extract to fn, move componenet panel?

	if QC_Product.Blend != nil {
		view.component_panel.UpdateBlend(QC_Product.Blend)
		return
	}

	var recipe_data blender.ProductRecipe
	// proc_name := "RecipeProduct.GetRecipes"
	proc_name := "FrictionReducerPanelView.GetRecipes"

	DB.Forall(proc_name,
		func() {},
		func(row *sql.Rows) {
			if err := row.Scan(
				&recipe_data.Recipe_id,
			); err != nil {
				log.Fatal("Crit: [RecipeProduct GetRecipes]: ", proc_name, err)
			}
		},
		DB.DB_Select_product_recipe, QC_Product.Product_id)

	recipe_data.GetComponents()
	view.component_panel.UpdateRecipe(&recipe_data)
}

func (view *QCWindow) SetBlend(QC_Product *product.QCProduct) {
	QC_Product.SetBlend(view.component_panel.Get())
	QC_Product.SaveBlend()
}

func (view *QCWindow) ComponentsEnable() {
	view.component_panel.Enable()
}

func (view *QCWindow) ComponentsDisable() {
	view.component_panel.Disable()
}
