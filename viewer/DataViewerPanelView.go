package viewer

import (
	"database/sql"
	"log"

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

	clear_product(e *windigo.Event)
	clear_lot(e *windigo.Event)
	ClearFilters(e *windigo.Event)

	update_lot(id int, name string)
	update_sample(id int, name string)

	// listeners
	product_field_OnChange(e *windigo.Event)
	lot_field_OnChange(e *windigo.Event)
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

	product_data map[string]int
	lot_data     []string

	SQLFilters *SQLFilterList

	LotFilter *SQLFilterDiscrete

	lot_field *GUI.SearchBox

	prod_panel, lot_panel *windigo.AutoPanel

	product_field, moniker_field,
	sample_field *GUI.ComboBox

	product_clear_button, lot_clear_button,
	filter_button, search_button, clear_button,
	reprint_sample_button, regen_sample_button, export_json_button *windigo.PushButton
}

// parent

func NewDataViewerPanelView(mainWindow windigo.Controller) *DataViewerPanelView {

	view := new(DataViewerPanelView)

	view.product_data = make(map[string]int)
	// lot_data := make(map[string]int)
	// var lot_data []string

	// LotFilter := &SQLFilterDiscrete{COL_KEY_LOT,nil}
	view.SQLFilters = NewSQLFilterList()

	view.LotFilter = NewSQLFilterDiscrete(
		COL_KEY_LOT,
		nil,
	)

	product_text := "Product"
	lot_text := "Lot Number"
	sample_text := "Sample Point"
	moniker_text := "Product Moniker"

	clear_label := "Clear"
	export_json_label := "Export JSON"

	product_panel := windigo.NewAutoPanel(mainWindow)

	//TODO array db_select_all_product

	prod_panel := windigo.NewAutoPanel(product_panel)

	product_field := GUI.NewListComboBox(prod_panel, product_text)

	product_clear_button := windigo.NewPushButton(prod_panel)
	product_clear_button.SetText("-")

	moniker_field := GUI.NewListComboBox(prod_panel, moniker_text)

	prod_panel.Dock(product_field, windigo.Left)
	prod_panel.Dock(product_clear_button, windigo.Left)
	prod_panel.Dock(moniker_field, windigo.Left)

	lot_panel := windigo.NewAutoPanel(product_panel)

	lot_field := GUI.NewLabeledListSearchBox(lot_panel, lot_text)
	lot_clear_button := windigo.NewPushButton(lot_panel)
	lot_clear_button.SetText("-")
	sample_field := GUI.NewListComboBox(lot_panel, sample_text)

	lot_panel.Dock(lot_field, windigo.Left)
	lot_panel.Dock(lot_clear_button, windigo.Left)
	lot_panel.Dock(sample_field, windigo.Left)

	filter_button := windigo.NewPushButton(product_panel)
	filter_button.SetText("Filter")

	search_button := windigo.NewPushButton(product_panel)
	search_button.SetText("Search")

	clear_button := windigo.NewPushButton(product_panel)
	clear_button.SetText(clear_label)

	reprint_sample_button := windigo.NewPushButton(product_panel)
	reprint_sample_button.SetText("Reprint Sample")

	regen_sample_button := windigo.NewPushButton(product_panel)
	regen_sample_button.SetText("Regen Sample")

	export_json_button := windigo.NewPushButton(product_panel)
	export_json_button.SetText(export_json_label)

	product_panel.Dock(prod_panel, windigo.Top)
	product_panel.Dock(lot_panel, windigo.Top)
	product_panel.Dock(filter_button, windigo.Left)
	product_panel.Dock(search_button, windigo.Left)
	product_panel.Dock(clear_button, windigo.Left)

	product_panel.Dock(reprint_sample_button, windigo.Left)
	product_panel.Dock(regen_sample_button, windigo.Left)
	product_panel.Dock(export_json_button, windigo.Left)

	view.AutoPanel = product_panel
	view.prod_panel = prod_panel
	view.lot_panel = lot_panel

	view.product_field = product_field
	view.lot_field = lot_field
	view.moniker_field = moniker_field
	view.sample_field = sample_field

	view.product_clear_button = product_clear_button
	view.lot_clear_button = lot_clear_button

	view.filter_button = filter_button
	view.search_button = search_button
	view.clear_button = clear_button

	view.reprint_sample_button = reprint_sample_button
	view.regen_sample_button = regen_sample_button

	view.export_json_button = export_json_button

	// functionality

	// clear := func() {
	// 	clear_product()
	// 	clear_lot()
	// }

	// combobox
	GUI.Fill_combobox_from_query_rows(product_field, func(row *sql.Rows) error {
		var (
			id                   int
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		view.product_data[name] = id

		product_field.AddItem(name)
		return nil
	},
		DB.DB_Select_product_info)

	lot_field.Fill_FromFnQuery(func(id int, name string) {
		// lot_data[name] = id
		view.lot_data = append(view.lot_data, name)
		view.update_lot(id, name)
	}, DB.DB_Select_product_lot_all)

	GUI.Fill_combobox_from_query(moniker_field, DB.DB_Select_all_product_moniker)
	GUI.Fill_combobox_from_query_fn(sample_field, view.update_sample, DB.DB_Select_all_sample_points)

	// listeners
	product_field.OnSelectedChange().Bind(view.product_field_OnChange)
	lot_field.OnSelectedChange().Bind(view.lot_field_OnChange)
	moniker_field.OnSelectedChange().Bind(view.moniker_field_OnChange)

	filter_button.OnClick().Bind(view.filter_button_OnClick)
	search_button.OnClick().Bind(view.search_button_OnClick)
	product_clear_button.OnClick().Bind(view.clear_product)
	lot_clear_button.OnClick().Bind(view.clear_lot)
	clear_button.OnClick().Bind(view.ClearFilters)

	reprint_sample_button.OnClick().Bind(view.reprint_sample_button_OnClick)
	regen_sample_button.OnClick().Bind(view.regen_sample_button_OnClick)
	export_json_button.OnClick().Bind(view.export_json_button_OnClick)

	return view
}

func (view *DataViewerPanelView) SetFont(font *windigo.Font) {

	view.product_field.SetFont(font)
	view.moniker_field.SetFont(font)
	view.lot_field.SetFont(font)
	view.sample_field.SetFont(font)

	view.product_clear_button.SetFont(font)
	view.lot_clear_button.SetFont(font)

	// view.AutoPanel.SetFont(font)
	// view.lot_panel.SetFont(font)
	// view.prod_panel.SetFont(font)

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
	view.prod_panel.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.lot_panel.SetSize(GUI.HPANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.product_clear_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.moniker_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.lot_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.lot_clear_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.sample_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	view.lot_panel.SetMarginTop(GUI.INTER_SPACER_HEIGHT)
	view.lot_panel.SetMarginLeft(GUI.HPANEL_MARGIN)
	view.sample_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)

	view.moniker_field.SetMarginLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.prod_panel.SetMarginLeft(GUI.HPANEL_MARGIN)
	view.prod_panel.SetMarginTop(GUI.TOP_SPACER_HEIGHT)

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

// TODO combine all this (CLEAR, OnChange) into filter
// conditional serch, refine (browes)
// TODO ADD BUTTON TO MOVE TO LOWER FILTERS
func (view *DataViewerPanelView) clear_product(e *windigo.Event) {

	view.mainWindow.SetTable(select_samples())

	view.product_field.SetSelectedItem(-1)

	COL_ITEMS_LOT = nil
	view.lot_field.Update(nil)

	for id, name := range view.lot_data {
		view.update_lot(id, name)
	}

	COL_ITEMS_SAMPLE_PT = nil
	GUI.Fill_combobox_from_query_fn(view.sample_field, view.update_sample, DB.DB_Select_all_sample_points)
	view.mainWindow.UpdateFilterListView()

}

func (view *DataViewerPanelView) clear_lot(e *windigo.Event) {

	product_id := view.product_data[view.product_field.GetSelectedItem()]
	view.mainWindow.clear_lot(product_id)
	view.lot_field.SetSelectedItem(-1)
}

func (view *DataViewerPanelView) clear_moniker(e *windigo.Event) {

	view.moniker_field.SetSelectedItem(-1)
}

func (view *DataViewerPanelView) ClearFilters(e *windigo.Event) {
	view.mainWindow.ClearFilters()
}

func (view *DataViewerPanelView) update_lot(id int, name string) {
	view.lot_field.AddEntry(name)
	COL_ITEMS_LOT = append(COL_ITEMS_LOT, name)
}
func (view *DataViewerPanelView) update_sample(id int, name string) {
	view.sample_field.AddItem(name)
	COL_ITEMS_SAMPLE_PT = append(COL_ITEMS_SAMPLE_PT, name)
}

// listeners
func (view *DataViewerPanelView) product_field_OnChange(e *windigo.Event) {
	// product_field_pop_data(product_field.GetSelectedItem())

	product_id := view.product_data[view.product_field.GetSelectedItem()]
	samples := select_product_samples(product_id)
	view.mainWindow.SetTable(samples)

	COL_ITEMS_LOT = nil
	view.lot_field.Fill_FromFnQuery(view.update_lot, DB.DB_Select_product_lot_product, product_id)

	COL_ITEMS_SAMPLE_PT = nil
	GUI.Fill_combobox_from_query_fn(view.sample_field, view.update_sample, DB.DB_Select_product_sample_points, product_id)

	view.mainWindow.UpdateFilterListView()

	view.moniker_field.SetSelectedItem(-1)

}

func (view *DataViewerPanelView) lot_field_OnChange(e *windigo.Event) {
	view.LotFilter.Set([]string{view.lot_field.GetSelectedItem()})
	view.SQLFilters.Filters = []SQLFilter{view.LotFilter}
	rows, err := QC_DB.Query(SAMPLE_SELECT_STRING + view.SQLFilters.Get())
	view.mainWindow.SetTable(_select_samples(rows, err, "search samples"))

	view.moniker_field.SetSelectedItem(-1)
}

func (view *DataViewerPanelView) moniker_field_OnChange(e *windigo.Event) {

	view.clear_product(nil)
	// view.product_field.SetSelectedItem(-1)
	// clear_lot
	view.lot_field.SetSelectedItem(-1)

	product_moniker := view.moniker_field.GetSelectedItem()

	samples := select_product_moniker(product_moniker)
	view.mainWindow.SetTable(samples)
}

func (view *DataViewerPanelView) filter_button_OnClick(e *windigo.Event) {
	view.mainWindow.ToggleFilterListView()
}

func (view *DataViewerPanelView) search_button_OnClick(e *windigo.Event) {

	//TODO
	log.Println("FilterListView:", view.mainWindow.FilterListView.Get().Get())
	log.Println("SAMPLE_SELECT_STRING", SAMPLE_SELECT_STRING+view.mainWindow.FilterListView.Get().Get())
	rows, err := QC_DB.Query(SAMPLE_SELECT_STRING + view.mainWindow.FilterListView.Get().Get())
	view.mainWindow.SetTable(_select_samples(rows, err, "search samples"))

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
