package product

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

func Release_testing_lot(lot_number string) error {
	proc_name := "product.Release_testing_lot"
	err := DB.Update(proc_name,
		// TODO only the BSQL?
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

func (product *BaseProduct) ResetLot() {
	product.Lot_number = ""
	product.Product_Lot_id = DB.DEFAULT_LOT_ID
	product.Lot_id = DB.DEFAULT_LOT_ID
	product.Product_name_customer = ""
	product.Product_name_customer_id = nullable.NullInt64Default()
	product.Blend = nil
}

// TODO product.get_coa_name()
func (product BaseProduct) get_base_filename(extension string) string {
	// if (product.Sample_point.Valid) {
	Product_name := strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.Product_name)), " ", "_")
	Lot_number := strings.ReplaceAll(strings.ToUpper(product.Lot_number), ":", "_")
	if product.Sample_point != "" {
		return fmt.Sprintf("%s-%s-%s.%s", Lot_number, product.Sample_point, Product_name, extension)
	}

	return fmt.Sprintf("%s-%s.%s", Lot_number, Product_name, extension)
}

func (product BaseProduct) get_pdf_name() string {
	return fmt.Sprintf("%s/%s", config.LABEL_PATH, product.get_base_filename("pdf"))
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
	proc_name := "BaseProduct-SetBlend"

	if product.Blend == Blend {
		return
	}

	product.Blend = Blend
	if product.Blend == nil ||
		product.Product_Lot_id == DB.DEFAULT_LOT_ID ||
		product.Blend.Recipe_id == DB.INVALID_ID {
		return
	}

	DB.Update(proc_name,
		DB.DB_Update_lot_recipe,
		Blend.Recipe_id,
		product.Product_Lot_id,
	)
}

func (product BaseProduct) SaveBlend() {
	if product.Blend == nil ||
		product.Product_Lot_id == DB.DEFAULT_LOT_ID ||
		product.Blend.Recipe_id == DB.INVALID_ID {
		return
	}

	product.Blend.Save(product.Product_Lot_id)
}

// TODO product.NewBin
// add button to call to account for unlogged samples

/*
 * NewStorageBin
 *
 * Print label for storage bin
 *
 * Returns next storage bin
 *
 */
func (measured_product BaseProduct) NewStorageBin() int {
	var qc_sample_storage_id int
	// get data for printing bin label, new bin creation
	proc_name := "BaseProduct-NewStorageBin"

	var (
		product_sample_storage_id, qc_sample_storage_offset,
		start_date, end_date, retain_storage_duration int
		qc_sample_storage_name, product_moniker_name string
	)
	if err := DB.Select_Error(proc_name,
		DB.DB_Select_gen_product_sample_storage.QueryRow(
			measured_product.Product_id,
		),
		&product_sample_storage_id, &product_moniker_name, &qc_sample_storage_offset, &qc_sample_storage_name,
		&start_date, &end_date, &retain_storage_duration,
	); err != nil {
		log.Println("Crit: ", proc_name, measured_product, err)
		panic(err)
	}

	// print
	PrintStorageBin(qc_sample_storage_name, product_moniker_name, start_date, end_date, retain_storage_duration)

	// product_moniker_id == product_sample_storage_id. update if this is no longer true
	qc_sample_storage_id, qc_sample_storage_offset = InsertStorageBin(product_sample_storage_id, qc_sample_storage_offset, product_moniker_name)

	// update qc_sample_storage_offset
	proc_name = "BaseProduct-NewStorageBin.Update"
	DB.Update(proc_name,
		DB.DB_Update_product_sample_storage_qc_sample,
		product_sample_storage_id,
		qc_sample_storage_id, qc_sample_storage_offset)
	return qc_sample_storage_id
}

//TODO move
/*
 * PrintStorageBin
 *
 * Gen new based on qc_sample_storage_offset,
 * product_moniker_id or product_sample_storage_id
 *
 * Print label for storage bin
 *
 * Returns next storage bin, updated qc_sample_storage_offset
 *
 */
