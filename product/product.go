package product

import (
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

func (measured_product Product) Save() int64 {

	proc_name := "Product.Save.Sample_point"
	DB.Insert(proc_name, DB.DB_Insel_sample_point, measured_product.Sample_point)

	// ??TODO compare
	proc_name = "Product.Save.Tester"
	log.Println("TRACE: ", proc_name)

	// DB.DB_insert_qc_tester.Exec(product.Tester)

	// DB.Forall_err(proc_name,
	// 	func() {},
	// 	func(row *sql.Rows) error {
	// 		var (
	// 			id   int64
	// 			name string
	// 		)
	//
	// 		if err := row.Scan(
	// 			&id, &name,
	// 		); err != nil {
	// 			return err
	// 		}
	//
	// 		log.Println("DEBUG: ", proc_name, id, name)
	//
	// 		return nil
	// 	},
	// 	DB.DB_Insel_qc_tester, measured_product.Tester)

	// var (
	// 	qc_tester_id   int64
	// 	qc_tester_name string
	// )
	//
	// if err := DB.Select_Error(proc_name,
	// 	DB.DB_Insel_qc_tester.QueryRow(
	// 		measured_product.Tester,
	// 	),
	// 	&qc_tester_id, &qc_tester_name,
	// ); err != nil {
	// 	log.Println("WARN: []%S]:", proc_name, err)
	//
	// }
	// log.Println("TRACE: ", proc_name, "Select_Error", qc_tester_id, qc_tester_name)

	//TODO DB.Select_Nullable_Error(

	DB.Insert(proc_name, DB.DB_Insel_qc_tester, measured_product.Tester)

	var qc_id int64
	proc_name = "Product.Save.All"
	if err := DB.Select_Error(proc_name,
		DB.DB_insert_measurement.QueryRow(
			measured_product.Lot_id, measured_product.Sample_point, measured_product.Tester, time.Now().UTC().UnixNano(), measured_product.PH, measured_product.SG, measured_product.String_test, measured_product.Viscosity,
		),
		&qc_id,
	); err != nil {
		log.Println("error[]%S]:", proc_name, err)
		threads.Show_status("Sample Recording Failed")
		return DB.INVALID_ID

	}
	proc_name = "Product.Save.Lot-status"

	if err := DB.Update(proc_name,
		DB.DB_Update_lot_list__status, measured_product.Lot_id, Status_TESTED,
	); err != nil {
		log.Println("error[]%S]:", proc_name, err)
		threads.Show_status("Sample Recording Failed")
		return DB.INVALID_ID

	}
	threads.Show_status("Sample Recorded")
	return qc_id
}

//TODO product.NewBin
// add button to call to account for unlogged samples

/*
 * NewStorageBin
 *
 * Print label for storage bin
 *
 * Returns next storage bin
 *
 */
func (measured_product Product) NewStorageBin() int {
	var (
		qc_sample_storage_id int
	)

	// get data for printing bin label, new bin creation
	proc_name := "Product.NewStorageBin"

	var (
		product_sample_storage_id, qc_sample_storage_offset, start_date_, end_date_ int64
		retain_storage_duration                                                     int
		qc_sample_storage_name, product_moniker_name                                string
	)
	if err := DB.Select_Error(proc_name,
		DB.DB_Select_gen_product_sample_storage.QueryRow(
			measured_product.Product_id,
		),
		&product_sample_storage_id, &product_moniker_name, &qc_sample_storage_offset, &qc_sample_storage_name,
		&start_date_, &end_date_, &retain_storage_duration,
	); err != nil {
		log.Println("Crit: ", proc_name, measured_product, err)
		panic(err)
	}
	start_date := time.Unix(0, start_date_)
	end_date := time.Unix(0, end_date_)
	retain_date := end_date.AddDate(0, retain_storage_duration, 0)
	// + retain_storage_duration*time.Month

	// print
	measured_product.PrintOldStorage(qc_sample_storage_name, product_moniker_name, &start_date, &end_date, &retain_date)

	// gen new based on qc_sample_storage_offset,
	// product_moniker_id or product_sample_storage_id
	proc_name = "Product.NewStorageBin.Insert"
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
		log.Println("Crit: ", proc_name, measured_product, err)
		panic(err)
	}
	measured_product.PrintNewStorage(qc_sample_storage_name, product_moniker_name, &start_date, &end_date, &retain_date)

	//update qc_sample_storage_offset
	proc_name = "Product.NewStorageBin.Update"
	DB.Update(proc_name,
		DB.DB_Update_product_sample_storage_qc_sample,
		product_sample_storage_id,
		qc_sample_storage_id, qc_sample_storage_offset)
	return qc_sample_storage_id
}

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
func (measured_product Product) GetStorage(numSamples int) int {

	proc_name := "Product.GetStorage"

	var (
		qc_sample_storage_id, storage_capacity int
	)
	if err := DB.Select_Error(proc_name,
		DB.DB_Select_product_sample_storage_capacity.QueryRow(
			measured_product.Product_id,
		),
		&qc_sample_storage_id, &storage_capacity,
	); err != nil {
		log.Println("Crit: ", proc_name, measured_product, err)
		panic(err)
	}

	// Check storage
	if storage_capacity < numSamples {
		qc_sample_storage_id = measured_product.NewStorageBin()
	}
	return qc_sample_storage_id

}

