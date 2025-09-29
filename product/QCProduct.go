package product

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/whatsupdocx/docx"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"

	"github.com/samuel-jimenez/qc_data_entry/datatypes"
	// "github.com/samuel-jimenez/whatsupdocx/docx"
)

/*
 * QCProduct
 *
 */

type QCProduct struct {
	Product
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
	empty_product.Product = qc_product.Product
	empty_product.UpdateFN = qc_product.UpdateFN
	*qc_product = empty_product
}

func (qc_product *QCProduct) Select_product_details() {
	proc_name := "Select_product_details"
	DB.Select_Error(
		proc_name,
		DB.DB_Select_product_details.QueryRow(qc_product.Product_id),
		&qc_product.Product_type, &qc_product.Container_type, &qc_product.Appearance,
		&qc_product.PH.Min, &qc_product.PH.Target, &qc_product.PH.Max,
		&qc_product.SG.Min, &qc_product.SG.Target, &qc_product.SG.Max,
		&qc_product.Density.Min, &qc_product.Density.Target, &qc_product.Density.Max,
		&qc_product.String_test.Min, &qc_product.String_test.Target, &qc_product.String_test.Max,
		&qc_product.Viscosity.Min, &qc_product.Viscosity.Target, &qc_product.Viscosity.Max,
	)
}

func (qc_product *QCProduct) Select_product_coa_details() {
	proc_name := "Select_product_coa_details"
	DB.Select_ErrNoRows(
		proc_name,
		DB.DB_Select_product_coa_details.QueryRow(qc_product.Product_id),
		&qc_product.Appearance,
		&qc_product.PH.Min, &qc_product.PH.Target, &qc_product.PH.Max,
		&qc_product.SG.Min, &qc_product.SG.Target, &qc_product.SG.Max,
		&qc_product.Density.Min, &qc_product.Density.Target, &qc_product.Density.Max,
		&qc_product.String_test.Min, &qc_product.String_test.Target, &qc_product.String_test.Max,
		&qc_product.Viscosity.Min, &qc_product.Viscosity.Target, &qc_product.Viscosity.Max,
	)
}

func (qc_product QCProduct) _upsert(db_upsert_statement *sql.Stmt) {

	DB.DB_Insert_appearance.Exec(qc_product.Appearance)
	DB.DB_Upsert_product_type.Exec(qc_product.Product_id, qc_product.Product_type)
	_, err := db_upsert_statement.Exec(
		qc_product.Product_id, qc_product.Product_type, qc_product.Appearance,
		qc_product.PH.Min, qc_product.PH.Target, qc_product.PH.Max,
		qc_product.SG.Min, qc_product.SG.Target, qc_product.SG.Max,
		qc_product.Density.Min, qc_product.Density.Target, qc_product.Density.Max,
		qc_product.String_test.Min, qc_product.String_test.Target, qc_product.String_test.Max,
		qc_product.Viscosity.Min, qc_product.Viscosity.Target, qc_product.Viscosity.Max,
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
func (qc_product QCProduct) Upsert()     { qc_product._upsert(DB.DB_Upsert_product_details) }
func (qc_product QCProduct) Upsert_coa() { qc_product._upsert(DB.DB_Upsert_product_coa_details) }

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

		Appearance_title = "Appearance"
		// Appearance_title = "Clarity/Color"

		viscosity_title = "Viscosity"
		string_title    = "String"
		ph_title        = "pH"
		sg_title        = "Specific Gravity"
		Density_title   = "Density"
	)

	var (
		// visual        = "PASS"
		visual = "FAIL"
	)
	if qc_product.Visual {
		visual = "PASS"
	}
	Format_sg := func(sg float64) string {
		return formats.Format_sg(sg, qc_product.Product.PH.Valid)
	}

	// TODO: use QCProduct so this is not required
	qc_product.Select_product_coa_details()

	write_CoA_row(table, Appearance_title, Appearance_units, qc_product.Appearance.String, visual)
	write_CoA_row_fmt(table, ph_title, "", qc_product.PH, qc_product.Product.PH, formats.Format_ph)
	write_CoA_row_fmt(table, sg_title, formats.SG_UNITS, qc_product.SG, qc_product.Product.SG, Format_sg)
	write_CoA_row_fmt(table, Density_title, formats.DENSITY_UNITS, qc_product.Density, qc_product.Product.Density, formats.Format_density)
	// write_CoA_row_fmt(table, string_title, formats.STRING_UNITS, product.String_test, product.Product.String_test, formats.Format_string_test)
	write_CoA_row_fmt_int64(table, string_title, formats.STRING_UNITS, qc_product.String_test, qc_product.Product.String_test, formats.Format_string_test)
	// write_CoA_row_fmt(table, viscosity_title, formats.VISCOSITY_UNITS, product.Viscosity, product.Product.Viscosity, formats.Format_viscosity)
	write_CoA_row_fmt_int64(table, viscosity_title, formats.VISCOSITY_UNITS, qc_product.Viscosity, qc_product.Product.Viscosity, formats.Format_viscosity)

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

func (qc_product QCProduct) Check(data Product) bool {
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
