package blender

import (
	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

type ComponentTypeViewer interface {
	windigo.Pane
	Update_products(product_list []string)
	SetFont(font *windigo.Font)
	RefreshSize()
}

type ComponentTypeView struct {
	*windigo.AutoPanel
	controls            []windigo.Controller
	component_add_field *GUI.SearchBox
}

func NewComponentTypeView(parent *RecipeView) *ComponentTypeView {

	view := new(ComponentTypeView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.component_add_field = GUI.NewSearchBox(view.AutoPanel)
	component_accept_button := windigo.NewPushButton(view.AutoPanel)
	component_accept_button.SetText("OK")
	component_accept_button.SetSize(ACCEPT_BUTTON_WIDTH, GUI.OFF_AXIS)
	component_cancel_button := windigo.NewPushButton(view.AutoPanel)
	component_cancel_button.SetText("Cancel")
	component_cancel_button.SetSize(CANCEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.controls = append(view.controls, view.component_add_field)
	view.controls = append(view.controls, component_accept_button)
	view.controls = append(view.controls, component_cancel_button)

	// Dock
	view.AutoPanel.Dock(view.component_add_field, windigo.Left)
	view.AutoPanel.Dock(component_accept_button, windigo.Left)
	view.AutoPanel.Dock(component_cancel_button, windigo.Left)

	component_accept_button.OnClick().Bind(func(e *windigo.Event) {
		component := NewComponentType(view.component_add_field.Text())
		Product_id := parent.Product_data[component.Component_name]
		if Product_id != DB.INVALID_ID {
			component.AddProduct(Product_id)
		} else {
			component.AddInbound()
		}
		view.AutoPanel.Hide()
		parent.Update_component_types()
	})
	component_cancel_button.OnClick().Bind(func(e *windigo.Event) {
		view.AutoPanel.Hide()
	})
	return view
}

func (view *ComponentTypeView) Update_products(product_list []string) {
	view.component_add_field.Update(product_list)
}

func (view *ComponentTypeView) SetFont(font *windigo.Font) {
	for _, control := range view.controls {
		control.SetFont(font)
	}
}

func (view *ComponentTypeView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_HEIGHT)
	view.component_add_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

}
