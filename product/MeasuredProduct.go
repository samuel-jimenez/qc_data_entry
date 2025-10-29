package product

import (
	"log"
	"strings"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/config"
	"github.com/samuel-jimenez/qc_data_entry/nullable"
	"github.com/samuel-jimenez/qc_data_entry/threads"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type MeasuredProduct struct {
	QCProduct
	PH          nullable.NullFloat64
	SG          nullable.NullFloat64
	Density     nullable.NullFloat64
	String_test nullable.NullInt64
	Viscosity   nullable.NullInt64
}

func MeasuredProduct_from_new() *MeasuredProduct {
	return new(MeasuredProduct)
}

func (measured_product MeasuredProduct) Save() int64 {

	proc_name := "MeasuredProduct-Save.Sample_point"
	DB.Insert(proc_name, DB.DB_Insel_sample_point, measured_product.Sample_point)

	// ??TODO compare
	proc_name = "MeasuredProduct-Save.Tester"
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
	proc_name = "MeasuredProduct-Save.All"
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
	proc_name = "MeasuredProduct-Save.Lot-status"

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

func Store(products ...MeasuredProduct) {
	numSamples := len(products)
	if numSamples <= 0 {
		return
	}
	measured_product := products[0]
	proc_name := "MeasuredProduct-Store"
	qc_sample_storage_id := measured_product.GetStorage(numSamples)

	for _, product := range products {
		qc_id := product.Save()
		// assign storage to sample
		log.Println("DEBUG: ", proc_name, qc_id, qc_sample_storage_id)

		DB.Update(proc_name,
			DB.DB_Update_qc_samples_storage,
			qc_id, qc_sample_storage_id)

	}

	// update capacity
	DB.Update(proc_name,
		DB.DB_Update_dec_product_sample_storage_capacity,
		qc_sample_storage_id, numSamples)

}

func (measured_product MeasuredProduct) Check_data() bool {
	return measured_product.Valid
}

func (measured_product MeasuredProduct) Printout() error {
	return measured_product.print()
}

func (measured_product MeasuredProduct) Output() error {
	if err := measured_product.Printout(); err != nil {
		log.Printf("Error: [%s]: %q\n", "Output", err)
		log.Printf("Debug: %q: %v\n", err, measured_product)
		return err
	}
	measured_product.format_sample()
	return measured_product.export_CoA()

}

func (measured_product MeasuredProduct) Output_sample() error {
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

func (measured_product MeasuredProduct) Save_xl() error {
	isomeric_product := strings.Split(measured_product.Product_name, " ")[1]
	valence_product := strings.Split(measured_product.Product_name_customer, " ")[1]
	return updateExcel(config.RETAIN_FILE_NAME, config.RETAIN_WORKSHEET_NAME, measured_product.Lot_number, isomeric_product, valence_product)
}

func (measured_product MeasuredProduct) print() error {

	pdf_path, err := measured_product.export_label_pdf()
	if err != nil {
		return err
	}
	threads.Show_status("Label Created")

	Print_PDF(pdf_path)

	return err
}
