package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/windigo"
)

type QCProduct struct {
	Product
	product_type ProductType
	SG           Range
	PH           Range
	Density      Range
	String_test  Range
	Viscosity    Range
}

type ProductType struct {
	sql.NullInt32
}

func ProductTypeDefault() ProductType {
	return ProductType{}
}

func ProductTypeFromIndex(index int) ProductType {
	return ProductType{sql.NullInt32{Int32: int32(index) + 1, Valid: true}}
}

func (product_type ProductType) toIndex() int {
	return int(product_type.Int32 - 1)

}

type Range struct {
	Min    NullFloat64
	Target NullFloat64
	Max    NullFloat64
}

func NewRange(
	min NullFloat64,
	target NullFloat64,
	max NullFloat64,
) Range {
	return Range{min, target, max}
}

type ProductTypeView struct {
	*windigo.GroupPanel
	Get   func() ProductType
	Error func()
}

func NewProductTypeView(panel *windigo.GroupPanel, get func() ProductType, err func()) ProductTypeView {
	return ProductTypeView{panel, get, err}
}

func BuildNewProductTypeView(parent windigo.Controller, group_text string, field_data ProductType, labels []string) ProductTypeView {

	var buttons []*windigo.RadioButton
	panel := windigo.NewGroupPanel(parent)
	panel.SetSize(50, 50)
	panel.SetText(group_text)

	dock := windigo.NewSimpleDock(panel)
	dock.SetMargins(15)

	for _, label_text := range labels {
		label := windigo.NewRadioButton(panel)
		label.SetSize(150, 25)
		label.SetMargins(10)
		label.SetText(label_text)
		buttons = append(buttons, label)
		dock.Dock(label, windigo.Left)

	}

	if field_data.Valid {

		buttons[field_data.toIndex()].SetChecked(true)
	}

	// Get()
	get := func() ProductType {
		for i, button := range buttons {
			if button.Checked() {
				// one-based
				return ProductTypeFromIndex(i)
			}
		}
		return ProductTypeDefault()
	}

	// Error()
	err := func() {
		log.Println("NO")
		//TODO
	}

	return NewProductTypeView(panel, get, err)
}

type RangeView struct {
	*windigo.Panel
	Get func() Range
}

func NewRangeView(panel *windigo.Panel, get func() Range) RangeView {
	return RangeView{panel, get}
}

func BuildNewRangeView(parent windigo.Controller, field_text string, field_data Range, format func(float64) string) RangeView {
	panel := windigo.NewPanel(parent)
	panel.SetSize(50, 50)

	dock := windigo.NewSimpleDock(panel)
	dock.SetMargins(10)
	label := windigo.NewLabel(panel)
	label.SetText(field_text)
	min_field := BuildNewNullFloat64View(panel, field_data.Min, format)
	target_field := BuildNewNullFloat64View(panel, field_data.Target, format)

	max_field := BuildNewNullFloat64View(panel, field_data.Max, format)

	dock.Dock(label, windigo.Left)
	dock.Dock(min_field, windigo.Left)
	dock.Dock(target_field, windigo.Left)

	dock.Dock(max_field, windigo.Left)

	get := func() Range {
		return NewRange(
			min_field.Get(),
			target_field.Get(),
			max_field.Get(),
		)
	}

	return NewRangeView(panel, get)
}

type NullFloat64View struct {
	*windigo.Edit
	Get func() NullFloat64
}

func NewNullFloat64View(control *windigo.Edit, get func() NullFloat64) NullFloat64View {
	return NullFloat64View{control, get}
}

func BuildNewNullFloat64View(parent windigo.Controller, field_data NullFloat64, format func(float64) string) NullFloat64View {
	edit_field := windigo.NewEdit(parent)
	if field_data.Valid {
		edit_field.SetText(format(field_data.Float64))
	}
	edit_field.OnKillFocus().Bind(func(e *windigo.Event) {
		edit_field.SetText(strings.TrimSpace(edit_field.Text()))
	})

	get := func() NullFloat64 {
		field_text := edit_field.Text()
		var value float64
		valid := field_text != ""
		if valid {
			value, err = strconv.ParseFloat(field_text, 64)
			valid = err == nil
		}
		//only valid if we have some text which parses correctly
		return NewNullFloat64(value, valid)
	}

	return NewNullFloat64View(edit_field, get)
}

