package qc

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/formats"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/util/math"
	"github.com/samuel-jimenez/windigo"
)

type FrictionReducerProduct struct {
	product.BaseProduct
	sg          float64
	string_test int64
	viscosity   int64
}

func (fr_product FrictionReducerProduct) toProduct() product.Product {
	return product.Product{
		BaseProduct: fr_product.Base(),
		PH:          nullable.NewNullFloat64(0, false),
		SG:          nullable.NewNullFloat64(fr_product.sg, true),
		Density:     nullable.NewNullFloat64(formats.Density_from_sg(fr_product.sg), true),
		String_test: nullable.NewNullInt64(fr_product.string_test),
		Viscosity:   nullable.NewNullInt64(fr_product.viscosity),
	}
}

func newFrictionReducerProduct(base_product product.BaseProduct, viscosity, mass, string_test float64) product.Product {

	sg := formats.SG_from_mass(mass)

	return FrictionReducerProduct{base_product, sg, int64(string_test), int64(viscosity)}.toProduct()

}

func (product FrictionReducerProduct) Check_data() bool {
	return true
}

type FrictionReducerProductView struct {
	*windigo.AutoPanel
	Get     func(base_product product.BaseProduct, replace_sample_point bool) product.Product
	Clear   func()
	SetFont func(font *windigo.Font)
	Refresh func()
}

func BuildNewFrictionReducerProductView(parent *windigo.AutoPanel, sample_point string, ranges_panel FrictionReducerProductRangesView) FrictionReducerProductView {

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	string_text := "String"
	mass_text := "Mass"

	group_panel := windigo.NewGroupAutoPanel(parent)
	group_panel.SetText(sample_point)

	visual_field := views.NewBoolCheckboxView(group_panel, visual_text)

	viscosity_field := views.NewNumberEditViewWithChange(group_panel, viscosity_text, ranges_panel.viscosity_field)

	mass_field := views.NewNumberEditView(group_panel, mass_text)

	string_field := views.NewNumberEditViewWithChange(group_panel, string_text, ranges_panel.string_field)

	density_field := views.NewMassDataView(group_panel, ranges_panel, mass_field)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(viscosity_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Top)
	group_panel.Dock(string_field, windigo.Top)

	get := func(base_product product.BaseProduct, replace_sample_point bool) product.Product {
		base_product.Visual = visual_field.Checked()
		if replace_sample_point {
			base_product.Sample_point = sample_point
		}
		return newFrictionReducerProduct(base_product, viscosity_field.Get(), density_field.Get(), string_field.Get())

	}
	clear := func() {
		visual_field.SetChecked(false)
		viscosity_field.Clear()

		density_field.Clear()

		string_field.Clear()

		ranges_panel.Clear()

	}

	setFont := func(font *windigo.Font) {
		group_panel.SetFont(font)
		visual_field.SetFont(font)
		viscosity_field.SetFont(font)
		density_field.SetFont(font) //?TODO
		string_field.SetFont(font)
	}
	refresh := func() {
		group_panel.SetSize(GROUP_WIDTH, GROUP_HEIGHT)
		group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, BTM_SPACER_WIDTH, BTM_SPACER_HEIGHT)

		visual_field.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
		viscosity_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)
		density_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.DATA_SUBFIELD_WIDTH, DATA_UNIT_WIDTH, GUI.EDIT_FIELD_HEIGHT)
		string_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.DATA_FIELD_WIDTH, GUI.EDIT_FIELD_HEIGHT)

	}

	return FrictionReducerProductView{group_panel, get, clear, setFont, refresh}

}

type FrictionReducerProductRangesView struct {
	*windigo.AutoPanel
	*views.MassRangesView

	viscosity_field,
	// mass_field,
	// sg_field,
	// density_field,
	string_field *views.RangeROView

	Update  func(qc_product *product.QCProduct)
	SetFont func(font *windigo.Font)
	Refresh func()
}

func (data_view FrictionReducerProductRangesView) Clear() {
	data_view.MassRangesView.Clear()
	data_view.viscosity_field.Clear()
	data_view.string_field.Clear()
}

