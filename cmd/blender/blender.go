package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/blender"
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
		log.Fatalf("Crit: error opening file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.Println("Info: Logging Started")

	//open_db
	// qc_db, err := sql.Open("sqlite3", DB_FILE)
	qc_db, err = sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal("Crit: error opening database: ", err)
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
	// show_Formulator()

}

var (
	qc_db *sql.DB
	// DB_Select_component_type_product,
	DB_Insert_internal_product_component_type,
	DB_Insert_inbound_product_component_type *sql.Stmt
)

func dbinit(db *sql.DB) {

	DB.Check_db(db)
	DB.DBinit(db)

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
	GUI.PRODUCT_FIELD_WIDTH = 15 * font_size
	GUI.PRODUCT_FIELD_HEIGHT = font_size*16/10 + 8
	FIELD_HEIGHT = 24

}

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
	object.Component_id = DB.Insel(DB.DB_Insert_component_types, DB.DB_Select_name_component_types, "Debug: ComponentType.Insel", object.Component_name)
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
// 		log.Printf("Err: [%s]: %q\n", proc_name, err)
// 		return recipe_data
// 	}
// 	insert_id, err := result.LastInsertId()
// 	if err != nil {
// 		log.Printf("Err: [%s]: %q\n", proc_name, err)
// 		return recipe_data
// 	}
// 	recipe_data = new(ProductComponent)
//
// 	recipe_data.Component_id = insert_id
// 	recipe_data.Product_id = object.Product_id
// 	object.Components = append(object.Components, recipe_data)
// 	return recipe_data
// }

func show_window() {

	log.Println("Info: Process started")

	window_title := "QC Data Blender"

	product_text := "Product"
	recipe_text := ""
	// component_text := "Component"

	product_data := make(map[string]int64)
	component_types_data := make(map[string]int64)

	add_button_width := 20
	accept_button_width := 50
	// cancel_button_width := 50

	// Blend_product := new(BlendProduct)
	Recipe_product := blender.NewRecipeProduct()
	var (
		product_list,
		component_types_list []string
	)

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)
	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)
	recipe_panel := windigo.NewAutoPanel(mainWindow)

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

	recipe_accept_button := windigo.NewPushButton(recipe_panel)
	recipe_accept_button.SetText("OK")
	recipe_accept_button.SetSize(accept_button_width, GUI.OFF_AXIS)

	Recipe_View := blender.NewRecipeView(mainWindow)
	//TODO BlendView

	// Dock
	product_panel.Dock(product_field, windigo.Left)
	product_panel.Dock(product_add_button, windigo.Left)
	recipe_panel.Dock(recipe_sel_field, windigo.Left)
	recipe_panel.Dock(recipe_add_button, windigo.Left)
	recipe_panel.Dock(recipe_accept_button, windigo.Left)

	dock.Dock(product_panel, windigo.Top)
	dock.Dock(recipe_panel, windigo.Top)
	dock.Dock(Recipe_View, windigo.Top)

	// combobox

	DB.Forall("fill_combobox_from_query",
		func() {
			product_field.DeleteAllItems()
		},
		func(row *sql.Rows) {
			var (
				id                   int64
				internal_name        string
				product_moniker_name string
			)
			if err := row.Scan(&id, &internal_name, &product_moniker_name); err == nil {
				name := product_moniker_name + " " + internal_name
				product_data[name] = id

				product_list = append(product_list, name)

				product_field.AddItem(name)
			} else {
				log.Printf("error: [%s]: %q\n", "fill_combobox_from_query", err)
				// return -1
			}
		},

		DB.DB_Select_product_info)

	// functionality

	// sizing
	refresh_vars := func(font_size int) {
		refresh_globals(font_size)

	}
	refresh := func(font_size int) {
		refresh_vars(font_size)

		product_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		recipe_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		Recipe_View.RefreshSize()

		product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		recipe_sel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

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
		Recipe_View.SetFont(windigo.DefaultFont)
		recipe_sel_field.SetFont(windigo.DefaultFont)
		recipe_add_button.SetFont(windigo.DefaultFont)
		recipe_accept_button.SetFont(windigo.DefaultFont)
		// threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}

	//TODO move
	update_component_types := func() {

		DB.Forall("update_component_types",
			func() { component_types_list = nil },
			func(row *sql.Rows) {
				var (
					id   int64
					name string
				)

				if err := row.Scan(&id, &name); err == nil {
					component_types_data[name] = id
					log.Println("DEBUG: update_component_types nsme", id, name)

					component_types_list = append(component_types_list, name)

					// component_field.AddItem(name)
				} else {
					log.Printf("error: [%s]: %q\n", "fgjifjofdgkjfokgf", err)
					// return -1
				}
			},
			DB.DB_Select_all_component_types)
		log.Println("DEBUG: update_component_types", component_types_list)
		Recipe_View.Update_component_types(component_types_list, component_types_data)
		// Recipe_View.Update_component_types(component_types_list)
	}
	update_component_types()

	set_font(GUI.BASE_FONT_SIZE)

	//event handling
	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		//TODO move to view
		Recipe_product = blender.NewRecipeProduct()
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

			recipe_sel_field.AddItem(strconv.Itoa(numRecipes))
			recipe_sel_field.SetSelectedItem(numRecipes)
			Recipe_View.Update(Recipe_product.NewRecipe())

			// log.Println("product_field", product_field.SelectedItem())

		}
	})
	recipe_accept_button.OnClick().Bind(func(e *windigo.Event) {

		log.Println("ERR: DEBUG: Recipe_View Get:", Recipe_View.Get())

	})

	recipe_sel_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		i, err := strconv.Atoi(recipe_sel_field.GetSelectedItem())
		if err != nil {
			log.Println("ERR: recipe_field strconv", err)
			Recipe_View.Update(nil)
			return
		}
		Recipe_View.Update(Recipe_product.Recipes[i])

		// Recipe_product = NewRecipeProduct()
		// name := product_field.GetSelectedItem()
		// Recipe_product.Set(name, product_data[name])
		// // Recipe_product.GetRecipes()
		// Recipe_product.LoadRecipeCombo(recipe_field)
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

