package views

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/blender/inbound"
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

	Recipe_Component *blender.RecipeComponent
	// BlendComponent       *BlendComponent
	Component_name_field  *windigo.Label
	Component_field       *GUI.SearchBox
	Amount_required_field *GUI.NumbEditView

	// amount_field   *GUI.NumbEditView

	component_types_list []string
	Component_types_data map[string]blender.BlendComponent

	Controls []windigo.Controller

	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
}

// TODO multiComp00
// mutliple compents in one
// e.g. recipe ingrediant x@100 : a@75 b@25

// TODO blendAmount00
// we need to capture to desired amount somewhere
// and then enter actual amounts
// TODO no amount_field
// NewBlendComponentView

func New_Bare_QCBlendComponentView_from_RecipeComponent_com(parent windigo.Controller, recipeComponent *blender.RecipeComponent) *QCBlendComponentView {
	// TODO multiComp00
	// TODO blendAmount00
	// list Component, amount required
	// capture lot, amount

	view := new(QCBlendComponentView)
	view.Recipe_Component = recipeComponent

	view.Component_types_data = make(map[string]blender.BlendComponent)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	// view.component_field = GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)
	// view.component_field = GUI.NewListSearchBoxWithLabels(view.AutoPanel, view.component_types_list)

	view.Component_name_field = windigo.NewLabel(view.AutoPanel)
	view.Amount_required_field = GUI.NumbEditView_from_new(view.AutoPanel)

	view.Component_field = GUI.NewListSearchBox(view.AutoPanel)
	// view.amount_field = windigo.NewEdit(view.AutoPanel)

	view.AutoPanel.Dock(view.Component_name_field, windigo.Left)
	view.AutoPanel.Dock(view.Amount_required_field, windigo.Left)
	view.AutoPanel.Dock(view.Component_field, windigo.Left)
	// view.AutoPanel.Dock(order_field, windigo.Left)

	// view.component_field.Label().SetText(RecipeComponent.Component_name)
	view.Component_name_field.SetText(recipeComponent.Component_name)
	// view.component_field.SetText(RecipeComponent.Component_name)
	// view.amount_required_field.Set(RecipeComponent.Component_amount)

	view.Update_component_types()

	return view
}

func NewQCBlendComponentView_from_RecipeComponent(parent windigo.Controller, recipeComponent *blender.RecipeComponent) *QCBlendComponentView {
	view := New_Bare_QCBlendComponentView_from_RecipeComponent_com(parent, recipeComponent)
	view.Amount_required_field.SetEnabled(false)
	view.Amount_required_field.SetInt(recipeComponent.Component_amount)

	// RefreshSize
	view.RefreshSize()

	DEL_BUTTON_WIDTH := 20

	component_add_button := windigo.NewPushButton(view.AutoPanel)
	component_add_button.SetText("+")
	component_add_button.OnClick().Bind(func(e *windigo.Event) {
		// TODO multiComp00
		// TODO mutliple compents in one

	})

	view.Controls = append(view.Controls, component_add_button)

	// toodo todo TODO FIX SIZES
	component_add_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(component_add_button, windigo.Left)

	return view
}

func NewQCBlendComponentView_from_BlendComponent(parent windigo.Controller, BlendComponent *blender.BlendComponent) *QCBlendComponentView {
	view := NewQCBlendComponentView_from_RecipeComponent(parent, &BlendComponent.RecipeComponent)
	view.Component_field.SetText(fmt.Sprintf("%s : %s", BlendComponent.Lot_name, BlendComponent.Container_name))
	return view
}

func (view *QCBlendComponentView) Get() *blender.BlendComponent {
	Component_name := view.Component_field.Text()

	BlendComponent := blender.BlendComponent_from_new()
	*BlendComponent = view.Component_types_data[Component_name]

	if BlendComponent.Lot_id == DB.INVALID_ID {
		// TODO make error
		return nil
	}
	log.Println("DEBUG: QCBlendComponentView update_component_types", BlendComponent, view.Component_field.GetSelectedItem(), view.Component_field.SelectedItem(), view.Component_types_data[view.Component_field.Text()])

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
			view.Component_types_data[name] = blender.BlendComponent{
				RecipeComponent: *view.Recipe_Component,
				Lot_id:          id,
				Inboundp:        Inboundp,
				Lot_name:        lot_name,
				Container_name:  container_name,
			}
			view.component_types_list = append(view.component_types_list, name)
			return nil
		},
		DB.DB_Select_component_type_product, view.Recipe_Component.Component_type_id,
		product.Status_TESTED,
		inbound.Status_UNAVAILABLE,
	)

	text := view.Component_field.Text()
	view.Component_field.Update(view.component_types_list)
	view.Component_field.SetText(text)
}

func (view *QCBlendComponentView) SetFont(font *windigo.Font) {
	view.Component_field.SetFont(font)
	view.Amount_required_field.SetFont(font)
	view.Component_name_field.SetFont(font)

	for _, control := range view.Controls {
		control.SetFont(font)
	}
}

func (view *QCBlendComponentView) RefreshSize() {
	view.SetSize(GUI.OFF_AXIS, GUI.PRODUCT_FIELD_HEIGHT)

	view.Component_name_field.SetSize(GUI.SOURCES_LABEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Amount_required_field.SetSize(GUI.DATA_SUBFIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Component_field.SetLabeledSize(GUI.OFF_AXIS, GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	// TODO GUI.SOURCES_MARGIN_WIDTH
	// view.component_field.SetLabeledSize(GUI.SOURCES_MARGIN_WIDTH, GUI.SOURCES_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
}
