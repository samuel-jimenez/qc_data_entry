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
	panels     []*windigo.AutoPanel
	buttons    []*windigo.CheckBox
	baseHeight int
	ready_p    bool
}

func DiscreteMultiView_from_new(parent windigo.Controller, labels []string) *DiscreteMultiView {
	view := new(DiscreteMultiView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	view.baseHeight = view.Height()
	check_labels := make([]CheckItem, len(labels))

	for i, label := range labels {
		check_labels[i] = CheckItem{Label: label}
	}
	view.build(check_labels)
	view.ready_p = true
	return view
}

type CheckItem struct {
	Label   string
	Checked bool
}

func (view *DiscreteMultiView) rebuild() {
	labels := make([]CheckItem, len(view.buttons))

	for i, label := range view.buttons {
		labels[i] = CheckItem{label.Text(), label.Checked()}
		label.Close()
	}
	for _, panel := range view.panels {
		panel.Close()
	}
	view.buttons = nil
	view.panels = nil
	view.build(labels)
}

func (view *DiscreteMultiView) build(labels []CheckItem) {
	height := view.baseHeight
	delta_height := view.baseHeight
	width := view.Width()
	curr_width := width
	margin := 10
	checkbox_size := 20 + view.GetSizeOfText(" ").Width

	panel := windigo.NewAutoPanel(view)
	panel.SetSize(GUI.OFF_AXIS, delta_height)
	view.Dock(panel, windigo.Top)
	view.panels = append(view.panels, panel)

	for _, label_text := range labels {
		text_width := panel.GetSizeOfText(label_text.Label).Width + checkbox_size
		curr_width -= margin
		curr_width -= text_width
		if curr_width < 0 {
			panel = windigo.NewAutoPanel(view)
			panel.SetSize(GUI.OFF_AXIS, delta_height)
			view.Dock(panel, windigo.Top)
			view.panels = append(view.panels, panel)

			height += delta_height
			curr_width = width
			curr_width -= margin
			curr_width -= text_width
		}
		label := windigo.NewCheckBox(panel)
		label.SetMarginLeft(margin)
		label.SetText(label_text.Label)
		label.SetChecked(label_text.Checked)
		label.SetSize(text_width, delta_height)
		panel.Dock(label, windigo.Left)
		view.buttons = append(view.buttons, label)

	}
	view.AutoPanel.SetSize(width, height)
	//
	if view.ready_p {
		view.Parent().Parent().RefreshSize()
	}
}

func (view *DiscreteMultiView) RefreshSize() {
	height := view.baseHeight
	delta_height := view.baseHeight
	width := view.Width()
	curr_width := width
	margin := 10
	checkbox_size := 20
	panel_index := 0
	var changed bool
	panel := view.panels[panel_index]
	panel.SetSize(width, delta_height)
	for _, label := range view.buttons {
		text_width := label.TextSize().Width + checkbox_size
		curr_width -= margin
		curr_width -= text_width

		if label.Parent() != panel && curr_width < 0 {
			panel_index++
			panel = view.panels[panel_index]
			panel.SetSize(width, delta_height)
			height += delta_height

			curr_width = width
			curr_width -= margin
			curr_width -= text_width

		} else if label.Parent() != panel || curr_width < 0 {
			changed = true
			break

		}
		label.SetSize(text_width, delta_height)

	}
	if changed {
		view.rebuild()
		return
	}
	view.AutoPanel.SetForcedSize(width, height)
}

func (view *DiscreteMultiView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	for _, control := range view.buttons {
		control.SetFont(font)
	}
}

func (view *DiscreteMultiView) SetBaseSize(width, height int) {
	view.baseHeight = height
	view.AutoPanel.SetSize(width, height*len(view.panels))
}

func (view *DiscreteMultiView) SetSize(width, height int) {
	view.SetBaseSize(width, height/len(view.panels))
}

func (view *DiscreteMultiView) SetForcedSize(width, height int) {
	view.AutoPanel.SetForcedSize(width, height)
	if width == 0 || height == 0 {
		return
	}
	view.RefreshSize()
}

func (view *DiscreteMultiView) Close() {
	view.buttons = nil // don't Refresh if we're closing
	view.AutoPanel.Close()
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

func (view *DiscreteMultiView) Clear() {
	for _, control := range view.buttons {
		control.SetChecked(false)
	}
}
