package product

import (
	"path/filepath"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/xuri/excelize/v2"
)

// TODO xl-1501-3
// TODO: move these to own file

func write_cell_fmt(xl_file *excelize.File, sheet, cell string, result nullable.NullFloat64, format_fn func(float64) string) {
	if result.Valid {
		xl_file.SetCellValue(sheet, cell, format_fn(result.Float64))
	}
}

func write_cell_fmt_int64(xl_file *excelize.File, sheet, cell string, result nullable.NullInt64) {
	if result.Valid {
		xl_file.SetCellValue(sheet, cell, formats.FormatInt(result.Int64))
	}
}

func updateBSExcel(file_path string, products ...*MeasuredProduct) error {
	numSamples := len(products)
	if numSamples <= 0 {
		return nil
	}
	measured_product := products[0]
	Format_sg := func(sg float64) string {
		return formats.Format_sg(sg, measured_product.QCProduct.SG.Method.String == METHOD_DMA)
	}

	xl_file_glon := file_path + "/" + measured_product.Lot_number + "*"

	files, err := filepath.Glob(xl_file_glon)
	if err != nil {
		return err
	}

	xl_file := file_path + "/" + measured_product.Lot_number + ".xlsx"
	if len(files) > 0 {
		xl_file = files[0]
	}

	return WithOpenFile(xl_file, func(xl_file *excelize.File) error {
		for _, sheet_name := range xl_file.GetSheetList() {

			visible, _ := xl_file.GetSheetVisible(sheet_name)
			if visible {
				// cell, _ := xl_file.GetCellValue( sheet_name)
				rows, err := xl_file.GetRows(sheet_name)
				if err != nil {
					return err
				}

				if updateBSExcel_ROWS(xl_file, rows, sheet_name, Format_sg, numSamples, products) {
					break
				}
			}
		}
		return xl_file.Save()
	})
}

func updateBSExcel_ROWS(xl_file *excelize.File, rows [][]string, sheet_name string, Format_sg func(float64) string, numSamples int, products []*MeasuredProduct) bool {
	max_COL := 16
	var result_col, btm_result_col int
	found_rsult := false
	double_sample := false

ROWS:

	// 			Clarity
	// Color
	// Density (g/ml)
	// Density (lb/gal)
	// String Test at 0.5gpt
	// Neat Viscosity

	for i, row := range rows {
		for j := 0; j < len(row) && j < max_COL; j++ {
			cell := row[j]
			if strings.Contains(cell, "Result") {
				if !found_rsult {
					result_col = j + 1
					found_rsult = true

				} else if numSamples > 1 {
					btm_result_col = j + 1
					double_sample = true

					continue ROWS
				}
			}
			if found_rsult {
				// todo appreaecne
				if strings.Contains(cell, "pH") {
					coords, _ := excelize.CoordinatesToCellName(result_col, i+1)

					write_cell_fmt(xl_file, sheet_name, coords, products[0].PH, formats.Format_ph)
					if double_sample {
						coords, _ := excelize.CoordinatesToCellName(btm_result_col, i+1)
						write_cell_fmt(xl_file, sheet_name, coords, products[1].PH, formats.Format_ph)

					}
					continue ROWS
				}
				if strings.Contains(cell, "Specific Gravity") || strings.Contains(cell, "Density (g/ml)") {
					coords, _ := excelize.CoordinatesToCellName(result_col, i+1)

					// xl_file.SetCellValue(sheet_name, coords, products[0].SG)
					write_cell_fmt(xl_file, sheet_name, coords, products[0].SG, Format_sg)

					if double_sample {
						coords, _ := excelize.CoordinatesToCellName(btm_result_col, i+1)

						write_cell_fmt(xl_file, sheet_name, coords, products[1].SG, Format_sg)

					}
					continue ROWS
				}
				if strings.Contains(cell, "Density (lb/gal)") {
					coords, _ := excelize.CoordinatesToCellName(result_col, i+1)

					write_cell_fmt(xl_file, sheet_name, coords, products[0].Density, formats.Format_density)
					if double_sample {
						coords, _ := excelize.CoordinatesToCellName(btm_result_col, i+1)
						write_cell_fmt(xl_file, sheet_name, coords, products[1].Density, formats.Format_density)

					}
					continue ROWS
				}
				if strings.Contains(cell, "String") {
					coords, _ := excelize.CoordinatesToCellName(result_col, i+1)

					write_cell_fmt_int64(xl_file, sheet_name, coords, products[0].String_test)
					if double_sample {
						coords, _ := excelize.CoordinatesToCellName(btm_result_col, i+1)
						write_cell_fmt_int64(xl_file, sheet_name, coords, products[1].String_test)

					}
					continue ROWS
				}
				if strings.Contains(cell, "Viscosity") {
					coords, _ := excelize.CoordinatesToCellName(result_col, i+1)

					write_cell_fmt_int64(xl_file, sheet_name, coords, products[0].Viscosity)
					if double_sample {
						coords, _ := excelize.CoordinatesToCellName(btm_result_col, i+1)
						write_cell_fmt_int64(xl_file, sheet_name, coords, products[1].Viscosity)

					}
					continue ROWS
				}
			}

		}
	}
	return found_rsult
}
