package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/config"
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

func (product BaseProduct) toBaseProduct() BaseProduct {
	return product
}

//TODO product.get_coa_name()

func (product BaseProduct) get_base_filename(extension string) string {
	// if (product.Sample_point.Valid) {
	if product.Sample_point != "" {
		return fmt.Sprintf("%s-%s-%s.%s", strings.ToUpper(product.Lot_number), product.Sample_point, strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.Product_type)), " ", "_"), extension)
	}

	return fmt.Sprintf("%s-%s.%s", strings.ToUpper(product.Lot_number), strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(product.Product_type)), " ", "_"), extension)
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

func (product *BaseProduct) insel_product_id(product_name string) {
	product.product_id = insel_product_id(product_name)
	log.Println("Debug: insel_product_id", product)
}

func (product *BaseProduct) insel_lot_id(lot_name string) {
	product.lot_id = insel_lot_id(lot_name, product.product_id)
}

func (product *BaseProduct) insel_product_self() *BaseProduct {
	product.insel_product_id(product.Product_type)
	return product

}

func (product *BaseProduct) insel_lot_self() *BaseProduct {
	product.insel_lot_id(product.Lot_number)
	return product

}

func (product BaseProduct) insel_all() *BaseProduct {
	return product.insel_product_self().insel_lot_self()
}
