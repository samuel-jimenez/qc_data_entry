package blender_ui

import (
	"database/sql"
	"log"
	"os/exec"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
	"github.com/samuel-jimenez/qc_data_entry/blender/blendsheet"
	"github.com/samuel-jimenez/windigo"
	"github.com/samuel-jimenez/windigo/w32"
	// "github.com/samuel-jimenez/windigo/w32"
)

/*
 * BlendStrappingProductViewer
 *
 */
type BlendStrappingProductViewer interface {
	windigo.Pane

	Operations_group_field_OnSelectedChange(*windigo.Event)
	Product_Field_OnSelectedChange(*windigo.Event)
	recipe_accept_button_OnClick(*windigo.Event)
	Blend_sel_field_OnSelectedChange(*windigo.Event)

	Get() *blender.BlendProduct
	Update_products(product_data map[string]int64)
	SetFont(font *windigo.Font)
	RefreshSize()
	SetHeelVolume(float64)
}

/*
 * BlendStrappingProductView
 *
 */
type BlendStrappingProductView struct {
	*windigo.AutoPanel

	Blend_Vessel *BlendVessel

	RecipeProduct *blender.RecipeProduct
	// BlendProduct                   *blender.BlendProduct
	Blend         *BlendView
	Product_data  map[string]int64
	Product_Field *GUI.SearchBox
	Operations_group_field,
	Customer_field,
	Blend_sel_field *GUI.ComboBox

	Tag_field, Seal_field *GUI.NumbestEditView
	Operators_field       *windigo.LabeledEdit

	panels   []*windigo.AutoPanel
	controls []windigo.Controller

	// Product_name string
	// Product_id     int64
}

