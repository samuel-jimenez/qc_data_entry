package main

import (
	"database/sql"
	"log"
	"os"

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

	RANGES_RO_PADDING,
	RANGES_RO_FIELD_WIDTH,
	RANGES_RO_SPACER_WIDTH,

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
	DB_Select_recipe_components *sql.Stmt
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
	select component_type_name, component_type_amount
		from bs.recipe_components
		join bs.component_types
		using (component_type_id)
		where recipe_list_id = ?
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
	Product_id int
	// Lot_id     int64
	Recipes []ProductRecipe
	// Product_name_customer_id nullable.NullInt64
	// Product_name_customer    string `json:"customer_product_name"`
}

func (r *RecipeProduct) Set(name string, i int) {
	r.Product_name = name
	r.Product_id = i
}

func (r *RecipeProduct) GetRecipes() {
	r.Recipes = DB_Select_product_recipe.Query(r.Product_id)
	//TODO
}

func NewRecipeProduct() *RecipeProduct {
	return new(RecipeProduct)
}

type ProductRecipe struct {
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

type RecipeComponent struct {
	Component_name   string `json:"product_name"`
	Component_amount float64
	// Lot_number               string `json:"lot_number"`
	// Sample_point             string
	// Visual                   bool
	Component_id int64
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

	product_data := make(map[string]int)

	// Blend_product := new(BlendProduct)
	Recipe_product := NewRecipeProduct()

	// build window
	mainWindow := windigo.NewForm(nil)
	mainWindow.SetText(window_title)
	dock := windigo.NewSimpleDock(mainWindow)

	product_panel := windigo.NewAutoPanel(mainWindow)

	product_field := GUI.NewComboBox(product_panel, product_text)

	// product_field := GUI.NewSizedListComboBox(prod_panel, label_width, field_width, field_height, product_text)

	// Dock

	product_panel.Dock(product_field, windigo.Left)
	dock.Dock(product_panel, windigo.Top)

	// combobox
	rows, err := DB.DB_Select_product_info.Query()
	GUI.Fill_combobox_from_query_rows(product_field, rows, err, func(rows *sql.Rows) {
		var (
			id                   int
			internal_name        string
			product_moniker_name string
		)
		if err := rows.Scan(&id, &internal_name, &product_moniker_name); err == nil {
			name := product_moniker_name + " " + internal_name
			product_data[name] = id

			product_field.AddItem(name)
		}
	})

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
		// threads.Status_bar.SetFont(windigo.DefaultFont)
		refresh(font_size)

	}
	set_font(GUI.BASE_FONT_SIZE)

	// functionality

	product_field.OnSelectedChange().Bind(func(e *windigo.Event) {
		Recipe_product = NewRecipeProduct()
		name := product_field.GetSelectedItem()
		Recipe_product.Set(name, product_data[name])
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(func(*windigo.Event) {
		windigo.Exit()
	})
	mainWindow.RunMainLoop() // Must call to start event loop.
}
