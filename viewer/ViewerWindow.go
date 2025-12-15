package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/product"
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

	newMonikerMenu_OnClick(*windigo.Event)
	reprint_sample_button_OnClick(*windigo.Event)
	regen_sample_button_OnClick(*windigo.Event)
	reprint_label_button_OnClick(*windigo.Event)
	regen_label_button_OnClick(*windigo.Event)
	reprint_bin_button_OnClick(*windigo.Event)
	regen_bin_button_OnClick(*windigo.Event)
}

/*
 * ViewerWindow
 *
 */
type ViewerWindow struct {
	*windigo.AutoForm

	selection_panel *DataViewerPanelView
	FilterListView  *SQLFilterListView

	table *QCDataView
}

func NewViewerWindow(parent windigo.Controller) *ViewerWindow {
	window_title := "QC Data Viewer"

	view := new(ViewerWindow)
	view.AutoForm = windigo.AutoForm_from_new(parent)
	view.SetSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	view.SetText(window_title)

	menu := view.NewMenu()
	// TODO settings

	// File
	fileMenu := menu.AddSubMenu("&File")
	newMenu := fileMenu.AddSubMenu_Shortcut("&New", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyN,
	})
	newMonikerMenu := newMenu.AddItem("&Moniker", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyM,
	})
	quitMenu := fileMenu.AddItem("&Quit", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyQ,
	})
	refMenu := fileMenu.AddItem("&Refresh", windigo.Shortcut{
		Key: windigo.KeyF5,
	})

	// Label
	labelMenu := menu.AddSubMenu("&Label")

	binReprintMenu := labelMenu.AddItem("Reprint &Bin", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyB,
	})
	labelReprintMenu := labelMenu.AddItem("Reprint &Label", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyL,
	})
	sampleReprintMenu := labelMenu.AddItem("Reprint Customer &Sample", windigo.Shortcut{
		Modifiers: windigo.ModControl,
		Key:       windigo.KeyS,
	})
	binRegenMenu := labelMenu.AddItem("Regen Bin", windigo.Shortcut{
		Modifiers: windigo.ModControl | windigo.ModShift,
		Key:       windigo.KeyB,
	})
	labelRegenMenu := labelMenu.AddItem("Regen Label", windigo.Shortcut{
		Modifiers: windigo.ModControl | windigo.ModShift,
		Key:       windigo.KeyL,
	})
	sampleRegenMenu := labelMenu.AddItem("Regen Customer Sample", windigo.Shortcut{
		Modifiers: windigo.ModControl | windigo.ModShift,
		Key:       windigo.KeyS,
	})

	newMonikerMenu.OnClick().Bind(view.newMonikerMenu_OnClick)
	quitMenu.OnClick().Bind(toplevel_ui.WndOnClose)
	refMenu.OnClick().Bind(view.search_button_OnClick)

	labelReprintMenu.OnClick().Bind(view.reprint_label_button_OnClick)
	labelRegenMenu.OnClick().Bind(view.regen_label_button_OnClick)
	sampleReprintMenu.OnClick().Bind(view.reprint_sample_button_OnClick)
	sampleRegenMenu.OnClick().Bind(view.regen_sample_button_OnClick)
	binReprintMenu.OnClick().Bind(view.reprint_bin_button_OnClick)
	binRegenMenu.OnClick().Bind(view.regen_bin_button_OnClick)

	// menu.Show() // TODO why is this sometimes needed? Maybe the SetSize call?

	selection_panel := NewDataViewerPanelView(view)

	FilterListView := NewSQLFilterListView(view)

	// product_panel.Dock(reprint_button, windigo.Left)

	// tabs := windigo.TabView_from_new(mainWindow)
	// tab_wb := tabs.AddAutoPanel("Water Based")
	// tab_oil := tabs.AddAutoPanel("Oil Based")
	// tab_fr := tabs.AddAutoPanel("Friction Reducer")

	table := NewQCDataView(view)
	table.Set(select_all_samples())
	table.Update()

	threads.Status_bar = windigo.NewStatusBar(view)
	view.SetStatusBar(threads.Status_bar)
	view.Dock(threads.Status_bar, windigo.Bottom)

	view.Dock(selection_panel, windigo.Top)
	view.Dock(FilterListView, windigo.Top)

	view.Dock(table, windigo.Fill)

	FilterListView.AddContinuousTime(COL_KEY_TIME, COL_LABEL_TIME)
	FilterListView.AddDiscreteSearch(COL_KEY_LOT, COL_LABEL_LOT, COL_ITEMS_LOT)
	FilterListView.AddDiscreteMulti(COL_KEY_SAMPLE_PT, COL_LABEL_SAMPLE_PT, COL_ITEMS_SAMPLE_PT)
	FilterListView.AddContinuous(COL_KEY_PH, COL_LABEL_PH)
	FilterListView.AddContinuous(COL_KEY_SG, COL_LABEL_SG)
	// TODO FilterListView.AddContinuous(COL_KEY_DENSITY, COL_LABEL_DENSITY)
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

func (view *ViewerWindow) newMonikerMenu_OnClick(*windigo.Event) {
	MonikerView := views.MonikerView_from_new(view)
	MonikerView.SetModal(false)
	MonikerView.Show()
	MonikerView.RefreshSize()
}

func (view *ViewerWindow) search_button_OnClick(*windigo.Event) {
	view.SetTable(select_samples(view.FilterListView.Get()))
	// TODO there can be only one
}

func (view *ViewerWindow) forall(fn func(QCData)) {
	for _, data := range view.GetTableSelected() {
		fn(data.(QCData))
	}
}

func Reprint_sample(data QCData) {
	data.Product().Reprint_sample()
}

func Reprint_label(data QCData) {
	data.Product().Reprint()
}

func Reprint_bin(data QCData) {
	if !data.Sample_bin.Valid {
		return
	}
	threads.LogError("ViewerWindow-Reprint_bin", "Error Printing Label", product.RePrintStorage(data.Sample_bin.String))
}

func Regen_sample(data QCData) {
	threads.LogError("ViewerWindow-Regen_sample", "Error Creating Label", data.Product().Regen_sample())
}

func Regen_label(data QCData) {
	threads.LogError("ViewerWindow-Regen_label", "Error Creating Label", data.Product().Print())
}

func Regen_bin(data QCData) {
	if !data.Sample_bin.Valid {
		return
	}
	product.RegenStorageBin(data.Sample_bin.String)
}

func (view *ViewerWindow) reprint_sample_button_OnClick(*windigo.Event) {
	view.forall(Reprint_sample)
}

func (view *ViewerWindow) reprint_label_button_OnClick(*windigo.Event) {
	view.forall(Reprint_label)
}

func (view *ViewerWindow) reprint_bin_button_OnClick(*windigo.Event) {
	view.forall(Reprint_bin)
}

func (view *ViewerWindow) regen_sample_button_OnClick(*windigo.Event) {
	view.forall(Regen_sample)
}

func (view *ViewerWindow) regen_label_button_OnClick(*windigo.Event) {
	view.forall(Regen_label)
}

func (view *ViewerWindow) regen_bin_button_OnClick(*windigo.Event) {
	view.forall(Regen_bin)
}
