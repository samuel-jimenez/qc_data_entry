package product

import (
	"database/sql"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
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

func (product_type Discrete) Index() int {
	return int(product_type.Int32 - 1)
}

func DiscreteFromInt(Int int) Discrete {
	return Discrete{sql.NullInt32{Int32: int32(Int), Valid: true}}
}

/*
 * DiscreteView
 *
 */
type DiscreteView struct {
	GUI.View
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

func (data_view *DiscreteView) Update(field_data Discrete) {
	if data_view.selected.Valid {
		data_view.buttons[data_view.selected.Index()].SetChecked(false)
	}

	data_view.selected = field_data

	if data_view.selected.Valid {
		data_view.buttons[data_view.selected.Index()].SetChecked(true)
	}
}

func (control *DiscreteView) OnSelectedChange() *windigo.EventManager {
	return &control.onSelectedChange
}

func (control DiscreteView) SetFont(font *windigo.Font) {
	control.View.ComponentFrame.(*windigo.AutoPanel).SetFont(font)
	for _, button := range control.buttons {
		button.SetFont(font)
	}
}

func (control DiscreteView) SetItemSize(width int) {
	for _, button := range control.buttons {
		button.SetSize(width, GUI.OFF_AXIS)
	}
}

func BuildNewDiscreteView(parent windigo.Controller, group_text string, field_data Discrete, labels []string) *DiscreteView {
	view := new(DiscreteView)

	panel := windigo.NewGroupAutoPanel(parent)
	panel.SetSize(GUI.DISCRETE_FIELD_WIDTH, GUI.DISCRETE_FIELD_HEIGHT)

	panel.SetText(group_text)
	panel.SetPaddingsAll(15)

	for i, label_text := range labels {
		label := windigo.NewRadioButton(panel)
		label.SetMarginLeft(10)
		label.SetText(label_text)

		label.OnClick().Bind(func(e *windigo.Event) {
			view.Set(i)
		})
		view.buttons = append(view.buttons, label)
		panel.Dock(label, windigo.Left)
	}
	// view.View.ComponentFrame = View{panel}
	view.View.ComponentFrame = panel

	view.Update(field_data)
	view.Ok()

	return view
}
