package views

import (
	"database/sql"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/*
 * ProductAppearanceViewer
 *
 */
type ProductAppearanceViewer interface {
	*windigo.LabeledEdit
	Get() product.ProductAppearance
	Set(field_data product.ProductAppearance)
}

/*
 * ProductAppearanceView
 *
 */
type ProductAppearanceView struct {
	*windigo.LabeledEdit
}

// func ProductAppearanceView_from_new(parent windigo.Controller,  label_width, control_width, height int, field_text string, field_data product.ProductAppearance) ProductAppearanceView {
func ProductAppearanceView_from_new(parent windigo.Controller, field_text string, field_data product.ProductAppearance) *ProductAppearanceView {

	view := new(ProductAppearanceView)
	view.LabeledEdit = windigo.LabeledEdit_with_size_from_new(parent, GUI.LABEL_WIDTH, GUI.OFF_AXIS, GUI.RANGES_FIELD_HEIGHT, field_text)
	view.LabeledEdit.SetPaddingsAll(GUI.RANGES_PADDING)
	if field_data.Valid {
		view.LabeledEdit.SetText(field_data.String)
	}
	view.LabeledEdit.OnKillFocus().Bind(func(e *windigo.Event) {
		view.Edit.SetText(strings.TrimSpace(view.LabeledEdit.Text()))
	})

	return view
}

func (view *ProductAppearanceView) Get() product.ProductAppearance {
	field_text := view.LabeledEdit.Text()
	Valid := field_text != ""
	return product.ProductAppearance{sql.NullString{String: field_text, Valid: Valid}}
}

func (view *ProductAppearanceView) Set(field_data product.ProductAppearance) {
	if field_data.Valid {
		view.LabeledEdit.SetText(field_data.String)
	} else {
		view.LabeledEdit.SetText("")
	}
}
