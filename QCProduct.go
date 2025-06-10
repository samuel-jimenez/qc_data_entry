package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/formats"
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

	rangeWindow.SetSize(800, 600) // (width, height)
	rangeWindow.SetText(WindowText)

	dock := windigo.NewSimpleDock(rangeWindow)
	dock.SetPaddingsAll(5)

	dock.SetPaddingTop(10)
	prod_label := windigo.NewLabel(rangeWindow)

	prod_label.SetText(WindowText)

	radio_dock := BuildNewProductTypeView(rangeWindow, "Type", product.product_type, []string{"Water Based", "Oil Based", "Friction Reducer"})

	coa_field := windigo.NewCheckBox(rangeWindow)
	coa_field.SetText("Save to COA published ranges")
	coa_field.SetMarginsAll(ERROR_MARGIN)

	appearance_dock := BuildNewProductAppearanceView(rangeWindow, "Appearance", product.Appearance)

	labels := build_text_dock(rangeWindow, []string{"", "Min", "Target", "Max"})
	ph_dock := BuildNewRangeView(rangeWindow, "pH", product.PH, formats.Format_ranges_ph)
	sg_dock := BuildNewRangeView(rangeWindow, "Specific Gravity", product.SG, formats.Format_ranges_sg)
	density_dock := BuildNewRangeView(rangeWindow, "Density", product.Density, formats.Format_ranges_density)
	string_dock := BuildNewRangeView(rangeWindow, "String Test \n\t at 0.5gpt", product.String_test, formats.Format_ranges_string_test)
	//TODO store string_amt "at 0.5gpt"
	visco_dock := BuildNewRangeView(rangeWindow, "Viscosity", product.Viscosity, formats.Format_ranges_viscosity)

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
	button_dock := build_button_dock(rangeWindow, []string{"OK", "Cancel"}, []func(){try_save, exit})

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
