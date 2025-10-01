package qc

import (
	"database/sql"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/windigo"
)

/*
 * StorageBinable
 *
 */
type StorageBinable interface {
	Save()
}

/*
 * StorageBin
 *
 */
type StorageBin struct {
	Product_sample_storage_id, Max_storage_capacity, QC_storage_capacity int
	QC_sample_storage_name, Product_moniker_name                         string
}

func New_StorageBin() *StorageBin {
	object := new(StorageBin)
	return object
}

func New_StorageBinFromRow(row *sql.Rows) (*StorageBin, error) {
	object := New_StorageBin()
	err := row.Scan(
		&object.Product_sample_storage_id, &object.Product_moniker_name, &object.QC_sample_storage_name, &object.Max_storage_capacity,
		&object.QC_storage_capacity,
	)

	//	if err != nil {
	//		log.Println("Crit: [BlendComponent-NewRecipeComponentfromSQL]: ", err)
	//	}
	return object, err
}

func (object *StorageBin) Save() {
	proc_name := "StorageBin-Save"
	// update capacity
	DB.Update(proc_name,
		DB.DB_Update_product_sample_storage_capacity,
		object.Product_sample_storage_id, object.QC_storage_capacity)
}

/*
 *     StorageBinViewer
 *
 */
type StorageBinViewer interface {
	windigo.Controller
	Save()
	Changedp() bool
	SetFont(font *windigo.Font)
	RefreshSize()
}

/*
 * StorageBinView
 *
 */
type StorageBinView struct {
	*windigo.AutoPanel

	Storage_Bin *StorageBin

	Product_name_field, Bin_name_field *windigo.Label

	Samples_field *GUI.NumbEditView
}

func New_StorageBinView(parent windigo.Controller, StorageBin *StorageBin) *StorageBinView {
	view := new(StorageBinView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.Storage_Bin = StorageBin

	view.Bin_name_field = windigo.NewLabel(view.AutoPanel)
	view.Product_name_field = windigo.NewLabel(view.AutoPanel)
	view.Samples_field = GUI.NewNumbEditView(view.AutoPanel)

	view.AutoPanel.Dock(view.Bin_name_field, windigo.Left)
	view.AutoPanel.Dock(view.Product_name_field, windigo.Left)
	view.AutoPanel.Dock(view.Samples_field, windigo.Left)

	view.Product_name_field.SetText(view.Storage_Bin.Product_moniker_name)
	view.Bin_name_field.SetText(view.Storage_Bin.QC_sample_storage_name)
	view.Samples_field.SetInt(float64(view.Storage_Bin.Max_storage_capacity - view.Storage_Bin.QC_storage_capacity))

	view.RefreshSize()

	return view
}

func (view *StorageBinView) Save() {
	if !(view.Changedp()) {
		return
	}
	view.Storage_Bin.QC_storage_capacity = view.Storage_Bin.Max_storage_capacity - int(view.Samples_field.Get())
	view.Storage_Bin.Save()
}

func (view *StorageBinView) Changedp() bool {
	return view.Storage_Bin.Max_storage_capacity-int(view.Samples_field.Get()) != view.Storage_Bin.QC_storage_capacity
}

func (view *StorageBinView) SetFont(font *windigo.Font) {
	view.Product_name_field.SetFont(font)
	view.Bin_name_field.SetFont(font)
	view.Samples_field.SetFont(font)
}

func (view *StorageBinView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_HEIGHT)
	view.Product_name_field.SetSize(GUI.SOURCES_LABEL_WIDTH, GUI.OFF_AXIS)
	view.Bin_name_field.SetSize(GUI.SOURCES_LABEL_WIDTH, GUI.OFF_AXIS)
	view.Samples_field.SetSize(GUI.SOURCES_LABEL_WIDTH, GUI.OFF_AXIS)
}

/*
 * StorageBinHeader
 *
 */
type StorageBinHeader struct {
	*GUI.TextDock
}

func New_StorageBinHeader(parent windigo.Controller) *StorageBinHeader {
	view := new(StorageBinHeader)
	view.TextDock = GUI.NewTextDock(parent, "Product", "Storage Bin", "Inventory")
	view.RefreshSize()
	return view
}

func (view *StorageBinHeader) RefreshSize() {
	view.SetDockSize(GUI.SOURCES_LABEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
}

/*
 * InventoryViewer
 *
 */
type InventoryViewer interface {
	windigo.Controller
	SetFont(font *windigo.Font)
	RefreshSize()

	Start()
	OnExit(arg *windigo.Event)
	Exit()
	TryExit()
	Save()
}

/*
 * InventoryView
 *
 */
type InventoryView struct {
	*windigo.Form

	StorageBin_Header *StorageBinHeader
	Components        []StorageBinView
	button_dock       *GUI.ButtonDock
}

func New_InventoryView() *InventoryView {
	window_title := "Sample Inventory"

	view := new(InventoryView)
	view.Form = windigo.NewForm(nil)
	view.SetText(window_title)

	dock := windigo.NewSimpleDock(view)
	view.StorageBin_Header = New_StorageBinHeader(view)
	dock.Dock(view.StorageBin_Header, windigo.Top)

	proc_name := "New_StorageBin"
	DB.Forall_err(proc_name,
		func() {},
		func(row *sql.Rows) error {
			sample_bin, err := New_StorageBinFromRow(row)
			if err != nil {
				return err
			}
			bin_view := New_StorageBinView(view, sample_bin)
			view.Components = append(view.Components, *bin_view)
			dock.Dock(bin_view, windigo.Top)
			return nil
		},
		DB.DB_Select_all_product_sample_storage)

	// view.button_dock = GUI.NewMarginalButtonDock(view, SUBMIT_CLEAR_BTN, []int{40, 0}, []func(){view.submit_cb, view.Clear})
	// view.button_dock = GUI.NewButtonDock(view, []string{"OK", "Cancel"}, []func(){view.Save, view.TryExit})
	view.button_dock = GUI.NewButtonDock(view, []string{"OK", "Cancel"}, []func(){view.Save, view.Exit})
	dock.Dock(view.button_dock, windigo.Top)

	return view
}

func (view *InventoryView) SetFont(font *windigo.Font) {
	view.StorageBin_Header.SetFont(font)
	for _, Component := range view.Components {
		Component.SetFont(font)
	}
}

func (view *InventoryView) RefreshSize() {
	view.StorageBin_Header.RefreshSize()
	for _, Component := range view.Components {
		Component.RefreshSize()
	}
}

func (view *InventoryView) Start() {
	view.Center()
	view.Show()
	view.OnClose().Bind(view.OnExit)
	view.RunMainLoop()
}

func (view *InventoryView) OnExit(arg *windigo.Event) {
	view.Exit()
}

func (view *InventoryView) Exit() {
	view.Close()
	windigo.Exit()
}

func (view *InventoryView) TryExit() {
	canExit := true
	// for _, Component := range view.Components && canExit {
	for i := 0; i < len(view.Components) && canExit; i++ {
		Component := view.Components[i]
		canExit = canExit && !(Component.Changedp())
	}
	if canExit {
		view.Exit()
	} else { //popup save? exit cancel dialog

	}
}

func (view *InventoryView) Save() {
	for _, Component := range view.Components {
		Component.Save()
	}
	view.Exit()
}