func (product *QCProduct) reset() {
	var empty_product QCProduct
	empty_product.Product = product.Product
	*product = empty_product
}

func (product *QCProduct) select_product_details() {
	var (
		product_name_customer *string

		product_name_customer_default string
	)

	product_name_customer_default = ""

	err := db_select_product_details.QueryRow(product.product_id).Scan(&product_name_customer, &product.product_type,
		&product.SG.Min, &product.SG.Max, &product.SG.Target,
		&product.PH.Min, &product.PH.Max, &product.PH.Target,
		&product.String_test.Min, &product.String_test.Max, &product.String_test.Target,
		&product.Viscosity.Min, &product.Viscosity.Max, &product.Viscosity.Target,
	)
	if err != nil {
		log.Printf("%q: %s\n", err, "select_product_details")

	}
	product.Product_name_customer = ValidOr(product_name_customer, product_name_customer_default)

}

func (product QCProduct) upsert() {

	_, err := db_upsert_product_details.Exec(product.product_id, product.product_type,
		product.SG.Min, product.SG.Max, product.SG.Target,
		product.PH.Min, product.PH.Max, product.PH.Target,
		product.String_test.Min, product.String_test.Max, product.String_test.Target,
		product.Viscosity.Min, product.Viscosity.Max, product.Viscosity.Target,
	)
	if err != nil {
		log.Printf("%q: %s\n", err, "upsert")
	}

	//TODO?
	// id, err := result.LastInsertId()
	// log.Println("upsert", id, err)

	// result.LastInsertId()
	// return product_type_id_default, product_name_customer_default

}

func (product *QCProduct) edit(
	product_type ProductType,
	SG,
	PH,
	// Density,
	String_test,
	Viscosity Range,
) {
	product.product_type = product_type
	product.SG = SG
	product.PH = PH
	// product.Density = Density
	product.String_test = String_test
	product.Viscosity = Viscosity
}

func (product *QCProduct) show_ranges_window() {

	rangeWindow := windigo.NewForm(nil)
	var WindowText string

	if product.Product_name_customer != "" {
		WindowText = fmt.Sprintf("%s (%s)", product.Product_type, product.Product_name_customer)
	} else {
		WindowText = product.Product_type
	}

	rangeWindow.SetSize(800, 600) // (width, height)
	rangeWindow.SetText(WindowText)

	dock := windigo.NewSimpleDock(rangeWindow)
	dock.SetMargins(5)

	dock.SetMarginTop(10)
	prod_label := windigo.NewLabel(rangeWindow)

	prod_label.SetText(WindowText)

	radio_dock := BuildNewProductTypeView(rangeWindow, "Type", product.product_type, []string{"Water Based", "Oil Based", "Friction Reducer"})

	labels := build_text_dock(rangeWindow, []string{"", "Min", "Target", "Max"})
	ph_dock := BuildNewRangeView(rangeWindow, "pH", product.PH, format_ranges_ph)
	sg_dock := BuildNewRangeView(rangeWindow, "Specific Gravity", product.SG, format_ranges_sg)
	string_dock := BuildNewRangeView(rangeWindow, "String Test \n\t at 0.5gpt", product.String_test, format_ranges_string_test)
	//TODO store string_amt "at 0.5gpt"
	visco_dock := BuildNewRangeView(rangeWindow, "Viscosity", product.Viscosity, format_ranges_viscosity)

	exit := func() {
		rangeWindow.Close()
		windigo.Exit()
	}
	save := func() {
		product.edit(
			radio_dock.Get(),
			sg_dock.Get(),
			ph_dock.Get(),
			string_dock.Get(),
			visco_dock.Get(),
		)
		product.upsert()
		show_status_bar("\t\tQC Data Updated")
		exit()
	}
	try_save := func() {
		if radio_dock.Get().Valid {
			save()
		} else {
			radio_dock.Error()
		}
	}
	button_dock := build_button_dock(rangeWindow, []string{"OK", "Cancel"}, []func(){try_save, exit})

	dock.Dock(prod_label, windigo.Top)
	dock.Dock(radio_dock, windigo.Top)
	dock.Dock(labels, windigo.Top)
	dock.Dock(ph_dock, windigo.Top)
	dock.Dock(sg_dock, windigo.Top)
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
