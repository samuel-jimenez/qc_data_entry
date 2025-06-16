package main

import "github.com/samuel-jimenez/windigo"

type TextDock struct {
	windigo.Pane
	labels []*windigo.Label
}

func (control TextDock) SetDockSize(width, height int) {
	control.SetSize(width, height)
	for _, label := range control.labels {
		label.SetSize(width, height)
	}
}

func NewTextDock(parent windigo.Controller, labels []string) *TextDock {
	control := new(TextDock)

	panel := windigo.NewAutoPanel(parent)

	for _, label_text := range labels {
		label := windigo.NewLabel(panel)
		label.SetText(label_text)

		panel.Dock(label, windigo.Left)
		control.labels = append(control.labels, label)
	}

	control.Pane = panel
	return control

}
