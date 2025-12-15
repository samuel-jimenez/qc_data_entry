package views

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/windigo"
)

type InboundLotViewer interface {
	windigo.Pane
	Update_products(product_list []string)
	SetFont(font *windigo.Font)
	RefreshSize()
}

type InboundLotView struct {
	*windigo.AutoPanel
	labeled              []windigo.Labelable
	controls             []windigo.Controller
	product_field        *GUI.SearchBox
	component_types_list []string
	component_types_data map[string]int64
}

func NewInboundLotView(parent windigo.Controller, RecipeComponent *blender.RecipeComponent) *InboundLotView {

	inbound_lot_label := "Lot #"
	container_label := "Container"

	view := new(InboundLotView)
	view.AutoPanel = windigo.NewAutoPanel(parent)
	DB.Forall_err("NewInboundLotView",
		func() {
			view.component_types_list = nil
			view.component_types_data = make(map[string]int64)
		},
		func(row *sql.Rows) error {
			var (
				id   int64
				name string
			)
			if err := row.Scan(
				&id, &name,
			); err != nil {
				return err
			}
			view.component_types_data[name] = id
			log.Println("DEBUG: NewInboundLotView nsme", id, name)
			view.component_types_list = append(view.component_types_list, name)
			return nil
		},
		DB.DB_Select_inbound_product_component_type_id, RecipeComponent.Component_type_id)

	view.product_field = GUI.NewListSearchBoxWithLabels(view.AutoPanel, view.component_types_list)
	// inbound_lot_name
	inbound_lot_field := windigo.LabeledEdit_from_new(view.AutoPanel, inbound_lot_label)
	// container_list
	// container_field = GUI.NewListSearchBoxFromQuery(view.AutoPanel, view.component_types_list)
	container_field := GUI.NewSearchBoxFromQuery(view.AutoPanel, DB.DB_Select_container_all)
	// TODO inbound_provider_list

	accept_button := windigo.NewPushButton(view.AutoPanel)
	accept_button.SetText("OK")
	accept_button.SetSize(GUI.ACCEPT_BUTTON_WIDTH, GUI.OFF_AXIS)
	cancel_button := windigo.NewPushButton(view.AutoPanel)
	cancel_button.SetText("Cancel")
	cancel_button.SetSize(GUI.CANCEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.product_field.Label().SetText(RecipeComponent.Component_name)
	container_field.Label().SetText(container_label)

	view.controls = append(view.controls, view.product_field)
	view.controls = append(view.controls, inbound_lot_field)
	view.controls = append(view.controls, container_field)
	view.controls = append(view.controls, accept_button)
	view.controls = append(view.controls, cancel_button)

	view.labeled = append(view.labeled, view.product_field)
	view.labeled = append(view.labeled, inbound_lot_field)
	view.labeled = append(view.labeled, container_field)

	// Dock
	view.AutoPanel.Dock(view.product_field, windigo.Left)

	view.AutoPanel.Dock(inbound_lot_field, windigo.Left)
	view.AutoPanel.Dock(container_field, windigo.Left)
	view.AutoPanel.Dock(accept_button, windigo.Left)
	view.AutoPanel.Dock(cancel_button, windigo.Left)

	accept_button.OnClick().Bind(func(e *windigo.Event) {

		log.Println("DEBUG: NewInboundLotView product_field", view.product_field.Text())
		view.Save(
			inbound_lot_field.Text(),
			view.product_field.Text(),
			// TODO inbound_provider_list
			"SNF",
			container_field.Text(),
		)
		view.OnClose().Fire(nil)
		view.Close()
	})
	cancel_button.OnClick().Bind(func(e *windigo.Event) {
		// view.AutoPanel.Hide()
		view.OnClose().Fire(nil)
		view.Close()
		// view.Exit()

	})
	return view
}
func (view *InboundLotView) Save(inbound_lot_name, inbound_product_name, inbound_provider_name, container_name string) {
	log.Println("DEBUG: [InboundLotView Save] product_field", inbound_lot_name, inbound_product_name, inbound_provider_name, container_name)

	inbound_product_id := view.component_types_data[inbound_product_name]
	if inbound_product_id == DB.INVALID_ID {
		log.Println("Err: [InboundLotView Save] Invalid Product: ", inbound_lot_name, inbound_product_name, inbound_provider_name, container_name)

	}
	inbound_provider_id := DB.Insel("InboundLotView Save insel inbound_provider", DB.DB_Insert_inbound_provider, DB.DB_Select_inbound_provider_id, inbound_provider_name)
	container_id := DB.Insel("InboundLotView Save insel container", DB.DB_Insert_container, DB.DB_Select_container_id, container_name)
	DB.DB_Insert_inbound_lot.Exec(inbound_lot_name, inbound_product_id, inbound_provider_id, container_id)
}

func (view *InboundLotView) Update_products(product_list []string) {
	view.product_field.Update(product_list)
}

func (view *InboundLotView) SetFont(font *windigo.Font) {
	for _, control := range view.controls {
		control.SetFont(font)
	}
}

func (view *InboundLotView) RefreshSize() {
	view.SetSize((GUI.LABEL_WIDTH+GUI.PRODUCT_FIELD_WIDTH)*len(view.labeled)+GUI.ACCEPT_BUTTON_WIDTH+GUI.CANCEL_BUTTON_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	for _, control := range view.labeled {
		control.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

	}

}
