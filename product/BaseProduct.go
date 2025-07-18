package product

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

type BaseProduct struct {
	Product_name             string `json:"product_name"`
	Lot_number               string `json:"lot_number"`
	Sample_point             string
	Visual                   bool
	Product_id               int64
	Lot_id                   int64
	Product_name_customer_id nullable.NullInt64
	Product_name_customer    string `json:"customer_product_name"`
}

func (product BaseProduct) Base() BaseProduct {
	return product
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

func (product BaseProduct) get_json_filename(path string, base_name string) string {
	return fmt.Sprintf("%s/%d-%s", path, time.Now().UTC().UnixNano(), base_name)
}

func (product BaseProduct) get_json_names() []string {
	var json_names []string
	base_name := product.get_base_filename("json")

	for _, JSON_PATH := range config.JSON_PATHS {

		json_names = append(json_names, product.get_json_filename(JSON_PATH, base_name))
	}
	return json_names
}

func (product *BaseProduct) Insel_product_self() *BaseProduct {
	product.Product_id = DB.Insel_product_id(product.Product_name)
	return product

}

func (product *BaseProduct) Update_lot(lot_number, product_name_customer string) *BaseProduct {

	log.Println("Debug: Update_lot Product_id", product.Product_id, lot_number, product_name_customer)
	if product_name_customer != "" && product.Product_id != DB.INVALID_ID {
		log.Println("Debug: Update_lot product_name_customer", product.Product_name_customer)
		product.Product_name_customer = product_name_customer
		product.Product_name_customer_id = nullable.NewNullInt64(DB.Insel_product_name_customer(product.Product_name_customer, product.Product_id))
	} else {
		product.Product_name_customer = ""
		product.Product_name_customer_id = nullable.NullInt64Default()
	}
	log.Println("Debug: Update_lot Product_name_customer_id", product.Product_name_customer, product.Product_name_customer_id)

	if lot_number != "" && product.Product_id != DB.INVALID_ID {
		product.Lot_number = lot_number
		product.Lot_id = DB.Insel_lot_id(product.Lot_number, product.Product_id)
		log.Println("Debug: Update_lot Lot_id", product.Lot_number, product.Lot_id)

		DB.DB_Update_lot_customer.Exec(product.Product_name_customer_id, product.Lot_id)
	} else {
		product.Lot_number = ""
		product.Lot_id = DB.DEFAULT_LOT_ID
	}
	return product
}

func (product *BaseProduct) Insel_lot_self() *BaseProduct {
	product.Insel_product_name_customer()
	return product

}

func (product *BaseProduct) Insel_product_name_customer() *BaseProduct {
	return product

}

func (product BaseProduct) Insel_all() *BaseProduct {
	return product.Insel_product_self().Insel_lot_self()
}
