package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/qc_data_entry/GUI/views/toplevel_ui"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/qc_data_entry/viewer"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// TODO
// Dynamic-column pivot table
// pivot_vtab

func main() {
	//load config
	config.Main_config = config.Load_config_viewer("qc_data_viewer")

	// log to file
	log_file, err := os.OpenFile(config.LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Crit: error opening file: %v", err)
	}
	defer log_file.Close()
	log.Println("Info: Logging to logfile:", config.LOG_FILE)

	log.SetOutput(log_file)
	log.Println("Info: Using config:", config.Main_config.ConfigFileUsed())

	//open_db
	// viewer.QC_DB, err := sql.Open("sqlite3", DB_FILE)
	viewer.QC_DB, err = sql.Open("sqlite3", ":memory:")
	viewer.QC_DB.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal("Crit: error opening database: ", err)
	}
	defer viewer.QC_DB.Close()
	log.Println("Info: Using db:", config.DB_FILE)
	viewer.DBinit(viewer.QC_DB)

	//setup print goroutine
	threads.PRINT_QUEUE = make(chan string, 4)
	defer close(threads.PRINT_QUEUE)
	go threads.Do_print_queue(threads.PRINT_QUEUE)

	//setup status_bar goroutine
	threads.STATUS_QUEUE = make(chan string, 16)
	defer close(threads.STATUS_QUEUE)
	go threads.Do_status_queue(threads.STATUS_QUEUE)

	viewer.Refresh_globals(GUI.BASE_FONT_SIZE)

	//show main window
	toplevel_ui.Show_window(viewer.NewViewerWindow(nil))
}
