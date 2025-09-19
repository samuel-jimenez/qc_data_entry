package subpanels

import (
	"database/sql"
	"log"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/product"
	"github.com/samuel-jimenez/qc_data_entry/util/math"
	"github.com/samuel-jimenez/windigo"
)

type FrictionReducerPanelView struct {
	Update          func(qc_product *product.QCProduct)
	ChangeContainer func(qc_product *product.QCProduct)
	SetFont         func(font *windigo.Font)
	RefreshSize     func()
}

// create table product_line (product_id integer not null primary key, product_name text);
func Show_fr(parent *windigo.AutoPanel, qc_product *product.QCProduct, create_new_product_cb func() product.BaseProduct) *FrictionReducerPanelView {
	DELTA_DIFF_VISCO := 200.

	top_text := "Top"
	// bottom_text := "Bottom"
	bottom_text := "Btm"

	component_panel := views.NewQCBlendView(parent)

	panel := windigo.NewAutoPanel(parent)

	ranges_panel := BuildNewFrictionReducerProductRangesView(panel, qc_product)

	top_group := BuildNewFrictionReducerProductView(panel, top_text, ranges_panel)
	bottom_group := BuildNewFrictionReducerProductView(panel, bottom_text, ranges_panel)

	top_group.viscosity_field.Entangle(bottom_group.viscosity_field, ranges_panel.viscosity_field, DELTA_DIFF_VISCO)

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
		// TODO blend013 do only if base_product.Blend != nil?
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
		if qc_product.Container_type == product.CONTAINER_RAILCAR {
			bottom_group.Show()
			button_dock_cars.Show()
			button_dock_totes.Hide()
		} else {
			bottom_group.Hide()
			button_dock_cars.Hide()
			// // TODO
			// if qc_product.Container_type == product.CONTAINER_TOTE {
			button_dock_totes.Show()
			// } else { // product.CONTAINER_SAMPLE
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

		panel.SetSize(GUI.OFF_AXIS, GUI.GROUP_HEIGHT)
		panel.SetMargins(GUI.GROUP_MARGIN, GUI.GROUP_MARGIN, 0, 0)

		button_dock_totes.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)
		button_dock_cars.SetDockSize(GUI.BUTTON_WIDTH, GUI.BUTTON_HEIGHT)

		top_group.RefreshSize()
		bottom_group.RefreshSize()

		ranges_panel.RefreshSize()
		component_panel.RefreshSize()
	}

	return &FrictionReducerPanelView{update, changeContainer, setFont, refresh}

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

		// * Check storage
		bottom_product.CheckStorage()

	} else { // TODO show confirm box
		log.Println("ERROR: Viscosity", top_product.Lot_number, top_product.Product_name, top_product.Viscosity, bottom_product.Viscosity)

	}

}
