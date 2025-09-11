package viewer

import (
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/* SQLFilterViewHeaderer
 *
 */
type SQLFilterViewHeaderer interface {
	windigo.Controller
	RefreshSize()
	SetFont(font *windigo.Font)
	SetHidePanel(panel *windigo.AutoPanel)
	SetShowPanel(panel *windigo.AutoPanel)
	AddItem(entry string)
	DelItem(entry string)
	Clear()
}

/* SQLFilterViewHeader
 *
 */
type SQLFilterViewHeader struct {
	*windigo.AutoPanel
	child  *windigo.AutoPanel
	Hidden bool
	parent SQLFilterViewable

	label       *windigo.Label
	hide_button *windigo.PushButton

	Options map[string]*SQLFilterViewHeaderOptionLabel
}

// func NewSQLFilterViewHeader(parent *SQLFilterView, label string) *SQLFilterViewHeader {
func NewSQLFilterViewHeader(parent SQLFilterViewable, label string) *SQLFilterViewHeader {
	header := new(SQLFilterViewHeader)
	header.Options = make(map[string]*SQLFilterViewHeaderOptionLabel)

	panel := windigo.NewAutoPanel(parent)
	panel_label := windigo.NewLabel(panel)
	panel_label.SetText(label)
	hide_button := windigo.NewPushButton(panel)
	hide_button.SetText("+")
	panel.Dock(panel_label, windigo.Left)
	panel.Dock(hide_button, windigo.Left)
	hide_button.OnClick().Bind(func(e *windigo.Event) {
		grandma := parent.Parent()
		if header.child != nil {
			header.Hidden = !header.Hidden
			if header.Hidden {
				grandma.SetSize(grandma.Width(), grandma.Height()-header.child.Height())
				hide_button.SetText("+")
				header.child.Hide()
				parent.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)
			} else {
				grandma.SetSize(grandma.Width(), grandma.Height()+header.child.Height())
				hide_button.SetText("-")
				header.child.Show()
				parent.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT+header.child.Height())
			}
		}
	})

	header.AutoPanel = panel
	header.parent = parent
	// header.parent = parent.(*SQLFilterView)
	header.label = panel_label
	header.hide_button = hide_button

	return header
}

func (header *SQLFilterViewHeader) RefreshSize() {
	header.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)
	header.label.SetSize(GUI.LABEL_WIDTH, HEADER_HEIGHT)
	if header.Hidden {
		header.parent.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)
	} else {
		header.parent.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT+header.child.Height())
	}
	header.hide_button.SetSize(GUI.SMOL_BUTTON_WIDTH, HEADER_HEIGHT)
	for _, control := range header.Options {
		control.RefreshSize()
	}
}

func (header *SQLFilterViewHeader) SetFont(font *windigo.Font) {
	header.AutoPanel.SetFont(font)
	header.label.SetFont(font)
	header.hide_button.SetFont(font)
	for _, control := range header.Options {
		control.SetFont(font)
	}
}

func (header *SQLFilterViewHeader) SetHidePanel(panel *windigo.AutoPanel) {
	header.child = panel
	header.Hidden = true
	header.child.Hide()
}

func (header *SQLFilterViewHeader) SetShowPanel(panel *windigo.AutoPanel) {
	header.child = panel
	header.Hidden = false
	header.child.Show()
}

func (header *SQLFilterViewHeader) AddItem(entry string) {
	// header.Options[entry] = NewSQLFilterViewHeaderOptionLabel(header, entry)

	panel_label := NewSQLFilterViewHeaderOptionLabel(header, entry)
	header.Dock(panel_label, windigo.Left)
	header.Options[entry] = panel_label

}

func (header *SQLFilterViewHeader) DelItem(entry string) {
	delete(header.Options, entry)
	header.parent.DelItem(entry)
}

func (header *SQLFilterViewHeader) Clear() {
	for entry, control := range header.Options {
		control.Close()
		delete(header.Options, entry)
	}
}
