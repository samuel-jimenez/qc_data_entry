package product

import (
	"database/sql"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * ProductAppearance
 *
 */

type ProductAppearance struct {
	sql.NullString
}

/*
 * ProductAppearanceViewer
 *
 */
type ProductAppearanceViewer interface {
	*windigo.LabeledEdit
	Get() ProductAppearance
	Set(field_data ProductAppearance)
}

/*
 * ProductAppearanceView
 *
 */
type ProductAppearanceView struct {
	*windigo.LabeledEdit
}

// func BuildNewProductAppearanceView(parent windigo.Controller,  label_width, control_width, height int, field_text string, field_data ProductAppearance) ProductAppearanceView {
func BuildNewProductAppearanceView(parent windigo.Controller, field_text string, field_data ProductAppearance) *ProductAppearanceView {

	view := new(ProductAppearanceView)
	view.LabeledEdit = windigo.NewSizedLabeledEdit(parent, GUI.LABEL_WIDTH, GUI.OFF_AXIS, GUI.RANGES_FIELD_HEIGHT, field_text)
	view.LabeledEdit.SetPaddingsAll(GUI.RANGES_PADDING)
	if field_data.Valid {
		view.LabeledEdit.SetText(field_data.String)
	}

	return view
}

func (view *ProductAppearanceView) Get() ProductAppearance {
	field_text := view.LabeledEdit.Text()
	Valid := field_text != ""
	return ProductAppearance{sql.NullString{String: field_text, Valid: Valid}}
}

func (view *ProductAppearanceView) Set(field_data ProductAppearance) {
	if field_data.Valid {
		view.LabeledEdit.SetText(field_data.String)
	} else {
		view.LabeledEdit.SetText("")
	}
}

type ProductAppearanceROViewer interface {
	*GUI.View
	windigo.BaseController
	Update(field_data ProductAppearance)
	SetFont(font *windigo.Font)
	RefreshSize()
}

/*
 * ProductAppearanceROView
 *
 */

type ProductAppearanceROView struct {
	*GUI.View
	data_field *windigo.LabeledLabel
}

func BuildNewProductAppearanceROView(parent windigo.Controller, field_text string, field_data ProductAppearance) *ProductAppearanceROView {
	data_field := windigo.NewLabeledLabel(parent, "")

	//TODO toolti[p]
	// label := windigo.NewLabel(panel)
	// label.SetText(field_text)

	return &ProductAppearanceROView{&GUI.View{ComponentFrame: data_field}, data_field}
}

func (view *ProductAppearanceROView) Update(field_data ProductAppearance) {
	if field_data.Valid {
		view.data_field.SetText(field_data.String)
	} else {
		view.data_field.SetText("")
	}
}
func (view *ProductAppearanceROView) SetFont(font *windigo.Font) { view.data_field.SetFont(font) }
func (view *ProductAppearanceROView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.RANGES_RO_FIELD_HEIGHT)
	view.SetPaddingsAll(GUI.ERROR_MARGIN)
}
