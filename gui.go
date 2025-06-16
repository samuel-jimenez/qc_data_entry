package main

import (
	"github.com/samuel-jimenez/windigo"
	"github.com/samuel-jimenez/windigo/w32"
)

var (
	WINDOW_WIDTH,
	WINDOW_HEIGHT,
	WINDOW_FUDGE_MARGIN,

	GROUPBOX_CUSHION,
	TOP_SPACER_WIDTH,
	TOP_SPACER_HEIGHT,
	INTER_SPACER_HEIGHT,
	BTM_SPACER_WIDTH,
	BTM_SPACER_HEIGHT,

	TOP_PANEL_WIDTH,
	REPRINT_BUTTON_WIDTH,
	TOP_PANEL_INTER_SPACER_WIDTH,

	RANGE_WIDTH,
	PRODUCT_TYPE_WIDTH,

	LABEL_WIDTH,
	PRODUCT_FIELD_WIDTH,
	PRODUCT_FIELD_HEIGHT,
	DATA_FIELD_WIDTH,
	DATA_SUBFIELD_WIDTH,
	DATA_UNIT_WIDTH,
	DATA_MARGIN,
	FIELD_HEIGHT,
	NUM_FIELDS,

	DISCRETE_FIELD_WIDTH,
	DISCRETE_FIELD_HEIGHT,

	RANGES_WINDOW_WIDTH,
	RANGES_WINDOW_HEIGHT,
	RANGES_WINDOW_PADDING,
	RANGES_PADDING,
	RANGES_FIELD_HEIGHT,
	RANGES_FIELD_SMALL_HEIGHT,
	RANGES_FIELD_WIDTH,
	RANGES_BUTTON_WIDTH,
	RANGES_BUTTON_HEIGHT,
	RANGES_BUTTON_MARGIN,

	OFF_AXIS,
	RANGES_RO_PADDING,
	RANGES_RO_FIELD_WIDTH,
	RANGES_RO_SPACER_WIDTH,
	RANGES_RO_FIELD_HEIGHT,

	BUTTON_WIDTH,
	BUTTON_HEIGHT,
	BUTTON_MARGIN,

	ERROR_MARGIN,

	GROUP_WIDTH,
	GROUP_HEIGHT,
	GROUP_MARGIN int

	erroredPen = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSolidColorBrush(windigo.RGB(255, 0, 64)))
	okPen      = windigo.NewPen(w32.PS_GEOMETRIC, 2, windigo.NewSystemColorBrush(w32.COLOR_BTNFACE))
)

