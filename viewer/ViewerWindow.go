package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"
)

/*
 * ViewerWinder
 *
 */
type ViewerWinder interface {
	clear_lot(product_id int)
	ToggleFilterListView()
	UpdateFilterListView()
	ClearFilters()
	SetTable(samples []QCData)
	GetTable() []windigo.ListItem
	SetFont(font *windigo.Font)
	RefreshSize()
	AddShortcuts()
	Set_font_size()
	Increase_font_size() bool
	Decrease_font_size() bool
}

//TODO add new moniker
//
// insert product_moniker_id :=
//
// INSERT INTO product_moniker
// (product_moniker_name)
//
//insel
// insert into qc_sample_storage_list values(27,'00-018-0000',18);
// Product.NewStorageBin.Insert
//
//
// insert into product_sample_storage values(18,18,12,24,27,0,24);
// offset:=0
// space:=max_cap
// product_moniker_id, duration,max_cap,storagbin_id,offset,space
//TODO end add new moniker

/*
 * ViewerWindow
 *
 */
type ViewerWindow struct {
	*windigo.Form

	selection_panel *DataViewerPanelView
	FilterListView  *SQLFilterListView

	table *QCDataView
}

func NewViewerWindow(parent windigo.Controller) *ViewerWindow {
	window_title := "QC Data Viewer"

	view := new(ViewerWindow)
	view.Form = windigo.NewForm(parent)
	view.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	view.SetText(window_title)

	// menu := view.NewMenu()
	// TODO settings
	// 	menu := mainWindow.NewMenu()
	//
	// 	fileMn := menu.AddSubMenu("File")
	// 	fileMn.AddItem("New", windigo.Shortcut{windigo.ModControl, windigo.KeyN})
	// 	editMn := menu.AddSubMenu("Edit")
	// 	cutMn := editMn.AddItem("Cut", windigo.Shortcut{windigo.ModControl, windigo.KeyX})
	//
	dock := windigo.NewSimpleDock(view)

	selection_panel := NewDataViewerPanelView(view)
	selection_panel.SetMainWindow(view)

	FilterListView := NewSQLFilterListView(view)

	// product_panel.Dock(reprint_button, windigo.Left)

	// tabs := windigo.NewTabView(mainWindow)
	// tab_wb := tabs.AddAutoPanel("Water Based")
	// tab_oil := tabs.AddAutoPanel("Oil Based")
	// tab_fr := tabs.AddAutoPanel("Friction Reducer")

	table := NewQCDataView(view)
	table.Set(select_all_samples())
	table.Update()

	threads.Status_bar = windigo.NewStatusBar(view)
	view.SetStatusBar(threads.Status_bar)
	dock.Dock(threads.Status_bar, windigo.Bottom)

	dock.Dock(selection_panel, windigo.Top)
	dock.Dock(FilterListView, windigo.Top)

	dock.Dock(table, windigo.Fill)

	FilterListView.AddContinuousTime(COL_KEY_TIME, COL_LABEL_TIME)
	FilterListView.AddDiscreteSearch(COL_KEY_LOT, COL_LABEL_LOT, COL_ITEMS_LOT)
	FilterListView.AddDiscreteMulti(COL_KEY_SAMPLE_PT, COL_LABEL_SAMPLE_PT, COL_ITEMS_SAMPLE_PT)
	FilterListView.AddContinuous(COL_KEY_PH, COL_LABEL_PH)
	FilterListView.AddContinuous(COL_KEY_SG, COL_LABEL_SG)
	//TODO FilterListView.AddContinuous(COL_KEY_DENSITY, COL_LABEL_DENSITY)
	FilterListView.AddContinuous(COL_KEY_STRING, COL_LABEL_STRING)
	FilterListView.AddContinuous(COL_KEY_VISCOSITY, COL_LABEL_VISCOSITY)
	FilterListView.Hide()

	// build object
	view.selection_panel = selection_panel
	view.FilterListView = FilterListView
	view.table = table

	return view
}

func (view *ViewerWindow) ToggleFilterListView() {
	if view.FilterListView.Visible() {
		view.FilterListView.Hide()
	} else {
		view.FilterListView.Show()
	}
}

func (view *ViewerWindow) ShowFilterListView() {
	view.FilterListView.Show()
}

func (view *ViewerWindow) UpdateFilterListView() {
	view.FilterListView.Update(COL_KEY_LOT, COL_ITEMS_LOT)
	view.FilterListView.Update(COL_KEY_SAMPLE_PT, COL_ITEMS_SAMPLE_PT)
}
func (view *ViewerWindow) ClearFilters() {
	view.FilterListView.Clear()
}

func (view *ViewerWindow) AddItem(key string, entry string) {
	view.FilterListView.AddItem(key, entry)
}

func (view *ViewerWindow) SetTable(samples []QCData) {
	view.table.Set(samples)
	view.table.Update()
}
func (view *ViewerWindow) GetTableSelected() []windigo.ListItem {
	return view.table.SelectedItems()
}
func (view *ViewerWindow) GetTableAllSelected() (entries []windigo.ListItem) {
	entries = view.table.SelectedItems()
	if len(entries) == 0 {
		for _, entry := range view.table.Get() {
			entries = append(entries, entry)
		}
	}
	return
}

func (view *ViewerWindow) SetFont(font *windigo.Font) {
	view.selection_panel.SetFont(font)
	view.FilterListView.SetFont(font)
	view.table.SetFont(font)
}

func (view *ViewerWindow) RefreshSize() {
	Refresh_globals(config.BASE_FONT_SIZE)
	view.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	view.selection_panel.RefreshSize()
	view.FilterListView.RefreshSize()
	view.table.RefreshSize()
}

func (view *ViewerWindow) AddShortcuts() {
	toplevel_ui.AddShortcuts(view)
}

func (mainWindow *ViewerWindow) Set_font_size() {
	toplevel_ui.Set_font_size()
	mainWindow.SetFont(windigo.DefaultFont)
	mainWindow.RefreshSize()
}
func (view *ViewerWindow) Increase_font_size() bool {
	config.BASE_FONT_SIZE += 1
	view.Set_font_size()
	return true
}
func (view *ViewerWindow) Decrease_font_size() bool {
	config.BASE_FONT_SIZE -= 1
	view.Set_font_size()
	return true
}
