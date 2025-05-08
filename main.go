package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/spf13/viper"
)

var main_config *viper.Viper
var (
	DB_PATH,
	DB_FILE,
	LABEL_PATH string
	LOG_FILE string

	JSON_PATHS []string

	print_queue       chan string
	qr_done           chan bool
	qr_sync_waitgroup sync.WaitGroup
)

func main() {
	//load config
	main_config = load_config()

	// log to file
	log_file, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)
	log.Println("Logging Started")

	//open_db
	// qc_db, err := sql.Open("sqlite3", DB_FILE)
	qc_db, err := sql.Open("sqlite3", ":memory:")
	qc_db.Exec("attach ? as 'bs'", DB_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer qc_db.Close()

	dbinit(qc_db)

	//setup print goroutine
	print_queue = make(chan string, 4)
	defer close(print_queue)
	go do_print_queue(print_queue)

	//setup qr goroutine
	defer qr_sync_waitgroup.Wait()
	qr_done = make(chan bool)
	// defer close(qr_done)

	//show main window
	show_window()

	close(qr_done)
}

func pdf_print(pdf_path string) error {

	app := "./PDFtoPrinter"
	cmd := exec.Command(app, pdf_path)
	err = cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()

	return err

}

func do_print_queue(print_queue chan string) {
	for {
		select {
		case new_file, ok := <-print_queue:
			if ok {
				log.Println("printing: ", new_file)
				err := pdf_print(new_file)
				if err != nil {
					log.Println(err)
				}
			} else {
				return
			}
		}
	}
}

type QRJson struct {
	Product_type string `json:"product_name"`
	Lot_number   string `json:"lot_number"`
}

func do_read_qr(
	qr_pop_data func(QRJson),
	qr_done chan bool) {

	defer qr_sync_waitgroup.Done()

	var webcam_waitgroup sync.WaitGroup
	defer webcam_waitgroup.Wait()

	var product QRJson

	webcam_text := make(chan string)

	webcam_done := make(chan bool)
	defer close(webcam_done)

	webcam_waitgroup.Add(1)
	go DoReadFromWebcam(&webcam_waitgroup, webcam_text, qr_done)

	for {
		select {
		case qr_json, ok := <-webcam_text:

			if ok {
				log.Println("ReadFromWebcam: ", qr_json)
				err := json.Unmarshal([]byte(qr_json), &product)
				if err == nil {
					qr_pop_data(product)
				}
			} else {
				log.Printf("do_read_qr qr_get exit qr_done: \n")
				return
			}
		case <-qr_done:
			return
		}
	}
}

func load_config() *viper.Viper {
	viper_config := viper.New()
	viper_config.SetConfigName("config") // name of config file (without extension)
	viper_config.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper_config.AddConfigPath(".")      // optionally look for config in the working directory
	// viper_config.AddConfigPath("/etc/appname/")  // path to look for the config file in
	// viper.AddConfigPath("$HOME/.config/qc_data_entry") // call multiple times to add many search paths
	viper_config.AddConfigPath("$HOME/.config/qc_data_entry") // call multiple times to add many search paths
	err := viper_config.ReadInConfig()                        // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error if desired

		viper_config.Set("db_path", ".")
		viper_config.Set("label_path", ".")
		viper_config.Set("log_file", "./qc_data_entry.log")
		viper_config.Set("json_paths", []string{"."})

		log.Println(viper_config.WriteConfigAs("config.toml"))
	} else if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	DB_FILE = viper_config.GetString("db_path") + "/qc.sqlite3"
	LABEL_PATH = viper_config.GetString("label_path")
	JSON_PATHS = viper_config.GetStringSlice("json_paths")
	LOG_FILE = viper_config.GetString("log_file")

	return viper_config
}
