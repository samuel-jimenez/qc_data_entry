package views

import (
	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

var (
	SMOL_BOTTLE_SIZE = "8OZ"
	SMOL_BOTTLE_CAP  = 24
	BIG_BOTTLE_SIZE  = "500ML"
	BIG_BOTTLE_CAP   = 20
)

/*
 * MonikerView
 *
 */
type MonikerView struct {
	*windigo.AutoPanel

	product_moniker_name_field    *windigo.LabeledEdit
	retain_storage_duration_field *NumbestEditView
	max_storage_capacity_field    *GUI.ComboBox

	max_storage_capacity int
}

func MonikerView_from_new(parent windigo.Controller) *MonikerView {

	product_moniker_name_text := "Product Moniker Name"
	retain_storage_duration_text := "Retain Storage Duration"

	// max_storage_capacity_text := "Max Storage Capacity"
	max_storage_capacity_text := "Sample Size"

	view := new(MonikerView)

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.product_moniker_name_field = windigo.NewLabeledEdit(view, product_moniker_name_text)

	view.retain_storage_duration_field = NumbestEditView_from_new(view, retain_storage_duration_text)
	view.max_storage_capacity_field = GUI.List_ComboBox_from_new(view, max_storage_capacity_text)

	max_storage_capacity_accept_button := windigo.NewPushButton(view)
	max_storage_capacity_accept_button.SetText("OK")
	// max_storage_capacity_accept_button.SetSize(GUI.ACCEPT_BUTTON_WIDTH, GUI.OFF_AXIS)
	max_storage_capacity_accept_button.SetSize(200, 200)

	// Dock
	view.Dock(view.product_moniker_name_field, windigo.Left)
	view.Dock(view.retain_storage_duration_field, windigo.Left)
	view.Dock(view.max_storage_capacity_field, windigo.Left)
	view.Dock(max_storage_capacity_accept_button, windigo.Left)

	// combobox
	for _, name := range []string{"BSFR", "BSWB", "BSWH"} {
		view.max_storage_capacity_field.AddItem(name)
	}

	//event handling
	view.max_storage_capacity_field.OnSelectedChange().Bind(view.max_storage_capacity_field_OnSelectedChange)
	max_storage_capacity_accept_button.OnClick().Bind(view.max_storage_capacity_accept_button_OnClick)

	return view
}

// add new moniker
func (view *MonikerView) AddMoniker() {
	proc_name := "MonikerView-AddMoniker"
	product_moniker_name := view.product_moniker_name_field.Text()
	retain_storage_duration := view.retain_storage_duration_field.Get()
	max_storage_capacity := view.max_storage_capacity

	// [compiler] (IncompatibleAssign) cannot use qc_sample_storage_offset (variable of type int) as int64 value in argument to measured_product.InsertStorageBin
	// yknow, this, combined with the whole yOU HaVe UNuSed VaRIBLEs thing makes prototyping really annoying
	// qc_sample_storage_offset.(int64) := 0
	// qc_sample_storage_offset := (int64) 0
	// qc_sample_storage_offset := 0
	// qc_sample_storage_offset int64:= 0
	// qc_sample_storage_offset := 0 int64
	// qc_sample_storage_offset := int64(0)
	var qc_sample_storage_offset int64 // lol it's 0.

	qc_storage_capacity := max_storage_capacity
	measured_product := product.MeasuredProduct_from_new()
	//
	product_moniker_id := DB.Insert(proc_name,
		DB.DB_Insert_product_moniker,
		product_moniker_name,
	)

	// product_moniker_id == product_sample_storage_id. update if this is no longer true
	product_sample_storage_id := product_moniker_id

	qc_sample_storage_id, qc_sample_storage_offset := measured_product.InsertStorageBin(product_sample_storage_id, qc_sample_storage_offset, product_moniker_name)

	DB.Insert(proc_name,
		DB.DB_Insert_product_sample_storage,
		product_sample_storage_id, product_moniker_id, retain_storage_duration, max_storage_capacity, qc_sample_storage_id, qc_sample_storage_offset, qc_storage_capacity)
}

func (view *MonikerView) max_storage_capacity_field_OnSelectedChange(e *windigo.Event) {
	switch view.max_storage_capacity_field.GetSelectedItem() {
	case SMOL_BOTTLE_SIZE:
		view.max_storage_capacity = SMOL_BOTTLE_CAP
	case BIG_BOTTLE_SIZE:
		view.max_storage_capacity = BIG_BOTTLE_CAP
	}
}

func (view *MonikerView) max_storage_capacity_accept_button_OnClick(e *windigo.Event) {
	if view.max_storage_capacity == 0 {
		// /?TODO display ERROR
		return
	}
	view.AddMoniker()
}