func BuildNewFrictionReducerProductRangesView(parent *windigo.AutoPanel, qc_product *product.QCProduct) FrictionReducerProductRangesView {

	visual_text := "Visual Inspection"
	viscosity_text := "Viscosity"
	mass_text := "Mass"
	string_text := "String"
	sg_text := "Specific Gravity"
	density_text := "Density"

	group_panel := windigo.NewAutoPanel(parent)

	visual_field := product.BuildNewProductAppearanceROView(group_panel, visual_text, qc_product.Appearance)

	viscosity_field := views.BuildNewRangeROView(group_panel, viscosity_text, qc_product.Viscosity, formats.Format_ranges_viscosity)

	string_field := views.BuildNewRangeROView(group_panel, string_text, qc_product.SG, formats.Format_ranges_string_test)

	mass_field := views.BuildNewRangeROViewMap(group_panel, mass_text, qc_product.SG, formats.Format_mass, formats.Mass_from_sg)

	sg_field := views.BuildNewRangeROView(group_panel, sg_text, qc_product.SG, formats.Format_ranges_sg)
	density_field := views.BuildNewRangeROView(group_panel, density_text, qc_product.Density, formats.Format_ranges_density)

	group_panel.Dock(visual_field, windigo.Top)
	group_panel.Dock(viscosity_field, windigo.Top)
	group_panel.Dock(mass_field, windigo.Top)
	group_panel.Dock(string_field, windigo.Top)
	group_panel.Dock(density_field, windigo.Bottom)
	group_panel.Dock(sg_field, windigo.Bottom)

	update := func(qc_product *product.QCProduct) {
		visual_field.Update(qc_product.Appearance)
		viscosity_field.Update(qc_product.Viscosity)
		string_field.Update(qc_product.String_test)

		mass_field.Update(qc_product.SG)
		sg_field.Update(qc_product.SG)
		density_field.Update(qc_product.Density)
	}

	setFont := func(font *windigo.Font) {
		visual_field.SetFont(font)
		viscosity_field.SetFont(font)
		mass_field.SetFont(font)
		string_field.SetFont(font)
		sg_field.SetFont(font)
		density_field.SetFont(font)
	}
	refresh := func() {
		group_panel.SetSize(GUI.DATA_FIELD_WIDTH, GROUP_HEIGHT)
		group_panel.SetPaddings(TOP_SPACER_WIDTH, TOP_SPACER_HEIGHT, GUI.RANGES_RO_PADDING, BTM_SPACER_HEIGHT)
		visual_field.Refresh()
		viscosity_field.Refresh()
		mass_field.Refresh()
		string_field.Refresh()
		sg_field.Refresh()
		density_field.Refresh()

	}

	return FrictionReducerProductRangesView{group_panel,
		&views.MassRangesView{Mass_field: &mass_field,
			SG_field:      &sg_field,
			Density_field: &density_field},
		&viscosity_field,
		&string_field,
		update, setFont, refresh}

}

// TODO
func check_dual_data(top_product, bottom_product product.Product) {
	// DELTA_DIFF_VISCO := 200
	var DELTA_DIFF_VISCO int64 = 200 //go sucks

	if math.Abs(top_product.Viscosity.Diff(bottom_product.Viscosity)) <= DELTA_DIFF_VISCO &&
		top_product.Check_data() && bottom_product.Check_data() {

		log.Println("debug: Check_data", top_product)
		// top_product.Save()
		// bottom_product.Save()

		// TODO blend013 ensurethis works with testing blends
		product.Store(top_product, bottom_product)

		// qc_product.Store(9)

		err := top_product.Printout()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "top_product.Printout", err)
		}

		log.Println("debug: Check_data", bottom_product)
		err = bottom_product.Printout()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "bottom_product.Printout", err)
		}
		//TODO find closest: RMS?
		err = bottom_product.Output_sample()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "bottom_product.Output_sample", err)
		}
		err = bottom_product.Save_xl()
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "bottom_product.Save_xl", err)
		}

	} else { // TODO show confirm box
		log.Println("ERROR: Viscosity", top_product.Lot_number, top_product.Product_name, top_product.Viscosity, bottom_product.Viscosity)

	}

}

type FrictionReducerPanelView struct {
	Update          func(qc_product *product.QCProduct)
	ChangeContainer func(qc_product *product.QCProduct)
	SetFont         func(font *windigo.Font)
	Refresh         func()
}

