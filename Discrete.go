package main

import (
	"database/sql"

	"github.com/samuel-jimenez/windigo"
)

/*
 * Discrete
 *
 */

type Discrete struct {
	sql.NullInt32
}

func DiscreteDefault() Discrete {
	return Discrete{}
}

func DiscreteFromIndex(index int) Discrete {
	return Discrete{sql.NullInt32{Int32: int32(index) + 1, Valid: true}}
}

func (product_type Discrete) toIndex() int {
	return int(product_type.Int32 - 1)

}

/*
 * DiscreteView
 *
 */

type DiscreteView struct {
	windigo.AutoPanel
	buttons          []*windigo.RadioButton
	selected         Discrete
	onSelectedChange windigo.EventManager
}

func (data_view DiscreteView) Get() Discrete {
	return data_view.selected
}

func (data_view *DiscreteView) Set(index int) {
	data_view.selected = DiscreteFromIndex(index)
	data_view.onSelectedChange.Fire(nil)
}

func (view *DiscreteView) Ok() {
	view.AutoPanel.SetBorder(okPen)
}

func (view *DiscreteView) Error() {
	view.AutoPanel.SetBorder(erroredPen)
}

func (data_view *DiscreteView) Update(field_data Discrete) {
	if data_view.selected.Valid {

		data_view.buttons[data_view.selected.toIndex()].SetChecked(false)
	}

	data_view.selected = field_data

	if data_view.selected.Valid {
		data_view.buttons[data_view.selected.toIndex()].SetChecked(true)
	}
}

func (control *DiscreteView) OnSelectedChange() *windigo.EventManager {
	return &control.onSelectedChange
}

func BuildNewDiscreteView(parent windigo.Controller, width, height int, group_text string, field_data Discrete, labels []string) *DiscreteView {
	view := new(DiscreteView)

	view.AutoPanel = windigo.NewGroupAutoPanel(parent)
	view.AutoPanel.SetSize(OFF_AXIS, height)
	view.AutoPanel.SetText(group_text)
	view.AutoPanel.SetPaddingsAll(15)

	for i, label_text := range labels {
		label := windigo.NewRadioButton(view.AutoPanel)
		label.SetSize(width, OFF_AXIS)
		label.SetMarginLeft(10)
		label.SetText(label_text)
		label.OnClick().Bind(func(e *windigo.Event) {
			view.Set(i)
		})
		view.buttons = append(view.buttons, label)
		view.AutoPanel.Dock(label, windigo.Left)
	}

	view.Update(field_data)
	view.Ok()

	return view
}

func BuildNewProductTypeView(parent windigo.Controller, group_text string, field_data Discrete, labels []string) *DiscreteView {
	return BuildNewDiscreteView(parent, PRODUCT_TYPE_WIDTH, RANGES_FIELD_HEIGHT, group_text, field_data, labels)
}
