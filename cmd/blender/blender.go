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
	DB_Select_recipe_components, DB_Insert_recipe_component *sql.Stmt
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
	object.Recipes = nil
	rows, err := DB_Select_product_recipe.Query(object.Product_id)
	fn := "GetRecipes"
	if err != nil {
		log.Printf("error: %q: %s\n", err, fn)
		// return -1
	}

	// data := make([]ProductRecipe, 0)
	for rows.Next() {
		var (
			recipe_data ProductRecipe
		)

		if err := rows.Scan(&recipe_data.Recipe_id); err != nil {
			log.Fatal(err)
		}
		recipe_data.Product_id = object.Product_id
		log.Println("DEBUG: GetRecipes qc_data", recipe_data)
		object.Recipes = append(object.Recipes, &recipe_data)

	}

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
	object.Components = nil
	rows, err := DB_Select_recipe_components.Query(object.Recipe_id)
	fn := "GetComponents"
	if err != nil {
		log.Printf("error: %q: %s\n", err, fn)
		// return -1
	}

	// data := make([]ProductRecipe, 0)
	for rows.Next() {
		Recipe_component := new(RecipeComponent)

		if err := rows.Scan(&Recipe_component.Component_name, &Recipe_component.Component_id, &Recipe_component.Component_amount, &Recipe_component.Add_order); err != nil {
			log.Fatal(err)
		}
		log.Println("DEBUG: GetComponents qc_data", Recipe_component)
		object.Components = append(object.Components, Recipe_component)
	}

}

// TODO
func (object *ProductRecipe) AddComponent() *RecipeComponent {
	proc_name := "ProductRecipe.AddComponent"
	var (
		component_data *RecipeComponent
	)
	result, err := DB_Insert_recipe_component.Exec(object.Product_id)
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return component_data
	}
	insert_id, err := result.LastInsertId()
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return component_data
	}
	component_data = new(ProductComponent)

	component_data.Component_id = insert_id
	component_data.Product_id = object.Product_id
	object.Components = append(object.Components, component_data)
	return component_data
}

type RecipeComponent struct {
	Component_name   string `json:"product_name"`
	Component_amount float64
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Component_id int64
	Add_order    int64
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

func show_window() {

	log.Println("Info: Process started")

	window_title := "QC Data Blender"

	product_text := "Product"
	recipe_text := ""
	component_text := "Component"

	product_data := make(map[string]int64)

	add_button_width := 20

	// Blend_product := new(BlendProduct)
	Recipe_product := NewRecipeProduct()
	var (
		currentRecipe *ProductRecipe
	)

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)
	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)
	recipe_panel := windigo.NewAutoPanel(mainWindow)
	component_panel := windigo.NewAutoPanel(mainWindow)

	product_field := GUI.NewComboBox(product_panel, product_text)

	// product_field := GUI.NewSizedListComboBox(prod_panel, label_width, field_width, field_height, product_text)
	//
	product_add_button := windigo.NewPushButton(product_panel)
	product_add_button.SetText("+")
	product_add_button.SetSize(add_button_width, GUI.OFF_AXIS)

	recipe_field := GUI.NewComboBox(recipe_panel, recipe_text)
	recipe_add_button := windigo.NewPushButton(recipe_panel)
	recipe_add_button.SetText("+")
	recipe_add_button.SetSize(add_button_width, GUI.OFF_AXIS)

	component_field := GUI.NewComboBox(component_panel, component_text)
	component_add_button := windigo.NewPushButton(component_panel)
	component_add_button.SetText("+")
	component_add_button.SetSize(add_button_width, GUI.OFF_AXIS)

	// Dock
	product_panel.Dock(product_field, windigo.Left)
	product_panel.Dock(product_add_button, windigo.Left)
	recipe_panel.Dock(recipe_field, windigo.Left)
	recipe_panel.Dock(recipe_add_button, windigo.Left)
	component_panel.Dock(component_field, windigo.Left)
	component_panel.Dock(component_add_button, windigo.Left)
	dock.Dock(product_panel, windigo.Top)
	dock.Dock(recipe_panel, windigo.Top)
	dock.Dock(component_panel, windigo.Top)

	// combobox
	rows, err := DB.DB_Select_product_info.Query()
	GUI.Fill_combobox_from_query_rows(product_field, rows, err, func(rows *sql.Rows) {
		var (
			id                   int64
			internal_name        string
			product_moniker_name string
		)
		if err := rows.Scan(&id, &internal_name, &product_moniker_name); err == nil {
			name := product_moniker_name + " " + internal_name
			product_data[name] = id

			product_field.AddItem(name)
		}
	})

	// functionality

	// sizing
	refresh_vars := func(font_size int) {
		refresh_globals(font_size)
	}
	refresh := func(font_size int) {
		refresh_vars(font_size)

		product_panel.SetSize(TOP_PANEL_WIDTH, PRODUCT_FIELD_HEIGHT)

		product_field.SetLabeledSize(GUI.LABEL_WIDTH, PRODUCT_FIELD_WIDTH, PRODUCT_FIELD_HEIGHT)

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
		// threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}

	update_components := func(object *ProductRecipe) {
		for _, component := range object.Components {
			log.Println("DEBUG: update_components", component)

			c_view := windigo.NewLabel(component_panel)
			c_view.SetText(component.Component_name)
		}

	}

	set_font(GUI.BASE_FONT_SIZE)

	//event handling
	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		Recipe_product = NewRecipeProduct()
		name := product_field.GetSelectedItem()
		Recipe_product.Set(name, product_data[name])
		// Recipe_product.GetRecipes()
		Recipe_product.LoadRecipeCombo(recipe_field)
		Recipe_product.GetComponents()
		currentRecipe = nil
	})

	product_add_button.OnClick().Bind(func(e *windigo.Event) {
		// log.Println("product_field", product_field.SelectedItem())

		// clear()
		// clear_product()
	})
	recipe_add_button.OnClick().Bind(func(e *windigo.Event) {
		if Recipe_product.Product_id != DB.INVALID_ID {
			//?TODO add  recip button
			//recip := Recipe_product.NewRecipe()
			currentRecipe = Recipe_product.NewRecipe()
			numRecipes := len(Recipe_product.Recipes)

			recipe_field.AddItem(strconv.Itoa(numRecipes))
			recipe_field.SetSelectedItem(numRecipes)

			// log.Println("product_field", product_field.SelectedItem())

		}
	})

	recipe_field.OnSelectedChange().Bind(func(e *windigo.Event) {

		i, err := strconv.Atoi(recipe_field.GetSelectedItem())
		if err != nil {
			log.Println("ERR: recipe_field strconv", err)
		}
		currentRecipe = Recipe_product.Recipes[i]
		update_components(currentRecipe)

		// Recipe_product = NewRecipeProduct()
		// name := product_field.GetSelectedItem()
		// Recipe_product.Set(name, product_data[name])
		// // Recipe_product.GetRecipes()
		// Recipe_product.LoadRecipeCombo(recipe_field)
	})

	component_add_button.OnClick().Bind(func(e *windigo.Event) {
		if currentRecipe != nil {
			//?TODO
			// currentRecipe.AddComponent

			// component_field.AddItem(strconv.Itoa(len(Recipe_product.Recipes)))
			// log.Println("product_field", product_field.SelectedItem())

		}
	})

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
