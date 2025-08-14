package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views"
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
)

func dbinit(db *sql.DB) {

	DB.Check_db(db)
	DB.DBinit(db)

}

func refresh_globals(font_size int) {

	GUI.GROUPBOX_CUSHION = font_size * 3 / 2
	TOP_SPACER_WIDTH = 7
	TOP_SPACER_HEIGHT = GUI.GROUPBOX_CUSHION + 2
	INTER_SPACER_HEIGHT = 2
	BTM_SPACER_WIDTH = 2
	BTM_SPACER_HEIGHT = 2
	TOP_PANEL_INTER_SPACER_WIDTH = 30

	GUI.DATA_MARGIN = 10

	GUI.LABEL_WIDTH = 10 * font_size
	GUI.PRODUCT_FIELD_WIDTH = 15 * font_size
	GUI.PRODUCT_FIELD_HEIGHT = font_size*16/10 + 8
	GUI.EDIT_FIELD_HEIGHT = 24
	GUI.EDIT_FIELD_HEIGHT = 3*font_size - 2

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

	// component_text := "Component"

	product_data := make(map[string]int64)

	ADD_BUTTON_WIDTH := 20

	// Blend_product := new(BlendProduct)
	var (
		product_list []string
	)

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)
	dock := windigo.NewSimpleDock(mainWindow)

	component_panel := windigo.NewAutoPanel(mainWindow)
	component_add_panel := windigo.NewAutoPanel(mainWindow)

	// product_field := GUI.NewSizedListComboBox(prod_panel, label_width, field_width, field_height, product_text)
	//

	Blend_product_view := views.NewBlendProductView(mainWindow)

	// Blend_View := blender.NewBlendView(mainWindow)

	// TODO move this too to view
	component_new_button := windigo.NewPushButton(component_panel)
	component_new_button.SetText("+")
	component_new_button.SetSize(ADD_BUTTON_WIDTH, ADD_BUTTON_WIDTH)

	component_add_field := GUI.NewSearchBox(component_add_panel)
	component_accept_button := windigo.NewPushButton(component_add_panel)
	component_accept_button.SetText("OK")
	component_accept_button.SetSize(GUI.ACCEPT_BUTTON_WIDTH, GUI.OFF_AXIS)
	component_cancel_button := windigo.NewPushButton(component_add_panel)
	component_cancel_button.SetText("Cancel")
	component_cancel_button.SetSize(GUI.CANCEL_BUTTON_WIDTH, GUI.OFF_AXIS)
	component_add_panel.Hide()

	// Dock

	component_panel.Dock(component_new_button, windigo.Top)
	component_add_panel.Dock(component_add_field, windigo.Left)
	component_add_panel.Dock(component_accept_button, windigo.Left)
	component_add_panel.Dock(component_cancel_button, windigo.Left)

	// dock.Dock(Blend_View, windigo.Top)
	// dock.Dock(component_panel, windigo.Top)
	dock.Dock(component_add_panel, windigo.Top)
	dock.Dock(Blend_product_view, windigo.Fill)

	// combobox

	GUI.Fill_combobox_from_query_rows(Blend_product_view.Product_Field, func(row *sql.Rows) error {
		var (
			id                   int64
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		product_data[name] = id

		product_list = append(product_list, name)

		Blend_product_view.Product_Field.AddItem(name)
		return nil
	},

		DB.DB_Select_product_info)
	// DB.Forall("fill_combobox_from_query",
	// 	func() {
	// 		product_field.DeleteAllItems()
	// 	},
	// 	func(row *sql.Rows) {
	// 		var (
	// 			id                   int64
	// 			internal_name        string
	// 			product_moniker_name string
	// 		)
	// 		if err := row.Scan(&id, &internal_name, &product_moniker_name); err == nil {
	// 			name := product_moniker_name + " " + internal_name
	// 			product_data[name] = id
	//
	// 			product_list = append(product_list, name)
	//
	// 			product_field.AddItem(name)
	// 		} else {
	// 			log.Printf("error: [%s]: %q\n", "fill_combobox_from_query", err)
	// 			// return -1
	// 		}
	// 	},
	//
	// 	DB.DB_Select_product_info)

	//TODO fixme
	// component_field.Update(product_list)
	component_add_field.Update(product_list)
	Blend_product_view.Update_products(product_data)

	// functionality

	// sizing
	refresh_vars := func(font_size int) {
		refresh_globals(font_size)

	}
	refresh := func(font_size int) {
		refresh_vars(font_size)
		//TODO grow
		// component_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		// component_panel.SetSize(TOP_PANEL_WIDTH, 2*GUI.PRODUCT_FIELD_HEIGHT)
		// component_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT+ADD_BUTTON_WIDTH)

		//TODO grow
		// Blend_View.SetSize(GUI.LABEL_WIDTH+TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		// Blend_View.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		Blend_product_view.RefreshSize()
		// Blend_View.RefreshSize()

		component_add_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)

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
		Blend_product_view.SetFont(windigo.DefaultFont)

		component_new_button.SetFont(windigo.DefaultFont)
		component_add_field.SetFont(windigo.DefaultFont)
		component_accept_button.SetFont(windigo.DefaultFont)
		component_cancel_button.SetFont(windigo.DefaultFont)
		// threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}

	set_font(GUI.BASE_FONT_SIZE)

	//event handling

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

	// component_text := "Component"

	product_data := make(map[string]int64)

	// Blend_product := new(BlendProduct)
	var (
		product_list []string
	)

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)
	dock := windigo.NewSimpleDock(mainWindow)

	Recipe_product_view := views.NewRecipeProductView(mainWindow)

	dock.Dock(Recipe_product_view, windigo.Fill)

	// combobox

	GUI.Fill_combobox_from_query_rows(Recipe_product_view.Product_Field, func(row *sql.Rows) error {
		var (
			id                   int64
			internal_name        string
			product_moniker_name string
		)
		if err := row.Scan(&id, &internal_name, &product_moniker_name); err != nil {
			return err
		}
		name := product_moniker_name + " " + internal_name
		product_data[name] = id

		product_list = append(product_list, name)

		Recipe_product_view.Product_Field.AddItem(name)
		return nil
	},

		DB.DB_Select_product_info)
	// DB.Forall("fill_combobox_from_query",
	// 	func() {
	// 		product_field.DeleteAllItems()
	// 	},
	// 	func(row *sql.Rows) {
	// 		var (
	// 			id                   int64
	// 			internal_name        string
	// 			product_moniker_name string
	// 		)
	// 		if err := row.Scan(&id, &internal_name, &product_moniker_name); err == nil {
	// 			name := product_moniker_name + " " + internal_name
	// 			product_data[name] = id
	//
	// 			product_list = append(product_list, name)
	//
	// 			product_field.AddItem(name)
	// 		} else {
	// 			log.Printf("error: [%s]: %q\n", "fill_combobox_from_query", err)
	// 			// return -1
	// 		}
	// 	},
	//
	// 	DB.DB_Select_product_info)

	Recipe_product_view.Update_products(product_list, product_data)

	// functionality

	// sizing
	refresh_vars := func(font_size int) {
		refresh_globals(font_size)

	}
	refresh := func(font_size int) {
		refresh_vars(font_size)
		//TODO grow
		// component_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		// component_panel.SetSize(TOP_PANEL_WIDTH, 2*GUI.PRODUCT_FIELD_HEIGHT)
		// component_panel.SetSize(TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT+ADD_BUTTON_WIDTH)

		//TODO grow
		// Recipe_View.SetSize(GUI.LABEL_WIDTH+TOP_PANEL_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		// Recipe_View.SetLabeledSize(GUI.LABEL_WIDTH, GUI.PRODUCT_FIELD_WIDTH, GUI.PRODUCT_FIELD_HEIGHT)
		Recipe_product_view.RefreshSize()
		// Recipe_View.RefreshSize()

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
		Recipe_product_view.SetFont(windigo.DefaultFont)

		// threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}
	set_font(GUI.BASE_FONT_SIZE)

	//event handling

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(func(*windigo.Event) {
		windigo.Exit()
	})
	// Must call to start event loop.
	mainWindow.RunMainLoop()
}
