package product

import (
	"log"

	"codeberg.org/go-pdf/fpdf"
	"github.com/samuel-jimenez/whatsupdocx/docx"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"

	"github.com/samuel-jimenez/qc_data_entry/datatypes"
	// "github.com/samuel-jimenez/whatsupdocx/docx"
)

var (
	QC_TEST_WIDTH,
	QC_SPEC_WIDTH,
	QC_RESULT_WIDTH float64

	EMPTY string
)

/*
 * QCProduct
 *
 */

type QCProduct struct {
	MeasuredProduct
	Appearance     ProductAppearance
	Product_type   Discrete
	Container_type ProductContainerType // bs.container_types
	PH             datatypes.Range
	SG             datatypes.Range
	Density        datatypes.Range
	String_test    datatypes.Range
	Viscosity      datatypes.Range
	UpdateFN       func(*QCProduct)
}

func NewQCProduct() *QCProduct {
	qc_product := new(QCProduct)
	qc_product.Product_id = DB.INVALID_ID
	qc_product.Product_Lot_id = DB.DEFAULT_LOT_ID
	qc_product.Lot_id = DB.DEFAULT_LOT_ID
	return qc_product
}

func (qc_product *QCProduct) ResetQC() {
	var empty_product QCProduct
	empty_product.MeasuredProduct = qc_product.MeasuredProduct
	empty_product.UpdateFN = qc_product.UpdateFN
	*qc_product = empty_product
}

func (qc_product *QCProduct) Select_product_details() {
	proc_name := "Select_product_details"
	DB.Select_ErrNoRows(
		proc_name,
		DB.DB_Select_product_details.QueryRow(qc_product.Product_id),
		&qc_product.Product_type, &qc_product.Container_type, &qc_product.Appearance,
		&qc_product.PH.Valid, &qc_product.PH.Publish_p, &qc_product.PH.Min, &qc_product.PH.Target, &qc_product.PH.Max,
		&qc_product.SG.Valid, &qc_product.SG.Publish_p, &qc_product.SG.Min, &qc_product.SG.Target, &qc_product.SG.Max,
		&qc_product.Density.Valid, &qc_product.Density.Publish_p, &qc_product.Density.Min, &qc_product.Density.Target, &qc_product.Density.Max,
		&qc_product.String_test.Valid, &qc_product.String_test.Publish_p, &qc_product.String_test.Min, &qc_product.String_test.Target, &qc_product.String_test.Max,
		&qc_product.Viscosity.Valid, &qc_product.Viscosity.Publish_p, &qc_product.Viscosity.Min, &qc_product.Viscosity.Target, &qc_product.Viscosity.Max,
	)
}

func (qc_product QCProduct) Upsert() {

	DB.DB_Insert_appearance.Exec(qc_product.Appearance)
	DB.DB_Upsert_product_type.Exec(qc_product.Product_id, qc_product.Product_type)
	_, err := DB.DB_Upsert_product_details.Exec(
		qc_product.Product_id, qc_product.Product_type, qc_product.Appearance,
		qc_product.PH.Valid, qc_product.PH.Publish_p, qc_product.PH.Min, qc_product.PH.Target, qc_product.PH.Max,
		qc_product.SG.Valid, qc_product.SG.Publish_p, qc_product.SG.Min, qc_product.SG.Target, qc_product.SG.Max,
		qc_product.Density.Valid, qc_product.Density.Publish_p, qc_product.Density.Min, qc_product.Density.Target, qc_product.Density.Max,
		qc_product.String_test.Valid, qc_product.String_test.Publish_p, qc_product.String_test.Min, qc_product.String_test.Target, qc_product.String_test.Max,
		qc_product.Viscosity.Valid, qc_product.Viscosity.Publish_p, qc_product.Viscosity.Min, qc_product.Viscosity.Target, qc_product.Viscosity.Max,
	)
	if err != nil {
		log.Printf("Error: [%s]: %q\n", "upsert", err)
	}
	//TODO?
	// id, err := result.LastInsertId()
	// log.Println("upsert", id, err)

	// result.LastInsertId()
	// return product_type_id_default, product_name_customer_default

}

func write_CoA_cell(row *docx.Row, value string) {
	cell := row.AddCell()
	cell.AddParagraph(value)
}

func write_CoA_row(table *docx.Table, title, units, spec, result string) {
	row := table.AddRow()
	write_CoA_cell(row, title)
	write_CoA_cell(row, units)
	write_CoA_cell(row, spec)
	write_CoA_cell(row, result)
}

func write_CoA_row_fmt(table *docx.Table, title, units string, spec datatypes.Range, result nullable.NullFloat64, format_fn func(float64) string) {
	if result.Valid {
		write_CoA_row(table, title, units, spec.CoA(format_fn), format_fn(result.Float64))
	}
}

func write_CoA_row_fmt_int64(table *docx.Table, title, units string, spec datatypes.Range, result nullable.NullInt64, format_fn func(float64) string) {
	if result.Valid {
		write_CoA_row(table, title, units, spec.CoA(format_fn), formats.FormatInt(result.Int64))
	}
}

