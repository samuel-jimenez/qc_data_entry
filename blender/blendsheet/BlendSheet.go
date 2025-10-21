package blendsheet

import (
	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

/*
 * BlendSheetable
 *
 */
type BlendSheetable interface {
}

/*
 * BlendSheet
 *
 */
type BlendSheet struct {
	// blender.BlendProduct //TODO need??
	// TODO maybe ProductBlend instead with BlendProduct.Save or Next_Lot_Number or equiv
	blender.ProductBlend
	// TODO   BlendProduct.Save

	// product.BaseProduct
	// product.QCProduct (Product_id)
	// Select_product_details
	// *QCProduct

	//TODO ??
	// 	we need: #
	//
	// 	type BaseProduct struct {
	// #		Product_name             string `json:"product_name"`
	// #		Lot_number               string `json:"lot_number"`
	// 		Sample_point             string
	// 		Tester                   nullable.NullString `json:"Tester"`
	// 		Visual                   bool
	// +		Product_id               int64
	// +		Lot_id                   int64
	// +		Product_Lot_id           int64
	// +		Product_name_customer_id nullable.NullInt64
	// #		Product_name_customer    string `json:"customer_product_name"`
	// #		Blend                    *blender.ProductBlend
	// 		Valid                    bool
	// 	}
	//
	// plus maybe some ids +

	// TODO: split out a new smaller type, then rename and subclass BaseProduct

	// 	PRODUCT INFO:
	// 	PRODUCT
	Product_name string
	Lot_number   string
	// 	CUSTNANE

	Product_id     int64
	Lot_id         int64
	Product_Lot_id int64

	Product_name_customer_id nullable.NullInt64
	Product_name_customer    string

	// Blend                    *blender.ProductBlend

	// TODO
	// 	WHO
	// Operators []string
	Operators string

	// 	WhERE
	// 	BATCH VESSEL
	Vessel   string
	Capacity string
	// Capacity float64

	// Capacity,
	Strap, HeelVolume, HeelMass float64
	Total_Component             *blender.BlendComponent

	Tag,
	Seal string

	// 	WHEN

	//empty:
	// TODO
	// 	TAG SEAL
	// 	FINAL
	//
	//
	// 	QC INFO
	//
	// 	TEST SPEC RESULT
	// 	TESTER
	//
	//

	// 	BLEND INFO
	//
	// 	COMPONENTE AMOUNTS
	// 	//ProductBlend, BlendProduct???
	//
	//

	// 	PROCEDURE
	// 	// TODO

	// 	QR
	// 	// TODO
	//

	// Operators []string
	// Capacity
	// Strap

	// Tag
	// Seal

	// Final Strap
	// /density?
	// TODO cerate blend sheet type
	// TODO print blend sheet
	// TODO gen QR
	// TODO cap:
	// operators
	// TAG
	// SEAL

	//frorm BlendVessel:
	// Vessel_field, Capacity_field *GUI.SearchBox
	// Strap_field

	// TODO Sheet-03988
}

//
// TODO Sheet-03988
/*/
new tables:
qc_ blend
masurement
spec
coa-p

As Is
≤ 60 seconds

String Test at 0.5gpt	60 ≤ seconds
Neat Viscosity	As Is'

procedure
product id
step-no
step_text


*/
func BlendSheet_from_new(ProductBlend *blender.ProductBlend, Product_name, operations_group, Product_name_customer string, Product_id int64, BlendVessel blender.BlendWessel, Total_Component *blender.BlendComponent,

	// func BlendSheet_from_new(ProductBlend *blender.ProductBlend, Product_name, operations_group, Product_name_customer string, Product_id int64, BlendVessel *BlendWessel, Total_Component *blender.BlendComponent,
	Tag, Seal, Operators string) *BlendSheet {

	proc_name := "BlendSheet_from_new"

	object := new(BlendSheet)
	object.ProductBlend = *ProductBlend
	object.Product_name = Product_name
	object.Product_id = Product_id

	// BlendProduct.Save :
	if Product_name_customer != "" {
		object.Product_name_customer = Product_name_customer
		object.Product_name_customer_id = nullable.NewNullInt64(DB.Insel_product_name_customer(object.Product_name_customer, object.Product_id))
	}

	object.Lot_id, object.Lot_number, _ = blender.Next_Both_Lot_Id_Number(operations_group)
	// object.Lot_number, _ = blender.Next_Lot_Number(operations_group)
	object.Product_Lot_id = DB.Insert(proc_name, DB.DB_Insert_blend_lot, object.Lot_id, object.Product_id, object.Product_name_customer_id, object.Recipe_id)

	object.ProductBlend.Save(object.Product_Lot_id)

	object.Vessel,
		object.Capacity,
		object.Strap, object.HeelVolume, object.HeelMass = BlendVessel.Get()
	object.Total_Component = Total_Component

	object.Tag, object.Seal,
		object.Operators =
		Tag, Seal,
		Operators

	return object
}
