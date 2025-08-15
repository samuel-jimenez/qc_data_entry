package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/threads"
	"github.com/samuel-jimenez/whatsupdocx"
	"github.com/samuel-jimenez/whatsupdocx/docx"
	"github.com/xuri/excelize/v2"
)

var (
	COA_LOT_TITLE = "Batch/Lot#"
)

type Product struct {
	BaseProduct
	PH          nullable.NullFloat64
	SG          nullable.NullFloat64
	Density     nullable.NullFloat64
	String_test nullable.NullInt64
	Viscosity   nullable.NullInt64
}

func (product Product) Save() int64 {

	// DB.DB_insert_sample_point.Exec(product.Sample_point)
	proc_name := "Product.Save.Sample_point"
	log.Println("DEBUG: Product.Save.Sample_point:", DB.Insert(proc_name, DB.DB_insert_sample_point, product.Sample_point))

	// DB.DB_insert_qc_tester.Exec(product.Tester)
	proc_name = "Product.Save.Tester"
	DB.Forall_err(proc_name,
		func() {},
		func(row *sql.Rows) error {
			var (
				id   int64
				name string
			)

			if err := row.Scan(
				&id, &name,
			); err != nil {
				return err
			}

			log.Println("DEBUG: ", proc_name, id, name)

			return nil
		},
		DB.DB_insert_qc_tester, product.Tester)

	var qc_id int64
	proc_name = "Product.Save.All"
	if err := DB.Select_Error(proc_name,
		DB.DB_insert_measurement.QueryRow(
			product.Lot_id, product.Sample_point, product.Tester, time.Now().UTC().UnixNano(), product.PH, product.SG, product.String_test, product.Viscosity,
		),
		&qc_id,
	); err != nil {
		log.Println("error[]%S]:", proc_name, err)
		threads.Show_status("Sample Recording Failed")
		return DB.INVALID_ID

	}
	threads.Show_status("Sample Recorded")
	return qc_id
}

// updateSstaorgr
/*
   create table bs.product_moniker (
   	product_moniker_id integer not null,
   	product_moniker_name text unique not null,
   primary key (product_moniker_id)
   );*/

// 		create table bs.product_sample_storage (
//
// product_sample_storage_id integer not null,
// product_moniker_id not null,
// retain_storage_duration integer not null,
// max_storage_capacity integer not null,
//
// qc_sample_storage_id not null,
// qc_sample_storage_offset integer not null,
// qc_storage_capacity integer not null,

// create table bs.qc_sample_storage_list (
// 	qc_sample_storage_id integer not null,
// 	qc_sample_storage_name text not null,
// 	product_moniker_id not null,
//
// foreign key (product_moniker_id) references product_moniker,
// unique (qc_sample_storage_name),
// primary key (qc_sample_storage_id)
// );
//
// reate table bs.qc_samples (
// 	qc_id integer not null,
// 	lot_id integer not null,
// 	sample_point_id integer,
// 	qc_tester_id integer,
// 	qc_sample_storage_id integer,
// 	time_stamp integer,
// 	ph real,
// 	specific_gravity real,
// 	string_test real,
// 	viscosity real,
// foreign key (lot_id) references lot_list,
// foreign key (sample_point_id) references product_sample_points,
// foreign key (qc_sample_storage_id) references qc_sample_storage_list,
// foreign key (qc_tester_id) references qc_tester_list,
// primary key (qc_id)
// );

/*
 * GetStorage
 *
 * Check storage
 *
 * Returns next storage bin with capacity for product
 * products come in sequentially
 * keep entire batches together
 *
 */
