package blender

import "github.com/samuel-jimenez/qc_data_entry/DB"

type ComponentType struct {
	Component_name string
	Component_id   int64
}

func ComponentType_from_new(Component_name string) *ComponentType {
	component := new(ComponentType)
	component.Component_name = Component_name
	component.Insel()
	return component
}

func (object *ComponentType) Insel() {
	object.Component_id = DB.Insel("ComponentType.Insel", DB.DB_Insert_component_types, DB.DB_Select_name_component_types, object.Component_name)
}
func (object *ComponentType) AddProduct(Product_id int64) {
	DB.DB_Insert_internal_product_component_type.Exec(object.Component_id, Product_id)
}

// ??TODO
func (object *ComponentType) AddInbound() {
	//TODO maybe pick name
	Product_id := DB.Insert("ComponentType.AddInbound", DB.DB_Insert_inbound_product, object.Component_name)
	DB.DB_Insert_inbound_product_component_type.Exec(object.Component_id, Product_id)
}
