package main

import "github.com/samuel-jimenez/windigo"

type ButtonDock struct {
	windigo.Pane
	buttons []*windigo.PushButton
}

func (control ButtonDock) SetFont(font *windigo.Font) {
	for _, button := range control.buttons {
		button.SetFont(font)
	}
}

func (control ButtonDock) SetDockSize(width, height int) {
	control.SetSize(width, height)
	for _, button := range control.buttons {
		button.SetSize(width, height)
	}
}

func NewButtonDock(parent windigo.Controller, labels []string, onclicks []func()) *ButtonDock {
	assertEqual(len(labels), len(onclicks))
	control := new(ButtonDock)
	panel := windigo.NewAutoPanel(parent)

	for i, label_text := range labels {
		button := windigo.NewPushButton(panel)
		button.SetMarginsAll(RANGES_BUTTON_MARGIN)

		button.OnClick().Bind(func(e *windigo.Event) {
			onclicks[i]()

		})
		button.SetText(label_text)
		panel.Dock(button, windigo.Left)
		control.buttons = append(control.buttons, button)

	}

	control.Pane = panel
	return control
}

func NewMarginalButtonDock(parent windigo.Controller, labels []string, margins []int, onclicks []func()) *ButtonDock {
	assertEqual(len(labels), len(margins))
	assertEqual(len(labels), len(onclicks))
	control := new(ButtonDock)
	panel := windigo.NewAutoPanel(parent)

	for i, label_text := range labels {
		button := windigo.NewPushButton(panel)
		button.SetMarginLeft(margins[i])

		button.OnClick().Bind(func(e *windigo.Event) {
			onclicks[i]()

		})
		button.SetText(label_text)
		control.buttons = append(control.buttons, button)

		panel.Dock(button, windigo.Left)

	}

	control.Pane = panel
	return control
}