func (product Product) GetStorage(numSamples int) int {

	proc_name := "Product.GetStorage"
	log.Println("DFEBUG: ", proc_name, numSamples, product)

	var (
		qc_sample_storage_id int
	)
	storage_capacity := 0
	if err := DB.Select_Error(proc_name,
		DB.DB_Select_product_sample_storage_capacity.QueryRow(product.Product_id),
		&qc_sample_storage_id, &storage_capacity,
	); err != nil {
		log.Println("Crit: ", proc_name, product, err)
		return int(DB.INVALID_ID)
	}

	log.Println("DFEBUG: ", proc_name, "strogaE", storage_capacity, numSamples, qc_sample_storage_id, product.Product_id)

	// Check storage
	if storage_capacity < numSamples {
		// get data for printing bin label, new bin creation
		proc_name = "Product.GetStorage.New"

		var (
			product_sample_storage_id, qc_sample_storage_offset, start_date_, end_date_ int64
			retain_storage_duration                                                     int
			qc_sample_storage_name, product_moniker_name                                string
		)
		if err := DB.Select_Error(proc_name,
			DB.DB_Select_gen_product_sample_storage.QueryRow(product.Product_id),
			&product_sample_storage_id, &product_moniker_name, &qc_sample_storage_offset, &qc_sample_storage_name,
			&start_date_, &end_date_, &retain_storage_duration,
		); err != nil {
			log.Println("Crit: ", proc_name, product, err)
			return int(DB.INVALID_ID)
		}
		start_date := time.Unix(0, start_date_)
		end_date := time.Unix(0, end_date_)
		retain_date := end_date.AddDate(0, retain_storage_duration, 0)
		// + retain_storage_duration*time.Month

		// print
		product.PrintStorage(qc_sample_storage_name, product_moniker_name, start_date, end_date, retain_date)

		// gen new based on qc_sample_storage_offset,
		// product_moniker_id or product_sample_storage_id
		proc_name = "Product.GetStorage.Insert"
		location := 0
		// product_moniker_id == product_sample_storage_id. update if this is no longer true
		product_moniker_id := product_sample_storage_id
		qc_sample_storage_offset += 1
		qc_sample_storage_name = fmt.Sprintf("%02d-%03d-%04d", location, product_moniker_id, qc_sample_storage_offset)

		if err := DB.Select_Error(proc_name,
			DB.DB_Insert_sample_storage.QueryRow(
				qc_sample_storage_name, product_moniker_id,
			),
			&qc_sample_storage_id,
		); err != nil {
			log.Println("Crit: ", proc_name, product, err)
			return int(DB.INVALID_ID)
		}

		//update qc_sample_storage_offset
		proc_name = "Product.GetStorage.Update"
		DB.Update(proc_name,
			DB.DB_Update_product_sample_storage_qc_sample,
			product_sample_storage_id,
			qc_sample_storage_id, qc_sample_storage_offset)

	}
	return qc_sample_storage_id

}

func Store(products ...Product) {
	numSamples := len(products)
	if numSamples <= 0 {
		return
	}
	product := products[0]
	proc_name := "Product.Store"
	log.Println("DFEBUG: ", proc_name)
	qc_sample_storage_id := product.GetStorage(numSamples)
	log.Println("DFEBUG: ", proc_name, numSamples, qc_sample_storage_id)

	for _, product := range products {
		qc_id := product.Save()
		// assign storage to sample
		log.Println("DFEBUG: ", proc_name, qc_id, qc_sample_storage_id)

		DB.Update(proc_name,
			DB.DB_Update_qc_samples_storage,
			qc_id, qc_sample_storage_id)

	}

	// update capacity
	DB.Update(proc_name,
		DB.DB_Update_product_sample_storage_capacity,
		qc_sample_storage_id, numSamples)

	// * Check storage
	product.GetStorage(0)

}

func (product Product) Export_json() {
	output_files := product.get_json_names()
	bytestring, err := json.MarshalIndent(product, "", "\t")

	if err != nil {
		log.Println("error [Product.Export_json]:", err)
	}

	for _, output_file := range output_files {
		if err := os.WriteFile(output_file, bytestring, 0666); err != nil {
			log.Fatal("Crit: Product.Export_json: ", err)
		}
	}
}

func (product Product) get_coa_template() string {
	//TODO db
	// return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, product.COA_TEMPLATE)
	// Product_type split
	//TODO fix this
	product_moniker := strings.Split(product.Product_name, " ")[0]
	if product_moniker == "PETROFLO" {
		return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA-PETROFLO.docx")
	}
	return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA.docx")
}

func (product Product) get_coa_name() string {
	return fmt.Sprintf("%s/%s", config.COA_FILEPATH, product.get_base_filename("docx"))
}

func (product Product) export_CoA() error {
	var (
		p_title   = "[PRODUCT_NAME]"
		Coa_title = "Parameter"
	)
	terms := []string{
		COA_LOT_TITLE,
	}

	template_file := product.get_coa_template()
	output_file := product.get_coa_name()

	doc, err := whatsupdocx.OpenDocument(template_file)
	if err != nil {
		return err
	}

	product_name := product.Product_name_customer
	if product_name == "" {
		product_name = product.Product_name
	}
	for _, item := range doc.Document.Body.Children {
		if para := item.Paragraph; para != nil {
			product.searchCOAPara(para, p_title, product_name)
		}

		if table := item.Table; table != nil {
			product.searchCOATable(table, Coa_title, terms)
		}
	}

	// save to file
	return doc.SaveTo(output_file)
}

