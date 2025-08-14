package blender

import (
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
)

type BlendComponent struct {
	RecipeComponent
	Lot_id   int64
	Inboundp bool
	Lot_name,
	Container_name string
}

func NewBlendComponent() *BlendComponent {
	return new(BlendComponent)
}

func (object *BlendComponent) Save(Product_Lot_id, Recipe_id int64) {
	proc_name := "BlendComponent.Save"

	log.Println("DEBUG: BlendComponent Save", object)
	if object.Inboundp {
		object.Component_id = DB.Insel(proc_name, DB.DB_Insert_inbound_blend_component, DB.DB_Select_inbound_blend_component, object.Component_type_id, object.Lot_id)
	} else {
		object.Component_id = DB.Insel(proc_name, DB.DB_Insert_internal_blend_component, DB.DB_Select_internal_blend_component, object.Component_type_id, object.Lot_id)
	}

	DB.Insert(proc_name, DB.DB_Insert_Product_blend, Product_Lot_id, Recipe_id, object.Component_id, object.Component_amount)

}
