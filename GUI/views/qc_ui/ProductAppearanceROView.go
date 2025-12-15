package qc_ui

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

/*
 * ProductAppearanceROViewer
 *
 */
type ProductAppearanceROViewer interface {
	*GUI.View
	windigo.BaseController
	Update(field_data product.ProductAppearance)
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

func ProductAppearanceROView_from_new(parent windigo.Controller, field_text string, field_data product.ProductAppearance) *ProductAppearanceROView {
	data_field := windigo.LabeledLabel_from_new(parent, "")

	//TODO toolti[p]
	// label := windigo.NewLabel(panel)
	// label.SetText(field_text)

	return &ProductAppearanceROView{&GUI.View{ComponentFrame: data_field}, data_field}
}

func (view *ProductAppearanceROView) Update(field_data product.ProductAppearance) {
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
