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

// func (product Product) MarshalJSON() ([]byte, error) {
// 	_, err := db_insert_measurement.Exec(product.lot_id, product.sample_point, time.Now().UTC().UnixNano(), product.sg, product.ph, product.string_test, product.viscosity)
// 	return err
// }

func (product Product) save() error {
	_, err := db_insert_measurement.Exec(product.lot_id, product.Sample_point, time.Now().UTC().UnixNano(), product.SG, product.PH, product.String_test, product.Viscosity)
	return err
}

func (product Product) export_json() {
	output_files := product.get_json_names()
	// bytestring, err := json.Marshal(product)
	bytestring, err := json.MarshalIndent(product, "", "\t")

	if err != nil {
		log.Println("error:", err)
	}

	// os.Stdout.Write(bytestring)
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
	return product.export_label_pdf()
}

func (product Product) output_sample() error {
	product.Product_type = product.Product_name_customer
	product.Sample_point = ""
	return product.export_label_pdf()
}

func (product Product) export_label_pdf() error {
	var label_width, label_height,
		field_width, field_height,
		label_col,
		// field_col,
		product_row,
		curr_row,
		curr_row_delta,
		lot_row float64

	label_width = 40
	label_height = 10

	field_width = 40
	field_height = 10

	label_col = 10
	// field_col = 120

	product_row = 0
	lot_row = 45

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(label_col, product_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.Product_type))

	if product.Density.Valid {
		curr_row = 5
		curr_row_delta = 5

	} else {
		curr_row = 10
		curr_row_delta = 10

	}
	//TODO unit

	if product.SG.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "SG")
		pdf.Cell(field_width, field_height, format_sg(product.SG.Float64))

	}

	if product.PH.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "pH")
		pdf.Cell(field_width, field_height, format_ph(product.PH.Float64))

	}

	if product.Density.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "DENSITY")
		pdf.Cell(field_width, field_height, format_density(product.Density.Float64))
	}

	if product.String_test.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "STRING")
		pdf.Cell(field_width, field_height, format_string_test(product.String_test.Float64))

	}
	if product.Viscosity.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "VISCOSITY")
		pdf.Cell(field_width, field_height, format_viscosity(product.Viscosity.Float64))

	}

	// log.Println(curr_row)

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.Lot_number))
	pdf.CellFormat(field_width, field_height, strings.ToUpper(product.Sample_point), "", 0, "R", false, 0, "")

	log.Println("saving to: ", product.get_pdf_name())
	err := pdf.OutputFileAndClose(product.get_pdf_name())
	return err
}
