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
	Container_type Discrete
	PH             datatypes.Range
	SG             datatypes.Range
	Density        datatypes.Range
	String_test    datatypes.Range
	Viscosity      datatypes.Range
	Update         func()
}

func NewQCProduct() *QCProduct {
	qc_product := new(QCProduct)
	qc_product.Product_id = DB.INVALID_ID
	qc_product.Lot_id = DB.DEFAULT_LOT_ID
	return qc_product
}

func (product *QCProduct) Reset() {
	var empty_product QCProduct
	empty_product.Product = product.Product
	empty_product.Update = product.Update
	*product = empty_product
}

func (product *QCProduct) Select_product_details() {

	err := DB.DB_Select_product_details.QueryRow(product.Product_id).Scan(
		&product.Product_type, &product.Container_type, &product.Appearance,
		&product.PH.Min, &product.PH.Target, &product.PH.Max,
		&product.SG.Min, &product.SG.Target, &product.SG.Max,
		&product.Density.Min, &product.Density.Target, &product.Density.Max,
		&product.String_test.Min, &product.String_test.Target, &product.String_test.Max,
		&product.Viscosity.Min, &product.Viscosity.Target, &product.Viscosity.Max,
	)
	if err != nil {
		log.Printf("Error: [%s]: %q\n",  "Select_product_details",  err)
	}
}

func (product *QCProduct) Select_product_coa_details() {

	err := DB.DB_Select_product_coa_details.QueryRow(product.Product_id).Scan(
		&product.Appearance,
		&product.PH.Min, &product.PH.Target, &product.PH.Max,
		&product.SG.Min, &product.SG.Target, &product.SG.Max,
		&product.Density.Min, &product.Density.Target, &product.Density.Max,
		&product.String_test.Min, &product.String_test.Target, &product.String_test.Max,
		&product.Viscosity.Min, &product.Viscosity.Target, &product.Viscosity.Max,
	)
	if err != nil {
		log.Printf("Error: [%s]: %q\n",  "select_product_coa_details",  err)
	}
}

func (product QCProduct) _upsert(db_upsert_statement *sql.Stmt) {

	DB.DB_Insert_appearance.Exec(product.Appearance)
	DB.DB_Upsert_product_type.Exec(product.Product_id, product.Product_type)
	_, err := db_upsert_statement.Exec(
		product.Product_id, product.Product_type, product.Appearance,
		product.PH.Min, product.PH.Target, product.PH.Max,
		product.SG.Min, product.SG.Target, product.SG.Max,
		product.Density.Min, product.Density.Target, product.Density.Max,
		product.String_test.Min, product.String_test.Target, product.String_test.Max,
		product.Viscosity.Min, product.Viscosity.Target, product.Viscosity.Max,
	)
	if err != nil {
		log.Printf("Error: [%s]: %q\n",  "upsert",  err)
	}
	//TODO?
	// id, err := result.LastInsertId()
	// log.Println("upsert", id, err)

	// result.LastInsertId()
	// return product_type_id_default, product_name_customer_default

}
func (product QCProduct) Upsert()     { product._upsert(DB.DB_Upsert_product_details) }
func (product QCProduct) Upsert_coa() { product._upsert(DB.DB_Upsert_product_coa_details) }

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

func (product QCProduct) write_CoA_rows(table *docx.Table) {
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
	if product.Visual {
		visual = "PASS"
	}
	Format_sg := func(sg float64) string {
		return formats.Format_sg(sg, product.Product.PH.Valid)
	}

	// TODO: use QCProduct so this is not required
	product.Select_product_coa_details()

	write_CoA_row(table, Appearance_title, Appearance_units, product.Appearance.String, visual)
	write_CoA_row_fmt(table, ph_title, "", product.PH, product.Product.PH, formats.Format_ph)
	write_CoA_row_fmt(table, sg_title, formats.SG_UNITS, product.SG, product.Product.SG, Format_sg)
	write_CoA_row_fmt(table, Density_title, formats.DENSITY_UNITS, product.Density, product.Product.Density, formats.Format_density)
	write_CoA_row_fmt(table, string_title, formats.STRING_UNITS, product.String_test, product.Product.String_test, formats.Format_string_test)
	write_CoA_row_fmt(table, viscosity_title, formats.VISCOSITY_UNITS, product.Viscosity, product.Product.Viscosity, formats.Format_viscosity)
}

func (product *QCProduct) Edit(
	product_type Discrete,
	Appearance ProductAppearance,
	PH,
	SG,
	Density,
	String_test,
	Viscosity datatypes.Range,
) {
	product.Product_type = product_type
	product.Appearance = Appearance
	product.PH = PH
	product.SG = SG
	product.Density = Density
	product.String_test = String_test
	product.Viscosity = Viscosity
}

func (product QCProduct) Check(data Product) bool {
	return product.PH.Check(data.PH.Float64) && product.SG.Check(data.SG.Float64)
}
