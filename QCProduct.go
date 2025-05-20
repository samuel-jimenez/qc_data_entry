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
	Update       func()
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

// func (field_data Range) Check(data NullFloat64) bool {
func (field_data Range) Check(data float64) bool {
	return (!field_data.Min.Valid ||
		field_data.Min.Float64 <= data) && (!field_data.Max.Valid ||
		data <= field_data.Max.Float64)
}

func (field_data Range) Map(data_map func(float64) float64) Range {
	return Range{field_data.Min.Map(data_map),
		field_data.Target.Map(data_map),
		field_data.Max.Map(data_map)}
}

type ProductTypeView struct {
	windigo.AutoPanel
	Get   func() ProductType
	Error func()
}

func BuildNewProductTypeView(parent windigo.Controller, group_text string, field_data ProductType, labels []string) ProductTypeView {

	var buttons []*windigo.RadioButton
	panel := windigo.NewGroupAutoPanel(parent)
	panel.SetSize(50, 50)
	panel.SetText(group_text)
	panel.SetPaddingsAll(15)

	for _, label_text := range labels {
		label := windigo.NewRadioButton(panel)
		label.SetSize(150, 25)
		label.SetMarginLeft(10)
		label.SetText(label_text)
		buttons = append(buttons, label)
		panel.Dock(label, windigo.Left)

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

	return ProductTypeView{panel, get, err}
}

type RangeView struct {
	windigo.AutoPanel
	Get func() Range
}

func BuildNewRangeView(parent windigo.Controller, field_text string, field_data Range, format func(float64) string) RangeView {

	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(50, 50)
	// panel.SetMarginsAll(10)
	panel.SetPaddingsAll(10)
	label := windigo.NewLabel(panel)
	label.SetText(field_text)

	min_field := BuildNewNullFloat64View(panel, field_data.Min, format)
	target_field := BuildNewNullFloat64View(panel, field_data.Target, format)
	max_field := BuildNewNullFloat64View(panel, field_data.Max, format)

	panel.Dock(label, windigo.Left)
	panel.Dock(min_field, windigo.Left)
	panel.Dock(target_field, windigo.Left)

	panel.Dock(max_field, windigo.Left)

	get := func() Range {
		return NewRange(
			min_field.Get(),
			target_field.Get(),
			max_field.Get(),
		)
	}

	return RangeView{panel, get}
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

type RangeROView struct {
	windigo.AutoPanel
	field_data Range
	min_field,
	min_field_spacer,
	target_field,
	max_field_spacer,
	max_field NullFloat64ROView
	data_map func(float64) float64
}

func (data_view *RangeROView) Update(update_data Range) {
	if data_view.data_map != nil {
		update_data = update_data.Map(data_view.data_map)
	}

	data_view.field_data = update_data
	data_view.min_field.Update(update_data.Min)
	data_view.min_field_spacer.Update(update_data.Min)
	data_view.target_field.Update(update_data.Target)
	data_view.max_field.Update(update_data.Max)
	data_view.max_field_spacer.Update(update_data.Max)
}

func (data_view RangeROView) Check(data float64) bool {
	return data_view.field_data.Check(data)
}

func BuildNewRangeROViewMap(parent windigo.Controller, field_text string, field_data Range, format func(float64) string, data_map func(float64) float64) RangeROView {

	panel := windigo.NewAutoPanel(parent)
	panel.SetSize(22, 22)
	panel.SetMarginTop(5)
	//TODO toolti[p]
	// label := windigo.NewLabel(panel)
	// label.SetText(field_text)

	spacer_format := "<"
	min_field := BuildNewNullFloat64ROView(panel, field_data.Min, format)
	min_field_spacer := BuildNullFloat64SpacerView(panel, field_data.Min, spacer_format)
	target_field := BuildNewNullFloat64ROView(panel, field_data.Target, format)
	max_field := BuildNewNullFloat64ROView(panel, field_data.Max, format)
	max_field_spacer := BuildNullFloat64SpacerView(panel, field_data.Max, spacer_format)

	// panel.Dock(label, windigo.Left)
	panel.Dock(min_field, windigo.Left)
	panel.Dock(min_field_spacer, windigo.Left)
	panel.Dock(target_field, windigo.Left)
	panel.Dock(max_field_spacer, windigo.Left)
	panel.Dock(max_field, windigo.Left)

	return RangeROView{panel, field_data, min_field, min_field_spacer,
		target_field,
		max_field,
		max_field_spacer, data_map}
}

func BuildNewRangeROView(parent windigo.Controller, field_text string, field_data Range, format func(float64) string) RangeROView {
	return BuildNewRangeROViewMap(parent, field_text, field_data, format, nil)
}

type NullFloat64ROView struct {
	*windigo.Label
	Update func(field_data NullFloat64)
}

func BuildNewNullFloat64ROView(parent windigo.Controller, field_data NullFloat64, format func(float64) string) NullFloat64ROView {
	data_field := windigo.NewLabel(parent)
	data_field.SetSize(40, 22)

	update := func(field_data NullFloat64) {
		if field_data.Valid {
			data_field.SetText(format(field_data.Float64))
		} else {
			data_field.SetText("")
		}
	}
	update(field_data)

	return NullFloat64ROView{data_field, update}
}

func BuildNullFloat64SpacerView(parent windigo.Controller, field_data NullFloat64, format string) NullFloat64ROView {
	data_field := windigo.NewLabel(parent)
	data_field.SetSize(20, 22)

	update := func(field_data NullFloat64) {
		if field_data.Valid {
			data_field.SetText(format)
		} else {
			data_field.SetText("")
		}
	}
	update(field_data)

	return NullFloat64ROView{data_field, update}
}

func (product *QCProduct) reset() {
	var empty_product QCProduct
	empty_product.Product = product.Product
	empty_product.Update = product.Update
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
	dock.SetPaddingsAll(5)

	dock.SetPaddingTop(10)
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
		product.Update()
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

func (product QCProduct) Check(data Product) bool {
	return product.PH.Check(data.PH.Float64) && product.SG.Check(data.SG.Float64)
}
