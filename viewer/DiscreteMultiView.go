package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * DiscreteMultiView
 *
 */
type DiscreteMultiView struct {
	*windigo.AutoPanel
	buttons []*windigo.CheckBox
}

func (data_view DiscreteMultiView) Get() []string {
	var selected []string
	for _, button := range data_view.buttons {
		if button.Checked() {
			selected = append(selected, button.Text())
		}
	}
	return selected
}

func BuildNewDiscreteMultiView(parent windigo.Controller, labels []string) *DiscreteMultiView {
	view := new(DiscreteMultiView)
	height := FIELD_HEIGHT
	delta_height := FIELD_HEIGHT
	// 	log.Println("ClientWidth",parent.Width())
	// width = parent.Width()
	width := WINDOW_WIDTH
	curr_width := width
	// log.Println("ClientWidth", WINDOW_WIDTH)
	// log.Println("ClientWidth", parent.Width())
	// margin := 0
	margin := 10
	checkbox_size := 10

	overpanel := windigo.NewAutoPanel(parent)
	// panel.SetPaddingsAll(15)
	panel := windigo.NewAutoPanel(overpanel)
	panel.SetSize(GUI.OFF_AXIS, delta_height)
	overpanel.Dock(panel, windigo.Top)

	for _, label_text := range labels {
		text_width := GUI.BASE_FONT_SIZE*len(label_text)*6/5 + checkbox_size
		// text_width := checkbox_size
		curr_width -= margin
		curr_width -= text_width
		if curr_width < 0 {
			panel = windigo.NewAutoPanel(overpanel)
			panel.SetSize(GUI.OFF_AXIS, delta_height)
			overpanel.Dock(panel, windigo.Top)
			height += delta_height

			curr_width = width
			curr_width -= margin
			curr_width -= text_width
		}
		label := windigo.NewCheckBox(panel)
		label.SetMarginLeft(margin)
		label.SetText(label_text)
		label.SetSize(text_width, FIELD_HEIGHT)

		// label.OnClick().Bind(func(e *windigo.Event) {
		// 	view.Set(i)
		// })
		view.buttons = append(view.buttons, label)
		panel.Dock(label, windigo.Left)

	}
	overpanel.SetSize(GUI.OFF_AXIS, height)
	view.AutoPanel = overpanel

	return view
}
