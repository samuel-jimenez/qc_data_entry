package blender

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/windigo"
)

type BlendViewer interface {
	windigo.Pane
	// Get() *ProductBlend
	// Update(blend *ProductBlend)
	// AddComponent()
	SetFont(font *windigo.Font)
	RefreshSize()
}

type BlendView struct {
	*windigo.AutoPanel
	Recipe *ProductRecipe
	Blend  *ProductBlend
	panel  *windigo.AutoPanel

	amount_field *views.NumberEditView
	Components   []*BlendComponentView
	// Product_id int64
	// Blend_id  int64
}

func NewBlendView(parent windigo.Controller) *BlendView {
	view := new(BlendView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	amount_text := "Amount"
	// amount_text := "Quantity"

	//TODO normalize amounts

	view.panel = windigo.NewAutoPanel(view.AutoPanel)

	view.amount_field = views.NewNumberEditView(view.panel, amount_text)

	// TODO/ RecipeHeader from  RecipeComponent Component Amount

	view.panel.Dock(view.amount_field, windigo.Left)
	view.AutoPanel.Dock(view.panel, windigo.Top)

	view.amount_field.OnChange().Bind(func(e *windigo.Event) {
		view.SetAmount(view.amount_field.Get())
	})

	return view
}

func (view *BlendView) Get() *ProductBlend {
	if view.Recipe == nil {
		return nil
	}
	view.Blend = NewProductBlendFromRecipe(view.Recipe)
	for _, component := range view.Components {
		blendComponent := component.Get()
		log.Println("DEBUG: BlendView Get", blendComponent)
		if blendComponent == nil {
			return nil
		}
		log.Println("DEBUG: BlendView Get", blendComponent)
		view.Blend.AddComponent(*blendComponent)
	}
	return view.Blend
}

// func (view *BlendView) Update(blend *ProductBlend) {
// 	if view.Blend == blend {
// 		return
// 	}
//
// 	// height := GUI.EDIT_FIELD_HEIGHT
// 	// delta_height := GUI.EDIT_FIELD_HEIGHT
// 	// 	total_height := height
// 	// delta_height := height
// 	// 	height := GUI.EDIT_FIELD_HEIGHT
// 	// delta_height := GUI.EDIT_FIELD_HEIGHT
//
// 	for _, component := range view.Components {
// 		component.Close()
// 	}
// 	view.Components = nil
// 	view.Blend = blend
// 	log.Println("BlendView Update", view.Blend)
// 	if view.Blend == nil {
// 		return
// 	}
// 	//TODO
//
// 	// for _, component := range view.Blend.Components {
// 	// 	log.Println("DEBUG: BlendView update_components", component)
// 	// 	// view.AddComponent()
// 	// 	//TODO
// 	// 	component_view := view.AddComponent()
// 	// 	component_view.Update(&component)
// 	// }
//
// }

func (view *BlendView) Update(recipe *ProductRecipe) {
	if view.Recipe == recipe {
		return
	}

	// height := GUI.EDIT_FIELD_HEIGHT
	// delta_height := GUI.EDIT_FIELD_HEIGHT
	// 	total_height := height
	// delta_height := height
	// 	height := GUI.EDIT_FIELD_HEIGHT
	// delta_height := GUI.EDIT_FIELD_HEIGHT

	for _, component := range view.Components {
		component.Close()
	}
	view.Components = nil
	view.Recipe = recipe
	log.Println("BlendView Update", view.Recipe)
	if view.Recipe == nil {
		return
	}

	for _, component := range view.Recipe.Components {
		log.Println("DEBUG: BlendView update_components", component)
		view.AddComponent(&component)
	}

}

func (view *BlendView) AddComponent(RecipeComponent *RecipeComponent) *BlendComponentView {
	log.Println("DEBUG: BlendView AddComponent", RecipeComponent)

	if view.Recipe == nil {
		return nil
	}

	height := view.ClientHeight()

	component_view := NewBlendComponentView(view, RecipeComponent)
	view.Components = append(view.Components, component_view)
	view.AutoPanel.Dock(component_view, windigo.Top)

	view.SetSize(view.ClientWidth(), height+component_view.ClientHeight())

	return component_view
}

func (view *BlendView) SetAmount(amount float64) {
	view.amount_field.Set(amount)
	for _, component := range view.Components {
		component.SetAmount(amount)
	}
}

func (view *BlendView) SetFont(font *windigo.Font) {
	view.amount_field.SetFont(font)
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
	// height := GUI.PRODUCT_FIELD_HEIGHT+GUI.ERROR_MARGIN
	height := GUI.EDIT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	log.Println("Crit: DEBUG: BlendView RefreshSize", height, GUI.PRODUCT_FIELD_HEIGHT, len(view.Components))

	view.SetSize(GUI.OFF_AXIS, height+len(view.Components)*delta_height)
	view.panel.SetSize(GUI.TOP_PANEL_WIDTH, height)
	view.amount_field.SetLabeledSize(GUI.LABEL_WIDTH-GUI.ERROR_MARGIN, GUI.PRODUCT_FIELD_WIDTH, height)

	for _, component := range view.Components {

		component.SetSize(width, delta_height)
		// component.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, delta_height)
	}

}
