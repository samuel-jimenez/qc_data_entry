package views

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/windigo"
)

/*
 * QCBlendViewer
 *
 */
type QCBlendViewer interface {
	windigo.Pane
	Get() *blender.ProductBlend
	UpdateRecipe(blend *blender.ProductRecipe)
	SetFont(font *windigo.Font)
	RefreshSize()
}

//TODO combine BlendView
/*
 * QCBlendView
 *
 */
type QCBlendView struct {
	*windigo.AutoPanel
	Recipe *blender.ProductRecipe
	Blend  *blender.ProductBlend

	//TODO to fully combine with BlendView this needs to be split out and given its own size fn
	Panel *windigo.AutoPanel

	Components []QCBlendComponentViewer
}

func QCBlendView_from_new(parent windigo.Controller) *QCBlendView {
	view := new(QCBlendView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.Panel = windigo.NewAutoPanel(view.AutoPanel)

	// TODO/ RecipeHeader from  RecipeComponent Component Amount

	view.AutoPanel.Dock(view.Panel, windigo.Top)
	return view
}

func (view *QCBlendView) Get() *blender.ProductBlend {
	if view.Blend != nil {
		return view.Blend
	}
	if view.Recipe == nil {
		return nil
	}
	view.Blend = blender.NewProductBlend_from_Recipe(view.Recipe)
	for _, component := range view.Components {
		blendComponent := component.Get()

		log.Println("DEBUG: QCBlendView-Get", blendComponent)
		if blendComponent == nil {
			// return nil
			// TODO return once qwe get everything nailed down?
			continue
		}
		// blendComponent
		log.Println("DEBUG: QCBlendView-Get-2", blendComponent)
		view.Blend.AddComponent(*blendComponent)
	}
	return view.Blend
}

func (view *QCBlendView) UpdateBlend(blend *blender.ProductBlend) {
	if view.Blend == blend {
		return
	}

	for _, component := range view.Components {
		component.Close()
	}
	view.Components = nil
	view.Blend = blend
	log.Println("QCBlendView-UpdateBlend", view.Blend)
	if view.Blend == nil {
		return
	}

	height := GUI.EDIT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	for _, component := range view.Blend.Components {

		component_view := NewQCBlendComponentView_from_BlendComponent(view, &component)
		view.Components = append(view.Components, component_view)
		view.AutoPanel.Dock(component_view, windigo.Top)
		height += delta_height

	}

	view.SetSize(width, height)
}

func (view *QCBlendView) UpdateRecipe(recipe *blender.ProductRecipe) {
	SAMPLE_SIZE := 100.

	// TODO blendAmount00
	// we need to capture to desired amount somewhere
	// and then enter actual amounts

	// [compiler] (MismatchedTypes) invalid operation: component.Component_amount *= SAMPLE_SIZE (mismatched types float64 and int)
	// go == shit

	if view.Recipe == recipe {
		log.Println("TRACE: QCBlendView-UpdateRecipe-return", view.Recipe, recipe)

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
	view.Blend = nil
	view.Recipe = recipe
	if view.Recipe == nil {
		return
	}

	//TODO fix height
	// height := view.ClientHeight()
	height := GUI.EDIT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT
	width := view.ClientWidth()

	for _, component := range view.Recipe.Components {
		// component.Component_amount *= SAMPLE_SIZE
		// [compiler] (MismatchedTypes) invalid operation: component.Component_amount *= SAMPLE_SIZE (mismatched types float64 and int)
		component.Component_amount *= SAMPLE_SIZE

		component_view := NewQCBlendComponentView_from_RecipeComponent(view, &component)
		view.Components = append(view.Components, component_view)
		view.AutoPanel.Dock(component_view, windigo.Top)
		height += delta_height

	}

	view.SetSize(width, height)
}

func (view *QCBlendView) SetFont(font *windigo.Font) {
	for _, component := range view.Components {
		component.SetFont(font)
	}
}

func (view *QCBlendView) RefreshSize() {
	height := GUI.EDIT_FIELD_HEIGHT
	delta_height := GUI.PRODUCT_FIELD_HEIGHT

	view.SetSize(GUI.OFF_AXIS, height+len(view.Components)*delta_height)
	view.Panel.SetSize(GUI.TOP_PANEL_WIDTH, height)

	for _, component := range view.Components {
		component.RefreshSize()
	}

}
