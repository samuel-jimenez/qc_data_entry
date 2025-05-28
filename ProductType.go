package main

import (
	"database/sql"

	"github.com/samuel-jimenez/windigo"
)

/*
 * ProductType
 *
 */

type ProductType struct {
	sql.NullInt32
}

func ProductTypeDefault() ProductType {
	return ProductType{}
}

func ProductTypeFromIndex(index int) ProductType {
	return ProductType{sql.NullInt32{Int32: int32(index) + 1, Valid: true}}
}

func (product_type ProductType) toIndex() int {
	return int(product_type.Int32 - 1)

}

/*
 * ProductTypeView
 *
 */

type ProductTypeView struct {
	windigo.AutoPanel
	Get   func() ProductType
	Ok    func()
	Error func()
}

func BuildNewProductTypeView(parent windigo.Controller, group_text string, field_data ProductType, labels []string) ProductTypeView {
	var buttons []*windigo.RadioButton
	panel := windigo.NewGroupAutoPanel(parent)
	panel.SetSize(OFF_AXIS, RANGES_FIELD_HEIGHT)
	panel.SetText(group_text)
	panel.SetPaddingsAll(15)

	for _, label_text := range labels {
		label := windigo.NewRadioButton(panel)
		label.SetSize(150, OFF_AXIS)
		label.SetMarginLeft(10)
		label.SetText(label_text)
		buttons = append(buttons, label)
		panel.Dock(label, windigo.Left)

	}

	if field_data.Valid {

		buttons[field_data.toIndex()].SetChecked(true)
	}

	// Get()
	get := func() ProductType {
		for i, button := range buttons {
			if button.Checked() {
				// one-based
				return ProductTypeFromIndex(i)
			}
		}
		return ProductTypeDefault()
	}
	// Ok()
	ok := func() {
		panel.SetBorder(okPen)
	}

	// Error()
	err := func() {
		panel.SetBorder(erroredPen)
	}

	ok()

	return ProductTypeView{panel, get, ok, err}
}
