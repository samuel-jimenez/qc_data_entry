package viewer

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

//TODO viewerf_refactor: split datatypes, views

/*
 * DataViewerPanelViewer
 *
 */
type DataViewerPanelViewer interface {
	windigo.Controller
	SetFont(font *windigo.Font)
	RefreshSize()
	SetMainWindow(mainWindow *ViewerWindow)

	ClearFilters(e *windigo.Event)

	update_product(id int, name string)
	update_lot(id int, name string)
	update_sample(id int, name string)

	clear_moniker()
	clear_product()
	clear_lot()

	SetTable()

	// listeners

	moniker_field_OnChange(e *windigo.Event)
	moniker_clear_button_OnClick(e *windigo.Event)

	product_field_OnChange(e *windigo.Event)
	product_clear_button_OnClick(e *windigo.Event)

	lot_field_OnChange(e *windigo.Event)
	lot_clear_button_OnClick(e *windigo.Event)
	lot_add_button_OnClick(e *windigo.Event)

	filter_button_OnClick(e *windigo.Event)
	search_button_OnClick(e *windigo.Event)
	reprint_sample_button_OnClick(e *windigo.Event)
	regen_sample_button_OnClick(e *windigo.Event)
	export_json_button_OnClick(e *windigo.Event)
}

/*
 * DataViewerPanelView
 *
 */
type DataViewerPanelView struct {
	*windigo.AutoPanel

	mainWindow *ViewerWindow

	product_map  map[string]int
	moniker_map  map[string][]string
	product_data []string

	SQLFilters *SQLFilterList

	ProductFilter,
	LotFilter,
	MonikerFilter *SQLFilterDiscrete

	lot_field *GUI.SearchBox

	product_panel, lot_panel *windigo.AutoPanel

	product_field, moniker_field,
	sample_field *GUI.ComboBox

	product_clear_button, lot_clear_button, lot_add_button, moniker_clear_button,
	filter_button, search_button, clear_button,
	reprint_sample_button, regen_sample_button, export_json_button *windigo.PushButton
}

// parent

