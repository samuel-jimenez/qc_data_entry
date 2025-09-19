package viewer

import (
	"slices"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 *
 * DiscreteSearchViewer
 *
 */
type DiscreteSearchViewer interface {
	windigo.Controller
	SetFont(font *windigo.Font)
	RefreshSize()

	Get() []string
	Update(set []string)

	box_OnSelectedChange(e *windigo.Event)
	Contains(entry string) bool

	AddItem(entry string)
	DelItem(entry string)

	Clear()
}

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

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.box = GUI.NewListSearchBoxWithLabels(view.AutoPanel, labels)
	view.box.OnSelectedChange().Bind(view.box_OnSelectedChange)

	view.AutoPanel.Dock(view.box, windigo.Left)

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

func (view DiscreteSearchView) Get() []string {
	return view.chosen
}

func (view *DiscreteSearchView) Update(set []string) {
	view.box.Update(set)
}

func (view *DiscreteSearchView) box_OnSelectedChange(e *windigo.Event) {
	entry := view.box.GetSelectedItem()
	view.parent.AddItem(entry)
}

func (view *DiscreteSearchView) Contains(entry string) bool {
	i := slices.Index(view.chosen, entry)
	return i >= 0
}

func (view *DiscreteSearchView) AddItem(entry string) {
	view.chosen = append(view.chosen, entry)
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