func BlendStrappingProductView_from_new(parent windigo.Controller) *BlendStrappingProductView {

	ops_group_text := "Ops Group"
	product_text := "Product"
	// recipe_text := ""
	recipe_text := "Recipe"
	customer_text := "Customer Name"

	Tag_text := "Tag"

	Seal_text := "Seal"

	Operators_text := "Operators"

	view := new(BlendStrappingProductView)

	view.RecipeProduct = blender.RecipeProduct_from_new()

	view.Product_data = make(map[string]int64)

	view.AutoPanel = windigo.NewAutoPanel(parent)

	ops_group_panel := windigo.NewAutoPanel(view)
	product_panel := windigo.NewAutoPanel(view)
	tag_seal_panel := windigo.NewAutoPanel(view)
	ops_panel := windigo.NewAutoPanel(view)
	recipe_panel := windigo.NewAutoPanel(view)
	view.Blend_Vessel = BlendVessel_from_new(view)

	view.panels = append(view.panels, ops_group_panel)
	view.panels = append(view.panels, product_panel)
	view.panels = append(view.panels, ops_panel)
	view.panels = append(view.panels, tag_seal_panel)
	view.panels = append(view.panels, recipe_panel)

	view.Operations_group_field = GUI.List_ComboBox_from_new(ops_group_panel, ops_group_text)

	view.Product_Field = GUI.NewLabeledListSearchBox(product_panel, product_text)

	view.Customer_field = GUI.ComboBox_from_new(product_panel, customer_text)
	view.Blend_sel_field = GUI.List_ComboBox_from_new(product_panel, recipe_text)

	view.Tag_field = GUI.NumbestEditView_from_new(tag_seal_panel, Tag_text)
	view.Seal_field = GUI.NumbestEditView_from_new(tag_seal_panel, Seal_text)
	view.Operators_field = windigo.NewLabeledEdit(ops_panel, Operators_text)

	recipe_accept_button := windigo.NewPushButton(recipe_panel)
	recipe_accept_button.SetText("OK")
	// recipe_accept_button.SetSize(GUI.ACCEPT_BUTTON_WIDTH, GUI.OFF_AXIS)
	recipe_accept_button.SetSize(200, 200)

	view.Blend = BlendView_from_new(view)

	view.controls = append(view.controls, recipe_accept_button)

	// Dock

	view.Dock(ops_group_panel, windigo.Top)
	view.Dock(product_panel, windigo.Top)
	view.Dock(tag_seal_panel, windigo.Top)
	view.Dock(ops_panel, windigo.Top)

	view.Dock(view.Blend_Vessel, windigo.Top)
	view.Dock(view.Blend, windigo.Top)
	view.Dock(recipe_panel, windigo.Top)

	ops_group_panel.Dock(view.Operations_group_field, windigo.Left)

	product_panel.Dock(view.Product_Field, windigo.Left)
	product_panel.Dock(view.Customer_field, windigo.Left)
	product_panel.Dock(view.Blend_sel_field, windigo.Left)
	tag_seal_panel.Dock(view.Tag_field, windigo.Left)
	tag_seal_panel.Dock(view.Seal_field, windigo.Left)
	ops_panel.Dock(view.Operators_field, windigo.Left)
	recipe_panel.Dock(recipe_accept_button, windigo.Left)

	// combobox
	GUI.Fill_combobox_from_query_rows(view.Product_Field, func(row *sql.Rows) error {
		var (
			id                   int64
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		view.Product_data[name] = id
		return nil
	},
		DB.DB_Select_product_info_all)

	for _, name := range []string{"BSFR", "BSWB", "BSWH"} {
		view.Operations_group_field.AddItem(name)
	}
	view.Operations_group_field.SetSelectedItem(0)
	view.Operations_group_field_OnSelectedChange(nil)
	//event handling

	view.Operations_group_field.OnSelectedChange().Bind(view.Operations_group_field_OnSelectedChange)

	view.Product_Field.OnSelectedChange().Bind(view.Product_Field_OnSelectedChange)
	recipe_accept_button.OnClick().Bind(view.recipe_accept_button_OnClick)
	view.Blend_sel_field.OnSelectedChange().Bind(view.Blend_sel_field_OnSelectedChange)

	return view
}

// functionality

func (view *BlendStrappingProductView) SetFont(font *windigo.Font) {
	view.Blend_Vessel.SetFont(windigo.DefaultFont)
	view.Blend.SetFont(font)

	view.Operations_group_field.SetFont(font)
	view.Product_Field.SetFont(font)
	view.Customer_field.SetFont(font)
	view.Blend_sel_field.SetFont(font)
	view.Tag_field.SetFont(font)
	view.Seal_field.SetFont(font)
	view.Operators_field.SetFont(font)

	for _, control := range view.controls {
		control.SetFont(font)
	}

}

func (view *BlendStrappingProductView) RefreshSize() {
	view.Blend_Vessel.RefreshSize()

	for _, panel := range view.panels {
		panel.SetSize(GUI.TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	}

	view.Operations_group_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Product_Field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Customer_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Customer_field.SetPaddingLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.Blend_sel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Blend_sel_field.SetPaddingLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.Tag_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Seal_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Seal_field.SetPaddingLeft(GUI.TOP_PANEL_INTER_SPACER_WIDTH)
	view.Operators_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH+GUI.LABEL_WIDTH+GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
	view.Blend.RefreshSize()

}

func (view *BlendStrappingProductView) Operations_group_field_OnSelectedChange(e *windigo.Event) {
	var (
		select_statement *sql.Stmt
		args             []any
	)
	switch view.Operations_group_field.GetSelectedItem() {

	case "BSFR":
		select_statement = DB.DB_Select_product_info_moniker

		args = []any{
			"PETROFLO"}

	// case "BSWB":
	// case "BSWH":
	default:

		select_statement = DB.DB_Select_product_info_all

	}
	GUI.Fill_combobox_from_query_rows(view.Product_Field, func(row *sql.Rows) error {
		var (
			id                   int64
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		view.Product_Field.AddItem(name)
		return nil
	},
		select_statement,
		args...)
}

func (view *BlendStrappingProductView) Product_Field_OnSelectedChange(e *windigo.Event) {
	// view.BlendProduct = NewBlendProduct()
	name := view.Product_Field.GetSelectedItem()
	Product_id := view.Product_data[name]
	view.RecipeProduct = blender.RecipeProduct_from_Product_id(Product_id)
	view.RecipeProduct.GetRecipes(view.Blend_sel_field)
	view.RecipeProduct.GetComponents()
	view.Blend_sel_field.OnSelectedChange().Fire(nil)
	GUI.Fill_combobox_from_query(view.Customer_field, DB.DB_Select_product_customer_info, Product_id)
}

func (view *BlendStrappingProductView) recipe_accept_button_OnClick(e *windigo.Event) {
	sheet := view.Get()
	// log.Println(sheet.Export_pdf())

	pdf_path, _ := sheet.Export_pdf()

	err := w32.ShellExecute(0, "", pdf_path, "", "", 0)
	if err != nil {
		log.Println("ERROROROR:", err)
		return
	}

	//
	// 		app := "./chrome"
	app := "C:/Program Files/Google/Chrome/Application/chrome.exe"
	// app := "C:/bin/chrome-win/chrome.exe"
	// app := "firefox"

	cmd := exec.Command(app, pdf_path)
	// // cmd := exec.Command(pdf_path)
	err = cmd.Start()
	if err != nil {
		log.Println("ERROROROR:", err)
		return
	}
	cmd.Wait()
}

func (view *BlendStrappingProductView) Blend_sel_field_OnSelectedChange(e *windigo.Event) {
	view.Blend.UpdateRecipe(view.RecipeProduct.Recipes[view.Blend_sel_field.GetSelectedItem()])
}

// return BlendSheet
// print blend sheet
// TODO
// operations_group
func (view *BlendStrappingProductView) Get() *blendsheet.BlendSheet {

	//TODO
	ProductBlend := view.Blend.Get()
	if ProductBlend == nil {
		return nil
	}

	// save defaults
	DB.Update("BlendStrappingProductView-Get",
		DB.DB_Update_product_defaults,
		view.RecipeProduct.Product_id,
		ProductBlend.Total)

	Blend_Sheet := blendsheet.BlendSheet_from_new(ProductBlend,
		view.Product_Field.GetSelectedItem(),
		view.Operations_group_field.GetSelectedItem(),
		view.Customer_field.GetSelectedItem(),
		view.RecipeProduct.Product_id,
		view.Blend_Vessel.Get(),
		view.Blend.Total_Component.Get(),
		view.Tag_field.Text(),
		view.Seal_field.Text(),
		view.Operators_field.Text(),
	)

	// BlendProduct := blender.BlendProduct_from_Recipe(view.RecipeProduct)
	//
	// if BlendProduct.Product_id == DB.INVALID_ID {
	// 	//TODO make error
	// 	return nil
	// }
	//
	// BlendProduct.Recipe = view.Blend.Recipe
	// BlendProduct.Blend = ProductBlend
	// BlendProduct.Save()	// TODO NOT HERE
	return Blend_Sheet

}

func (view *BlendStrappingProductView) Update_products(product_data map[string]int64) {
	view.Product_data = product_data
}

func (view *BlendStrappingProductView) SetHeelVolume(heel float64) {
	view.Blend.SetHeelVolume(heel)
}

func (view *BlendStrappingProductView) SetHeel(heel float64) {
	view.Blend.SetHeel(heel)
}

func (view *BlendStrappingProductView) SetStrap(volume float64) {
	view.Blend_Vessel.SetStrap(volume)
}

func (view *BlendStrappingProductView) GetStrap(volume float64) float64 {
	return view.Blend_Vessel.GetStrap(volume)
}