func (qc_product QCProduct) write_CoA_rows(table *docx.Table) {
	var (
		Appearance_units = "Pass/fail"
	)

	var (
		// visual        = "PASS"
		visual = "FAIL"
	)
	if qc_product.Visual {
		visual = "PASS"
	}
	Format_sg := func(sg float64) string {
		return formats.Format_sg(sg, qc_product.MeasuredProduct.PH.Valid)
	}

	// TODO: use QCProduct so this is not required
	qc_product.Select_product_details()

	write_CoA_row(table, formats.APPEARANCE_TEXT, Appearance_units, qc_product.Appearance.String, visual)
	write_CoA_row_fmt(table, formats.PH_TEXT, "", qc_product.PH, qc_product.MeasuredProduct.PH, formats.Format_ph)
	write_CoA_row_fmt(table, formats.SG_TEXT, formats.SG_UNITS, qc_product.SG, qc_product.MeasuredProduct.SG, Format_sg)
	write_CoA_row_fmt(table, formats.DENSITY_TEXT, formats.DENSITY_UNITS, qc_product.Density, qc_product.MeasuredProduct.Density, formats.Format_density)
	// write_CoA_row_fmt(table, string_title, formats.STRING_UNITS, product.String_test, product.Product.String_test, formats.Format_string_test)
	write_CoA_row_fmt_int64(table, formats.STRING_TEXT_MINI, formats.STRING_UNITS, qc_product.String_test, qc_product.MeasuredProduct.String_test, formats.Format_string_test)
	// write_CoA_row_fmt(table, viscosity_title, formats.VISCOSITY_UNITS, product.Viscosity, product.Product.Viscosity, formats.Format_viscosity)
	write_CoA_row_fmt_int64(table, formats.VISCOSITY_TEXT, formats.VISCOSITY_UNITS, qc_product.Viscosity, qc_product.MeasuredProduct.Viscosity, formats.Format_viscosity)

}

// TODO Method
func write_QC_row(pdf *fpdf.Fpdf, height float64, title, units, spec, result_0, result_1 string) {

	pdf.CellFormat(QC_TEST_WIDTH, height, title, "1", 0, "", false, 0, "")
	// pdf.CellFormat(qc_test_width, qc_header_height, units, "1",0, "", false, 0, "")
	pdf.CellFormat(QC_SPEC_WIDTH, height, spec, "1", 0, "", false, 0, "")
	pdf.CellFormat(QC_RESULT_WIDTH, height, result_0, "1", 0, "", false, 0, "")
	pdf.CellFormat(QC_RESULT_WIDTH, height, result_1, "1", 1, "", false, 0, "")

}

func write_QC_row_fmt(pdf *fpdf.Fpdf, height float64, title, units string, spec datatypes.Range, format_fn func(float64) string) {
	if spec.Valid {
		write_QC_row(pdf, height, title, units, spec.QC(format_fn), "", "")
	}
}

func (qc_product QCProduct) Write_QC_rows(pdf *fpdf.Fpdf, qc_header_height,
	qc_height float64) {

	var (
		Appearance_units = "Pass/fail"
	)

	//TODO
	Format_sg := func(sg float64) string {
		return formats.Format_sg(sg, qc_product.PH.Valid)
	}

	pdf.SetFontSize(12)
	pdf.CellFormat(QC_TEST_WIDTH+QC_SPEC_WIDTH+QC_RESULT_WIDTH+QC_RESULT_WIDTH, qc_header_height, "QC Testing", "1", 1, "C", false, 0, "")

	pdf.SetFontSize(10)
	// TODO Method
	write_QC_row(pdf, qc_header_height, "Test", "units", "Spec", "Result - Top", "Result - Bottom")

	pdf.SetFontStyle("")
	write_QC_row(pdf, qc_height, formats.APPEARANCE_TEXT, Appearance_units, qc_product.Appearance.String, EMPTY, EMPTY)
	write_QC_row_fmt(pdf, qc_height, formats.PH_TEXT, "", qc_product.PH, formats.Format_ph)
	write_QC_row_fmt(pdf, qc_height, formats.SG_TEXT, formats.SG_UNITS, qc_product.SG, Format_sg)
	write_QC_row_fmt(pdf, qc_height, formats.DENSITY_TEXT, formats.DENSITY_UNITS, qc_product.Density, formats.Format_density)
	write_QC_row_fmt(pdf, qc_height, formats.STRING_TEXT_MINI, formats.STRING_UNITS, qc_product.String_test, formats.Format_string_test)
	write_QC_row_fmt(pdf, qc_height, formats.VISCOSITY_TEXT, formats.VISCOSITY_UNITS, qc_product.Viscosity, formats.Format_viscosity)
}

func (qc_product *QCProduct) Edit(
	product_type Discrete,
	Appearance ProductAppearance,
	PH,
	SG,
	Density,
	String_test,
	Viscosity datatypes.Range,
) {
	qc_product.Product_type = product_type
	qc_product.Appearance = Appearance
	qc_product.PH = PH
	qc_product.SG = SG
	qc_product.Density = Density
	qc_product.String_test = String_test
	qc_product.Viscosity = Viscosity
}

func (qc_product QCProduct) Check(data MeasuredProduct) bool {
	return qc_product.PH.Check(data.PH.Float64) && qc_product.SG.Check(data.SG.Float64)
}

func (qc_product *QCProduct) SetUpdate(UpdateFN func(*QCProduct)) {
	qc_product.UpdateFN = UpdateFN
}

func (qc_product *QCProduct) Update() {
	if qc_product.UpdateFN == nil {
		return
	}
	qc_product.UpdateFN(qc_product)
}
