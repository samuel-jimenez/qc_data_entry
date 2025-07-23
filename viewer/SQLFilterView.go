package viewer

import (
	"slices"
	"strings"

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

	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(GUI.LABEL_WIDTH+SMOL_BUTTON_EDGE, FIELD_HEIGHT)

	close_button := windigo.NewPushButton(panel)
	close_button.SetText("X")
	close_button.SetSize(SMOL_BUTTON_EDGE, SMOL_BUTTON_EDGE)
	// close_button.SetPos(GUI.LABEL_WIDTH, GUI.OFF_AXIS)
	close_button.SetPos(GUI.OFF_AXIS, GUI.OFF_AXIS)
	panel.SetMarginLeft(SMOL_BUTTON_EDGE)
	panel.SetPaddingLeft(SMOL_BUTTON_EDGE)

	panel_label := windigo.NewLabel(panel)
	panel_label.SetText(entry)

	// GUI.LABEL_WIDTH

	close_button.OnClick().Bind(func(e *windigo.Event) {
		panel.Close()
		parent.parent.DelItem(entry)
	})

	panel.Dock(panel_label, windigo.Left)
	// panel.Dock(panel_label, windigo.Fill)
	// panel.Dock(close_button, windigo.Left)

	// label.ComponentFrame = panel
	// label.Label = panel_label
	// return &LabeledLabel{panel, label}
	label.LabeledLabel = &windigo.LabeledLabel{ComponentFrame: panel, Label: panel_label}
	label.close_button = close_button
	return label

}

/* SQLFilterViewHeader
 *
 */
type SQLFilterViewHeader struct {
	*windigo.AutoPanel
	child  *windigo.AutoPanel
	Hidden bool
	parent *SQLFilterView
	//TODO // Options map[string]SQLFilterViewHeaderOptionLabel

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
	// panel_label := windigo.NewLabel(header)
	panel_label := NewSQLFilterViewHeaderOptionLabel(header, entry)

	header.Dock(panel_label, windigo.Left)
}

func NewSQLFilterViewHeader(parent *SQLFilterView, label string) *SQLFilterViewHeader {
	header := new(SQLFilterViewHeader)
	panel := windigo.NewAutoPanel(parent)
	panel_label := windigo.NewLabel(panel)
	panel_label.SetText(label)
	hide_button := windigo.NewPushButton(panel)
	hide_button.SetText("+")
	hide_button.SetSize(SMOL_BUTTON_EDGE, SMOL_BUTTON_EDGE)
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
				parent.SetSize(GUI.OFF_AXIS, SMOL_BUTTON_EDGE)
			} else {
				grandma.SetSize(grandma.Width(), grandma.Height()+header.child.Height())
				hide_button.SetText("-")
				header.child.Show()
				parent.SetSize(GUI.OFF_AXIS, SMOL_BUTTON_EDGE+header.child.Height())
			}
		}
	})

	header.AutoPanel = panel
	header.parent = parent

	return header
}

// TODO
// /*
//  * SQLFilterViewable
//  *
//  */
// type SQLFilterViewable interface {
// 	windigo.AutoPane
// 	windigo.Editable
// 	windigo.DiffLabelable
// 	Get() SQLFilter
// 	Update(set []string)
// }

/* SQLFilterView
 *
 */
type SQLFilterView struct {
	*windigo.AutoPanel
	// Key() string
	Get     func() SQLFilter
	Update  func(set []string)
	AddItem func(entry string)
	DelItem func(entry string)
}

// func (s SQLFilterView) AddItem(entry string) {
// }

// func NewSQLFilterView(parent windigo.Controller) SQLFilterView {
// 	panel := windigo.NewAutoPanel(parent)
// 	return SQLFilterView{panel}
// }

// SQLFilter

func NewSQLFilterViewContinuous(parent windigo.Controller, key, label string) *SQLFilterView {
	view := new(SQLFilterView)
	panel := windigo.NewAutoPanel(parent)
	view.AutoPanel = panel
	panel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	panel_label := NewSQLFilterViewHeader(view, label)
	panel_label.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	field_panel := windigo.NewAutoPanel(panel)
	panel_label.SetHidePanel(field_panel)

	field_panel.SetSize(GUI.OFF_AXIS, FIELD_HEIGHT)
	min_field := NewNullStringView(field_panel)
	max_field := NewNullStringView(field_panel)
	field_panel.Dock(min_field, windigo.Left)
	field_panel.Dock(max_field, windigo.Left)
	panel.Dock(panel_label, windigo.Top)
	panel.Dock(field_panel, windigo.Top)

	get := func() SQLFilter {
		return SQLFilterContinuous{key,
			min_field.Get(),
			max_field.Get()}
	}

	view.Get = get
	return view
}

func NewSQLFilterViewDiscreteMulti(parent windigo.Controller, key, label string, set []string) *SQLFilterView {
	view := new(SQLFilterView)
	panel := windigo.NewAutoPanel(parent)
	view.AutoPanel = panel
	panel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	panel_label := NewSQLFilterViewHeader(view, label)
	panel_label.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	selection_options := BuildNewDiscreteMultiView(panel, set)
	panel_label.SetHidePanel(selection_options.AutoPanel)

	panel.Dock(panel_label, windigo.Top)
	panel.Dock(selection_options, windigo.Top)

	get := func() SQLFilter {
		return SQLFilterDiscrete{key,
			selection_options.Get()}
	}
	update := func(set []string) {
		selection_options.Close()
		selection_options = BuildNewDiscreteMultiView(panel, set)
		if panel_label.Hidden {
			panel_label.SetHidePanel(selection_options.AutoPanel)
		} else {
			panel_label.SetShowPanel(selection_options.AutoPanel)
			panel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT+selection_options.Height())
		}
		panel.Dock(selection_options, windigo.Top)
	}

	view.Get = get
	view.Update = update
	return view

}

