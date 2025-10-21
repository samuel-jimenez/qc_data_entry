package viewer

import (
	"github.com/samuel-jimenez/windigo"
)

/*
 * SQLFilterViewable
 *
 */
type SQLFilterViewable interface {
	// windigo.AutoPane
	windigo.Controller
	// windigo.Editable
	// windigo.DiffLabelable
	Get() string
	Update(set []string)
	AddItem(entry string)
	DelItem(entry string)
	Clear()

	RefreshSize()
	SetFont(font *windigo.Font)
}

/* SQLFilterView
 *
 */
type SQLFilterView struct {
	*windigo.AutoPanel
	label *SQLFilterViewHeader
	key   string

	// Key() string

	// Get     func() SQLFilter
	// Update  func(set []string)
	// AddItem func(entry string)
	// DelItem func(entry string)
	//
	// RefreshSize func()
	// SetFont     func(font *windigo.Font)
}

func (view *SQLFilterView) Update(set []string) {}

func (view *SQLFilterView) AddItem(entry string) {
	view.label.AddItem(entry)
}
func (view *SQLFilterView) DelItem(entry string) {}

func (view *SQLFilterView) Clear() {
	view.label.Clear()
}

func (view *SQLFilterView) RefreshSize() {
	// view.AutoPanel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)
	view.label.RefreshSize()
}

func (view *SQLFilterView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	view.label.SetFont(font)
}
