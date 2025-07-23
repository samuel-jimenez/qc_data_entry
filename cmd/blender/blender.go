package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/windigo"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	WINDOW_WIDTH,
	WINDOW_HEIGHT,
	WINDOW_FUDGE_MARGIN,

	TOP_SPACER_WIDTH,
	TOP_SPACER_HEIGHT,
	INTER_SPACER_HEIGHT,
	BTM_SPACER_WIDTH,
	BTM_SPACER_HEIGHT,

	TOP_PANEL_WIDTH,
	REPRINT_BUTTON_WIDTH,
	TOP_PANEL_INTER_SPACER_WIDTH,

	RANGE_WIDTH,

	PRODUCT_FIELD_WIDTH,
	PRODUCT_FIELD_HEIGHT,
	DATA_FIELD_WIDTH,
	DATA_SUBFIELD_WIDTH,
	DATA_UNIT_WIDTH,
	DATA_MARGIN,
	FIELD_HEIGHT,
	NUM_FIELDS,

	BUTTON_WIDTH,
	BUTTON_HEIGHT,
	BUTTON_MARGIN,

	GROUP_WIDTH,
	GROUP_HEIGHT,
	GROUP_MARGIN int
)

func main() {
	//load config
	config.Main_config = config.Load_config_viewer("qc_data_blender")

	// log to file
	log_file, err := os.OpenFile(config.LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.Println("Info: Logging Started")

	//open_db
	// qc_db, err := sql.Open("sqlite3", DB_FILE)
	qc_db, err = sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer qc_db.Close()

	dbinit(qc_db)
	/*
		//setup print goroutine
		threads.PRINT_QUEUE = make(chan string, 4)
		defer close(threads.PRINT_QUEUE)
		go threads.Do_print_queue(threads.PRINT_QUEUE)*/

	//setup status_bar goroutine
	threads.STATUS_QUEUE = make(chan string, 16)
	defer close(threads.STATUS_QUEUE)
	go threads.Do_status_queue(threads.STATUS_QUEUE)

	//show main window
	show_window()

}

var (
	qc_db *sql.DB
	DB_Select_product_recipe, DB_Insert_product_recipe,
	DB_Select_recipe_components, DB_Insert_recipe_component,
	DB_Select_name_component_types, DB_Select_all_component_types, DB_Insert_component_types,
	// DB_Select_component_type_product,
	DB_Insert_internal_product_component_type,
	DB_Insert_inbound_product_component_type *sql.Stmt
)

func dbinit(db *sql.DB) {

	DB.Check_db(db)
	DB.DBinit(db)

	DB_Select_product_recipe = DB.PrepareOrElse(db, `
	select recipe_list_id
		from bs.recipe_list
		where product_id = ?
	`)

	DB_Insert_product_recipe = DB.PrepareOrElse(db, `
	insert into bs.recipe_list
		(product_id)
		values (?)
	returning recipe_list_id
	`)

	DB_Select_recipe_components = DB.PrepareOrElse(db, `
	select component_type_id, component_type_name, component_type_amount, component_add_order
		from bs.recipe_components
		join bs.component_types
		using (component_type_id)
		where recipe_list_id = ?
	`)

	DB_Insert_recipe_component = DB.PrepareOrElse(db, `
	insert into bs.recipe_components
		(recipe_list_id,component_type_id,component_type_amount,component_add_order)
		values (?,?,?,?)
	returning recipe_components_id
	`)

	DB_Select_name_component_types = DB.PrepareOrElse(db, `
	select component_type_id
		from bs.component_types
		where component_type_name = ?
	`)

	DB_Select_all_component_types = DB.PrepareOrElse(db, `
	select component_type_id, component_type_name
		from bs.component_types
		order by component_type_name
	`)

	DB_Insert_component_types = DB.PrepareOrElse(db, `
	insert into bs.component_types
		(component_type_name)
		values (?)
	returning component_type_id
	`)

	/*




				DB_Select_component_type_product = DB.PrepareOrElse(db, `
		select component_type_product_id, component_type_product_name, component_type_product_amount, component_add_order
			from bs.component_type_product
			where component_type_id = ?
		`)


	*/

	DB_Insert_internal_product_component_type = DB.PrepareOrElse(db, `
	insert into bs.component_type_product_internal
		(component_type_id,product_id)
		values (?,?)
	returning component_type_product_internal_id
	`)

	DB_Insert_inbound_product_component_type = DB.PrepareOrElse(db, `
	insert into bs.component_type_product_inbound
		(component_type_id,inbound_product_id)
		values (?,?)
	returning component_type_product_inbound_id
	`)

}

func refresh_globals(font_size int) {

	GUI.GROUPBOX_CUSHION = font_size * 3 / 2
	TOP_SPACER_WIDTH = 7
	TOP_SPACER_HEIGHT = GUI.GROUPBOX_CUSHION + 2
	INTER_SPACER_HEIGHT = 2
	BTM_SPACER_WIDTH = 2
	BTM_SPACER_HEIGHT = 2
	TOP_PANEL_INTER_SPACER_WIDTH = 30

	GUI.LABEL_WIDTH = 10 * font_size
	PRODUCT_FIELD_WIDTH = 15 * font_size
	PRODUCT_FIELD_HEIGHT = font_size*16/10 + 8
	FIELD_HEIGHT = 24

}

type RecipeProduct struct {
	Product_name string `json:"product_name"`
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Product_id int64
	// Lot_id     int64
	Recipes []*ProductRecipe
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

func (object *RecipeProduct) Set(name string, i int64) {
	object.Product_name = name
	object.Product_id = i
}

// func (r *RecipeProduct) GetRecipes() (rows *sql.Rows,error){
// 	return DB_Select_product_recipe.Query(r.Product_id)
// 	// r.Recipes =
//
// 	// rows, err :=
// 	//TODO
// }

func (object *RecipeProduct) GetRecipes() {
	DB.Forall("GetRecipes",
		func() { object.Recipes = nil },
		func(rows *sql.Rows) {
			var (
				recipe_data ProductRecipe
			)

			if err := rows.Scan(&recipe_data.Recipe_id); err != nil {
				log.Fatal(err)
			}
			recipe_data.Product_id = object.Product_id
			log.Println("DEBUG: GetRecipes qc_data", recipe_data)
			object.Recipes = append(object.Recipes, &recipe_data)

		},
		DB_Select_product_recipe, object.Product_id)
}

func (object *RecipeProduct) LoadRecipeCombo(combo_field *GUI.ComboBox) {
	object.Recipes = nil
	rows, err := DB_Select_product_recipe.Query(object.Product_id)
	i := 0
	GUI.Fill_combobox_from_query_rows(combo_field, rows, err, func(rows *sql.Rows) {

		var (
			recipe_data ProductRecipe
		)

		if err := rows.Scan(&recipe_data.Recipe_id); err != nil {
			log.Fatal(err)
		}
		recipe_data.Product_id = object.Product_id
		log.Println("DEBUG: GetRecipes qc_data", recipe_data)
		object.Recipes = append(object.Recipes, &recipe_data)
		// combo_field.AddItem(strconv.FormatInt(i, 10))
		combo_field.AddItem(strconv.Itoa(i))
		i++
	})
	if i != 0 {
		combo_field.SetSelectedItem(0)
	}
}

func (object *RecipeProduct) NewRecipe() *ProductRecipe {
	proc_name := "RecipeProduct.NewRecipe"
	var (
		recipe_data *ProductRecipe
	)
	result, err := DB_Insert_product_recipe.Exec(object.Product_id)
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return recipe_data
	}
	insert_id, err := result.LastInsertId()
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return recipe_data
	}
	recipe_data = new(ProductRecipe)

	recipe_data.Recipe_id = insert_id
	recipe_data.Product_id = object.Product_id
	object.Recipes = append(object.Recipes, recipe_data)
	return recipe_data
}