/*
 * CheckStorage
 *
 * Check storage bin, print label if full
 *
 */
func (measured_product Product) CheckStorage() {
	measured_product.GetStorage(1)
}

func Store(products ...Product) {
	numSamples := len(products)
	if numSamples <= 0 {
		return
	}
	measured_product := products[0]
	proc_name := "Product.Store"
	log.Println("DFEBUG: ", proc_name)
	qc_sample_storage_id := measured_product.GetStorage(numSamples)
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

}

func (measured_product Product) Export_json() {
	output_files := measured_product.get_json_names()
	bytestring, err := json.MarshalIndent(measured_product, "", "\t")

	if err != nil {
		log.Println("error [Product.Export_json]:", err)
	}

	for _, output_file := range output_files {
		if err := os.WriteFile(output_file, bytestring, 0666); err != nil {
			log.Fatal("Crit: Product.Export_json: ", err)
		}
	}
}

func (measured_product Product) get_coa_template() string {
	//TODO db
	// return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, product.COA_TEMPLATE)
	// Product_type split
	//TODO fix this
	product_moniker := strings.Split(measured_product.Product_name, " ")[0]
	if product_moniker == "PETROFLO" {
		return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA-PETROFLO.docx")
	}
	return fmt.Sprintf("%s/%s", config.COA_TEMPLATE_PATH, "CoA.docx")
}

func (product Product) get_coa_name() string {
	return fmt.Sprintf("%s/%s", config.COA_FILEPATH, product.get_base_filename("docx"))
}

func (measured_product Product) export_CoA() error {
	var (
		p_title   = "[PRODUCT_NAME]"
		Coa_title = "Parameter"
	)
	terms := []string{
		COA_LOT_TITLE,
	}

	template_file := measured_product.get_coa_template()
	output_file := measured_product.get_coa_name()

	doc, err := whatsupdocx.OpenDocument(template_file)
	if err != nil {
		return err
	}

	product_name := measured_product.Product_name_customer
	if product_name == "" {
		product_name = measured_product.Product_name
	}
	for _, item := range doc.Document.Body.Children {
		if para := item.Paragraph; para != nil {
			measured_product.searchCOAPara(para, p_title, product_name)
		}

		if table := item.Table; table != nil {
			measured_product.searchCOATable(table, Coa_title, terms)
		}
	}

	// save to file
	return doc.SaveTo(output_file)
}

func (measured_product Product) searchCOAPara(para *docx.Paragraph, p_title, product_name string) {
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

func (measured_product Product) searchCOATable(table *docx.Table, Coa_title string, terms []string) {
	for _, row := range table.RowContents {
		if measured_product.searchCOARow(table, row.Row, Coa_title, terms) {
			return
		}
	}
}

func (measured_product Product) searchCOARow(table *docx.Table, row *docx.Row, Coa_title string, terms []string) (finished bool) {
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
						QCProduct{Product: measured_product}.write_CoA_rows(table)
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
					field.AddText(measured_product.Lot_number)
					return false
				}
			}
		}
	}
	return false
}

func (measured_product Product) toProduct() Product {
	return measured_product
}

func (measured_product Product) Check_data() bool {
	return measured_product.Valid
}

func (measured_product Product) Printout() error {
	//TODO test
	measured_product.Export_json()
	return measured_product.print()
}

func (measured_product Product) Output() error {
	if err := measured_product.Printout(); err != nil {
		log.Printf("Error: [%s]: %q\n", "Output", err)
		log.Printf("Debug: %q: %v\n", err, measured_product)
		return err
	}
	measured_product.format_sample()
	return measured_product.export_CoA()

}

func (measured_product *Product) format_sample() {
	if measured_product.Product_name_customer != "" {
		measured_product.Product_name = measured_product.Product_name_customer
	}
	measured_product.Sample_point = ""
}

func (measured_product Product) Output_sample() error {
	log.Println("DEBUG: Output_sample ", measured_product)
	measured_product.format_sample()
	log.Println("DEBUG: Output_sample formatted", measured_product)
	if err := measured_product.export_CoA(); err != nil {
		log.Printf("Error: [%s]: %q\n", "Output_sample", err)
		log.Printf("Debug: %q: %v\n", err, measured_product)
		return err
	}
	return measured_product.print()
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

func (measured_product Product) Save_xl() error {
	isomeric_product := strings.Split(measured_product.Product_name, " ")[1]
	valence_product := strings.Split(measured_product.Product_name_customer, " ")[1]
	return updateExcel(config.RETAIN_FILE_NAME, config.RETAIN_WORKSHEET_NAME, measured_product.Lot_number, isomeric_product, valence_product)
}

func (measured_product Product) print() error {

	pdf_path, err := measured_product.export_label_pdf()
	if err != nil {
		return err
	}
	threads.Show_status("Label Created")

	_print(pdf_path)

	return err
}

func (measured_product Product) Reprint() {
	_print(measured_product.get_pdf_name())
}

func (measured_product Product) Reprint_sample() {
	measured_product.format_sample()
	log.Println("DEBUG: Reprint_sample formatted", measured_product)
	measured_product.Reprint()
}
