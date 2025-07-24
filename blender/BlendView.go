package blender

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

type BlendViewer interface {
	windigo.Pane
	// Get() *ProductBlend
	// Update(blend *ProductBlend)
	// Update_component_types(component_types_list []string, component_types_data map[string]int64)
	// AddComponent()
	SetFont(font *windigo.Font)
	RefreshSize()
}

type BlendView struct {
	*windigo.AutoPanel
	Recipe               *ProductRecipe
	Blend                *ProductBlend
	Components           []*BlendComponentView
	component_types_list []string
	component_types_data map[string]int64
	// Product_id int64
	// Blend_id  int64
}

func NewBlendView(parent windigo.Controller) *BlendView {
	view := new(BlendView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	//TODO normalize amounts

	return view
}

/*
func (view *BlendView) Get() *ProductBlend {
	if view.Blend == nil {
		return nil
	}
	//TODO
	view.Blend.Components = nil
	for i, component := range view.Components {
		blendComponent := component.Get()
		log.Println("DEBUG: BlendView Get", blendComponent)
		if blendComponent != nil {
			blendComponent.Add_order = i
			log.Println("DEBUG: BlendView Get", blendComponent)
			view.Blend.AddComponent(*blendComponent)
		}
	}
	view.Blend.SaveComponents()
	return view.Blend

}

func (view *BlendView) Update(blend *ProductBlend) {
	if view.Blend == blend {
		return
	}

	// height := FIELD_HEIGHT
	// delta_height := FIELD_HEIGHT
	// 	total_height := height
	// delta_height := height
	// 	height := FIELD_HEIGHT
	// delta_height := FIELD_HEIGHT

	for _, component := range view.Components {
		component.Close()
	}
	view.Components = nil
	view.Blend = blend
	log.Println("BlendView Update", view.Blend)
	if view.Blend == nil {
		return
	}

	for _, component := range view.Blend.Components {
		log.Println("DEBUG: BlendView update_components", component)
		// view.AddComponent()
		//TODO
		component_view := view.AddComponent()
		component_view.Update(&component)
	}

}

func (view *BlendView) Update_component_types(component_types_list []string, component_types_data map[string]int64) {
	view.component_types_list = component_types_list
	view.component_types_data = component_types_data
	if view.Blend == nil {
		return
	}

	for _, component := range view.Components {
		log.Println("DEBUG: BlendView update_component_types", component)
		component.Update_component_types(component_types_list)
	}
}

func (view *BlendView) AddComponent() *BlendComponentView {
	if view.Blend == nil {
		return nil
	}

	height := view.ClientHeight()

	component_view := NewBlendComponentView(view)
	view.Components = append(view.Components, component_view)
	view.AutoPanel.Dock(component_view, windigo.Top)

	view.SetSize(view.ClientWidth(), height+component_view.ClientHeight())

	return component_view
}*/

func (view *BlendView) SetFont(font *windigo.Font) {
	for _, component := range view.Components {
		component.SetFont(font)
	}
}

// //TODO grow
// func (view *BlendView) SetLabeledSize(label_width, control_width, height int) {
// total_height := height
// delta_height := height
// 	view.SetSize(label_width+control_width, height)
// 	for _, component := range view.Components {
// 		component.SetSize(label_width+control_width, height)
// 	}
// }

func (view *BlendView) RefreshSize() {
	height := GUI.PRODUCT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	log.Println("Crit: DEBUG: BlendView RefreshSize", height, GUI.PRODUCT_FIELD_HEIGHT, len(view.Components))

	view.SetSize(GUI.OFF_AXIS, height+len(view.Components)*delta_height)
	for _, component := range view.Components {

		component.SetSize(width, delta_height)
		// component.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	}

}