// create table product_line (product_id integer not null primary key, product_name text);
func show_fr(parent *windigo.AutoPanel, qc_product *product.QCProduct, create_new_product_cb func() product.BaseProduct) *FrictionReducerPanelView {

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	component_panel := views.NewQCBlendView(parent)

	panel := windigo.NewAutoPanel(parent)

	ranges_panel := BuildNewFrictionReducerProductRangesView(panel, qc_product)

	top_group := BuildNewFrictionReducerProductView(panel, top_text, ranges_panel)
	bottom_group := BuildNewFrictionReducerProductView(panel, bottom_text, ranges_panel)

	submit_cb := func() {
		base_product := create_new_product_cb()

		// TODO blend012 ensurethis works with testing blends
		//component_panel.saVE
		log.Println("DEBUG: FrictionReducerPanelView.submit_cb base_product", base_product)
		base_product.SetBlend(component_panel.Get())
		log.Println("DEBUG: FrictionReducerPanelView.submit_cb base_product", base_product)
		//TODO make sure this is the only time it is saved
		base_product.SaveBlend()

		top_product := top_group.Get(base_product, true)
		bottom_product := bottom_group.Get(base_product, true)
		log.Println("debug: FrictionReducerPanelView.submit_cb.top", top_product)
		log.Println("debug: FrictionReducerPanelView.submit_cb.btm", bottom_product)
		check_dual_data(top_product, bottom_product)

	}

	clear_cb := func() {
		top_group.Clear()
		bottom_group.Clear()
	}

	tote_cb := func() {
		base_product := create_new_product_cb()
		base_product.SetBlend(component_panel.Get())
		log.Println("DEBUG: FrictionReducerPanelView.submit_cb.tote base_product", base_product)

		top_product := top_group.Get(base_product, false)
		if top_product.Check_data() {
			log.Println("debug: FrictionReducerPanelView.submit_cb.tote", top_product)
			top_product.Save()
			err := top_product.Output()
			if err != nil {
				log.Printf("Error: [%s]: %q\n", "top_product.Output", err)
			}
			//TODO component_panel.saVE

		}
	}

	button_dock_totes := GUI.NewMarginalButtonDock(parent, []string{"Submit", "Clear"}, []int{40, 0}, []func(){tote_cb, clear_cb})
	button_dock_cars := GUI.NewMarginalButtonDock(parent, []string{"Submit", "Clear"}, []int{40, 0}, []func(){submit_cb, clear_cb})
	button_dock_totes.Hide()

	panel.Dock(top_group, windigo.Left)
	panel.Dock(bottom_group, windigo.Left)
	// TODO
	panel.Dock(ranges_panel, windigo.Right)
	// panel.Dock(ranges_panel, windigo.Left)
	//	parent.Dock(component_panel, windigo.Top)
	// parent.Dock(component_panel, windigo.Top)

	parent.Dock(panel, windigo.Top)
	parent.Dock(button_dock_totes, windigo.Top)
	parent.Dock(button_dock_cars, windigo.Top)

	parent.Dock(component_panel, windigo.Top)

	update := func(qc_product *product.QCProduct) {
		ranges_panel.Update(qc_product)

		log.Println("TRACE: show_fr.update Blend", qc_product.Blend)
		// TODO recip00
		// extract to fn, move componenet panel?

		if qc_product.Blend != nil {
			component_panel.UpdateBlend(qc_product.Blend)
			return
		}

		var (
			recipe_data blender.ProductRecipe
		)
		// proc_name := "RecipeProduct.GetRecipes"
		proc_name := "FrictionReducerPanelView.GetRecipes"

		DB.Forall(proc_name,
			func() {},
			func(row *sql.Rows) {

				if err := row.Scan(
					&recipe_data.Recipe_id,
				); err != nil {
					log.Fatal("Crit: [RecipeProduct GetRecipes]: ", proc_name, err)
				}
				log.Println("DEBUG: GetRecipes qc_data", proc_name, recipe_data)

			},
			DB.DB_Select_product_recipe, qc_product.Product_id)

		recipe_data.GetComponents()
		component_panel.UpdateRecipe(&recipe_data)

	}

	changeContainer := func(qc_product *product.QCProduct) {
		if int(qc_product.Container_type.Int32) == CONTAINER_RAILCAR {
			bottom_group.Show()
			button_dock_cars.Show()
			button_dock_totes.Hide()
		} else {
			bottom_group.Hide()
			button_dock_cars.Hide()
			// // TODO
			// if int(qc_product.Container_type.Int32) == CONTAINER_TOTE {
			button_dock_totes.Show()
			// } else { // CONTAINER_SAMPLE
			// 					NO COA
			// no storage
			// TODO blend013 ensurethis works with testing blends

			// }
		}

	}

	setFont := func(font *windigo.Font) {
		ranges_panel.SetFont(font)      //?TODO
		top_group.SetFont(font)         //?TODO
		bottom_group.SetFont(font)      //?TODO
		button_dock_totes.SetFont(font) //?TODO
		button_dock_cars.SetFont(font)  //?TODO
		component_panel.SetFont(font)
	}
	refresh := func() {

		panel.SetSize(GUI.OFF_AXIS, GROUP_HEIGHT)
		panel.SetMargins(GROUP_MARGIN, GROUP_MARGIN, 0, 0)

		button_dock_totes.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)
		button_dock_cars.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)

		top_group.Refresh()
		bottom_group.Refresh()
		ranges_panel.Refresh()
		component_panel.RefreshSize()
	}

	return &FrictionReducerPanelView{update, changeContainer, setFont, refresh}

}

// c.f. blender.BlendComponentView
