package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/qc_data_entry/config"
)

func main() {
	//load config
	config.Main_config = config.Load_config()

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
	qc_db, err := sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", config.DB_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer qc_db.Close()

	dbinit(qc_db)

	//setup print goroutine
	print_queue = make(chan string, 4)
	defer close(print_queue)
	go do_print_queue(print_queue)

	//setup status_bar goroutine
	status_queue = make(chan string, 16)
	defer close(status_queue)
	go do_status_queue(status_queue)

	//show main window
	show_window()

}