func show_Formulator() {

	log.Println("Info: Process started")

	window_title := "QC Data Formulator"

	product_text := "Product"
	recipe_text := ""
	// component_text := "Component"

	product_data := make(map[string]int64)
	component_types_data := make(map[string]int64)

	add_button_width := 20
	accept_button_width := 50
	cancel_button_width := 50

	// Blend_product := new(BlendProduct)
	Recipe_product := blender.NewRecipeProduct()
	var (
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

	recipe_accept_button := windigo.NewPushButton(recipe_panel)
	recipe_accept_button.SetText("OK")
	recipe_accept_button.SetSize(accept_button_width, GUI.OFF_AXIS)

	Recipe_View := blender.NewRecipeView(mainWindow)

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
	component_add_panel.Hide()

	// Dock
	product_panel.Dock(product_field, windigo.Left)
	product_panel.Dock(product_add_button, windigo.Left)
	recipe_panel.Dock(recipe_sel_field, windigo.Left)
	recipe_panel.Dock(recipe_add_button, windigo.Left)
	recipe_panel.Dock(recipe_accept_button, windigo.Left)

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
		func(row *sql.Rows) {
			var (
				id                   int64
				internal_name        string
				product_moniker_name string
			)
			if err := row.Scan(&id, &internal_name, &product_moniker_name); err == nil {
				name := product_moniker_name + " " + internal_name
				product_data[name] = id

				product_list = append(product_list, name)

				product_field.AddItem(name)
			} else {
				log.Printf("error: [%s]: %q\n", "fill_combobox_from_query", err)
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

		product_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		recipe_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		//TODO grow
		// component_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		// component_panel.SetSize(TOP_PANEL_WIDTH, 2*GUI.PRODUCT_FIELD_HEIGHT)
		component_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT+add_button_width)

		//TODO grow
		// Recipe_View.SetSize(GUI.LABEL_WIDTH+TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		// Recipe_View.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		Recipe_View.RefreshSize()

		component_add_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

		product_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		recipe_sel_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		component_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		component_add_field.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

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
		Recipe_View.SetFont(windigo.DefaultFont)
		recipe_sel_field.SetFont(windigo.DefaultFont)
		recipe_add_button.SetFont(windigo.DefaultFont)
		recipe_accept_button.SetFont(windigo.DefaultFont)
		component_field.SetFont(windigo.DefaultFont)
		component_add_button.SetFont(windigo.DefaultFont)
		component_new_button.SetFont(windigo.DefaultFont)
		component_add_field.SetFont(windigo.DefaultFont)
		component_accept_button.SetFont(windigo.DefaultFont)
		component_cancel_button.SetFont(windigo.DefaultFont)
		// threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}

	//TODO move
	update_component_types := func() {

		DB.Forall("update_component_types",
			func() { component_types_list = nil },
			func(row *sql.Rows) {
				var (
					id   int64
					name string
				)

				if err := row.Scan(&id, &name); err == nil {
					component_types_data[name] = id
					log.Println("DEBUG: update_component_types nsme", id, name)

					component_types_list = append(component_types_list, name)

					// component_field.AddItem(name)
				} else {
					log.Printf("error: [%s]: %q\n", "fgjifjofdgkjfokgf", err)
					// return -1
				}
			},
			DB.DB_Select_all_component_types)
		log.Println("DEBUG: update_component_types", component_types_list)
		component_field.Update(component_types_list)
		Recipe_View.Update_component_types(component_types_list, component_types_data)
		// Recipe_View.Update_component_types(component_types_list)
	}
	update_component_types()

	set_font(GUI.BASE_FONT_SIZE)

	//event handling
	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		//TODO move to view
		Recipe_product = blender.NewRecipeProduct()
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

			recipe_sel_field.AddItem(strconv.Itoa(numRecipes))
			recipe_sel_field.SetSelectedItem(numRecipes)
			Recipe_View.Update(Recipe_product.NewRecipe())

			// log.Println("product_field", product_field.SelectedItem())

		}
	})
	recipe_accept_button.OnClick().Bind(func(e *windigo.Event) {

		log.Println("ERR: DEBUG: Recipe_View Get:", Recipe_View.Get())

	})

	recipe_sel_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		i, err := strconv.Atoi(recipe_sel_field.GetSelectedItem())
		if err != nil {
			log.Println("ERR: recipe_field strconv", err)
			Recipe_View.Update(nil)
			return
		}
		Recipe_View.Update(Recipe_product.Recipes[i])

		// Recipe_product = NewRecipeProduct()
		// name := product_field.GetSelectedItem()
		// Recipe_product.Set(name, product_data[name])
		// // Recipe_product.GetRecipes()
		// Recipe_product.LoadRecipeCombo(recipe_field)
	})

	component_add_button.OnClick().Bind(func(e *windigo.Event) {
		if Recipe_View.Recipe != nil {
			Recipe_View.AddComponent()
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