func PrintStorageBin(qc_sample_storage_name, product_moniker_name string, start_date_, end_date_, retain_storage_duration int,
) {
	start_date := time.Unix(0, int64(start_date_))
	end_date := time.Unix(0, int64(end_date_))
	retain_date := end_date.AddDate(0, retain_storage_duration, 0)
	// + retain_storage_duration*time.Month

	// print
	PrintOldStorage(qc_sample_storage_name, product_moniker_name, &start_date, &end_date, &retain_date)
}

/*
 * RegenStorageBin
 *
 * Print label for storage bin
 *
 *
 */
func RegenStorageBin(qc_sample_storage_name string) {
	// get data for printing bin label
	proc_name := "BaseProduct-RePrintStorageBin"

	var (
		start_date, end_date, retain_storage_duration int
		product_moniker_name                          string
	)
	if err := DB.Select_Error(proc_name,
		DB.DB_Select_reprint_product_sample_storage.QueryRow(
			qc_sample_storage_name,
		),
		&product_moniker_name, &qc_sample_storage_name,
		&start_date, &end_date, &retain_storage_duration,
	); err != nil {
		log.Println("Crit: ", proc_name, err)
		panic(err)
	}

	// print
	PrintStorageBin(qc_sample_storage_name, product_moniker_name, start_date, end_date, retain_storage_duration)
}

//TODO move
/*
 * InsertStorageBin
 *
 * Gen new based on qc_sample_storage_offset,
 * product_moniker_id or product_sample_storage_id
 *
 * Print label for storage bin
 *
 * Returns next storage bin, updated qc_sample_storage_offset
 *
 */
func InsertStorageBin(product_moniker_id, qc_sample_storage_offset int, product_moniker_name string) (int, int) {
	var (
		qc_sample_storage_id int
		date                 time.Time
	)
	proc_name := "BaseProduct-InsertStorageBin"
	location := 0
	qc_sample_storage_offset += 1
	qc_sample_storage_name := fmt.Sprintf("%02d-%03d-%04d", location, product_moniker_id, qc_sample_storage_offset)

	if err := DB.Select_Error(proc_name,
		DB.DB_Insert_sample_storage.QueryRow(
			qc_sample_storage_name, product_moniker_id,
		),
		&qc_sample_storage_id,
	); err != nil {
		log.Println("Crit: ", proc_name, err)
		panic(err)
	}
	PrintNewStorage(qc_sample_storage_name, product_moniker_name, &date, &date, &date)
	return qc_sample_storage_id, qc_sample_storage_offset
}

/*
 * GetStorage
 *
 * Check storage
 *
 * Returns next storage bin with capacity for product
 * products come in sequentially
 * keep entire batches together
 *
 */
func (measured_product BaseProduct) GetStorage(numSamples int) int {
	proc_name := "BaseProduct-GetStorage"

	var qc_sample_storage_id, storage_capacity int
	if err := DB.Select_Error(proc_name,
		DB.DB_Select_product_sample_storage_capacity.QueryRow(
			measured_product.Product_id,
		),
		&qc_sample_storage_id, &storage_capacity,
	); err != nil {
		log.Println("Crit: ", proc_name, measured_product, err)
		panic(err)
	}

	// Check storage
	if storage_capacity < numSamples {
		qc_sample_storage_id = measured_product.NewStorageBin()
	}
	return qc_sample_storage_id
}

/*
 * CheckStorage
 *
 * Check storage bin, print label if full
 *
 */
func (measured_product BaseProduct) CheckStorage() {
	measured_product.GetStorage(1)
}

func (measured_product *BaseProduct) format_sample() {
	if measured_product.Product_name_customer != "" {
		measured_product.Product_name = measured_product.Product_name_customer
	}
	measured_product.Sample_point = ""
}

func (measured_product BaseProduct) Reprint() {
	Print_PDF(measured_product.get_pdf_name())
}

func (measured_product BaseProduct) Reprint_sample() {
	measured_product.format_sample()
	log.Println("DEBUG: Reprint_sample formatted", measured_product)
	measured_product.Reprint()
}