func NewDataViewerPanelView(mainWindow windigo.Controller) *DataViewerPanelView {

	view := new(DataViewerPanelView)

	view.product_map = make(map[string]int)
	view.moniker_map = make(map[string][]string)

	// LotFilter := &SQLFilterDiscrete{COL_KEY_LOT,nil}
	view.SQLFilters = NewSQLFilterList()

	view.ProductFilter = NewSQLFilterDiscrete(
		COL_KEY_PRODUCT,
		nil,
	)
	view.LotFilter = NewSQLFilterDiscrete(
		COL_KEY_LOT,
		nil,
	)
	view.MonikerFilter = NewSQLFilterDiscrete(
		COL_KEY_MONIKER,
		nil,
	)

	clear_field_label := "-"
	add_field_label := "+"

	filter_label := "Filter"
	search_label := "Search"
	clear_label := "Clear"

	reprint_sample_label := "Reprint Sample"
	regen_sample_label := "Regen Sample"
	export_json_label := "Export JSON"

	view.AutoPanel = windigo.NewAutoPanel(mainWindow)

	//TODO array db_select_all_product

	view.product_panel = windigo.NewAutoPanel(view.AutoPanel)

	view.product_field = GUI.NewListComboBox(view.product_panel, COL_LABEL_PRODUCT)

	view.product_clear_button = windigo.NewPushButton(view.product_panel)
	view.product_clear_button.SetText(clear_field_label)

	view.moniker_field = GUI.NewListComboBox(view.product_panel, COL_LABEL_MONIKER)
	view.moniker_clear_button = windigo.NewPushButton(view.product_panel)
	view.moniker_clear_button.SetText(clear_field_label)

	view.product_panel.Dock(view.product_field, windigo.Left)
	view.product_panel.Dock(view.product_clear_button, windigo.Left)
	view.product_panel.Dock(view.moniker_field, windigo.Left)
	view.product_panel.Dock(view.moniker_clear_button, windigo.Left)

	view.lot_panel = windigo.NewAutoPanel(view.AutoPanel)

	view.lot_field = GUI.NewLabeledListSearchBox(view.lot_panel, COL_LABEL_LOT)
	view.lot_clear_button = windigo.NewPushButton(view.lot_panel)
	view.lot_clear_button.SetText(clear_field_label)
	view.lot_add_button = windigo.NewPushButton(view.lot_panel)
	view.lot_add_button.SetText(add_field_label)

	view.sample_field = GUI.NewListComboBox(view.lot_panel, COL_LABEL_SAMPLE_PT)

	view.lot_panel.Dock(view.lot_field, windigo.Left)
	view.lot_panel.Dock(view.lot_clear_button, windigo.Left)
	view.lot_panel.Dock(view.lot_add_button, windigo.Left)
	view.lot_panel.Dock(view.sample_field, windigo.Left)

	view.filter_button = windigo.NewPushButton(view.AutoPanel)
	view.filter_button.SetText(filter_label)

	view.search_button = windigo.NewPushButton(view.AutoPanel)
	view.search_button.SetText(search_label)

	view.clear_button = windigo.NewPushButton(view.AutoPanel)
	view.clear_button.SetText(clear_label)

	view.reprint_sample_button = windigo.NewPushButton(view.AutoPanel)
	view.reprint_sample_button.SetText(reprint_sample_label)

	view.regen_sample_button = windigo.NewPushButton(view.AutoPanel)
	view.regen_sample_button.SetText(regen_sample_label)

	view.export_json_button = windigo.NewPushButton(view.AutoPanel)
	view.export_json_button.SetText(export_json_label)

	view.AutoPanel.Dock(view.product_panel, windigo.Top)
	view.AutoPanel.Dock(view.lot_panel, windigo.Top)
	view.AutoPanel.Dock(view.filter_button, windigo.Left)
	view.AutoPanel.Dock(view.search_button, windigo.Left)
	view.AutoPanel.Dock(view.clear_button, windigo.Left)

	view.AutoPanel.Dock(view.reprint_sample_button, windigo.Left)
	view.AutoPanel.Dock(view.regen_sample_button, windigo.Left)
	view.AutoPanel.Dock(view.export_json_button, windigo.Left)

	// functionality
	view.SQLFilters.Filters = []SQLFilter{view.ProductFilter, view.LotFilter, view.MonikerFilter}

	// clear := func() {
	// 	clear_product()
	// 	clear_lot()
	// }

	// combobox
	DB.Forall_err("DataViewerPanelView.fill-product-moniker",
		func() {
			view.product_field.DeleteAllItems()
		},
		func(row *sql.Rows) error {
			var (
				id                   int
				internal_name        string
				product_moniker_name string
			)
			if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
				return err
			}
			name := product_moniker_name + " " + internal_name

			view.product_map[name] = id
			view.product_data = append(view.product_data, name)

			view.moniker_map[product_moniker_name] = append(view.moniker_map[product_moniker_name], name)

			view.product_field.AddItem(name)
			return nil
		},
		DB.DB_Select_product_info_all)

	view.lot_field.Fill_FromFnQuery(
		view.update_lot,
		DB.DB_Select_product_lot_all)

	GUI.Fill_combobox_from_query(view.moniker_field, DB.DB_Select_all_product_moniker)
	GUI.Fill_combobox_from_query_fn(view.sample_field, view.update_sample, DB.DB_Select_all_sample_points)

	// listeners
	view.product_field.OnSelectedChange().Bind(view.product_field_OnChange)
	view.lot_field.OnSelectedChange().Bind(view.lot_field_OnChange)
	view.moniker_field.OnSelectedChange().Bind(view.moniker_field_OnChange)

	view.filter_button.OnClick().Bind(view.filter_button_OnClick)
	view.search_button.OnClick().Bind(view.search_button_OnClick)
	view.product_clear_button.OnClick().Bind(view.product_clear_button_OnClick)
	view.moniker_clear_button.OnClick().Bind(view.moniker_clear_button_OnClick)
	view.lot_clear_button.OnClick().Bind(view.lot_clear_button_OnClick)
	view.lot_add_button.OnClick().Bind(view.lot_add_button_OnClick)
	view.clear_button.OnClick().Bind(view.ClearFilters)

	view.reprint_sample_button.OnClick().Bind(view.reprint_sample_button_OnClick)
	view.regen_sample_button.OnClick().Bind(view.regen_sample_button_OnClick)
	view.export_json_button.OnClick().Bind(view.export_json_button_OnClick)

	return view
}

