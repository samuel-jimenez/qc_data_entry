package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"codeberg.org/go-pdf/fpdf"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Product struct {
	BaseProduct
	SG          NullFloat64
	PH          NullFloat64
	Density     NullFloat64
	String_test NullFloat64
	Viscosity   NullFloat64
}

func (product Product) save() {
	db_insert_sample_point.Exec(product.Sample_point)
	_, err := db_insert_measurement.Exec(product.lot_id, product.Sample_point, time.Now().UTC().UnixNano(), product.PH, product.SG, product.String_test, product.Viscosity)
	if err != nil {
		log.Println("error:", err)
		show_status("Sample Recording Failed")
	} else {
		show_status("Sample Recorded")
	}
}

func (product Product) export_json() {
	output_files := product.get_json_names()
	bytestring, err := json.MarshalIndent(product, "", "\t")

	if err != nil {
		log.Println("error:", err)
	}

	for _, output_file := range output_files {
		if err := os.WriteFile(output_file, bytestring, 0666); err != nil {
			log.Fatal(err)
		}
	}
}

func (product Product) toProduct() Product {
	return product
}

func (product Product) check_data() bool {
	return true
}

func (product Product) output() error {
	product.export_json()
	return product.print()
}

func (product *Product) format_sample() {
	if product.Product_name_customer != "" {
		product.Product_type = product.Product_name_customer
	}
	product.Sample_point = ""
}

func (product Product) output_sample() error {
	product.format_sample()
	return product.print()
}

func _print(pdf_path string) {
	print_queue <- pdf_path
	show_status("Label Printed")
}

func (product Product) print() error {

	pdf_path, err := product.export_label_pdf()
	if err != nil {
		return err
	}
	show_status("Label Created")

	_print(pdf_path)

	return err
}

func (product Product) reprint() {
	_print(product.get_pdf_name())
}

func (product Product) reprint_sample() {
	product.format_sample()
	product.reprint()
}

func (product Product) export_label_pdf() (string, error) {
	var label_width, label_height,
		field_width, field_height,
		unit_width, unit_height,
		label_col,
		// field_col,
		product_row,
		curr_row,
		curr_row_delta,
		lot_row float64

	label_width = 40
	label_height = 10

	field_width = 20
	field_height = 10

	unit_width = 40
	unit_height = 10

	label_col = 10
	// field_col = 120

	product_row = 0
	lot_row = 45

	file_path := product.get_pdf_name()

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(label_col, product_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.Product_type))

	if product.Density.Valid {
		curr_row = 5
		curr_row_delta = 6

	} else {
		curr_row = 10
		curr_row_delta = 10

	}

	var sg_derived bool
	if product.PH.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "pH")
		pdf.Cell(field_width, field_height, format_ph(product.PH.Float64))
		sg_derived = false
	} else {
		sg_derived = true
	}

	if product.SG.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "SG")
		pdf.Cell(field_width, field_height, format_sg(product.SG.Float64, !sg_derived))
		pdf.Cell(unit_width, unit_height, "g/mL")
	}

	if product.Density.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "DENSITY")
		pdf.Cell(field_width, field_height, format_density(product.Density.Float64))
		pdf.Cell(unit_width, unit_height, "lb/gal")
	}

	if product.String_test.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "STRING")
		pdf.Cell(field_width, field_height, format_string_test(product.String_test.Float64))
		pdf.Cell(unit_width, unit_height, "s")
	}

	if product.Viscosity.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "VISCOSITY")
		pdf.Cell(field_width, field_height, format_viscosity(product.Viscosity.Float64))
		pdf.Cell(unit_width, unit_height, "cP")
	}

	// log.Println(curr_row)

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(label_width, field_height, strings.ToUpper(product.Lot_number))
	pdf.CellFormat(unit_width, field_height, strings.ToUpper(product.Sample_point), "", 0, "R", false, 0, "")

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	return file_path, err
}
