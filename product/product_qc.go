package product

import (
	"codeberg.org/go-pdf/fpdf"
	"github.com/samuel-jimenez/qc_data_entry/datatypes"
	"github.com/samuel-jimenez/qc_data_entry/formats"
)

// TODO Method
func write_QC_row(pdf *fpdf.Fpdf, height float64, title, units, Method, spec, result_0, result_1 string) {
	pdf.CellFormat(QC_TEST_WIDTH, height, title, "1", 0, "", false, 0, "")
	// pdf.CellFormat(qc_test_width, qc_header_height, units, "1",0, "", false, 0, "")
	pdf.CellFormat(QC_SPEC_WIDTH, height, spec, "1", 0, "", false, 0, "")
	pdf.CellFormat(QC_RESULT_WIDTH, height, result_0, "1", 0, "", false, 0, "")
	pdf.CellFormat(QC_RESULT_WIDTH, height, result_1, "1", 1, "", false, 0, "")
}

func write_QC_row_fmt(pdf *fpdf.Fpdf, height float64, title, units string, spec datatypes.Range, format_fn func(float64) string) {
	if spec.Valid {
		write_QC_row(pdf, height, title, units, spec.Method.String, spec.QC(format_fn), EMPTY, EMPTY)
	}
}

func (qc_product QCProduct) Write_QC_rows(pdf *fpdf.Fpdf, qc_header_height,
	qc_height float64,
) {
	Appearance_units := "Pass/fail"
	Appearance_Method := "Visual"

	// TODO
	Format_sg := func(sg float64) string {
		return formats.Format_sg(sg, qc_product.SG.Method.String == METHOD_DMA)
	}

	pdf.SetFontSize(12)
	pdf.CellFormat(QC_TEST_WIDTH+QC_SPEC_WIDTH+QC_RESULT_WIDTH+QC_RESULT_WIDTH, qc_header_height, "QC Testing", "1", 1, "C", false, 0, "")

	pdf.SetFontSize(10)
	// TODO Method
	write_QC_row(pdf, qc_header_height, "Test", "Units", " Method", "Spec", "Result - Top", "Result - Bottom")

	pdf.SetFontStyle("")
	write_QC_row(pdf, qc_height, formats.APPEARANCE_TEXT, Appearance_units, Appearance_Method, qc_product.Appearance.String, EMPTY, EMPTY)
	write_QC_row_fmt(pdf, qc_height, formats.PH_TEXT, "", qc_product.PH, formats.Format_ph)
	write_QC_row_fmt(pdf, qc_height, formats.SG_TEXT, formats.SG_UNITS, qc_product.SG, Format_sg)
	write_QC_row_fmt(pdf, qc_height, formats.DENSITY_TEXT, formats.DENSITY_UNITS, qc_product.Density, formats.Format_density)
	write_QC_row_fmt(pdf, qc_height, formats.STRING_TEXT_MINI, formats.STRING_UNITS, qc_product.String_test, formats.Format_string_test)
	write_QC_row_fmt(pdf, qc_height, formats.VISCOSITY_TEXT, formats.VISCOSITY_UNITS, qc_product.Viscosity, formats.Format_viscosity)
}
