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
	*windigo.LabeledEdit
	Get func() ProductAppearance
	Set func(field_data ProductAppearance)
}

// func BuildNewProductAppearanceView(parent windigo.Controller,  label_width, control_width, height int, field_text string, field_data ProductAppearance) ProductAppearanceView {
func BuildNewProductAppearanceView(parent windigo.Controller, field_text string, field_data ProductAppearance) ProductAppearanceView {

	field := windigo.NewSizedLabeledEdit(parent, LABEL_WIDTH, OFF_AXIS, RANGES_FIELD_HEIGHT, field_text)
	field.SetPaddingsAll(RANGES_PADDING)
	if field_data.Valid {
		field.SetText(field_data.String)
	}

	get := func() ProductAppearance {
		field_text := field.Text()
		field_valid := field_text != ""
		return ProductAppearance{sql.NullString{String: field_text, Valid: field_valid}}
	}
	update := func(field_data ProductAppearance) {
		if field_data.Valid {
			field.SetText(field_data.String)
		} else {
			field.SetText("")
		}
	}

	return ProductAppearanceView{field, get, update}
}

/*
 * ProductAppearanceROView
 *
 */

type ProductAppearanceROView struct {
	*View
	Update  func(field_data ProductAppearance)
	SetFont func(*windigo.Font)
	Refresh func()
}

func BuildNewProductAppearanceROView(parent windigo.Controller, field_text string, field_data ProductAppearance) ProductAppearanceROView {
	data_field := windigo.NewLabeledLabel(parent, "")

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

	refresh := func() {
		data_field.SetSize(OFF_AXIS, RANGES_RO_FIELD_HEIGHT)
		data_field.SetPaddingsAll(ERROR_MARGIN)
	}

	return ProductAppearanceROView{&View{data_field}, update, data_field.SetFont, refresh}
}
