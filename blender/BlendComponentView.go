package blender

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/*
 * BlendComponentViewer
 *
 */
type BlendComponentViewer interface {
	Get() *BlendComponent
	Update(*RecipeComponent)
	Update_component_types(component_types_list []string)
	SetAmount(amount float64)
	// SetFont(font *windigo.Font)
}

/*
 * BlendComponentView
 *
 */
type BlendComponentView struct {
	*windigo.AutoPanel
	RecipeComponent *RecipeComponent
	// BlendComponent       *BlendComponent
	component_field      *GUI.SearchBox
	amount_field         *NumbEditView
	component_types_list []string
	component_types_data map[string]ComponentSource

	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
}

type ComponentSource struct {
	int64
	bool
}

func (componentSource ComponentSource) Unpack() (int64, bool) {
	return componentSource.int64, componentSource.bool
}

// TODO
func NewBlendComponentView(parent *BlendView, RecipeComponent *RecipeComponent) *BlendComponentView {
	DEL_BUTTON_WIDTH := 20

	view := new(BlendComponentView)
	view.RecipeComponent = RecipeComponent

	// view.BlendComponent = NewBlendComponent()
	// view.BlendComponent.Add_order = len(parent.Components)
	// view.BlendComponent.

	view.component_types_data = make(map[string]ComponentSource)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	// view.component_field = GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)
	// view.component_field = GUI.NewListSearchBoxWithLabels(view.AutoPanel, view.component_types_list)
	view.component_field = GUI.NewListSearchBox(view.AutoPanel)

	// cf NumberEditView
	// view.amount_field = windigo.NewEdit(view.AutoPanel)
	view.amount_field = NewNumbEditView(view.AutoPanel)

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

	view.AutoPanel.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.component_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.amount_field.SetSize(GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	lot_add_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(view.component_field, windigo.Left)
	view.AutoPanel.Dock(view.amount_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)
	view.AutoPanel.Dock(lot_add_button, windigo.Left)

	view.component_field.Label().SetText(RecipeComponent.Component_name)
	// view.component_field.SetText(RecipeComponent.Component_name)
	view.amount_field.Set(RecipeComponent.Component_amount)

	view.Update_component_types()

	view.amount_field.OnChange().Bind(func(e *windigo.Event) {
		parent.SetAmount(view.amount_field.Get() / view.RecipeComponent.Component_amount)
	})

	return view
}

func (view *BlendComponentView) Get() *BlendComponent {
	BlendComponent := NewBlendComponent()
	BlendComponent.RecipeComponent = *view.RecipeComponent

	Component_name := view.component_field.Text()
	BlendComponent.Lot_id, BlendComponent.Inboundp = view.component_types_data[Component_name].Unpack()
	if BlendComponent.Lot_id == DB.INVALID_ID {
		//TODO make error
		return nil
	}
	//TODO check for zero amount?
	BlendComponent.Component_amount = view.amount_field.Get()
	log.Println("DEBUG: BlendComponentView update_component_types", BlendComponent, view.component_field.GetSelectedItem(), view.component_field.SelectedItem(), view.component_types_data[view.component_field.Text()])

	return BlendComponent
}

func (view *BlendComponentView) Update(RecipeComponent *RecipeComponent) {
	log.Println("DEBUG: BlendComponentView Update", RecipeComponent)

	if view.RecipeComponent == RecipeComponent || RecipeComponent == nil {
		log.Println("Warn: BlendComponentView Update return", RecipeComponent, view.RecipeComponent)
		return
	}

	view.RecipeComponent = RecipeComponent

}

func (view *BlendComponentView) Update_component_types() {
	// search and fill from
	// RecipeComponent.Component_type_id
	DB.Forall("BlendComponentView",
		func() {
			view.component_types_list = nil
		},
		func(row *sql.Rows) {
			var (
				Inboundp     bool
				id           int64
				product_name string
				lot_name     string
			)

			if err := row.Scan(
				&Inboundp, &id,
				&product_name, &lot_name,
			); err == nil {
				name := fmt.Sprintf("%s - %s", lot_name, product_name)
				view.component_types_data[name] = ComponentSource{id, Inboundp}
				log.Println("DEBUG: update_component_types nsme", id, name)
				view.component_types_list = append(view.component_types_list, name)
			} else {
				log.Printf("error: [%s]: %q\n", "BlendComponentView", err)
			}
		},
		DB.DB_Select_component_type_product, view.RecipeComponent.Component_type_id)

	text := view.component_field.Text()
	log.Println("DEBUG: BlendComponentView update_component_types", text)
	view.component_field.Update(view.component_types_list)
	view.component_field.SetText(text)
}

func (view *BlendComponentView) SetAmount(amount float64) {
	view.amount_field.Set(amount * view.RecipeComponent.Component_amount)
}

/*
 * //TODO

func (view *BlendComponentView) SetFont(font *windigo.Font) {
	// for _, component := range view.Components {
	// 	component.SetFont(font)
	// }
}

// func (view *BlendComponentView) SetLabeledSize(label_width, control_width, height int) {
// 	view.SetSize(label_width+control_width, height)
// }


func (view *BlendComponentView) SetLabeledSize(label_width, control_width, height int) {
	view.SetSize(label_width+control_width, height)

	view.Label().SetSize(label_width, height)

}*/
