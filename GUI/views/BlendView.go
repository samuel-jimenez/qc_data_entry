package views

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlendViewer
 *
 */
type BlendViewer interface {
	windigo.Pane
	Get() *blender.ProductBlend
	Update(blend *blender.ProductBlend)
	SetAmount(float64)
	SetFont(font *windigo.Font)
	RefreshSize()
}

/*
 * BlendView
 *
 */
type BlendView struct {
	QCBlendView
	amount_field *NumberEditView
}

func NewBlendView(parent windigo.Controller) *BlendView {
	view := new(BlendView)
	view.QCBlendView = *NewQCBlendView(parent)

	amount_text := "Amount"
	// amount_text := "Quantity"

	//TODO normalize amounts

	view.amount_field = NewNumberEditView(view.panel, amount_text)

	// TODO/ RecipeHeader from  RecipeComponent Component Amount

	view.panel.Dock(view.amount_field, windigo.Left)

	view.amount_field.OnChange().Bind(func(e *windigo.Event) {
		view.SetAmount(view.amount_field.Get())
	})

	return view
}

func (view *BlendView) Update(recipe *blender.ProductRecipe) {
	if view.Recipe == recipe {
		return
	}

	for _, component := range view.Components {
		component.Close()
	}
	view.Components = nil
	view.Recipe = recipe
	log.Println("BlendView Update", view.Recipe)
	if view.Recipe == nil {
		return
	}

	//TODO fix height
	// height := view.ClientHeight()
	height := GUI.EDIT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	for _, component := range view.Recipe.Components {
		log.Println("DEBUG: BlendView update_components", component)

		component_view := NewBlendComponentView(view, &component)
		view.Components = append(view.Components, component_view)
		view.AutoPanel.Dock(component_view, windigo.Top)
		height += delta_height

	}

	view.SetSize(width, height)
	log.Println("Crit: DEBUG: QCBlendView.Update", width, height)

}

func (view *BlendView) SetAmount(amount float64) {
	view.amount_field.Set(amount)
	for _, component := range view.Components {
		component.(*BlendComponentView).SetAmount(amount)
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
