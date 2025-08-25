package views

import (
	"time"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
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
}

func (view *ClockTimerView) RefreshSize() {

	view.SetSize(GUI.CLOCK_WIDTH, GUI.OFF_AXIS)
}

func NewClockTimerView(parent windigo.Controller) *ClockTimerView {
	view := new(ClockTimerView)

	clock_timer_width := 35
	clock_timer_offset_h := 10
	clock_timer_offset_v := 30
	clock_font_size := 25
	clock_font_style_flags := windigo.FontNormal

	clock_panel := windigo.NewAutoPanel(parent)
	clock_panel.SetSize(GUI.CLOCK_WIDTH, GUI.OFF_AXIS)
	clock_display_now := windigo.NewSizedLabeledLabel(clock_panel, clock_timer_width, GUI.OFF_AXIS, "")
	clock_display_future := windigo.NewSizedLabeledLabel(clock_panel, clock_timer_width, GUI.OFF_AXIS, "")
	clock_display_past := windigo.NewSizedLabeledLabel(clock_panel, clock_timer_width, GUI.OFF_AXIS, "")
	clock_display_now.SetMarginTop(clock_timer_offset_v)
	clock_display_future.SetMarginLeft(clock_timer_offset_h)
	clock_display_future.SetMarginRight(clock_timer_offset_h)
	clock_display_past.SetMarginTop(clock_timer_offset_v)
	clock_display_past.SetMarginRight(clock_timer_offset_h)

	clock_font := windigo.NewFont(clock_display_now.Font().Family(), clock_font_size, clock_font_style_flags)

	clock_display_now.SetFont(clock_font)
	clock_display_future.SetFont(clock_font)
	clock_display_past.SetFont(clock_font)

	Clock_ticker := time.Tick(time.Second)

	// Dock
	clock_panel.Dock(clock_display_now, windigo.Left)
	clock_panel.Dock(clock_display_future, windigo.Left)
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
	return view
}