func (object *RecipeProduct) GetComponents() {
	for _, recipe_data := range object.Recipes {
		recipe_data.GetComponents()
	}
}

/*
func (object *RecipeProduct) AddRecipe() {
}*/
// func _select_samples(rows *sql.Rows, err error, fn string) []viewer.QCData {

//TODO forall(DB_Select_all_component_types,fn)

// 	DB_Select_all_component_types
// 	_list []string
// 	func forall(rows *sql.Rows, err error, fn string) []viewer.QCData {
// 			object.Recipes = nil
// rows, err := DB_Select_product_recipe.Query(object.Product_id)
// fn := "GetRecipes"
// if err != nil {
// 	log.Printf("error: %q: %s\n", err, fn)
// 	// return -1
// }
//
// // data := make([]ProductRecipe, 0)
// for rows.Next() {
// 	var (
// 		recipe_data ProductRecipe
// 	)
//
// 	if err := rows.Scan(&recipe_data.Recipe_id); err != nil {
// 		log.Fatal(err)
// 	}
// 	recipe_data.Product_id = object.Product_id
// 	log.Println("DEBUG: GetRecipes qc_data", recipe_data)
// 	object.Recipes = append(object.Recipes, &recipe_data)

// }}
/*
func forall(select_statement *sql.Stmt, start_fn, row_fn func() any, fn string, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: %q: %s\n", err, fn)
		return
	}
	start_fn()
	for rows.Next() {
		row_fn()
	}
}
component_types_list := forall(DB_Select_all_component_types,func(){},func(){}fn)
component_field.Update(component_types_list)*/

