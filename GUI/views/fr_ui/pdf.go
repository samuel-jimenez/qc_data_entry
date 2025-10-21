package fr_ui

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"codeberg.org/go-pdf/fpdf"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/QR"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/util"
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// TODO pdf-IG99
// TODO extract this stuff to package

/*
 * pdf_SetBaseline
 *
 * set base margin
 *
 */
func pdf_SetBaseline(pdf *fpdf.Fpdf, x float64) {
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
}

/*
 * pdf_MultiCell
 *
 * MultiCell
 *
 * fpdf, youre just terrible
 *
 */
func pdf_MultiCell(pdf *fpdf.Fpdf, w, h float64, txtStr, borderStr string, ln int,
	alignStr string, fill bool) {
	x, y := pdf.GetXY()
	pdf.MultiCell(w, h, txtStr, borderStr, alignStr, fill)
	switch ln {
	case 0:
		pdf.SetXY(x+w, y)
	case 1:
		// pdf.SetXY(0, y+h)
		pdf.SetXY(x+w, y)
		pdf.Ln(-1)

	}
}

/*
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

// TODO extract this stuff to package
// TODO pdf-IG99

func Increment_row_pdf(pdf *fpdf.Fpdf, x, y, delta_y float64, width, height float64, txtStr string) float64 {
pdf.SetXY(x, y)
pdf.Cell(width, height, txtStr)
return y + delta_y
}

pdf.SetLeftMargin(10)
pdf.SetX(10)
/*
 * Increment_row_pdf
 *
 * print a cell, move row by delta
 *
*/

/*
func Increment_row_pdf(pdf *fpdf.Fpdf, x, y, delta_y float64, width, height float64, txtStr string) float64 {
	pdf.SetXY(x, y)
	pdf.Cell(width, height, txtStr)
	return y + delta_y
}

pdf.SetLeftMargin(10)
pdf.SetX(10)

/*
 * Export_Storage_pdf
 *
 * Create pdf with given data
 *
*/
/*
func Export_Storage_pdf(file_path, qc_sample_storage_name, product_moniker_name string, start_date, end_date, retain_date *time.Time, printDates bool) error {
	proc_name := "Product.Export_Storage_pdf"

	curr_col := 10.
	curr_row := 5.
	Top_f_height := 10.

	cell_width := 50.
	cell_height := 10.

	pdf := fpdf.New("L", "mm", "A7", "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 32)

	curr_row = Increment_row_pdf(pdf, curr_col, curr_row, Top_f_height, cell_width, cell_height, qc_sample_storage_name)

	curr_row = Increment_row_pdf(pdf, curr_col, curr_row, Top_f_height, cell_width, cell_height, product_moniker_name)

	if printDates {
		pdf.SetFontSize(24)
		curr_row = Increment_row_pdf(pdf, curr_col, curr_row, Top_f_height, cell_width, cell_height, start_date.Format(time.DateOnly))

		curr_row = Increment_row_pdf(pdf, curr_col, curr_row, Top_f_height, cell_width, cell_height, end_date.Format(time.DateOnly))

		curr_row = Increment_row_pdf(pdf, curr_col, curr_row, Top_f_height, cell_width, cell_height, retain_date.Format(time.DateOnly))
	}

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	util.LogError(proc_name, err)
	return err

}

// TODO
/*


// operators
// TAG
// SEAL

// TODO fields:
// TAG
// SEAL

// FINAL STRAP
// /density?


product info
product
custnane
batch vessel
who
were
when
tag seal
final


QC info

test spec result
tester


blend info

componente amounts



procedure

QR



*/
func (measured_product BlendSheet) get_pdf_name() string {
	// BLENDSHEET_PATH := "."
	// BLENDSHEET_PATH := "C:/Users/QC Lab/Documents/golang/qc_data_entry"
	// return fmt.Sprintf("%s/%s-%s.pdf", BLENDSHEET_PATH,strings.ToUpper(measured_product.Lot_number), strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(measured_product.Product_name)), " ", "_"))
	// return fmt.Sprintf("%s/%s-%s.pdf", BLENDSHEET_PATH, strings.ToUpper(measured_product.Lot_number), strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(measured_product.Product_name)), " ", "_"))
	return fmt.Sprintf("%s/%s-%s.pdf", config.BLENDSHEET_PATH, strings.ToUpper(measured_product.Lot_number), strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(measured_product.Product_name)), " ", "_"), "-", "_"))
}

