package product

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

func Release_testing_lot(lot_number string) error {
	proc_name := "product.Release_testing_lot"
	err := DB.Update(proc_name,
		//TODO only the BSQL?
		DB.DB_Update_lot_list__component_status, lot_number, Status_SHIPPED,
	)
	if err != nil {
		log.Println("error[]%S]:", proc_name, err)
	}
	return err
}

type BaseProduct struct {
	Product_name             string `json:"product_name"`
	Lot_number               string `json:"lot_number"`
	Sample_point             string
	Tester                   nullable.NullString `json:"Tester"`
	Visual                   bool
	Product_id               int64
	Lot_id                   int64
	Product_Lot_id           int64
	Product_name_customer_id nullable.NullInt64
	Product_name_customer    string `json:"customer_product_name"`
	Blend                    *blender.ProductBlend
	Valid                    bool
}

func (product BaseProduct) Base() BaseProduct {
	return product
}

func (product *BaseProduct) ResetLot() {
	product.Lot_number = ""
	product.Product_Lot_id = DB.DEFAULT_LOT_ID
	product.Lot_id = DB.DEFAULT_LOT_ID
	product.Product_name_customer = ""
	product.Product_name_customer_id = nullable.NullInt64Default()
	product.Blend = nil
}

//TODO product.get_coa_name()

func (product BaseProduct) get_base_filename(extension string) string {
	// if (product.Sample_point.Valid) {
	if product.Sample_point != "" {
		return fmt.Sprintf("%s-%s-%s.%s", strings.ToUpper(product.Lot_number), product.Sample_point, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.Product_name)), " ", "_"), extension)
	}

	return fmt.Sprintf("%s-%s.%s", strings.ToUpper(product.Lot_number), strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.Product_name)), " ", "_"), extension)
}

func (product BaseProduct) get_pdf_name() string {
	return fmt.Sprintf("%s/%s", config.LABEL_PATH, product.get_base_filename("pdf"))
}

func (product Product) get_storage_pdf_name(qc_sample_storage_name string) string {
	return fmt.Sprintf("%s/%s.%s", config.LABEL_PATH, qc_sample_storage_name, "pdf")
}

func (product *BaseProduct) Insel_product_self() *BaseProduct {
	product.Product_id = DB.Insel_product_id(product.Product_name)
	return product

}

func (product *BaseProduct) Update_lot(lot_number, product_name_customer string) *BaseProduct {
	product.ResetLot()
	if product.Product_id == DB.INVALID_ID {
		return product
	}
	if product_name_customer != "" {
		product.Product_name_customer = product_name_customer
		product.Product_name_customer_id = nullable.NewNullInt64(DB.Insel_product_name_customer(product.Product_name_customer, product.Product_id))
	}

	if lot_number != "" {
		// we cannot guarantee uniqueness of batch numbers
		product.Lot_number, product.Lot_id, product.Product_Lot_id = DB.Select_product_lot_name(lot_number, product.Product_id)

		DB.DB_Update_lot_customer.Exec(product.Product_name_customer_id, product.Product_Lot_id)
	}
	return product
}

//nned to export:
// 			Product_name             string `json:"product_name"`
// 	Lot_number               string `json:"lot_number"`
// 		Product_name_customer    string `json:"customer_product_name"`
//
// *	Sample_point             string
// *	Tester                   nullable.NullString `json:"Tester"`
// *		PH          nullable.NullFloat64
// *	SG          nullable.NullFloat64
//* 	String_test nullable.NullInt64
// *	Viscosity   nullable.NullInt64
// 	Visual                   bool
//
// 	Blend                    *blender.ProductBlend
// ??TODO
// 	glut
// 	quat
//

// need to save:
// Lot_id
// func (product *BaseProduct) Update_inbound_lot(lot_number, product_name_customer string) *BaseProduct {
func (base_product *BaseProduct) Update_testing_lot(lot_number string) {
	base_product.ResetLot()
	proc_name := "BaseProduct.Update_testing_lot"

	var Product_name_customer nullable.NullString
	if err := DB.Select_Error(proc_name,
		DB.DB_Select_product_lot_list_name.QueryRow(lot_number),
		&base_product.Lot_id,
		&base_product.Product_name, &base_product.Lot_number, &Product_name_customer,
	); err != nil {
		log.Println("Crit: ", proc_name, base_product, err)
		base_product.ResetLot()
		return
	}
	base_product.Product_name_customer = Product_name_customer.String

	Blend := blender.ProductBlend_from_new()
	base_product.Blend = Blend

	// TODO blend012 ensure doesn't break show_fr()

	// base_product.Insel_product_self()

	DB.Forall_exit(proc_name,
		func() {
		},
		func(row *sql.Rows) error {
			blendComponent, err := blender.BlendComponent_from_SQL(row)
			if err != nil {
				return err
			}
			Blend.AddComponent(*blendComponent)
			return nil
		},
		DB.DB_Select_product_lot_list_sources, base_product.Lot_number)
}

func (base_product *BaseProduct) Update_testing_tttToday_Junior() {
	proc_name := "BaseProduct.Update_testing_tttToday_Junior"

	operations_group := "BSQL"
	lot_number, err := blender.Next_Lot_Number(operations_group)
	if err != nil {
		return
	}
	if err = DB.Update(proc_name,
		DB.DB_Update_lot_list_name,
		base_product.Lot_id,
		lot_number,
	); err != nil {
		return
	}
	base_product.Update_testing_lot(lot_number)
}

func (product *BaseProduct) SetTester(Tester string) {
	product.Tester = nullable.NewNullString(Tester)

}

func (product *BaseProduct) SetBlend(Blend *blender.ProductBlend) {

	proc_name := "BaseProduct.SetBlend"

	log.Println("Debug: ", proc_name, product.Blend, Blend, product)
	if product.Blend == Blend {
		log.Println("TRACE: Debug: retrurning ", proc_name, product.Blend, Blend, product)

		return
	}
	product.Blend = Blend
	if Blend == nil {
		return
	}
	DB.Update(proc_name,
		DB.DB_Update_lot_recipe,
		Blend.Recipe_id,
		product.Product_Lot_id,
	)
}

func (product BaseProduct) SaveBlend() {
	//TODO make sure this is the only time it is saved
	if product.Blend != nil {
		product.Blend.Save(product.Product_Lot_id)
	}
}