/*
func forall(select_statement *sql.Stmt, start_fn, row_fn func() any, fn string, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: %q: %s\n", err, fn)
		return
	}
	start_fn()
	for rows.Next() {
		row_fn()
	}
}
// component_types_list :=
forall(DB_Select_all_component_types,func(){},func(){
	var (
			id   int
			name string
		)

		if err := rows.Scan(&id, &name); err == nil {
			fn(id, name)
		} else {
			log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
			// return -1
		}
},fn)
component_field.Update(component_types_list)*/

//TODO component__type_list := forall(DB_Select_all_component_types,func(){},func(){}fn)
// component_field.Update(component__type_list)
// var (
// 			id   int
// 			name string
// 		)
//
// 		if err := rows.Scan(&id, &name); err == nil {
// 			fn(id, name)
// 		} else {
// 			log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
// 			// return -1
// 		}

// 	//TODO genericize
// func _select_samples(rows *sql.Rows, err error, fn string) []viewer.QCData {
// 	if err != nil {
// 		log.Printf("error: %q: %s\n", err, fn)
// 		// return -1
// 	}
//
// 	data := make([]viewer.QCData, 0)
// 	for rows.Next() {
// 		var (
// 			qc_data              viewer.QCData
// 			_timestamp           int64
// 			product_moniker_name string
// 			internal_name        string
// 		)
//
// 		if err := rows.Scan(&product_moniker_name, &internal_name,
// 			&qc_data.Product_name_customer,
// 			&qc_data.Lot_name,
// 			&qc_data.Sample_point,
// 			&_timestamp,
// 			&qc_data.PH,
// 			&qc_data.Specific_gravity,
// 			&qc_data.String_test,
// 			&qc_data.Viscosity); err != nil {
// 			log.Fatal(err)
// 		}
// 		qc_data.Product_name = product_moniker_name + " " + internal_name
//
// 		qc_data.Time_stamp = time.Unix(0, _timestamp)
// 		data = append(data, qc_data)
// 	}
// 	return data
// }

func NewRecipeProduct() *RecipeProduct {
	return new(RecipeProduct)
}

type ProductRecipe struct {
	Components []*RecipeComponent
	// Product_name string `json:"product_name"`
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Product_id int64
	Recipe_id  int64
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

func (object *ProductRecipe) GetComponents() {
	DB.Forall("GetComponents",
		func() {
			object.Components = nil
		},
		func(rows *sql.Rows) {
			Recipe_component := new(RecipeComponent)

			if err := rows.Scan(&Recipe_component.Component_name, &Recipe_component.Component_id, &Recipe_component.Component_amount, &Recipe_component.Add_order); err != nil {
				log.Fatal(err)
			}
			log.Println("DEBUG: GetComponents qc_data", Recipe_component)
			object.Components = append(object.Components, Recipe_component)
		},
		DB_Select_recipe_components, object.Recipe_id)
}

// recipe_list_id integer not null,
// component_type_id not null,
// component_type_amount real,
// component_add_order not null,
// TODO
// func (object *ProductRecipe) AddComponent() *RecipeComponent {
// 	proc_name := "ProductRecipe.AddComponent"
// 	var (
// 		component_data *RecipeComponent
// 	)
// 	result, err := DB_Insert_recipe_component.Exec(object.Product_id)
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, proc_name)
// 		return component_data
// 	}
// 	insert_id, err := result.LastInsertId()
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, proc_name)
// 		return component_data
// 	}
// 	component_data = new(RecipeComponent)
//
// 	component_data.Component_id = insert_id
// 	component_data.Product_id = object.Product_id
// 	object.Components = append(object.Components, component_data)
// 	return component_data
// }

