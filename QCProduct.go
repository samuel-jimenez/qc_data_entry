package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/samuel-jimenez/whatsupdocx/docx"

	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"

	// "github.com/samuel-jimenez/whatsupdocx/docx"

	"github.com/samuel-jimenez/windigo"
)

/*
 * QCProduct
 *
 */

type QCProduct struct {
	Product
	Appearance     ProductAppearance
	product_type   Discrete
	container_type Discrete
	PH             Range
	SG             Range
	Density        Range
	String_test    Range
	Viscosity      Range
	Update         func()
}

func (product *QCProduct) reset() {
	var empty_product QCProduct
	empty_product.Product = product.Product
	empty_product.Update = product.Update
	*product = empty_product
}

func (product *QCProduct) select_product_details() {

	err := db_select_product_details.QueryRow(product.product_id).Scan(
		&product.product_type, &product.container_type, &product.Appearance,
		&product.PH.Min, &product.PH.Target, &product.PH.Max,
		&product.SG.Min, &product.SG.Target, &product.SG.Max,
		&product.Density.Min, &product.Density.Target, &product.Density.Max,
		&product.String_test.Min, &product.String_test.Target, &product.String_test.Max,
		&product.Viscosity.Min, &product.Viscosity.Target, &product.Viscosity.Max,
	)
	if err != nil {
		log.Printf("Error: %q: %s\n", err, "select_product_details")
	}
}

func (product *QCProduct) select_product_coa_details() {

	err := db_select_product_coa_details.QueryRow(product.product_id).Scan(
		&product.Appearance,
		&product.PH.Min, &product.PH.Target, &product.PH.Max,
		&product.SG.Min, &product.SG.Target, &product.SG.Max,
		&product.Density.Min, &product.Density.Target, &product.Density.Max,
		&product.String_test.Min, &product.String_test.Target, &product.String_test.Max,
		&product.Viscosity.Min, &product.Viscosity.Target, &product.Viscosity.Max,
	)
	if err != nil {
		log.Printf("Error: %q: %s\n", err, "select_product_coa_details")
	}
}

func (product QCProduct) _upsert(db_upsert_statement *sql.Stmt) {

	db_insert_appearance.Exec(product.Appearance)
	db_upsert_product_type.Exec(product.product_id, product.product_type)
	_, err := db_upsert_statement.Exec(
		product.product_id, product.product_type, product.Appearance,
		product.PH.Min, product.PH.Target, product.PH.Max,
		product.SG.Min, product.SG.Target, product.SG.Max,
		product.Density.Min, product.Density.Target, product.Density.Max,
		product.String_test.Min, product.String_test.Target, product.String_test.Max,
		product.Viscosity.Min, product.Viscosity.Target, product.Viscosity.Max,
	)
	if err != nil {
		log.Printf("Error: %q: %s\n", err, "upsert")
	}
	//TODO?
	// id, err := result.LastInsertId()
	// log.Println("upsert", id, err)

	// result.LastInsertId()
	// return product_type_id_default, product_name_customer_default

}
func (product QCProduct) upsert()     { product._upsert(db_upsert_product_details) }
func (product QCProduct) upsert_coa() { product._upsert(db_upsert_product_coa_details) }

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

func write_CoA_row_fmt(table *docx.Table, title, units string, spec Range, result nullable.NullFloat64, format_fn func(float64) string) {
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

	product.select_product_coa_details()
	write_CoA_row(table, Appearance_title, Appearance_units, product.Appearance.String, visual)
	write_CoA_row_fmt(table, ph_title, "", product.PH, product.Product.PH, formats.Format_ph)
	write_CoA_row_fmt(table, sg_title, formats.SG_UNITS, product.SG, product.Product.SG, Format_sg)
	write_CoA_row_fmt(table, Density_title, formats.DENSITY_UNITS, product.Density, product.Product.Density, formats.Format_density)
	write_CoA_row_fmt(table, string_title, formats.STRING_UNITS, product.String_test, product.Product.String_test, formats.Format_string_test)
	write_CoA_row_fmt(table, viscosity_title, formats.VISCOSITY_UNITS, product.Viscosity, product.Product.Viscosity, formats.Format_viscosity)
}

func (product *QCProduct) edit(
	product_type Discrete,
	Appearance ProductAppearance,
	PH,
	SG,
	Density,
	String_test,
	Viscosity Range,
) {
	product.product_type = product_type
	product.Appearance = Appearance
	product.PH = PH
	product.SG = SG
	product.Density = Density
	product.String_test = String_test
	product.Viscosity = Viscosity
}

