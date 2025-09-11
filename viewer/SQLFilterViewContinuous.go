package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/* SQLFilterViewContinuous
 *
 */
type SQLFilterViewContinuous struct {
	SQLFilterView

	field_panel          *windigo.AutoPanel
	min_field, max_field NullStringView
}

func NewSQLFilterViewContinuous(parent windigo.Controller, key, label string) *SQLFilterViewContinuous {
	view := new(SQLFilterViewContinuous)
	panel := windigo.NewAutoPanel(parent)
	view.AutoPanel = panel
	panel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	panel_label := NewSQLFilterViewHeader(view, label)
	panel_label.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	field_panel := windigo.NewAutoPanel(panel)
	panel_label.SetHidePanel(field_panel)

	field_panel.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	min_field := NewNullStringView(field_panel)
	max_field := NewNullStringView(field_panel)
	field_panel.Dock(min_field, windigo.Left)
	field_panel.Dock(max_field, windigo.Left)
	panel.Dock(panel_label, windigo.Top)
	panel.Dock(field_panel, windigo.Top)

	view.label = panel_label
	view.field_panel = field_panel
	view.min_field = min_field
	view.max_field = max_field
	view.key = key
	return view
}

func (view *SQLFilterViewContinuous) Get() SQLFilter {
	return SQLFilterContinuous{view.key,
		view.min_field.Get(),
		view.max_field.Get()}
}

func (view *SQLFilterViewContinuous) RefreshSize() {
	view.SQLFilterView.RefreshSize()
	view.field_panel.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
}

func (view *SQLFilterViewContinuous) SetFont(font *windigo.Font) {
	view.SQLFilterView.SetFont(font)
	view.field_panel.SetFont(font)
	view.min_field.SetFont(font)
	view.max_field.SetFont(font)
}

func (view *SQLFilterViewContinuous) Clear() {
	view.SQLFilterView.Clear()
	view.min_field.Clear()
	view.max_field.Clear()
}