//
//

type ComponentType struct {
	Component_name string
	Component_id   int64
}

func NewComponentType(Component_name string) *ComponentType {
	component := new(ComponentType)
	component.Component_name = Component_name
	component.Insel()
	return component
}

func (object *ComponentType) Insel() {
	object.Component_id = DB.Insel(DB_Insert_component_types, DB_Select_name_component_types, "Debug: ComponentType.Insel", object.Component_name)
}
func (object *ComponentType) AddProduct(Product_id int64) {
	DB_Insert_internal_product_component_type.Exec(object.Component_id, Product_id)
}

// ??TODO
func (object *ComponentType) AddInbound(Product_id int64) {
	DB_Insert_inbound_product_component_type.Exec(object.Component_id, Product_id)
}

//
//new comp, new from prod
//TODO
// func NewComponent() *ProductComponent {
// 	proc_name := "NewComponent"
// 	var (
// 		recipe_data *ProductComponent
// 	)
// 	result, err := DB_Insert_product_recipe.Exec(object.Product_id)
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, proc_name)
// 		return recipe_data
// 	}
// 	insert_id, err := result.LastInsertId()
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, proc_name)
// 		return recipe_data
// 	}
// 	recipe_data = new(ProductComponent)
//
// 	recipe_data.Component_id = insert_id
// 	recipe_data.Product_id = object.Product_id
// 	object.Components = append(object.Components, recipe_data)
// 	return recipe_data
// }