func (product *QCProduct) show_ranges_window() {

	rangeWindow := windigo.NewForm(nil)
	var WindowText string
	// rangeWindow.SetTranslucentBackground()

	if product.Product_name_customer != "" {
		WindowText = fmt.Sprintf("%s (%s)", product.Product_type, product.Product_name_customer)
	} else {
		WindowText = product.Product_type
	}
	rangeWindow.SetSize(RANGES_WINDOW_WIDTH,
		RANGES_WINDOW_HEIGHT) // (width, height)
	rangeWindow.SetText(WindowText)

	dock := windigo.NewSimpleDock(rangeWindow)
	dock.SetPaddingsAll(RANGES_WINDOW_PADDING)
	dock.SetPaddingTop(RANGES_PADDING)

	prod_label := windigo.NewLabel(rangeWindow)
	prod_label.SetText(WindowText)
	prod_label.SetSize(OFF_AXIS, RANGES_FIELD_SMALL_HEIGHT)

	radio_dock := BuildNewDiscreteView(rangeWindow, "Type", product.product_type, []string{"Water Based", "Oil Based", "Friction Reducer"})
	radio_dock.SetSize(OFF_AXIS, DISCRETE_FIELD_HEIGHT)
	radio_dock.SetItemSize(PRODUCT_TYPE_WIDTH)
	radio_dock.SetPaddingsAll(GROUPBOX_CUSHION)

	coa_field := windigo.NewCheckBox(rangeWindow)
	coa_field.SetText("Save to COA published ranges")
	coa_field.SetMarginsAll(ERROR_MARGIN)
	coa_field.SetSize(OFF_AXIS, RANGES_FIELD_SMALL_HEIGHT)

	appearance_dock := BuildNewProductAppearanceView(rangeWindow, "Appearance", product.Appearance)

	labels := NewTextDock(rangeWindow, []string{"", "Min", "Target", "Max"})
	labels.SetMarginsAll(RANGES_PADDING)
	labels.SetDockSize(RANGES_FIELD_WIDTH, RANGES_FIELD_SMALL_HEIGHT)
	//TODO center
	//TODO layout split n

	ph_dock := BuildNewRangeView(rangeWindow, "pH", product.PH, formats.Format_ranges_ph)
	ph_dock.SetLabeledSize(LABEL_WIDTH, RANGES_FIELD_WIDTH, RANGES_FIELD_HEIGHT)

	sg_dock := BuildNewRangeView(rangeWindow, "Specific Gravity", product.SG, formats.Format_ranges_sg)
	sg_dock.SetLabeledSize(LABEL_WIDTH, RANGES_FIELD_WIDTH, RANGES_FIELD_HEIGHT)

	density_dock := BuildNewRangeView(rangeWindow, "Density", product.Density, formats.Format_ranges_density)
	density_dock.SetLabeledSize(LABEL_WIDTH, RANGES_FIELD_WIDTH, RANGES_FIELD_HEIGHT)

	string_dock := BuildNewRangeView(rangeWindow, "String Test \n\t at 0.5gpt", product.String_test, formats.Format_ranges_string_test)
	string_dock.SetLabeledSize(LABEL_WIDTH, RANGES_FIELD_WIDTH, RANGES_FIELD_HEIGHT)
	//TODO store string_amt "at 0.5gpt"

	visco_dock := BuildNewRangeView(rangeWindow, "Viscosity", product.Viscosity, formats.Format_ranges_viscosity)
	visco_dock.SetLabeledSize(LABEL_WIDTH, RANGES_FIELD_WIDTH, RANGES_FIELD_HEIGHT)

	exit := func() {
		rangeWindow.Close()
		windigo.Exit()
	}
	save := func() {
		product.edit(
			radio_dock.Get(),
			appearance_dock.Get(),
			ph_dock.Get(),
			sg_dock.Get(),
			density_dock.Get(),
			string_dock.Get(),
			visco_dock.Get(),
		)

		product.upsert()
		show_status("QC Data Updated")
		product.Update()
		exit()
	}
	save_coa := func() {
		var coa_product QCProduct
		coa_product.Product = product.Product
		coa_product.edit(
			radio_dock.Get(),
			appearance_dock.Get(),
			ph_dock.Get(),
			sg_dock.Get(),
			density_dock.Get(),
			string_dock.Get(),
			visco_dock.Get(),
		)

		coa_product.upsert_coa()
		show_status("COA Data Updated")
		product.select_product_details()
		product.Update()
		exit()
	}

	try_save := func() {
		if radio_dock.Get().Valid {
			radio_dock.Ok()
			if coa_field.Checked() {
				save_coa()
			} else {
				save()
			}
		} else {
			radio_dock.Error()
		}
	}

	load_coa := func() {
		product.reset()
		product.select_product_coa_details()
		appearance_dock.Set(product.Appearance)
		ph_dock.Set(product.PH)
		sg_dock.Set(product.SG)
		density_dock.Set(product.Density)
		string_dock.Set(product.String_test)
		visco_dock.Set(product.Viscosity)
	}

	button_dock := NewButtonDock(rangeWindow, []string{"OK", "Cancel", "Load CoA Data"}, []func(){try_save, exit, load_coa})
	button_dock.SetDockSize(RANGES_BUTTON_WIDTH, RANGES_BUTTON_HEIGHT)
	button_dock.SetMarginLeft(RANGES_PADDING)

	dock.Dock(prod_label, windigo.Top)
	dock.Dock(radio_dock, windigo.Top)
	dock.Dock(coa_field, windigo.Top)
	dock.Dock(appearance_dock, windigo.Top)

	dock.Dock(labels, windigo.Top)
	dock.Dock(ph_dock, windigo.Top)
	dock.Dock(sg_dock, windigo.Top)
	dock.Dock(density_dock, windigo.Top)
	dock.Dock(string_dock, windigo.Top)
	dock.Dock(visco_dock, windigo.Top)
	dock.Dock(button_dock, windigo.Top)

	rangeWindow.Center()
	rangeWindow.Show()
	rangeWindow.OnClose().Bind(
		func(arg *windigo.Event) {
			exit()
		})
	rangeWindow.RunMainLoop() // Must call to start event loop.
}

func (product QCProduct) Check(data Product) bool {
	return product.PH.Check(data.PH.Float64) && product.SG.Check(data.SG.Float64)
}
