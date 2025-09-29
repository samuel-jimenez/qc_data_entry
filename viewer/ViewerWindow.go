package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
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

	mainWindow := new(ViewerWindow)
	mainWindow.Form = windigo.NewForm(parent)
	mainWindow.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	mainWindow.SetText(window_title)

	dock := windigo.NewSimpleDock(mainWindow)

	selection_panel := NewDataViewerPanelView(mainWindow)
	selection_panel.SetMainWindow(mainWindow)

	FilterListView := NewSQLFilterListView(mainWindow)

	// product_panel.Dock(reprint_button, windigo.Left)

	// tabs := windigo.NewTabView(mainWindow)
	// tab_wb := tabs.AddAutoPanel("Water Based")
	// tab_oil := tabs.AddAutoPanel("Oil Based")
	// tab_fr := tabs.AddAutoPanel("Friction Reducer")

	table := NewQCDataView(mainWindow)
	table.Set(select_all_samples())
	table.Update()

	threads.Status_bar = windigo.NewStatusBar(mainWindow)
	mainWindow.SetStatusBar(threads.Status_bar)
	dock.Dock(threads.Status_bar, windigo.Bottom)

	dock.Dock(selection_panel, windigo.Top)
	dock.Dock(FilterListView, windigo.Top)

	dock.Dock(table, windigo.Fill)

	FilterListView.AddContinuous(COL_KEY_TIME, COL_LABEL_TIME)
	FilterListView.AddDiscreteSearch(COL_KEY_LOT, COL_LABEL_LOT, COL_ITEMS_LOT)
	FilterListView.AddDiscreteMulti(COL_KEY_SAMPLE_PT, COL_LABEL_SAMPLE_PT, COL_ITEMS_SAMPLE_PT)
	FilterListView.AddContinuous(COL_KEY_PH, COL_LABEL_PH)
	FilterListView.AddContinuous(COL_KEY_SG, COL_LABEL_SG)
	//TODO FilterListView.AddContinuous(COL_KEY_DENSITY, COL_LABEL_DENSITY)
	FilterListView.AddContinuous(COL_KEY_STRING, COL_LABEL_STRING)
	FilterListView.AddContinuous(COL_KEY_VISCOSITY, COL_LABEL_VISCOSITY)
	FilterListView.Hide()

	// build object
	mainWindow.selection_panel = selection_panel
	mainWindow.FilterListView = FilterListView
	mainWindow.table = table

	return mainWindow
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
func (view *ViewerWindow) GetTable() []windigo.ListItem {
	return view.table.SelectedItems()
}

func (view *ViewerWindow) SetFont(font *windigo.Font) {
	view.selection_panel.SetFont(font)
	view.FilterListView.SetFont(font)
	view.table.SetFont(font)
}

func (view *ViewerWindow) RefreshSize() {
	Refresh_globals(GUI.BASE_FONT_SIZE)
	view.selection_panel.RefreshSize()
	view.FilterListView.RefreshSize()
	// view.table.RefreshSize()
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
	GUI.BASE_FONT_SIZE += 1
	view.Set_font_size()
	return true
}
func (view *ViewerWindow) Decrease_font_size() bool {
	GUI.BASE_FONT_SIZE -= 1
	view.Set_font_size()
	return true
}
