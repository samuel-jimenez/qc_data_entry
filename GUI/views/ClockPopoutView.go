package views

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * ClockPopoutView
 *
 */
type ClockPopoutView struct {
	*windigo.AutoDialog

	clock_panel *ClockTimerView
}

func ClockPopoutView_from_new(parent windigo.Controller) *ClockPopoutView {
	window_title := "String Test"
	view := new(ClockPopoutView)

	view.AutoDialog = windigo.AutoDialog_from_new(parent)
	view.Dialog.SetText(window_title)

	view.clock_panel = NewClockTimerView(view)
	view.Dock(view.clock_panel, windigo.Fill)

	// view.clock_panel.SetOnLBDown(view.EventClose)
	view.clock_panel.SetOnLBUp(view.EventClose)

	view.RefreshSize()
	return view
}

func (view *ClockPopoutView) SetSize(w, h int) {
	view.Dialog.SetSize(w+GUI.WINDOW_FUDGE_MARGIN_W, h+GUI.WINDOW_FUDGE_MARGIN_H)
}

func (view *ClockPopoutView) RefreshSize() {
	// view.SetSize(GUI.CLOCK_WIDTH, qc.TOP_SUBPANEL_HEIGHT)
	// 	num_rows := 3
	//
	// 	TOP_SUBPANEL_HEIGHT = GUI.TOP_SPACER_HEIGHT + num_rows*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT) + GUI.BTM_SPACER_HEIGHT
	CLOCK_HEIGHT := GUI.TOP_SPACER_HEIGHT + 3*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT) + GUI.BTM_SPACER_HEIGHT // TODO

	// view.SetSize(GUI.CLOCK_WIDTH, GUI.TOP_SPACER_HEIGHT+3*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT)+GUI.BTM_SPACER_HEIGHT)
	// view.SetSize(GUI.CLOCK_WIDTH, GUI.TOP_SPACER_HEIGHT+3*(GUI.PRODUCT_FIELD_HEIGHT+GUI.INTER_SPACER_HEIGHT)+GUI.BTM_SPACER_HEIGHT+GUI.GROUPBOX_CUSHION)
	view.SetSize(GUI.CLOCK_WIDTH, CLOCK_HEIGHT+GUI.GROUPBOX_CUSHION)
	// + num_rows*+ GUI.PRODUCT_FIELD_HEIGHT
	view.clock_panel.RefreshSize()
}
