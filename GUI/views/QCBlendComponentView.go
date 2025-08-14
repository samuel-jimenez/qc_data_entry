package views

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/windigo"
)

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
	component_field      *GUI.SearchBox
	component_types_list []string
	component_types_data map[string]blender.BlendComponent

	controls []windigo.Controller

	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
}

//TODO mutliple compents in one
// TODO recipe ingrediant x@100 : a@75 b@25

// TODO no amount_field
func NewQCBlendComponentView(parent windigo.Controller, RecipeComponent *blender.RecipeComponent) *QCBlendComponentView {
	DEL_BUTTON_WIDTH := 20

	log.Println("DEBUG: NewQCBlendComponentView", RecipeComponent)

	view := new(QCBlendComponentView)
	view.recipeComponent = RecipeComponent

	view.component_types_data = make(map[string]blender.BlendComponent)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	// view.component_field = GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)
	// view.component_field = GUI.NewListSearchBoxWithLabels(view.AutoPanel, view.component_types_list)
	view.component_field = GUI.NewListSearchBox(view.AutoPanel)

	//??TODO make memnber?
	var InboundLotView *InboundLotView

	lot_add_button := windigo.NewPushButton(view.AutoPanel)
	lot_add_button.SetText("+")
	// TODO add lot
	lot_add_button.OnClick().Bind(func(e *windigo.Event) {
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

	view.controls = append(view.controls, lot_add_button)

	view.AutoPanel.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.component_field.SetLabeledSize(GUI.SOURCES_LABEL_WIDTH, GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	// toodo todo TODO FIX SIZES
	lot_add_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(view.component_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)
	view.AutoPanel.Dock(lot_add_button, windigo.Left)

	view.component_field.Label().SetText(RecipeComponent.Component_name)
	// view.component_field.SetText(RecipeComponent.Component_name)

	view.Update_component_types()

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
		DB.DB_Select_component_type_product, view.recipeComponent.Component_type_id)

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
	for _, control := range view.controls {
		control.SetFont(font)
	}
}

func (view *QCBlendComponentView) RefreshSize() {

	view.SetSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_HEIGHT)
	view.component_field.SetLabeledSize(GUI.SOURCES_LABEL_WIDTH, GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

}
