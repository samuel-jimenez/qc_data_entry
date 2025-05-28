package main

import (
	"database/sql"

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
 * ProductAppearanceView
 *
 */
type ProductAppearanceView struct {
	windigo.LabeledEdit
	Get func() ProductAppearance
}

// func BuildNewProductAppearanceView(parent windigo.Controller,  label_width, control_width, height int, field_text string, field_data ProductAppearance) ProductAppearanceView {
func BuildNewProductAppearanceView(parent windigo.Controller, field_text string, field_data ProductAppearance) ProductAppearanceView {

	field := windigo.NewLabeledEdit(parent, LABEL_WIDTH, OFF_AXIS, RANGES_FIELD_HEIGHT, field_text)
	field.SetPaddingsAll(RANGES_PADDING)
	if field_data.Valid {
		field.SetText(field_data.String)
	}

	get := func() ProductAppearance {
		field_text := field.Text()
		field_valid := field_text != ""
		return ProductAppearance{sql.NullString{String: field_text, Valid: field_valid}}
		// return ProductAppearance{sql.NullString{String: field.Text(), Valid: true}}
	}

	return ProductAppearanceView{field, get}
}

/*
 * ProductAppearanceROView
 *
 */

type ProductAppearanceROView struct {
	windigo.LabeledLabel
	Update func(field_data ProductAppearance)
}

func (view *ProductAppearanceROView) Ok() {
	view.SetBorder(nil)
}

func (view *ProductAppearanceROView) Error() {
	view.SetBorder(erroredPen)
}

func BuildNewProductAppearanceROView(parent windigo.Controller, field_text string, field_data ProductAppearance) ProductAppearanceROView {
	data_field := windigo.NewLabeledLabel(parent, 40, 22, "")
	data_field.SetPaddingsAll(RANGES_PADDING)
	//TODO toolti[p]
	// label := windigo.NewLabel(panel)
	// label.SetText(field_text)
	update := func(field_data ProductAppearance) {
		if field_data.Valid {
			data_field.SetText(field_data.String)
		} else {
			data_field.SetText("")
		}
	}
	update(field_data)

	return ProductAppearanceROView{data_field, update}
}
