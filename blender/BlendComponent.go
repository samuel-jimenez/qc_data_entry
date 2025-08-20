package blender

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
)

type BlendComponent struct {
	RecipeComponent
	Lot_id   int64
	Inboundp bool
	Lot_name,
	Container_name string
	// TODO blendAmount00
	// Component_actual_amount  float64

}

func NewBlendComponent() *BlendComponent {
	return new(BlendComponent)
}

func NewBlendComponentfromSQL(row *sql.Rows) (*BlendComponent, error) {
	blendComponent := NewBlendComponent()
	err := row.Scan(
		// &blendComponent.Component_id,
		&blendComponent.Component_name, &blendComponent.Lot_id, &blendComponent.Lot_name, &blendComponent.Container_name, &blendComponent.Component_amount, &blendComponent.Add_order,
		&blendComponent.Inboundp,
	)
	if err != nil {
		log.Println("Crit: [BlendComponent-NewRecipeComponentfromSQL]: ", err)
	}
	return blendComponent, err
}

func (object *BlendComponent) Save(Product_Lot_id int64) {
	proc_name := "BlendComponent.Save"
	if object.Component_id == 0 {
		log.Println("Crit: BlendComponent-Save uninitialized object.Component_id", object)
		panic(proc_name)
	}
	var Component_id int64

	log.Println("DEBUG: BlendComponent-Save", object)
	if object.Inboundp {
		Component_id = DB.Insel(proc_name, DB.DB_Insert_inbound_blend_component, DB.DB_Select_inbound_blend_component, object.Component_type_id, object.Lot_id)
	} else {
		Component_id = DB.Insel(proc_name, DB.DB_Insert_internal_blend_component, DB.DB_Select_internal_blend_component, object.Component_type_id, object.Lot_id)
	}

	DB.Insert(proc_name, DB.DB_Insert_Product_blend, Product_Lot_id, object.Component_id, Component_id, object.Component_amount)

}

func (object BlendComponent) Text() []string {
	return []string{
		object.Component_name,
		object.Lot_name,
		object.Container_name,
	}
}