func (measured_product BlendSheet) get_qr_name() string {
	return fmt.Sprintf("%s/%s.png", config.QR_PATH, strings.ToUpper(measured_product.Lot_number))
}

func (measured_product BlendSheet) export_pdf() (string, error) {
	var (
		label_width, Top_f_height,
		field_width, field_height,
		// unit_width, unit_height,
		data_height,
		logo_width,
		logo_col,
		top_col,
		// field_col,
		tester_col,
		tester_height,
		start_row,
		testing_margin,
		Procedure_margin,
		qc_row,
		qc_header_height,
		qc_height,
		// qc_test_width,
		// qc_spec_width,
		// qc_result_width,
		// qc_margin,
		Component_margin,
		Component_col,
		qc_col float64

		EMPTY string
	)

	product.QC_TEST_WIDTH = 30
	product.QC_SPEC_WIDTH = 40
	product.QC_RESULT_WIDTH = 40

	top_margin := 5.
	testing_margin = 5.
	start_row = 5
	Component_margin = 5
	Procedure_margin = 2
	btm_margin := 1.

	// lot_row float64

	label_width = 60
	field_width = 100

	Top_f_height = 7
	field_height = 9.5
	data_height = 5

	// logo_width = 30
	logo_width = 40
	// logo_width = 50
	// unit_width = 40
	// unit_height = 10

	// qr_col, qr_row, qr_width := 220., 0., 50.
	// qr_col, qr_row, qr_width := 0., logo_width*log_h/log_w+start_row, logo_width+logo_col
	// qr_col, qr_row, qr_width := 0., 0., logo_width

	logo_col = 3
	// logo_col = 0

	top_col = logo_col + logo_width + top_margin
	// field_col = 120
	// tester_col = top_col + label_width + field_width + 10

	tester_height = 20

	ORDER_width := 10.
	VAlue_width := 20.
	LABEL_height := 4.
	// LABEL_height := 10.

	// lot_row = 45
	qc_col = 15
	// qc_col = top_col
	// qc_col = tester_col

	tester_col = qc_col +
		product.QC_TEST_WIDTH +
		product.QC_SPEC_WIDTH +
		2*product.QC_RESULT_WIDTH + 10

	// qc_margin = 5
	// //
	qc_header_height = 5
	qc_height = field_height
	//
	// qc_test_width = 30
	// qc_spec_width = 30
	// qc_result_width = 30

	Component_col = 10

	Procedure_height := 5.
	Procedure_width := 250.

	file_path := measured_product.get_pdf_name()

	commanFormatter := message.NewPrinter(language.English)

	// product.QCProduct
	qc_product := product.NewQCProduct()
	qc_product.Product_id = measured_product.Product_id
	qc_product.Select_product_details()

	// pdf := fpdf.New("L", "mm", "A4", "")
	pdf := fpdf.New("L", "mm", "Letter", "C:/Windows/Fonts/")

	//
	pdf.AddUTF8Font("Arial", "", "")
	pdf.AddUTF8Font("Arial", "B", "arialbd.ttf")

	// pdf.SetAutoPageBreak(false, 0)
	pdf.SetAutoPageBreak(true, .5)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	//logo
	// infoPtr = pdf.RegisterImage(imageFileStr, "")
	LOGO_file := config.LOGO_PATH + "/Isomeric-logo.jpg"

	// pdf.ImageOptions(LOGO_file, logo_col, start_row, logo_width, 0, false, fpdf.ImageOptions{ImageType: "jpg", ReadDpi: true}, 0, "")
	pdf.ImageOptions(LOGO_file, logo_col, start_row, logo_width, 0, false, fpdf.ImageOptions{ImageType: "jpg", ReadDpi: true}, 0, "")
	// pdf.ImageOptions(LOGO_file, logo_col, 0, logo_width, 0, false, fpdf.ImageOptions{ImageType: "jpg", ReadDpi: true}, 0, "")
	// pdf.ImageOptions(LOGO_file, logo_col, qr_width, logo_width, 0, false, fpdf.ImageOptions{ImageType: "jpg", ReadDpi: true}, 0, "")

	// pdf.GetImageInfo(config.LOGO_PATH+"/Isomeric-logo.jpg")
	// log_w, log_h := pdf.GetImageInfo(LOGO_file).Extent()

	// log.Println(logo_width, log_w, log_h, logo_width*log_h/log_w)

	// log_w, log_h := pdf.GetImageInfo(LOGO_file).Extent()
	// qr_col, qr_row, qr_width := 0., logo_width*log_h/log_w, logo_width

	// qr_col, qr_row, qr_width := tester_col, 0., logo_width

	// log.Println(pdf.GetPageSize())
	// log.Println(pdf.GetMargins())

	page_width, page_height := pdf.GetPageSize()
	// left, top, right, bottom :=pdf.GetMargins()
	// _, _, right, _ := pdf.GetMargins()

	qr_col, qr_row, qr_width := page_width-logo_width, 0., logo_width
	// tester_col = qr_col

	// QR
	proc_name := "WriteQR"
	util.LogError(proc_name, qrcode.WriteFile(string(util.VALUES(json.Marshal(QR.QRJson{Product_type: measured_product.Product_name, Lot_number: measured_product.Lot_number}))), qrcode.Medium, 256, measured_product.get_qr_name()))

	pdf.ImageOptions(measured_product.get_qr_name(), qr_col, qr_row, qr_width, 0, false, fpdf.ImageOptions{ImageType: "png", ReadDpi: true}, 0, "")

	// operators
	// TAG
	// SEAL

	// TODO fields:
	// TAG
	// SEAL

	// FINAL STRAP
	// /density?

	// 	product info
	// 	product
	// 	custnane
	// 	batch vessel
	// 	who
	// 	were
	// 	when
	// 	tag seal
	// 	final
	//
	//

	//
	// TODO Sheet-03988

	pdf_SetBaseline(pdf, top_col)

	pdf.SetXY(top_col, start_row)
	pdf.CellFormat(label_width, Top_f_height, "ISOMERIC NAME:", "1", 0, "", false, 0, "")
	pdf.CellFormat(field_width, Top_f_height, strings.ToUpper(measured_product.Product_name), "1", 1, "", false, 0, "")

	pdf.CellFormat(label_width, Top_f_height, "CUSTOMER NAME:", "1", 0, "", false, 0, "")
	pdf.CellFormat(field_width, Top_f_height, strings.ToUpper(measured_product.Product_name_customer), "1", 1, "", false, 0, "")

	pdf.CellFormat(label_width, Top_f_height, "Lot Number", "1", 0, "", false, 0, "")
	pdf.CellFormat(field_width, Top_f_height, strings.ToUpper(measured_product.Lot_number), "1", 1, "", false, 0, "")

	pdf.CellFormat(label_width, Top_f_height, "Operators", "1", 0, "", false, 0, "")
	pdf.SetFontSize(12)
	pdf.CellFormat(field_width, Top_f_height, strings.ToUpper(measured_product.Operators), "1", 1, "", false, 0, "")

	pdf.SetFontSize(14)

	label_width = 40
	field_width = 40

	pdf.CellFormat(label_width, Top_f_height, "Date", "1", 0, "", false, 0, "")
	pdf.CellFormat(field_width, Top_f_height, time.Now().Format(time.DateOnly), "1", 0, "", false, 0, "")

	pdf.CellFormat(label_width, Top_f_height, "Time", "1", 0, "", false, 0, "")
	pdf.CellFormat(field_width, Top_f_height, EMPTY, "1", 1, "", false, 0, "")

	pdf.CellFormat(label_width, data_height, "Blend Vessel", "1", 0, "", false, 0, "")
	pdf.CellFormat(field_width, data_height, strings.ToUpper(measured_product.Vessel), "1", 0, "", false, 0, "")

	pdf.CellFormat(label_width, data_height, "Capacity", "1", 0, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, GUI.Format_int(float64(measured_product.Amount)), "1", 1, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, strconv.Itoa(measured_product.Amount), "1", 1, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, commanFormatter.Sprintf("%.0f", measured_product.Capacity), "1", 1, "", false, 0, "")
	pdf.CellFormat(field_width, data_height, strings.ToUpper(measured_product.Capacity), "1", 1, "", false, 0, "")

	pdf.SetFontSize(12)
	pdf.CellFormat(label_width, data_height, "Blend Quantity, lb", "1", 0, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, GUI.Format_int(float64(measured_product.Amount)), "1", 1, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, strconv.Itoa(measured_product.Amount), "1", 1, "", false, 0, "")
	pdf.CellFormat(field_width, data_height, commanFormatter.Sprintf("%d", measured_product.Amount), "1", 0, "", false, 0, "")

	pdf.CellFormat(label_width, data_height, "Blend Quantity, gal", "1", 0, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, GUI.Format_int(float64(measured_product.Amount)), "1", 1, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, strconv.Itoa(measured_product.Amount), "1", 1, "", false, 0, "")
	pdf.CellFormat(field_width, data_height, commanFormatter.Sprintf("%d", int(measured_product.Total_Component.Gallons-measured_product.HeelVolume)), "1", 1, "", false, 0, "")

	// pdf.CellFormat(label_width, data_height, "Blend Quantity", "1", 0, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, GUI.Format_int(float64(measured_product.Amount)), "1", 1, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, strconv.Itoa(measured_product.Amount), "1", 1, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, commanFormatter.Sprintf("%d", measured_product.Amount), "1", 1, "", false, 0, "")

	// pdf.CellFormat(label_width, data_height, "Final weight", "1", 0, "", false, 0, "")
	// pdf.CellFormat(field_width, data_height, GUI.Format_int(measured_product.Total_Component.Component_amount), "1", 1, "", false, 0, "")

	//qc_row
	pdf.Ln(testing_margin)
	qc_row = pdf.GetY()

	// // // -------------------------------------- SIDEBAR ------------------------------------------
	// pdf.SetDrawColor(255, 0, 0)
	// pdf.SetTextColor(255, 0, 0)
	// pdf.SetFillColor(0, 255, 0)
	pdf_SetBaseline(pdf, tester_col)
	pdf.SetXY(tester_col, qc_row)
	pdf.SetFontStyle("B")
	pdf.SetFontSize(10)
	// pdf.CellFormat(field_width, tester_height, "TESTED BY: ", "0", 2, "T", true, 0, "")
	pdf.CellFormat(field_width, tester_height, "TESTED BY: ", "0", 2, "T", false, 0, "")
	pdf.CellFormat(field_width, data_height, "TAG: "+strings.ToUpper(measured_product.Tag), "0", 2, "", false, 0, "")
	pdf.CellFormat(field_width, data_height, "SEAL: "+strings.ToUpper(measured_product.Seal), "0", 2, "", false, 0, "")
	pdf.CellFormat(field_width, tester_height, "FINAL STRAP: ", "0", 2, "T", false, 0, "")

	// // // -------------------------------------- QC ------------------------------------------

	pdf_SetBaseline(pdf, qc_col)
	pdf.SetY(qc_row)
	qc_product.Write_QC_rows(pdf, qc_header_height,
		qc_height)

	// // // -------------------------------------- Components ------------------------------------------

	pdf.SetFontStyle("B")
	pdf.SetFontSize(8)
	pdf.Ln(Component_margin)
	pdf_SetBaseline(pdf, Component_col)

	//youre just terrible
	pdf_MultiCell(pdf, ORDER_width, 2*LABEL_height, "Order", "1", 0, "", false)
	pdf_MultiCell(pdf, field_width, 2*LABEL_height, "", "1", 0, "", false)
	pdf_MultiCell(pdf, field_width, 2*LABEL_height, "Lot Number", "1", 0, "", false)
	pdf_MultiCell(pdf, field_width, 2*LABEL_height, "Railcar", "1", 0, "", false)
	pdf_MultiCell(pdf, VAlue_width, LABEL_height, "Weight Required", "1", 0, "", false)
	pdf_MultiCell(pdf, field_width, 2*LABEL_height, "Actual Weight", "1", 0, "", false)
	pdf_MultiCell(pdf, VAlue_width, 2*LABEL_height, "Gallons", "1", 0, "", false)
	pdf_MultiCell(pdf, VAlue_width, LABEL_height, "Density, lbs/gal", "1", 0, "", false)
	pdf_MultiCell(pdf, VAlue_width, 2*LABEL_height, "Inches\n", "1", 1, "", false)

	pdf.SetFontSize(10)
	//HEEL
	HeelDensity := "--"
	if measured_product.HeelVolume != 0 {
		HeelDensity = GUI.Format_float(measured_product.HeelMass / measured_product.HeelVolume)
	}
	if measured_product.HeelMass > 0 {

		pdf.CellFormat(ORDER_width, data_height, EMPTY, "1", 0, "", false, 0, "")
		pdf.CellFormat(3*field_width, data_height, "HEEL", "1", 0, "C", false, 0, "")
		// pdf.CellFormat(VAlue_width, data_height, GUI.Format_int(measured_product.HeelMass), "1", 0, "R", false, 0, "")
		// pdf.CellFormat(field_width, data_height, GUI.Format_int(measured_product.HeelMass), "1", 0, "C", false, 0, "")
		// pdf.CellFormat(VAlue_width, data_height, GUI.Format_int(measured_product.HeelVolume), "1", 0, "R", false, 0, "")
		pdf.CellFormat(VAlue_width, data_height, commanFormatter.Sprintf("%.0f", measured_product.HeelMass), "1", 0, "R", false, 0, "")
		pdf.CellFormat(field_width, data_height, commanFormatter.Sprintf("%.0f", measured_product.HeelMass), "1", 0, "C", false, 0, "")
		pdf.CellFormat(VAlue_width, data_height, commanFormatter.Sprintf("%.0f", measured_product.HeelVolume), "1", 0, "R", false, 0, "")
		pdf.CellFormat(VAlue_width, data_height, HeelDensity, "1", 0, "R", false, 0, "")
		pdf.CellFormat(VAlue_width, data_height, GUI.Format_float(measured_product.Strap), "1", 1, "R", false, 0, "")
	}

	for _, Component := range measured_product.ProductBlend.Components {
		pdf.CellFormat(ORDER_width, field_height, strconv.Itoa(Component.Add_order), "1", 0, "", false, 0, "")
		pdf.CellFormat(field_width, field_height, strings.ToUpper(Component.Component_name), "1", 0, "", false, 0, "")
		pdf.CellFormat(field_width, field_height, strings.ToUpper(Component.Lot_name), "1", 0, "", false, 0, "")
		pdf.CellFormat(field_width, field_height, strings.ToUpper(Component.Container_name), "1", 0, "", false, 0, "")
		// pdf.CellFormat(VAlue_width, field_height, GUI.Format_int(Component.Component_amount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(VAlue_width, field_height, commanFormatter.Sprintf("%.0f", Component.Component_amount), "1", 0, "R", false, 0, "")

		pdf.CellFormat(field_width, field_height, EMPTY, "1", 0, "", false, 0, "")
		// pdf.CellFormat(VAlue_width, field_height, GUI.Format_int(Component.Gallons), "1", 0, "R", false, 0, "")
		pdf.CellFormat(VAlue_width, field_height, commanFormatter.Sprintf("%.0f", Component.Gallons), "1", 0, "R", false, 0, "")
		pdf.CellFormat(VAlue_width, field_height, GUI.Format_float(Component.Density), "1", 0, "R", false, 0, "")
		pdf.CellFormat(VAlue_width, field_height, GUI.Format_float(Component.Strap), "1", 1, "R", false, 0, "")

	}
	// TOTAL
	pdf.CellFormat(ORDER_width, data_height, EMPTY, "1", 0, "", false, 0, "")
	pdf.CellFormat(3*field_width, data_height, "TOTAL", "1", 0, "C", false, 0, "")
	// pdf.CellFormat(VAlue_width, data_height, GUI.Format_int(measured_product.Total_Component.Component_amount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(VAlue_width, data_height, commanFormatter.Sprintf("%.0f", measured_product.Total_Component.Component_amount), "1", 0, "R", false, 0, "")
	pdf.CellFormat(field_width, data_height, EMPTY, "1", 0, "C", false, 0, "")
	// pdf.CellFormat(VAlue_width, data_height, GUI.Format_int(measured_product.Total_Component.Gallons), "1", 0, "R", false, 0, "")
	pdf.CellFormat(VAlue_width, data_height, commanFormatter.Sprintf("%.0f", measured_product.Total_Component.Gallons), "1", 0, "R", false, 0, "")
	pdf.CellFormat(VAlue_width, data_height, GUI.Format_float(measured_product.Total_Component.Density), "1", 1, "R", false, 0, "")
	// pdf.CellFormat(VAlue_width, data_height, EMPTY, "1", 1, "R", false, 0, "")

	// Procedure
	// // // -------------------------------------- Procedure ------------------------------------------

	BlendProcedure := measured_product.ProductBlend.GetProcedure()
	if Procedure_margin+float64(len(BlendProcedure)+1)*Procedure_height+pdf.GetY()+btm_margin > page_height {
		pdf.AddPage()
	}

	pdf.Ln(Procedure_margin)
	pdf.SetFontStyle("B")
	pdf.SetFontSize(12)
	pdf.CellFormat(Procedure_width, Procedure_height, measured_product.Product_name+" Procedure", "1", 1, "C", false, 0, "")

	// pdf.SetFontStyle("")
	pdf.SetFontStyle("B")
	pdf.SetFontSize(10)
	for i, Procedure_step := range BlendProcedure {
		pdf.CellFormat(ORDER_width, Procedure_height, strconv.Itoa(i+1), "", 0, "", false, 0, "")
		pdf.CellFormat(Procedure_width, Procedure_height, Procedure_step, "", 1, "", false, 0, "")

	}

	/*






		[compiler] (MismatchedTypes) invalid operation: len(BlendProcedure) * Procedure_height (mismatched types int and )







	*/
	// pdf.CellFormat(field_width, field_height, strings.ToUpper(measured_product.Operators), "1", 0, "", false, 0, "")
	//
	// pdf.Cell(label_width, label_height, "Operators")
	// pdf.CellFormat(field_width, field_height, strings.ToUpper(measured_product.Operators), "1", 1, "", false, 0, "")

	//
	// // // -------------------------------------- SIDEBAR ------------------------------------------

	// ISOMERIC NAME:    PETROFLO FR-1579
	//
	// CUSTOMER NAME: V-SLICK 450A-M
	//
	// BATCH NUMBER		BSFR250919A			Date		2025-09-19
	// Start/End Time
	// Blend Vessel		SLZU 287035-5			OPERATORS		DAMIAN,BRANDON,ZACH,MARC
	// Final wt In tote, lbs		0

	// Blend Density, lbs/gal (QAQC Input)			9.100
	// Heel, lbs		0
	// CONFIDENTIAL		Blend Quantity, lbs		43,000			Total Quantity, Gal			0
	// Blend Quantity, Gal			4755

	// 	QC Testing - Top Sample
	// Test	Spec		Result
	// Clarity	Milky
	// Color	Grayish-White
	// Density (g/ml)	1.060 - 1.101
	// // Density (lb/gal)	8.846 - 9.180
	// String Test at 0.5gpt	60 ≤ seconds
	// Neat Viscosity	As Is
	// Order of Addition	Component	Lot Number		Railcar		Weight Percent		Weight Required		Actual Weight		Reference Gallons		Density, lbs/gal		Inches
	// 			TAG: 5879
	//
	// 			SEAL:238338
	//
	// 			TESTED BY:
	//
	// 			FINAL STRAP:
	//
	//
	// Order of Addition	Component	Lot Number		Railcar		Weight Percent		Weight Required		Actual Weight		Reference Gallons		Density, lbs/gal		Inches
	// HEEL						0		0		0		9.100		0"
	// 1	PETROFLO FR2220					35.00%		15050				1650		9.122		27.75"
	// 2	FLOJET DRMAX E562					65.00%		27950				3106		9.000		62.25"
	// TOTAL					100.00%		43000		0		4755
	//

	//
	// Package		Net Weight
	// 330 GAL		2800
	//
	//
	/*
		pdf.SetXY(label_col, lot_row)
		pdf.Cell(label_width, field_height, strings.ToUpper(measured_product.Lot_number))
		pdf.CellFormat(unit_width, field_height, strings.ToUpper(measured_product.Sample_point), "", 0, "R", false, 0, "")*/

	log.Println("Info: Saving to: ", file_path)
	err := pdf.OutputFileAndClose(file_path)
	return file_path, err
}
