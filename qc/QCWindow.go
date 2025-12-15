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
	"github.com/samuel-jimenez/qc_data_entry/util"
	"github.com/samuel-jimenez/windigo"
)

const (
	PANEL_WATER_BASED = iota
	PANEL_OIL_BASED
	PANEL_FR
	NUM_PANELS
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

type TabPanelViewer interface {
	SetFont(*windigo.Font)
	RefreshSize()
	Update(*product.QCProduct)
	Clear()
	ChangeContainer(*product.QCProduct)
	Focus()
}

/*
 * QCWindow
 *
 */
type QCWindow struct {
	*windigo.AutoForm

	product_panel *TopPanelView

	tab_bar         *windigo.TabView
	tab_panels      []TabPanelViewer
	component_panel *views.QCBlendView

	keygrab *windigo.Edit
}

func QCWindow_from_new(parent windigo.Controller) *QCWindow {
	window_title := "QC Data Entry"

	view := new(QCWindow)
	view.AutoForm = windigo.AutoForm_from_new(parent)
	view.SetText(window_title)

	view.tab_panels = make([]TabPanelViewer, NUM_PANELS)

	// menu := view.NewMenu()
	//
	// popupMenu := windigo.NewContextMenu()
	// copyMenu := popupMenu.AddItem("Copy", windigo.Shortcut{
	// 	Modifiers: windigo.ModControl,
	// 	Key:       windigo.KeyC,
	// })

	view.keygrab = windigo.NewEdit(view)
	view.keygrab.Hide()

	view.product_panel = TopPanelView_from_new(view)

	view.tab_bar = windigo.TabView_from_new(view)
	tab_wb := view.tab_bar.AddAutoPanel(formats.BLEND_WB)
	tab_oil := view.tab_bar.AddAutoPanel(formats.BLEND_OIL)
	tab_fr := view.tab_bar.AddAutoPanel(formats.BLEND_FR)

	view.component_panel = views.QCBlendView_from_new(view)

	threads.Status_bar = windigo.NewStatusBar(view)
	view.SetStatusBar(threads.Status_bar)

	//
	//
	// Dock
	view.Dock(threads.Status_bar, windigo.Bottom)
	view.Dock(view.product_panel, windigo.Top)
	view.Dock(view.component_panel, windigo.Bottom)
	view.Dock(view.tab_bar, windigo.Top)           // view.tab_bar should prefer docking at the top
	view.Dock(view.tab_bar.Panels(), windigo.Fill) // tab panels dock just below view.tab_bar and fill area

	//
	//
	//
	//
	// functionality

	view.tab_panels[PANEL_WATER_BASED] = WaterBasedPanelView_from_new(tab_wb, view.product_panel)
	view.tab_panels[PANEL_OIL_BASED] = OilBasedPanelView_from_new(tab_oil, view.product_panel)
	view.tab_panels[PANEL_FR] = FrictionReducerPanelView_from_new(tab_fr, view.product_panel)

	view.AddShortcuts()

	return view
}

func (view *QCWindow) SetFont(font *windigo.Font) {
	view.Form.SetFont(font)

	view.product_panel.SetFont(font)
	view.tab_bar.SetFont(font)
	for _, panel := range view.tab_panels {
		panel.SetFont(font)
	}

	view.component_panel.SetFont(font)

	threads.Status_bar.SetFont(font)
}

func (view *QCWindow) RefreshSize() {
	Refresh_globals(config.BASE_FONT_SIZE)

	view.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)

	view.product_panel.RefreshSize()

	view.tab_bar.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	for _, panel := range view.tab_panels {
		panel.RefreshSize()
	}

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
	view.tab_panels[PANEL_FR].ChangeContainer(qc_product)
}

func (view *QCWindow) SetCurrentTab(i int) {
	view.tab_bar.SetCurrent(i)
}

func (view *QCWindow) UpdateProduct(QC_Product *product.QCProduct) {
	view.product_panel.UpdateProduct(QC_Product)
	for _, panel := range view.tab_panels {
		panel.Update(QC_Product)
	}

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
		util.NOOP,
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

func (view *QCWindow) FocusTab() bool {
	view.tab_panels[view.tab_bar.Current()].Focus()
	return true
}
