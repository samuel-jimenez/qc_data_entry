package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"codeberg.org/go-pdf/fpdf"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/whatsupdocx"
	"github.com/samuel-jimenez/whatsupdocx/docx"
	"github.com/xuri/excelize/v2"
)

var (
	COA_LOT_TITLE = "Batch/Lot#"
)

type Product struct {
	BaseProduct
	SG          nullable.NullFloat64
	PH          nullable.NullFloat64
	Density     nullable.NullFloat64
	String_test nullable.NullFloat64
	Viscosity   nullable.NullFloat64
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

func (product Product) get_coa_template() string {
	//TODO db
	// return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, product.COA_TEMPLATE)
	// Product_type split
	//TODO fix this
	product_moniker := strings.Split(product.Product_type, " ")[0]
	if product_moniker == "PETROFLO" {
		return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA-PETROFLO.docx")
	}
	return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA.docx")
}

func (product Product) get_coa_name() string {
	return fmt.Sprintf("%s/%s", config.COA_FILEPATH, product.get_base_filename("docx"))
}

func (product Product) export_CoA() error {
	var (
		p_title   = "[PRODUCT_NAME]"
		Coa_title = "Parameter"
	)
	terms := []string{
		COA_LOT_TITLE,
	}

	template_file := product.get_coa_template()
	output_file := product.get_coa_name()

	doc, err := whatsupdocx.OpenDocument(template_file)
	if err != nil {
		return err
	}

	product_name := product.Product_name_customer
	if product_name == "" {
		product_name = product.Product_type
	}
	for _, item := range doc.Document.Body.Children {
		if para := item.Paragraph; para != nil {
			product.searchCOAPara(para, p_title, product_name)
		}

		if table := item.Table; table != nil {
			product.searchCOATable(table, Coa_title, terms)
		}
	}

	// save to file
	return doc.SaveTo(output_file)
}

func (product Product) searchCOAPara(para *docx.Paragraph, p_title, product_name string) {
	if strings.Contains(para.String(), p_title) {
		for _, child := range para.Children {
			if run := child.Run; run != nil && strings.Contains(run.String(), p_title) {

				run.Clear()

				//Add product name
				run.AddText(product_name)
				return
			}
		}
	}
}

func (product Product) searchCOATable(table *docx.Table, Coa_title string, terms []string) {
	for _, row := range table.RowContents {
		if product.searchCOARow(table, row.Row, Coa_title, terms) {
			return
		}
	}
}

func (product Product) searchCOARow(table *docx.Table, row *docx.Row, Coa_title string, terms []string) (finished bool) {
	currentHeading := ""
ROW:
	for i, cell := range row.Contents {
		for _, cont := range cell.Cell.Contents {
			field := cont.Paragraph
			if field == nil {
				continue
			}
			switch i {
			case 0:
				if currentHeading == "" {
					field_string := field.String()
					if strings.Contains(field_string, Coa_title) {
						QCProduct{Product: product}.write_CoA_rows(table)
						return true
					}
					for _, term := range terms {
						if strings.Contains(field_string, term) {
							currentHeading = term
							continue ROW
						}
					}
				}
			case 1:
				switch currentHeading {
				case COA_LOT_TITLE:
					field.AddText(product.Lot_number)
					return false
				}
			}
		}
	}
	return false
}

func (product Product) toProduct() Product {
	return product
}

func (product Product) check_data() bool {
	return true
}

func (product Product) printout() error {
	product.export_json()
	return product.print()
}

func (product Product) output() error {
	err := product.export_CoA()
	if err != nil {
		return err
	}
	//TODO test
	return product.printout()
}

func (product *Product) format_sample() {
	if product.Product_name_customer != "" {
		product.Product_type = product.Product_name_customer
	}
	product.Sample_point = ""
}

func (product Product) output_sample() error {
	product.format_sample()
	err := product.export_CoA()
	if err != nil {
		return err
	}
	return product.print()
	//TODO clean
	// return nil

}

func withOpenFile(file_name string, FN func(*excelize.File) error) error {
	xl_file, err := excelize.OpenFile(file_name)
	if err != nil {
		log.Printf("Error: %q: %s\n", err, "withOpenFile")
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := xl_file.Close(); err != nil {
			log.Printf("Error: %q: %s\n", err, "withOpenFile")
		}
	}()
	return FN(xl_file)
}

func updateExcel(file_name, worksheet_name string, row ...string) error {

	return withOpenFile(file_name, func(xl_file *excelize.File) error {
		// Get all the rows in the worksheet.
		rows, err := xl_file.GetRows(worksheet_name)
		if err != nil {
			log.Printf("Error: %q: %s\n", err, "updateExcel")
			return err
		}
		startCell := fmt.Sprintf("A%v", len(rows)+1)

		err = xl_file.SetSheetRow(worksheet_name, startCell, &row)
		if err != nil {
			log.Printf("Error: %q: %s\n", err, "updateExcel")
			return err
		}
		return xl_file.Save()
	})
}

func (product Product) save_xl() error {
	isomeric_product := strings.Split(product.Product_type, " ")[1]
	valence_product := strings.Split(product.Product_name_customer, " ")[1]
	return updateExcel(config.RETAIN_FILE_NAME, config.RETAIN_WORKSHEET_NAME, product.Lot_number, isomeric_product, valence_product)
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
		pdf.Cell(field_width, field_height, formats.Format_string_test(product.String_test.Float64))
		pdf.Cell(unit_width, unit_height, formats.STRING_UNITS)
	}

	if product.Viscosity.Valid {
		curr_row += curr_row_delta
		pdf.SetXY(label_col, curr_row)
		pdf.Cell(label_width, label_height, "VISCOSITY")
		pdf.Cell(field_width, field_height, formats.Format_viscosity(product.Viscosity.Float64))
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