func refresh_globals(font_size int) {

	GROUPBOX_CUSHION = font_size * 3 / 2
	TOP_SPACER_WIDTH = 7
	TOP_SPACER_HEIGHT = GROUPBOX_CUSHION + 2
	INTER_SPACER_HEIGHT = 2
	BTM_SPACER_WIDTH = 2
	BTM_SPACER_HEIGHT = 2
	TOP_PANEL_INTER_SPACER_WIDTH = 30

	LABEL_WIDTH = 10 * font_size
	PRODUCT_FIELD_WIDTH = 15 * font_size
	PRODUCT_FIELD_HEIGHT = font_size*16/10 + 8
	DATA_FIELD_WIDTH = 6 * font_size
	DATA_SUBFIELD_WIDTH = 5 * font_size
	DATA_UNIT_WIDTH = 4 * font_size
	// DATA_UNIT_WIDTH = 3 * font_size
	DATA_MARGIN = 10

	FIELD_HEIGHT = 3*font_size - 2
	NUM_FIELDS = 6

	OFF_AXIS = 0
	RANGES_RO_PADDING = 10
	RANGES_RO_FIELD_WIDTH = 7*font_size/2 + 14
	RANGES_RO_SPACER_WIDTH = 20
	RANGES_RO_FIELD_HEIGHT = FIELD_HEIGHT

	DISCRETE_FIELD_WIDTH = 17 * font_size
	DISCRETE_FIELD_HEIGHT = font_size * 9 / 2
	PRODUCT_TYPE_WIDTH = 8*font_size + 90

	RANGES_PADDING = 10
	// RANGES_FIELD_HEIGHT =  DISCRETE_FIELD_HEIGHT
	RANGES_FIELD_HEIGHT = font_size*5/2 + 30
	RANGES_FIELD_SMALL_HEIGHT = font_size * 5 / 2
	RANGES_FIELD_WIDTH = 200
	RANGES_BUTTON_WIDTH = 200
	RANGES_BUTTON_HEIGHT = RANGES_FIELD_HEIGHT
	RANGES_BUTTON_MARGIN = 10
	RANGES_WINDOW_PADDING = 5

	BUTTON_WIDTH = 10 * font_size
	BUTTON_HEIGHT = 40

	REPRINT_BUTTON_WIDTH = 10 * font_size

	ERROR_MARGIN = 3
	GROUP_MARGIN = 5
	WINDOW_FUDGE_MARGIN = 16
	BUTTON_MARGIN = 5

	RANGES_WINDOW_WIDTH = 2*RANGES_WINDOW_PADDING + 2*RANGES_PADDING + LABEL_WIDTH + 3*RANGES_FIELD_WIDTH + WINDOW_FUDGE_MARGIN
	RANGES_WINDOW_HEIGHT = RANGES_WINDOW_PADDING + 3*RANGES_PADDING + 2*ERROR_MARGIN + DISCRETE_FIELD_HEIGHT + 2*RANGES_FIELD_SMALL_HEIGHT + 7*RANGES_FIELD_HEIGHT + RANGES_BUTTON_HEIGHT + WINDOW_FUDGE_MARGIN

	TOP_PANEL_WIDTH = 2*LABEL_WIDTH + PRODUCT_FIELD_WIDTH + TOP_PANEL_INTER_SPACER_WIDTH + 2*REPRINT_BUTTON_WIDTH + 4*BUTTON_MARGIN

	// GROUP_WIDTH = 20*font_size + 10
	GROUP_WIDTH = TOP_SPACER_WIDTH + LABEL_WIDTH + DATA_SUBFIELD_WIDTH + DATA_UNIT_WIDTH + DATA_MARGIN + BTM_SPACER_WIDTH + 2*ERROR_MARGIN
	GROUP_HEIGHT = TOP_SPACER_HEIGHT +
		NUM_FIELDS*FIELD_HEIGHT + BTM_SPACER_HEIGHT

	// 	WINDOW_WIDTH = max(65*font_size, 400)
	// WINDOW_HEIGHT = 20*font_size + 300
	// RANGE_WIDTH = WINDOW_WIDTH - 2*GROUP_WIDTH - GROUP_MARGIN

	RANGE_WIDTH = TOP_SPACER_WIDTH + 3*RANGES_RO_FIELD_WIDTH + 2*RANGES_RO_SPACER_WIDTH + RANGES_RO_PADDING
	// RANGE_WIDTH = 3*RANGES_RO_FIELD_WIDTH + 2*RANGES_RO_SPACER_WIDTH
	// RANGE_WIDTH = WINDOW_WIDTH - 2*GROUP_WIDTH - GROUP_MARGIN

	// WINDOW_WIDTH = max(65*font_size, 400)
	// WINDOW_WIDTH = GROUP_MARGIN + 2*GROUP_WIDTH + RANGE_WIDTH + WINDOW_FUDGE_MARGIN
	WINDOW_WIDTH = max(GROUP_MARGIN+2*GROUP_WIDTH+RANGE_WIDTH, TOP_PANEL_WIDTH) + WINDOW_FUDGE_MARGIN

	WINDOW_HEIGHT = 20*font_size + 300

}

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

func build_text_dock(parent windigo.Controller, labels []string) *TextDock {
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

func build_button_dock(parent windigo.Controller, labels []string, onclicks []func()) *ButtonDock {
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

func build_marginal_button_dock(parent windigo.Controller, labels []string, margins []int, onclicks []func()) *ButtonDock {
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