func (view *DataViewerPanelView) SetFont(font *windigo.Font) {

	view.product_field.SetFont(font)
	view.moniker_field.SetFont(font)
	view.lot_field.SetFont(font)
	view.sample_field.SetFont(font)

	view.product_clear_button.SetFont(font)
	view.moniker_clear_button.SetFont(font)
	view.lot_clear_button.SetFont(font)
	view.lot_add_button.SetFont(font)

	view.filter_button.SetFont(font)
	view.search_button.SetFont(font)
	view.clear_button.SetFont(font)

	view.reprint_sample_button.SetFont(font)
	view.regen_sample_button.SetFont(font)
	view.export_json_button.SetFont(font)

}

func (view *DataViewerPanelView) RefreshSize() {

	reprint_sample_button_margin := GUI.HPANEL_MARGIN + GUI.LABEL_WIDTH + GUI.PRODUCT_FIELD_WIDTH + GUI.SMOL_BUTTON_WIDTH + GUI.TOP_PANEL_INTER_SPACER_WIDTH - 3*GUI.BUTTON_WIDTH

	top_panel_height := GUI.TOP_SPACER_HEIGHT + GUI.INTER_SPACER_HEIGHT + 2*GUI.PRODUCT_FIELD_HEIGHT + GUI.BUTTON_HEIGHT

	view.AutoPanel.SetSize(GUI.TOP_PANEL_WIDTH, top_panel_height)
	view.product_panel.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.lot_panel.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.product_clear_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.moniker_clear_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.moniker_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.lot_clear_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.lot_add_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.sample_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.lot_panel.SetMarginTop(GUI.INTER_SPACER_HEIGHT)
	view.lot_panel.SetMarginLeft(GUI.HPANEL_MARGIN)
	view.sample_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH - GUI.SMOL_BUTTON_WIDTH)

	view.moniker_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.product_panel.SetMarginLeft(GUI.HPANEL_MARGIN)
	view.product_panel.SetMarginTop(GUI.TOP_SPACER_HEIGHT)

	view.filter_button.SetSize(GUI.BUTTON_WIDTH, GUI.OFF_AXIS)
	view.search_button.SetSize(GUI.BUTTON_WIDTH, GUI.OFF_AXIS)
	view.clear_button.SetSize(GUI.BUTTON_WIDTH, GUI.OFF_AXIS)

	view.reprint_sample_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.reprint_sample_button.SetMarginLeft(reprint_sample_button_margin)
	view.regen_sample_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.export_json_button.SetSize(GUI.REPRINT_BUTTON_WIDTH, GUI.OFF_AXIS)

}

func (view *DataViewerPanelView) SetMainWindow(mainWindow *ViewerWindow) {
	view.mainWindow = mainWindow
}

func (view *DataViewerPanelView) ClearFilters(e *windigo.Event) {
	view.mainWindow.ClearFilters()
}

func (view *DataViewerPanelView) update_product(id int, name string) {
	view.product_field.AddItem(name)
	// COL_ITEMS_LOT = append(COL_ITEMS_LOT, name)
}
func (view *DataViewerPanelView) update_lot(id int, name string) {
	view.lot_field.AddItem(name)
	COL_ITEMS_LOT = append(COL_ITEMS_LOT, name)
}
func (view *DataViewerPanelView) update_sample(id int, name string) {
	view.sample_field.AddItem(name)
	COL_ITEMS_SAMPLE_PT = append(COL_ITEMS_SAMPLE_PT, name)
}

func (view *DataViewerPanelView) clear_moniker() {
	view.moniker_field.SetSelectedItem(-1)
	view.MonikerFilter.Set(nil)
	view.clear_product()

	view.product_field.DeleteAllItems()
	for id, name := range view.product_data {
		view.update_product(id, name)
	}
}

