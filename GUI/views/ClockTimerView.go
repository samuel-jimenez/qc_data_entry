package views

import (
	"time"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/windigo"
)

/*
 * ClockTimerViewer
 *
 */
type ClockTimerViewer interface {
	windigo.Controller
	SetFont(font *windigo.Font)
	RefreshSize()
}

/*
 * ClockTimerView
 *
 */
type ClockTimerView struct {
	*windigo.AutoPanel
	clock_display_future,
	clock_display_now,
	clock_display_past *windigo.LabeledLabel
	clock_font *windigo.Font
}

func NewClockTimerView(parent windigo.Controller) *ClockTimerView {
	view := new(ClockTimerView)

	clock_font_size := 25
	clock_font_style_flags := windigo.FontNormal

	clock_panel := windigo.NewAutoPanel(parent)
	clock_display_future := windigo.NewLabeledLabel(clock_panel, "")
	clock_display_now := windigo.NewLabeledLabel(clock_panel, "")
	clock_display_past := windigo.NewLabeledLabel(clock_panel, "")

	view.clock_font = windigo.NewFont(clock_display_now.Font().Family(), clock_font_size, clock_font_style_flags)

	clock_display_future.SetFont(view.clock_font)
	clock_display_now.SetFont(view.clock_font)
	clock_display_past.SetFont(view.clock_font)

	clock_display_future.SetFGColor(windigo.RGB(128, 0, 0))
	clock_display_now.SetFGColor(windigo.RGB(0, 128, 0))
	clock_display_past.SetFGColor(windigo.RGB(0, 0, 128))

	Clock_ticker := time.Tick(time.Second)

	// Dock
	clock_panel.Dock(clock_display_future, windigo.Left)
	clock_panel.Dock(clock_display_now, windigo.Left)
	clock_panel.Dock(clock_display_past, windigo.Left)

	// Clock
	go func() {
		mix_time := 20 * time.Second
		for tock := range Clock_ticker {
			clock_display_now.SetText(tock.Format("05"))
			clock_display_future.SetText(tock.Add(mix_time).Format("05"))
			clock_display_past.SetText(tock.Add(-mix_time).Format("05"))
		}
	}()

	// build object
	view.AutoPanel = clock_panel
	view.clock_display_future = clock_display_future
	view.clock_display_now = clock_display_now
	view.clock_display_past = clock_display_past
	return view
}

func (view *ClockTimerView) SetFont(font *windigo.Font) {
	view.clock_display_future.SetFont(font)
	view.clock_display_now.SetFont(font)
	view.clock_display_past.SetFont(font)
}

func (view *ClockTimerView) RefreshSize() {
	clock_font_size := 3 * config.BASE_FONT_SIZE
	clock_font_style_flags := windigo.FontNormal

	old_font := view.clock_font
	view.clock_font = windigo.NewFont(old_font.Family(), clock_font_size, clock_font_style_flags)
	old_font.Dispose()
	view.SetFont(view.clock_font)

	clock_timer_offset_v := 3 * config.BASE_FONT_SIZE

	view.SetSize(GUI.CLOCK_WIDTH, GUI.OFF_AXIS)
	view.clock_display_future.SetSize(GUI.CLOCK_TIMER_WIDTH, GUI.OFF_AXIS)
	view.clock_display_now.SetSize(GUI.CLOCK_TIMER_WIDTH, GUI.OFF_AXIS)
	view.clock_display_past.SetSize(GUI.CLOCK_TIMER_WIDTH, GUI.OFF_AXIS)

	view.clock_display_future.SetMarginTop(clock_timer_offset_v)
	view.clock_display_future.SetMarginLeft(GUI.CLOCK_TIMER_OFFSET_H)
	view.clock_display_future.SetMarginRight(GUI.CLOCK_TIMER_OFFSET_H)
	view.clock_display_now.SetMarginRight(GUI.CLOCK_TIMER_OFFSET_H)
	view.clock_display_past.SetMarginTop(clock_timer_offset_v)
	view.clock_display_past.SetMarginRight(GUI.CLOCK_TIMER_OFFSET_H)
}