func NewSQLFilterViewDiscreteSearch(parent windigo.Controller, key, label string, set []string) *SQLFilterView {
	view := new(SQLFilterView)

	panel := windigo.NewAutoPanel(parent)
	view.AutoPanel = panel
	panel.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	panel_label := NewSQLFilterViewHeader(view, label)
	panel_label.SetSize(GUI.OFF_AXIS, HEADER_HEIGHT)

	selection_options := BuildNewDiscreteSearchView(view, set)
	panel_label.SetHidePanel(selection_options.AutoPanel)
	selection_options.SetSize(GUI.OFF_AXIS, FIELD_HEIGHT)

	panel.Dock(panel_label, windigo.Top)
	panel.Dock(selection_options, windigo.Top)

	get := func() SQLFilter {
		return SQLFilterDiscrete{key,
			selection_options.Get()}
	}

	view.Get = get
	view.Update = selection_options.Update
	view.AddItem = panel_label.AddItem
	view.DelItem = selection_options.DelItem

	return view
}

/*
 *
 * DiscreteSearchView

*
* TODO reconcile with SearchBox

GUI.SearchBox

* ComboBoxable ?
*/

/*
 *
 */
type DiscreteSearchView struct {
	*windigo.AutoPanel
	box *GUI.ComboBox
	entries,
	chosen []string
}

// TODO ?
// Filters map[string]SQLFilterView
func (data_view DiscreteSearchView) Get() []string {
	// var selected []string
	// for _, button := range data_view.buttons {
	// 	if button.Checked() {
	// 		selected = append(selected, button.Text())
	// 	}
	// }
	return data_view.chosen
	//TODO ?
}

func (data_view *DiscreteSearchView) Update(set []string) {
	data_view.entries = set
	data_view.box.DeleteAllItems()
	for _, name := range set {
		data_view.box.AddItem(name)
	}
}

func (data_view DiscreteSearchView) Search(terms []string) {
	data_view.box.DeleteAllItems()

	for _, entry := range data_view.entries {
		matched := true
		upcased := strings.ToUpper(entry)

		for _, term := range terms {
			if !strings.Contains(upcased, term) {
				matched = false
				break
			}
		}
		if matched {
			data_view.box.AddItem(entry)
		}
	}
}

func (data_view *DiscreteSearchView) DelItem(entry string) {
	i := slices.Index(data_view.chosen, entry)
	if i > -1 {
		data_view.chosen = slices.Delete(data_view.chosen, i, i+1)
	}
}

// func BuildNewDiscreteSearchView(parent windigo.Controller, labels []string) *DiscreteSearchView {
func BuildNewDiscreteSearchView(parent *SQLFilterView, labels []string) *DiscreteSearchView {
	data_view := new(DiscreteSearchView)
	// 	log.Println("ClientWidth",parent.Width())
	// width = parent.Width()

	panel := windigo.NewAutoPanel(parent)
	// panel.SetPaddingsAll(15)
	// panel := windigo.NewAutoPanel(overpanel)

	data_view.box = GUI.NewComboBox(panel, "")
	data_view.box.SetLabeledSize(GUI.OFF_AXIS, PRODUCT_FIELD_WIDTH, FIELD_HEIGHT)
	data_view.box.OnChange().Bind(func(e *windigo.Event) {

		start, end := data_view.box.Selected()
		text := strings.ToUpper(data_view.box.Text())

		terms := strings.Split(text, " ")
		data_view.Search(terms)

		data_view.box.SetText(text)
		data_view.box.SelectText(start, end)
		data_view.box.ShowDropdown(true)

	})
	data_view.box.OnSelectedChange().Bind(func(e *windigo.Event) {
		entry := data_view.box.GetSelectedItem()
		data_view.chosen = append(data_view.chosen, entry)
		parent.AddItem(entry)
	})

	// view.box.SetSize(GUI.OFF_AXIS, FIELD_HEIGHT)
	panel.SetSize(GUI.OFF_AXIS, FIELD_HEIGHT)

	panel.Dock(data_view.box, windigo.Left)

	// panel.Dock(view.box, windigo.Top)

	data_view.AutoPanel = panel

	data_view.Update(labels)

	return data_view
}

//TODO reconcile
// func NewSearchBox(parent windigo.Controller) *SearchBox {
// 	data_view := new(SearchBox)
//
// 	data_view.ComboBox = NewComboBox(parent, "")
// 	data_view.ComboBox.OnChange().Bind(func(e *windigo.Event) {
// 		log.Println("CRIT: DEBUG: NewSearchBox OnChange", data_view.Text(), data_view.SelectedItem())
//
// 		start, _ := data_view.Selected()
//
// 		text := strings.ToUpper(data_view.Text())
//
// 		terms := strings.Split(text, " ")
// 		data_view.Search(terms)
//
// 		data_view.ShowDropdown(true)
//
// 		data_view.SetText(text)
//
// 		// data_view.SelectText(start, end)
// 		data_view.SelectText(start, -1)
// 		data_view.onChange.Fire(e)
// 		// data_view.OnChange().Fire(e)
//
// 	})
//
// 	return data_view
// }
