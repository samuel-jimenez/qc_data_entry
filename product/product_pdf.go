package product

import (
	"log"
	"strings"
	"time"

	"codeberg.org/go-pdf/fpdf"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/threads"
)

func _print(pdf_path string) {
	if threads.PRINT_QUEUE != nil {
		threads.PRINT_QUEUE <- pdf_path
		threads.Show_status("Label Printed")
	} else {
		log.Println("Warn: Print queue not configured. Call threads.Do_print_queue() to set up.")

	}
}

/*
 * PrintStorage
 *
 * Create and print storage label
 *
 */
func (product Product) PrintStorage(qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date time.Time) error {
	file_path := product.get_storage_pdf_name(qc_sample_storage_name)

	if err := Export_Storage_pdf(file_path, qc_sample_storage_name, product_moniker_name, start_date, end_date, retain_date); err != nil {
		return err
	}
	threads.Show_status("Label Created")

	_print(file_path)

	return nil

}

/*
 * increment_row_pdf
 *
 * print a cell, move row by delta
 *
 */
func increment_row_pdf(pdf *fpdf.Fpdf, x, y, delta_y float64, width, height float64, txtStr string) float64 {
	pdf.SetXY(x, y)
	pdf.Cell(width, height, txtStr)
	return y + delta_y
}

/*
 * Export_Storage_pdf
 *
 * Create pdf with given data
 *
 */
func Export_Storage_pdf(file_path, qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date time.Time) error {
	proc_name := "Product.Export_Storage_pdf"

	curr_col := 10.
	curr_row := 5.
	curr_row_delta := 10.

	cell_width := 50.
	cell_height := 10.

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 32)

	curr_row = increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, qc_sample_storage_name)

	curr_row = increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, product_moniker_name)

	pdf.SetFontSize(24)
	curr_row = increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, start_date.Format(time.DateOnly))

	curr_row = increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, end_date.Format(time.DateOnly))

	curr_row = increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, retain_date.Format(time.DateOnly))

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
	}
	return err

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
	pdf.Cell(field_width, field_height, strings.ToUpper(product.Product_name))

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
		pdf.Cell(field_width, field_height, formats.Format_ph(product.PH.Float64))
		sg_derived = false
	} else {
		sg_derived = true
	}

	if product.SG.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "SG")
		pdf.Cell(field_width, field_height, formats.Format_sg(product.SG.Float64, !sg_derived))
		pdf.Cell(unit_width, unit_height, formats.SG_UNITS)
	}

	if product.Density.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "DENSITY")
		pdf.Cell(field_width, field_height, formats.Format_density(product.Density.Float64))
		pdf.Cell(unit_width, unit_height, formats.DENSITY_UNITS)
	}

	if product.String_test.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "STRING")
		// pdf.Cell(field_width, field_height, formats.Format_string_test(product.String_test.Int64))
		pdf.Cell(field_width, field_height, formats.FormatInt(product.String_test.Int64))
		pdf.Cell(unit_width, unit_height, formats.STRING_UNITS)
	}

	if product.Viscosity.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "VISCOSITY")
		// pdf.Cell(field_width, field_height, formats.Format_viscosity(product.Viscosity.Int64))
		pdf.Cell(field_width, field_height, formats.FormatInt(product.Viscosity.Int64))
		pdf.Cell(unit_width, unit_height, formats.VISCOSITY_UNITS)
	}

	// log.Println(curr_row)

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(label_width, field_height, strings.ToUpper(product.Lot_number))
	pdf.CellFormat(unit_width, field_height, strings.ToUpper(product.Sample_point), "", 0, "R", false, 0, "")

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	return file_path, err
}