func (product Product) searchCOAPara(para *docx.Paragraph, p_title, product_name string) {
	if strings.Contains(para.String(), p_title) {
		for _, child := range para.Children {
			if run := child.Run; run != nil && strings.Contains(run.String(), p_title) {

				run.Clear()

				//Add product name
				run.AddText(product_name)
				return
			}
		}
	}
}

func (product Product) searchCOATable(table *docx.Table, Coa_title string, terms []string) {
	for _, row := range table.RowContents {
		if product.searchCOARow(table, row.Row, Coa_title, terms) {
			return
		}
	}
}

func (product Product) searchCOARow(table *docx.Table, row *docx.Row, Coa_title string, terms []string) (finished bool) {
	currentHeading := ""
ROW:
	for i, cell := range row.Contents {
		for _, cont := range cell.Cell.Contents {
			field := cont.Paragraph
			if field == nil {
				continue
			}
			switch i {
			case 0:
				if currentHeading == "" {
					field_string := field.String()
					if strings.Contains(field_string, Coa_title) {
						QCProduct{Product: product}.write_CoA_rows(table)
						return true
					}
					for _, term := range terms {
						if strings.Contains(field_string, term) {
							currentHeading = term
							continue ROW
						}
					}
				}
			case 1:
				switch currentHeading {
				case COA_LOT_TITLE:
					field.AddText(product.Lot_number)
					return false
				}
			}
		}
	}
	return false
}

func (product Product) toProduct() Product {
	return product
}

func (product Product) Check_data() bool {
	return product.Valid
}

func (product Product) Printout() error {
	//TODO test
	product.Export_json()
	return product.print()
}

func (product Product) Output() error {
	if err := product.Printout(); err != nil {
		log.Printf("Error: [%s]: %q\n", "Output", err)
		log.Printf("Debug: %q: %v\n", err, product)
		return err
	}
	product.format_sample()
	return product.export_CoA()

}

func (product *Product) format_sample() {
	if product.Product_name_customer != "" {
		product.Product_name = product.Product_name_customer
	}
	product.Sample_point = ""
}

func (product Product) Output_sample() error {
	log.Println("DEBUG: Output_sample ", product)
	product.format_sample()
	log.Println("DEBUG: Output_sample formatted", product)
	if err := product.export_CoA(); err != nil {
		log.Printf("Error: [%s]: %q\n", "Output_sample", err)
		log.Printf("Debug: %q: %v\n", err, product)
		return err
	}
	return product.print()
	//TODO clean
	// return nil

}

// TODO: move these to own file
func withOpenFile(file_name string, FN func(*excelize.File) error) error {
	xl_file, err := excelize.OpenFile(file_name)
	if err != nil {
		log.Printf("Error: [%s]: %q\n", "withOpenFile", err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := xl_file.Close(); err != nil {
			log.Printf("Error: [%s]: %q\n", "withOpenFile", err)
		}
	}()
	return FN(xl_file)
}

func updateExcel(file_name, worksheet_name string, row ...string) error {

	return withOpenFile(file_name, func(xl_file *excelize.File) error {
		// Get all the rows in the worksheet.
		rows, err := xl_file.GetRows(worksheet_name)
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "updateExcel", err)
			return err
		}
		startCell := fmt.Sprintf("A%v", len(rows)+1)

		err = xl_file.SetSheetRow(worksheet_name, startCell, &row)
		if err != nil {
			log.Printf("Error: [%s]: %q\n", "updateExcel", err)
			return err
		}
		return xl_file.Save()
	})
}

func (product Product) Save_xl() error {
	isomeric_product := strings.Split(product.Product_name, " ")[1]
	valence_product := strings.Split(product.Product_name_customer, " ")[1]
	return updateExcel(config.RETAIN_FILE_NAME, config.RETAIN_WORKSHEET_NAME, product.Lot_number, isomeric_product, valence_product)
}

func (product Product) print() error {

	pdf_path, err := product.export_label_pdf()
	if err != nil {
		return err
	}
	threads.Show_status("Label Created")

	_print(pdf_path)

	return err
}

func (product Product) Reprint() {
	_print(product.get_pdf_name())
}

func (product Product) Reprint_sample() {
	product.format_sample()
	log.Println("DEBUG: Reprint_sample formatted", product)
	product.Reprint()
}
