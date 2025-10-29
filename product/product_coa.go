package product

import (
	"fmt"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/datatypes"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/whatsupdocx"
	"github.com/samuel-jimenez/whatsupdocx/docx"
)

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

func (qc_product MeasuredProduct) write_CoA_rows(table *docx.Table) {
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
		return formats.Format_sg(sg, qc_product.QCProduct.SG.Method.String == METHOD_DMA)
	}

	// TODO: ensure this is not required
	// qc_product.Select_product_details()

	write_CoA_row(table, formats.APPEARANCE_TEXT, Appearance_units, qc_product.Appearance.String, visual)
	write_CoA_row_fmt(table, formats.PH_TEXT, "", qc_product.QCProduct.PH, qc_product.PH, formats.Format_ph)
	write_CoA_row_fmt(table, formats.SG_TEXT, formats.SG_UNITS, qc_product.QCProduct.SG, qc_product.SG, Format_sg)
	write_CoA_row_fmt(table, formats.DENSITY_TEXT, formats.DENSITY_UNITS, qc_product.QCProduct.Density, qc_product.Density, formats.Format_density)
	// write_CoA_row_fmt(table, string_title, formats.STRING_UNITS, product.String_test, product.Product.String_test, formats.Format_string_test)
	write_CoA_row_fmt_int64(table, formats.STRING_TEXT_MINI, formats.STRING_UNITS, qc_product.QCProduct.String_test, qc_product.String_test, formats.Format_string_test)
	// write_CoA_row_fmt(table, viscosity_title, formats.VISCOSITY_UNITS, product.Viscosity, product.Product.Viscosity, formats.Format_viscosity)
	write_CoA_row_fmt_int64(table, formats.VISCOSITY_TEXT, formats.VISCOSITY_UNITS, qc_product.QCProduct.Viscosity, qc_product.Viscosity, formats.Format_viscosity)

}

func (measured_product MeasuredProduct) get_coa_template() string {
	//TODO db
	// return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, product.COA_TEMPLATE)
	// Product_type split
	//TODO fix this
	product_moniker := strings.Split(measured_product.Product_name, " ")[0]
	if product_moniker == "PETROFLO" {
		return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA-PETROFLO.docx")
	}
	return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA.docx")
}

func (measured_product MeasuredProduct) get_coa_name() string {
	return fmt.Sprintf("%s/%s", config.COA_FILEPATH, measured_product.get_base_filename("docx"))
}

func (measured_product MeasuredProduct) export_CoA() error {
	var (
		p_title   = "[PRODUCT_NAME]"
		Coa_title = "Parameter"
	)
	terms := []string{
		COA_LOT_TITLE,
	}

	template_file := measured_product.get_coa_template()
	output_file := measured_product.get_coa_name()

	doc, err := whatsupdocx.OpenDocument(template_file)
	if err != nil {
		return err
	}

	product_name := measured_product.Product_name_customer
	if product_name == "" {
		product_name = measured_product.Product_name
	}
	for _, item := range doc.Document.Body.Children {
		if para := item.Paragraph; para != nil {
			measured_product.searchCOAPara(para, p_title, product_name)
		}

		if table := item.Table; table != nil {
			measured_product.searchCOATable(table, Coa_title, terms)
		}
	}

	// save to file
	return doc.SaveTo(output_file)
}

func (measured_product MeasuredProduct) searchCOAPara(para *docx.Paragraph, p_title, product_name string) {
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

func (measured_product MeasuredProduct) searchCOATable(table *docx.Table, Coa_title string, terms []string) {
	for _, row := range table.RowContents {
		if measured_product.searchCOARow(table, row.Row, Coa_title, terms) {
			return
		}
	}
}

func (measured_product MeasuredProduct) searchCOARow(table *docx.Table, row *docx.Row, Coa_title string, terms []string) (finished bool) {
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
						measured_product.write_CoA_rows(table)
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
					field.AddText(measured_product.Lot_number)
					return false
				}
			}
		}
	}
	return false
}
