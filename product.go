package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"codeberg.org/go-pdf/fpdf"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Product struct {
	BaseProduct
	sg          sql.NullFloat64
	ph          sql.NullFloat64
	density     sql.NullFloat64
	string_test sql.NullFloat64
	viscosity   sql.NullFloat64
}

func (product Product) save() error {
	_, err := db_insert_measurement.Exec(product.lot_id, product.sample_point, time.Now().UTC().UnixNano(), product.sg, product.ph, product.string_test, product.viscosity)
	return err
}

func (product Product) toProduct() Product {
	return product
}

func (product Product) check_data() bool {
	return true
}

func (product Product) print() error {
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
	pdf.Cell(field_width, field_height, strings.ToUpper(product.product_type))

	if product.density.Valid {
		curr_row = 5
		curr_row_delta = 5

	} else {
		curr_row = 10
		curr_row_delta = 10

	}
	//TODO unit

	if product.sg.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "SG")
		pdf.Cell(field_width, field_height, format_sg(product.sg.Float64))

	}

	if product.ph.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "pH")
		pdf.Cell(field_width, field_height, format_ph(product.ph.Float64))

	}

	if product.density.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "DENSITY")
		pdf.Cell(field_width, field_height, format_density(product.density.Float64))
	}

	if product.string_test.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "STRING")
		pdf.Cell(field_width, field_height, format_string_test(product.string_test.Float64))

	}
	if product.viscosity.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "VISCOSITY")
		pdf.Cell(field_width, field_height, format_viscosity(product.viscosity.Float64))

	}

	fmt.Println(curr_row)

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(product.lot_number))
	pdf.CellFormat(field_width, field_height, strings.ToUpper(product.sample_point), "", 0, "R", false, 0, "")

	fmt.Println("saving to: ", product.get_pdf_name())
	err := pdf.OutputFileAndClose(product.get_pdf_name())
	return err
}
