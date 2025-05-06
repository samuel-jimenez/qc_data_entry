package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/winc"
)

type BaseProduct struct {
	Product_type          string `json:"product_name"`
	Lot_number            string `json:"lot_number"`
	Sample_point          string
	Visual                bool
	product_id            int64
	lot_id                int64
	Product_name_customer string `json:"customer_product_name"`
}

func NewBaseProduct(product_field winc.Controller, lot_field winc.Controller, sample_field winc.Controller) BaseProduct {
	return BaseProduct{strings.ToUpper(product_field.Text()), strings.ToUpper(lot_field.Text()), strings.ToUpper(sample_field.Text()), false, -1, -1, ""}
}

func (product BaseProduct) toBaseProduct() BaseProduct {
	return product
}

func (product BaseProduct) get_base_filename(extension string) string {
	// if (product.Sample_point.Valid) {
	if product.Sample_point != "" {
		return fmt.Sprintf("%s-%s-%s-%s.%s", strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.Product_type)), " ", "_"), strings.ToUpper(product.Lot_number), product.Sample_point, extension)
	}

	return fmt.Sprintf("%s-%s-%s.%s", strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.Product_type)), " ", "_"), strings.ToUpper(product.Lot_number), extension)
}

func (product BaseProduct) get_pdf_name() string {
	return fmt.Sprintf("%s/%s", LABEL_PATH, product.get_base_filename("pdf"))

}

func (product BaseProduct) get_json_filename(path string, base_name string) string {

	return fmt.Sprintf("%s/%s-%s", path, time.Now().UTC().UnixNano(), base_name)
}

func (product BaseProduct) get_json_names() []string {
	var json_names []string
	base_name := product.get_base_filename("json")

	for _, JSON_PATH := range JSON_PATHS {

		json_names = append(json_names, product.get_json_filename(JSON_PATH, base_name))
	}
	return json_names

}

func (product *BaseProduct) insel_product_id(product_name string) {
	log.Println("insel_product_id product_name", product_name)

	product.product_id = insel_product_id(product_name)
	log.Println("insel_product_id", product)
}

func (product *BaseProduct) insel_lot_id(lot_name string) {
	product.lot_id = insel_lot_id(lot_name, product.product_id)
	log.Println("lot_id", product.lot_id)

}

func (product *BaseProduct) insel_product_self() *BaseProduct {
	product.insel_product_id(product.Product_type)
	return product

}

func (product *BaseProduct) insel_lot_self() *BaseProduct {
	product.insel_lot_id(product.Lot_number)
	log.Println("insel_lot_self", product.lot_id)
	return product

}

func (product BaseProduct) insel_all() *BaseProduct {
	return product.insel_product_self().insel_lot_self()

}

func (product *BaseProduct) copy_ids(product_lot BaseProduct) {
	if product_lot.product_id > 0 {
		product.Product_name_customer = select_product_name_customer(product_lot.product_id)
		product.product_id = product_lot.product_id
		if product_lot.lot_id > 0 {
			product.lot_id = product_lot.lot_id
		} else {
			product.insel_lot_self()
		}
	} else {
		product.insel_all()
		product.Product_name_customer = select_product_name_customer(product.product_id)

	}

	log.Println("copy_ids lot_id", product.lot_id)

}