func (view *DataViewerPanelView) clear_product() {
	view.clear_lot()
	view.product_field.SetSelectedItem(-1)
	view.ProductFilter.Set(nil)

	COL_ITEMS_LOT = nil
	view.lot_field.Update(nil)
	select_lot(view.update_lot, view.SQLFilters.Get())

	COL_ITEMS_SAMPLE_PT = nil
	GUI.Fill_combobox_from_query_fn(view.sample_field, view.update_sample, DB.DB_Select_all_sample_points)
}

func (view *DataViewerPanelView) clear_lot() {
	view.lot_field.SetSelectedItem(-1)
	view.LotFilter.Set(nil)
}

func (view *DataViewerPanelView) SetTable() {
	view.mainWindow.SetTable(select_samples(view.SQLFilters.Get()))
}

// listeners

func (view *DataViewerPanelView) moniker_field_OnChange(e *windigo.Event) {

	product_moniker := view.moniker_field.GetSelectedItem()
	view.MonikerFilter.Set([]string{product_moniker})

	view.clear_product()

	view.product_field.DeleteAllItems()
	for id, name := range view.moniker_map[product_moniker] {
		view.update_product(id, name)
	}

	view.SetTable()
	view.mainWindow.UpdateFilterListView()

}

func (view *DataViewerPanelView) moniker_clear_button_OnClick(e *windigo.Event) {

	view.clear_moniker()
	view.SetTable()
	view.mainWindow.UpdateFilterListView()

}

func (view *DataViewerPanelView) product_field_OnChange(e *windigo.Event) {
	// product_field_pop_data(product_field.GetSelectedItem())

	product_id := view.product_map[view.product_field.GetSelectedItem()]
	view.ProductFilter.Set([]string{strconv.Itoa(product_id)})
	view.LotFilter.Set(nil)

	COL_ITEMS_LOT = nil
	view.lot_field.Fill_FromFnQuery(view.update_lot, DB.DB_Select_product_lot_product, product_id)

	COL_ITEMS_SAMPLE_PT = nil
	GUI.Fill_combobox_from_query_fn(view.sample_field, view.update_sample, DB.DB_Select_product_sample_points, product_id)

	view.SetTable()
	view.mainWindow.UpdateFilterListView()
}

func (view *DataViewerPanelView) product_clear_button_OnClick(e *windigo.Event) {

	view.clear_product()

	view.SetTable()
	view.mainWindow.UpdateFilterListView()

}

func (view *DataViewerPanelView) lot_field_OnChange(e *windigo.Event) {
	view.LotFilter.Set([]string{view.lot_field.GetSelectedItem()})
	view.SetTable()
}

func (view *DataViewerPanelView) lot_clear_button_OnClick(e *windigo.Event) {
	view.clear_lot()
	view.SetTable()
}

func (view *DataViewerPanelView) lot_add_button_OnClick(e *windigo.Event) {
	view.mainWindow.ShowFilterListView()
	view.mainWindow.AddItem(COL_KEY_LOT, view.lot_field.GetSelectedItem())
}

func (view *DataViewerPanelView) filter_button_OnClick(e *windigo.Event) {
	view.mainWindow.ToggleFilterListView()
}

func (view *DataViewerPanelView) search_button_OnClick(e *windigo.Event) {
	view.mainWindow.SetTable(select_samples(view.mainWindow.FilterListView.Get().Get()))
}

func (view *DataViewerPanelView) reprint_sample_button_OnClick(e *windigo.Event) {
	for _, data := range view.mainWindow.GetTable() {
		data.(QCData).Product().Reprint_sample()
	}

}

func (view *DataViewerPanelView) regen_sample_button_OnClick(e *windigo.Event) {
	for _, data := range view.mainWindow.GetTable() {
		if err := data.(QCData).Product().Output_sample(); err != nil {
			log.Printf("Error: [%s]: %q\n", "regen_sample_button", err)
			log.Printf("Debug: %q: %v\n", err, data)
			threads.Show_status("Error Creating Label")
		}
	}

}

func (view *DataViewerPanelView) export_json_button_OnClick(e *windigo.Event) {

	for _, data := range view.mainWindow.GetTable() {
		data.(QCData).Product().Export_json()
		log.Printf("Info: Exporting json: %v\n", data)
	}

}