type RecipeComponent struct {
	Component_name   string
	Component_amount float64
	Component_id     int64
	Add_order        int64
	// Lot_id     int64
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

type BlendComponent struct {
	RecipeComponent
	// Component_name string `json:"product_name"`
	// Component_amount float64
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	// Product_id int64
	Lot_id int64
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

/*
lot_provider
lot_isomeric*/

type ProductBlend struct {
	Components []BlendComponent
	// Product_name string `json:"product_name"`
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Product_id int64
	Recipe_id  int64
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

type BlendProduct struct {
	Product_name string `json:"product_name"`
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Product_id int64
	Lot_id     int64
	Recipe     ProductBlend
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

type RecipeViewer interface {
	Get() *ProductRecipe
	Update(recipe *ProductRecipe)
	Update_component_types(component_types_list []string)
	AddComponent()
}

type RecipeView struct {
	*windigo.AutoPanel
	Recipe               *ProductRecipe
	Components           []*RecipeComponentView
	component_types_list []string
	// Product_id int64
	// Recipe_id  int64
}

func (view *RecipeView) Get() *ProductRecipe {
	if view.Recipe == nil {
		return nil
	}
	// object.Recipe.AddComponent

	// Recipe =SQLFilterDiscrete{key,
	// 	selection_options.Get()}
	return view.Recipe

}

func (view *RecipeView) Update(recipe *ProductRecipe) {
	if view.Recipe == recipe {
		return
	}

	for _, component := range view.Components {
		component.Close()
	}
	view.Components = nil
	view.Recipe = recipe
	log.Println("RecipeView Update", view.Recipe)
	if view.Recipe == nil {
		return
	}

	for _, component := range view.Recipe.Components {
		log.Println("DEBUG: RecipeView update_components", component)
		component_view := NewRecipeComponentView(view)
		view.Components = append(view.Components, component_view)
		//TODO
		// component_view.Update(component)
	}

}

func (view *RecipeView) Update_component_types(component_types_list []string) {
	view.component_types_list = component_types_list
	if view.Recipe == nil {
		return
	}

	for _, component := range view.Components {
		log.Println("DEBUG: RecipeView update_component_types", component)
		component.Update_component_types(component_types_list)
	}
}

func (view *RecipeView) AddComponent() {
	if view.Recipe != nil {
		// object.Recipe.AddComponent
		//TODO
		// (object.Recipe_id,)
		component_data := NewRecipeComponentView(view)
		view.Components = append(view.Components, component_data)
	}
}

func NewRecipeView(parent windigo.Controller) *RecipeView {
	view := new(RecipeView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	return view
}

type RecipeComponentView struct {
	*windigo.AutoPanel
	RecipeComponent *RecipeComponent
	// Component_name   string
	// Component_amount float64
	// Component_id int64
	// Add_order    int64
	Get                    func() *RecipeComponent
	Update                 func(*RecipeComponent)
	Update_component_types func(component_types_list []string)
}

func NewRecipeComponentView(parent *RecipeView) *RecipeComponentView {
	DEL_BUTTON_WIDTH := 20

	view := new(RecipeComponentView)
	view.AutoPanel = windigo.NewAutoPanel(parent)

	component_field := GUI.NewSearchBoxWithLabels(view.AutoPanel, parent.component_types_list)

	component_del_button := windigo.NewPushButton(view.AutoPanel)
	component_del_button.SetText("+")
	component_del_button.SetSize(DEL_BUTTON_WIDTH, GUI.OFF_AXIS)

	view.AutoPanel.Dock(component_field, windigo.Left)
	view.AutoPanel.Dock(component_del_button, windigo.Left)

	view.Get = func() *RecipeComponent {
		// Recipe =SQLFilterDiscrete{key,
		// 	selection_options.Get()}
		return view.RecipeComponent

	}

	view.Update_component_types = func(component_types_list []string) {
		text := component_field.Text()
		log.Println("DEBUG: RecipeComponentView update_component_types", text)
		component_field.Update(component_types_list)
		component_field.SetText(text)
		// c_view := windigo.NewLabel(component_panel)
		// c_view.SetText(component.Component_name)
	}

	return view
}

func show_window() {

	log.Println("Info: Process started")

	window_title := "QC Data Blender"

	product_text := "Product"
	recipe_text := ""
	// component_text := "Component"

	product_data := make(map[string]int64)

	add_button_width := 20
	accept_button_width := 50
	cancel_button_width := 50

	// Blend_product := new(BlendProduct)
	Recipe_product := NewRecipeProduct()
	var (
		currentRecipe *ProductRecipe
		product_list,
		component_types_list []string
	)

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)
	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)
	recipe_panel := windigo.NewAutoPanel(mainWindow)
	component_panel := windigo.NewAutoPanel(mainWindow)
	component_add_panel := windigo.NewAutoPanel(mainWindow)

	product_field := GUI.NewComboBox(product_panel, product_text)

	// product_field := GUI.NewSizedListComboBox(prod_panel, label_width, field_width, field_height, product_text)
	//
	product_add_button := windigo.NewPushButton(product_panel)
	product_add_button.SetText("+")
	product_add_button.SetSize(add_button_width, GUI.OFF_AXIS)

	recipe_sel_field := GUI.NewComboBox(recipe_panel, recipe_text)
	recipe_add_button := windigo.NewPushButton(recipe_panel)
	recipe_add_button.SetText("+")
	recipe_add_button.SetSize(add_button_width, GUI.OFF_AXIS)

	Recipe_View := NewRecipeView(mainWindow)

	// component_field := GUI.NewComboBox(component_panel, component_text)
	component_field := GUI.NewSearchBox(component_panel)
	component_add_button := windigo.NewPushButton(component_panel)
	component_add_button.SetText("+")
	component_add_button.SetSize(add_button_width, GUI.OFF_AXIS)

	component_new_button := windigo.NewPushButton(component_panel)
	component_new_button.SetText("+")
	component_new_button.SetSize(add_button_width, add_button_width)

	component_add_field := GUI.NewSearchBox(component_add_panel)
	component_accept_button := windigo.NewPushButton(component_add_panel)
	component_accept_button.SetText("OK")
	component_accept_button.SetSize(accept_button_width, GUI.OFF_AXIS)
	component_cancel_button := windigo.NewPushButton(component_add_panel)
	component_cancel_button.SetText("Cancel")
	component_cancel_button.SetSize(cancel_button_width, GUI.OFF_AXIS)
	// component_add_panel.Hide()

	// Dock
	product_panel.Dock(product_field, windigo.Left)
	product_panel.Dock(product_add_button, windigo.Left)
	recipe_panel.Dock(recipe_sel_field, windigo.Left)
	recipe_panel.Dock(recipe_add_button, windigo.Left)
	component_panel.Dock(component_new_button, windigo.Top)
	component_panel.Dock(component_field, windigo.Left)
	component_panel.Dock(component_add_button, windigo.Left)
	component_add_panel.Dock(component_add_field, windigo.Left)
	component_add_panel.Dock(component_accept_button, windigo.Left)
	component_add_panel.Dock(component_cancel_button, windigo.Left)

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(recipe_panel, windigo.Top)
	dock.Dock(Recipe_View, windigo.Top)
	dock.Dock(component_panel, windigo.Top)
	dock.Dock(component_add_panel, windigo.Top)

	// combobox

	DB.Forall("fill_combobox_from_query",
		func() {
			product_field.DeleteAllItems()
		},
		func(rows *sql.Rows) {
			var (
				id                   int64
				internal_name        string
				product_moniker_name string
			)
			if err := rows.Scan(&id, &internal_name, &product_moniker_name); err == nil {
				name := product_moniker_name + " " + internal_name
				product_data[name] = id

				product_list = append(product_list, name)

				product_field.AddItem(name)
			} else {
				log.Printf("error: %q: %s\n", err, "fill_combobox_from_query")
				// return -1
			}
		},

		DB.DB_Select_product_info)

	//TODO fixme
	component_field.Update(product_list)
	component_add_field.Update(product_list)

	// functionality

	// sizing
	refresh_vars := func(font_size int) {
		refresh_globals(font_size)

	}
	refresh := func(font_size int) {
		refresh_vars(font_size)

		product_panel.SetSize(TOP_PANEL_WIDTH, PRODUCT_FIELD_HEIGHT)
		recipe_panel.SetSize(TOP_PANEL_WIDTH, PRODUCT_FIELD_HEIGHT)
		//TODO grow
		// component_panel.SetSize(TOP_PANEL_WIDTH, PRODUCT_FIELD_HEIGHT)
		// component_panel.SetSize(TOP_PANEL_WIDTH, 2*PRODUCT_FIELD_HEIGHT)
		component_panel.SetSize(TOP_PANEL_WIDTH, PRODUCT_FIELD_HEIGHT+add_button_width)

		component_add_panel.SetSize(TOP_PANEL_WIDTH, PRODUCT_FIELD_HEIGHT)

		product_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)
		recipe_sel_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)
		component_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)
		component_add_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)

	}

	//size
	set_font := func(font_size int) {
		GUI.BASE_FONT_SIZE = font_size

		config.Main_config.Set("font_size", GUI.BASE_FONT_SIZE)
		config.Write_config(config.Main_config)

		old_font := windigo.DefaultFont
		windigo.DefaultFont = windigo.NewFont(old_font.Family(), GUI.BASE_FONT_SIZE, 0)
		old_font.Dispose()

		mainWindow.SetFont(windigo.DefaultFont)
		product_field.SetFont(windigo.DefaultFont)
		product_add_button.SetFont(windigo.DefaultFont)
		recipe_sel_field.SetFont(windigo.DefaultFont)
		recipe_add_button.SetFont(windigo.DefaultFont)
		component_field.SetFont(windigo.DefaultFont)
		component_add_button.SetFont(windigo.DefaultFont)
		component_new_button.SetFont(windigo.DefaultFont)
		component_add_field.SetFont(windigo.DefaultFont)
		component_accept_button.SetFont(windigo.DefaultFont)
		component_cancel_button.SetFont(windigo.DefaultFont)
		// threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}

	//TODO
	update_components := func(object *ProductRecipe) {
		for _, component := range object.Components {
			log.Println("DEBUG: update_components", component)

			c_view := windigo.NewLabel(component_panel)
			c_view.SetText(component.Component_name)
		}

	}
	update_component_types := func() {

		DB.Forall("update_component_types",
			func() { component_types_list = nil },
			func(rows *sql.Rows) {
				var (
					id   int
					name string
				)

				if err := rows.Scan(&id, &name); err == nil {
					// component_types_data[name] = id
					log.Println("DEBUG: update_component_types nsme", id, name)

					component_types_list = append(component_types_list, name)

					// component_field.AddItem(name)
				} else {
					log.Printf("error: %q: %s\n", err, "fgjifjofdgkjfokgf")
					// return -1
				}
			},
			DB_Select_all_component_types)
		log.Println("DEBUG: update_component_types", component_types_list)
		component_field.Update(component_types_list)
		Recipe_View.Update_component_types(component_types_list)
	}
	update_component_types()

	set_font(GUI.BASE_FONT_SIZE)

	//event handling
	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		Recipe_product = NewRecipeProduct()
		name := product_field.GetSelectedItem()
		Recipe_product.Set(name, product_data[name])
		// Recipe_product.GetRecipes()
		Recipe_product.LoadRecipeCombo(recipe_sel_field)
		Recipe_product.GetComponents()
		recipe_sel_field.OnSelectedChange().Fire(nil)
	})

	product_add_button.OnClick().Bind(func(e *windigo.Event) {
		log.Println("TODO product_add_button")

		// clear()
		// clear_product()
	})
	recipe_add_button.OnClick().Bind(func(e *windigo.Event) {
		if Recipe_product.Product_id != DB.INVALID_ID {

			numRecipes := len(Recipe_product.Recipes)
			currentRecipe = Recipe_product.NewRecipe()

			recipe_sel_field.AddItem(strconv.Itoa(numRecipes))
			recipe_sel_field.SetSelectedItem(numRecipes)
			Recipe_View.Update(currentRecipe)

			// log.Println("product_field", product_field.SelectedItem())

		}
	})

	recipe_sel_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		// TODO currentRecipe
		currentRecipe = nil
		i, err := strconv.Atoi(recipe_sel_field.GetSelectedItem())
		if err != nil {
			log.Println("ERR: recipe_field strconv", err)
			Recipe_View.Update(currentRecipe)
			return
		}
		currentRecipe = Recipe_product.Recipes[i]
		update_components(currentRecipe)
		Recipe_View.Update(currentRecipe)

		// Recipe_product = NewRecipeProduct()
		// name := product_field.GetSelectedItem()
		// Recipe_product.Set(name, product_data[name])
		// // Recipe_product.GetRecipes()
		// Recipe_product.LoadRecipeCombo(recipe_field)
	})

	component_add_button.OnClick().Bind(func(e *windigo.Event) {
		if currentRecipe != nil {
			//?TODO
			Recipe_View.AddComponent()
			// currentRecipe.AddComponent

			// component_field.AddItem(strconv.Itoa(len(Recipe_product.Recipes)))
			// log.Println("product_field", product_field.SelectedItem())

		}
	})

	component_new_button.OnClick().Bind(func(e *windigo.Event) {
		component_add_panel.Show()
	})

	component_accept_button.OnClick().Bind(func(e *windigo.Event) {
		component := NewComponentType(component_add_field.Text())
		Product_id := product_data[component.Component_name]
		if Product_id != DB.INVALID_ID {
			component.AddProduct(Product_id)
		}
		component_add_panel.Hide()
		update_component_types()
	})
	component_cancel_button.OnClick().Bind(func(e *windigo.Event) {
		component_add_panel.Hide()
	})

	// component_add_field.OnChange().Bind(func(e *windigo.Event) {
	// 	log.Println("CRIT: DEBUG: component_add_field OnChange", component_add_field.Text(), component_add_field.SelectedItem())
	// 	component_add_field.ShowDropdown(true)
	//
	// 	// component_add_field.SetText(text)
	//
	// 	// data_view.SelectText(start, end)
	// 	// component_add_field.SelectText(start, -1)
	// })

	// component_field.OnSelectedChange().Bind(func(e *windigo.Event) {
	//
	// 	i, err := strconv.Atoi(component_field.GetSelectedItem())
	// 	if err != nil {
	// 		log.Println("ERR: component_field strconv", err)
	// 	}
	// 	update_components(Recipe_product.Recipes[i])
	//
	// 	// Recipe_product = NewRecipeProduct()
	// 	// name := product_field.GetSelectedItem()
	// 	// Recipe_product.Set(name, product_data[name])
	// 	// // Recipe_product.GetRecipes()
	// 	// Recipe_product.LoadRecipeCombo(component_field)
	// })

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(func(*windigo.Event) {
		windigo.Exit()
	})
	mainWindow.RunMainLoop() // Must call to start event loop.
}
