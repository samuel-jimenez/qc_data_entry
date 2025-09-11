package viewer

import (
	"slices"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 *
 * DiscreteSearchView
 *
 */
type DiscreteSearchView struct {
	*windigo.AutoPanel
	box    *GUI.SearchBox
	chosen []string
	parent *SQLFilterViewDiscreteSearch
}

func BuildNewDiscreteSearchView(parent *SQLFilterViewDiscreteSearch, labels []string) *DiscreteSearchView {
	view := new(DiscreteSearchView)
	view.parent = parent

	panel := windigo.NewAutoPanel(parent)
	view.AutoPanel = panel

	view.box = GUI.NewListSearchBoxWithLabels(panel, labels)
	view.box.OnSelectedChange().Bind(view.AddItem)

	panel.Dock(view.box, windigo.Left)

	view.RefreshSize()

	return view
}

func (view *DiscreteSearchView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_HEIGHT)
	view.box.SetLabeledSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
}

func (view *DiscreteSearchView) SetFont(font *windigo.Font) {
	view.AutoPanel.SetFont(font)
	view.box.SetFont(font)
}

// TODO ?
// Filters map[string]SQLFilterView
func (view DiscreteSearchView) Get() []string {
	// var selected []string
	// for _, button := range data_view.buttons {
	// 	if button.Checked() {
	// 		selected = append(selected, button.Text())
	// 	}
	// }
	return view.chosen
}

func (view *DiscreteSearchView) Update(set []string) {
	view.box.Update(set)
}

func (view *DiscreteSearchView) AddItem(e *windigo.Event) {
	entry := view.box.GetSelectedItem()
	if i := slices.Index(view.chosen, entry); i >= 0 {
		return
	}
	view.chosen = append(view.chosen, entry)
	view.parent.AddItem(entry)
}

func (view *DiscreteSearchView) DelItem(entry string) {
	i := slices.Index(view.chosen, entry)
	if i < 0 {
		return
	}
	view.chosen = slices.Delete(view.chosen, i, i+1)
}

func (view *DiscreteSearchView) Clear() {
	view.chosen = nil
}
