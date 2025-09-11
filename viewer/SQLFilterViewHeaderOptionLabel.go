package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/* SQLFilterViewHeaderOptionLabel
 *
 */
type SQLFilterViewHeaderOptionLabel struct {
	// *windigo.AutoPanel
	*windigo.LabeledLabel
	close_button *windigo.PushButton
	// base         *SQLFilterView
}

func NewSQLFilterViewHeaderOptionLabel(parent *SQLFilterViewHeader, entry string) *SQLFilterViewHeaderOptionLabel {
	label := new(SQLFilterViewHeaderOptionLabel)
	// panel := windigo.NewLabeledLabel(parent)

	panel := windigo.NewAutoPanel(parent.AutoPanel)
	panel.SetSize(GUI.LABEL_WIDTH+GUI.SMOL_BUTTON_WIDTH, GUI.EDIT_FIELD_HEIGHT)

	close_button := windigo.NewPushButton(panel)
	close_button.SetText("X")
	close_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.SMOL_BUTTON_WIDTH)
	// close_button.SetPos(GUI.LABEL_WIDTH, GUI.OFF_AXIS)
	close_button.SetPos(GUI.OFF_AXIS, GUI.OFF_AXIS)
	panel.SetMarginLeft(GUI.SMOL_BUTTON_WIDTH)
	panel.SetPaddingLeft(GUI.SMOL_BUTTON_WIDTH)

	panel_label := windigo.NewLabel(panel)
	panel_label.SetText(entry)
	// panel_label.SetSize(GUI.LABEL_WIDTH, GUI.OFF_AXIS)

	close_button.OnClick().Bind(func(e *windigo.Event) {
		panel.Close()
		parent.DelItem(entry)
	})

	// panel.Dock(panel_label, windigo.Left)
	panel.Dock(panel_label, windigo.Fill)

	// label.ComponentFrame = panel
	// label.Label = panel_label
	// return &LabeledLabel{panel, label}
	label.LabeledLabel = &windigo.LabeledLabel{ComponentFrame: panel, Label: panel_label}
	label.close_button = close_button
	return label

}

// TODO
func (view *SQLFilterViewHeaderOptionLabel) RefreshSize() {

	view.SetSize(GUI.LABEL_WIDTH+GUI.SMOL_BUTTON_WIDTH, GUI.OFF_AXIS)
	view.close_button.SetSize(GUI.SMOL_BUTTON_WIDTH, GUI.SMOL_BUTTON_WIDTH)

	// view.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)
	// view.Label.SetSize(GUI.LABEL_WIDTH, GUI.OFF_AXIS)
	// view.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)

	// view.SetLabeledSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	// for _, control := range view.Filters {
	// 	//TODO
	// 	// control.RefreshSize()
	// }
}

// TODO
func (view *SQLFilterViewHeaderOptionLabel) SetFont(font *windigo.Font) {
	view.LabeledLabel.SetFont(font)
	view.close_button.SetFont(font)
}
