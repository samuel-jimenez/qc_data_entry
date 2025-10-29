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
	"github.com/samuel-jimenez/qc_data_entry/util"
)

// TODO extract this stuff to package
// TODO print-pdf-019

func Print_PDF(pdf_path string) {
	if threads.PRINT_QUEUE != nil {
		threads.PRINT_QUEUE <- pdf_path
		threads.Show_status("Label Printed")
	} else {
		log.Println("Warn: Print queue not configured. Call threads.Do_print_queue() to set up.")

	}
}

/*
 * PrintOldStorage
 *
 * Create and print storage label with dates
 *
 */
func (measured_product BaseProduct) PrintOldStorage(qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date *time.Time) error {

	if err := measured_product.PrintStorage(qc_sample_storage_name, product_moniker_name, start_date, end_date, retain_date, true); err != nil {
		return err
	}
	return nil
}

/*
 * PrintNewStorage
 *
 * Create and print storage label
 * name and id only
 *
 */
func (measured_product BaseProduct) PrintNewStorage(qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date *time.Time) error {
	if err := measured_product.PrintStorage(qc_sample_storage_name, product_moniker_name, start_date, end_date, retain_date, false); err != nil {
		return err
	}
	return nil

}

/*
 * PrintStorage
 *
 * Create and print storage label
 *
 */
func (measured_product BaseProduct) PrintStorage(qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date *time.Time, printDates bool) error {
	file_path := measured_product.get_storage_pdf_name(qc_sample_storage_name)

	if err := Export_Storage_pdf(file_path, qc_sample_storage_name, product_moniker_name, start_date, end_date, retain_date, printDates); err != nil {
		return err
	}
	threads.Show_status("Label Created")

	Print_PDF(file_path)

	return nil

}

// TODO extract this stuff to package
// TODO pdf-IG99

/*
 * Increment_row_pdf
 *
 * print a cell, move row by delta
 *
 */
func Increment_row_pdf(pdf *fpdf.Fpdf, x, y, delta_y float64, width, height float64, txtStr string) float64 {
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
func Export_Storage_pdf(file_path, qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date *time.Time, printDates bool) error {
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

	curr_row = Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, qc_sample_storage_name)

	curr_row = Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, product_moniker_name)

	if printDates {
		pdf.SetFontSize(24)
		curr_row = Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, start_date.Format(time.DateOnly))

		curr_row = Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, end_date.Format(time.DateOnly))

		curr_row = Increment_row_pdf(pdf, curr_col, curr_row, curr_row_delta, cell_width, cell_height, retain_date.Format(time.DateOnly))
	}

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	util.LogError(proc_name, err)
	return err

}

func (measured_product MeasuredProduct) export_label_pdf() (string, error) {
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

	file_path := measured_product.get_pdf_name()

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(label_col, product_row)
	pdf.Cell(field_width, field_height, strings.ToUpper(measured_product.Product_name))

	if measured_product.Density.Valid {
		curr_row = 5
		curr_row_delta = 6

	} else {
		curr_row = 10
		curr_row_delta = 10

	}

	sg_derived := measured_product.QCProduct.SG.Method.String == METHOD_GARDCO_CUP
	if measured_product.PH.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "pH")
		pdf.Cell(field_width, field_height, formats.Format_ph(measured_product.PH.Float64))
	} else {
	}

	if measured_product.SG.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "SG")
		pdf.Cell(field_width, field_height, formats.Format_sg(measured_product.SG.Float64, !sg_derived))
		pdf.Cell(unit_width, unit_height, formats.SG_UNITS)
	}

	if measured_product.Density.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "DENSITY")
		pdf.Cell(field_width, field_height, formats.Format_density(measured_product.Density.Float64))
		pdf.Cell(unit_width, unit_height, formats.DENSITY_UNITS)
	}

	if measured_product.String_test.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "STRING")
		// pdf.Cell(field_width, field_height, formats.Format_string_test(product.String_test.Int64))
		pdf.Cell(field_width, field_height, formats.FormatInt(measured_product.String_test.Int64))
		pdf.Cell(unit_width, unit_height, formats.STRING_UNITS)
	}

	if measured_product.Viscosity.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "VISCOSITY")
		// pdf.Cell(field_width, field_height, formats.Format_viscosity(product.Viscosity.Int64))
		pdf.Cell(field_width, field_height, formats.FormatInt(measured_product.Viscosity.Int64))
		pdf.Cell(unit_width, unit_height, formats.VISCOSITY_UNITS)
	}

	// log.Println(curr_row)

	pdf.SetXY(label_col, lot_row)
	pdf.Cell(label_width, field_height, strings.ToUpper(measured_product.Lot_number))
	pdf.CellFormat(unit_width, field_height, strings.ToUpper(measured_product.Sample_point), "", 0, "R", false, 0, "")

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	return file_path, err
}
