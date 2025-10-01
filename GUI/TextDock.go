package GUI

import "github.com/samuel-jimenez/windigo"

type TextDock struct {
	windigo.Pane
	Labels []*windigo.Label
}

func NewTextDock(parent windigo.Controller, labels ...string) *TextDock {

	control := new(TextDock)

	panel := windigo.NewAutoPanel(parent)

	for _, label_text := range labels {
		label := windigo.NewLabel(panel)
		label.SetText(label_text)

		panel.Dock(label, windigo.Left)
		control.Labels = append(control.Labels, label)
	}

	control.Pane = panel
	return control

}

func (view TextDock) SetDockSize(width, height int) {
	view.SetSize(width, height)
	for _, label := range view.Labels {
		label.SetSize(width, height)
	}
}

func (view *TextDock) SetFont(font *windigo.Font) {
	for _, label := range view.Labels {
		label.SetFont(font)
	}
}
