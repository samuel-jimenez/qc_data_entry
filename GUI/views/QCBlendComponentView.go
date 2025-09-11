package views

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/blender/blendbound"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/windigo"
)

// c.f. BlendComponentView

/*
 * QCBlendComponentViewer
 *
 */
type QCBlendComponentViewer interface {
	windigo.Controller
	Get() *blender.BlendComponent
	Update_component_types()
	// SetAmount(amount float64)
	SetFont(font *windigo.Font)
	RefreshSize()
}

/*
 * QCBlendComponentView
 *
 */
type QCBlendComponentView struct {
	*windigo.AutoPanel

	recipeComponent *blender.RecipeComponent
	// BlendComponent       *BlendComponent
	Component_name_field  *windigo.Label
	component_field       *GUI.SearchBox
	amount_required_field *NumbEditView

	component_types_list []string
	component_types_data map[string]blender.BlendComponent

	controls []windigo.Controller

	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
}

//TODO multiComp00
// mutliple compents in one
// e.g. recipe ingrediant x@100 : a@75 b@25

// TODO blendAmount00
// we need to capture to desired amount somewhere
// and then enter actual amounts
// TODO no amount_field
// NewBlendComponentView
func NewQCBlendComponentView_from_RecipeComponent(parent windigo.Controller, RecipeComponent *blender.RecipeComponent) *QCBlendComponentView {
	//TODO multiComp00
	// TODO blendAmount00
	// list Component, amount required
	// capture lot, amount
	DEL_BUTTON_WIDTH := 20

	log.Println("DEBUG: NewQCBlendComponentView", RecipeComponent)

	view := new(QCBlendComponentView)
	view.recipeComponent = RecipeComponent

	view.component_types_data = make(map[string]blender.BlendComponent)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	// view.component_field = GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)
	// view.component_field = GUI.NewListSearchBoxWithLabels(view.AutoPanel, view.component_types_list)

	view.Component_name_field = windigo.NewLabel(view.AutoPanel)
	view.amount_required_field = NewNumbEditView(view.AutoPanel)
	view.amount_required_field.SetEnabled(false)

	view.component_field = GUI.NewListSearchBox(view.AutoPanel)

	//??TODO make memnber?
	var InboundLotView *InboundLotView

	// TODO add lot
	// lot_add_button := windigo.NewPushButton(view.AutoPanel)
	// 	lot_add_button.SetText("+")
	// 	lot_add_button.OnClick().Bind(func(e *windigo.Event) {
	// view.controls = append(view.controls, lot_add_button)

	component_add_button := windigo.NewPushButton(view.AutoPanel)
	component_add_button.SetText("+")
	component_add_button.OnClick().Bind(func(e *windigo.Event) {
		//TODO multiComp00
		//TODO mutliple compents in one
		return
		if InboundLotView != nil {
			return
		}
		InboundLotView = NewInboundLotView(view, RecipeComponent)
		InboundLotView.RefreshSize()
		view.AutoPanel.Dock(InboundLotView, windigo.Left)
		InboundLotView.OnClose().Bind(func(e *windigo.Event) {
			InboundLotView = nil
			view.Update_component_types()
		})
	})

	view.controls = append(view.controls, component_add_button)

	// RefreshSize
	view.RefreshSize()
	// toodo todo TODO FIX SIZES
	component_add_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(view.Component_name_field, windigo.Left)
	view.AutoPanel.Dock(view.amount_required_field, windigo.Left)
	view.AutoPanel.Dock(view.component_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)
	view.AutoPanel.Dock(component_add_button, windigo.Left)

	// view.component_field.Label().SetText(RecipeComponent.Component_name)
	view.Component_name_field.SetText(RecipeComponent.Component_name)
	// view.component_field.SetText(RecipeComponent.Component_name)
	// view.amount_required_field.Set(RecipeComponent.Component_amount)
	view.amount_required_field.SetInt(RecipeComponent.Component_amount)

	view.Update_component_types()

	return view
}

func NewQCBlendComponentView_from_BlendComponent(parent windigo.Controller, BlendComponent *blender.BlendComponent) *QCBlendComponentView {
	view := NewQCBlendComponentView_from_RecipeComponent(parent, &BlendComponent.RecipeComponent)
	view.component_field.SetText(fmt.Sprintf("%s : %s", BlendComponent.Lot_name, BlendComponent.Container_name))
	return view
}

func (view *QCBlendComponentView) Get() *blender.BlendComponent {
	Component_name := view.component_field.Text()

	BlendComponent := blender.NewBlendComponent()
	*BlendComponent = view.component_types_data[Component_name]

	if BlendComponent.Lot_id == DB.INVALID_ID {
		//TODO make error
		return nil
	}
	log.Println("DEBUG: BlendComponentView update_component_types", BlendComponent, view.component_field.GetSelectedItem(), view.component_field.SelectedItem(), view.component_types_data[view.component_field.Text()])

	return BlendComponent
}

func (view *QCBlendComponentView) Update_component_types() {
	// search and fill from
	// RecipeComponent.Component_type_id

	DB.Forall_err("BlendComponentView",
		func() {
			view.component_types_list = nil
		},
		func(row *sql.Rows) error {
			var (
				Inboundp       bool
				id             int64
				product_name   string
				lot_name       string
				container_name string
			)
			if err := row.Scan(
				&Inboundp, &id,
				&product_name, &lot_name, &container_name,
			); err != nil {
				return err
			}
			name := fmt.Sprintf("%s : %s - %s", lot_name, container_name, product_name)
			view.component_types_data[name] = blender.BlendComponent{
				RecipeComponent: *view.recipeComponent,
				Lot_id:          id,
				Inboundp:        Inboundp,
				Lot_name:        lot_name,
				Container_name:  container_name,
			}
			view.component_types_list = append(view.component_types_list, name)
			return nil
		},
		DB.DB_Select_component_type_product, view.recipeComponent.Component_type_id,
		product.Status_TESTED,
		blendbound.Status_UNAVAILABLE,
	)

	text := view.component_field.Text()
	view.component_field.Update(view.component_types_list)
	view.component_field.SetText(text)
}

/*

	// proc_name := "RecipeProduct.GetRecipes"

	DB.Forall(proc_name,
		func() {
			object.Recipes = nil
			combo_field.DeleteAllItems()
		},
		func(row *sql.Rows) {
			var (
				recipe_data ProductRecipe
			)

			if err := row.Scan(
				&recipe_data.Recipe_id,
			); err != nil {
				log.Fatal("Crit: [RecipeProduct GetRecipes]: ", proc_name, err)
			}
			log.Println("DEBUG: GetRecipes qc_data", proc_name, recipe_data)
			object.Recipes = append(object.Recipes, &recipe_data)
			combo_field.AddItem(strconv.Itoa(i))

		},
		DB.DB_Select_product_lot_components, object.Product_id)
*/

func (view *QCBlendComponentView) SetFont(font *windigo.Font) {

	view.component_field.SetFont(font)
	view.amount_required_field.SetFont(font)
	view.Component_name_field.SetFont(font)

	for _, control := range view.controls {
		control.SetFont(font)
	}
}

func (view *QCBlendComponentView) RefreshSize() {

	view.SetSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_HEIGHT)

	view.Component_name_field.SetSize(GUI.SOURCES_LABEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.amount_required_field.SetSize(GUI.DATA_SUBFIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.component_field.SetLabeledSize(GUI.OFF_AXIS, GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	// TODO GUI.SOURCES_MARGIN_WIDTH
	// view.component_field.SetLabeledSize(GUI.SOURCES_MARGIN_WIDTH, GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

}
