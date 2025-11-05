package viewer

import (
	"database/sql"
	"log"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/DB"
	"github.com/samuel-jimenez/qc_data_entry/util"
)

var QC_DB *sql.DB

/*
 * ??TODO
func _select_samples(proc_name string, select_statement *sql.Stmt, args ...any) []QCData {
	data := make([]QCData, 0)
	DB.Forall_exit(proc_name,
		func() {},
		func(row *sql.Rows) error {
			var (
				qc_data              QCData
				_timestamp           int64
				product_moniker_name string
				internal_name        string
			)

			if err := row.Scan(
				&product_moniker_name, &internal_name,
				&qc_data.Product_name_customer,
				&qc_data.Lot_name,
				&qc_data.Sample_point,
				&qc_data.Sample_bin,
				&_timestamp,
				&qc_data.PH,
				&qc_data.Specific_gravity,
				&qc_data.String_test,
				&qc_data.Viscosity,
			); err != nil {
				return err
			}
			qc_data.Product_name = product_moniker_name + " " + internal_name

			qc_data.Time_stamp = time.Unix(0, _timestamp)
			data = append(data, qc_data)
			return nil

		},
		select_statement, args...)
	return data
}

func select_samples() []QCData {
	return _select_samples("select_samples", DB_Select_samples)
}

func select_product_samples(product_id int) []QCData {
	return _select_samples("select_product_samples", DB_Select_product_samples, product_id)
}*/

func __select_samples(proc_name string, rows *sql.Rows, err error, query string) []QCData {
	if err != nil {
		log.Printf("error: [%s]: %q\n%v\n", proc_name, err, query)
		// return -1
	}

	data_0 := make([]QCData, 0)
	data_1 := make([]QCData, 0)
	for rows.Next() {
		var (
			qc_data              QCData
			_timestamp           int64
			product_moniker_name string
			internal_name        string
			internalp            bool
		)

		if err := rows.Scan(
			&internalp,
			&product_moniker_name, &internal_name,
			&qc_data.Product_name_customer,
			&qc_data.Lot_name,
			&qc_data.Sample_point,
			&qc_data.Sample_bin,
			&_timestamp,
			&qc_data.PH,
			&qc_data.Specific_gravity,
			&qc_data.Density,
			&qc_data.String_test,
			&qc_data.Viscosity,
		); err != nil {
			log.Fatalf("Crit: [%s]: %v", proc_name, err)
		}
		if internalp {
			qc_data.Product_name = product_moniker_name + " " + internal_name
		} else {
			qc_data.Product_name = internal_name
		}
		qc_data.Time_stamp = time.Unix(0, _timestamp)
		data_0 = append(data_0, qc_data)
	}

	// DB_Select_product_lot_list_sources
	for _, val := range data_0 {
		val.GetComponents()
		data_1 = append(data_1, val)
	}

	return data_1
}

func _select_samples(proc_name string, query string, args ...any) []QCData {
	rows, err := QC_DB.Query(query, args...)
	return __select_samples(proc_name, rows, err, query)
}

func select_all_samples() []QCData {
	rows, err := DB_Select_samples.Query()
	return __select_samples("select_all_samples", rows, err, "")
}

func select_samples(query string) []QCData {
	return _select_samples("select_samples", util.Concat(SAMPLE_SELECT_STRING, query, SAMPLE_ORDER_STRING))
}

func select_lot(fn func(int, string), query string) {
	proc_name := "select_lot"

	rows, err := QC_DB.Query(util.Concat(LOT_SELECT_STRING, query, LOT_ORDER_STRING))
	if err != nil {
		log.Printf("error: [%s]: %q\n", proc_name, err)
		return
	}
	for rows.Next() {
		var (
			id   int
			name string
		)
		if err := rows.Scan(
			&id, &name,
		); err != nil {
			log.Printf("error: [%s]: %q\n", proc_name, err)
		}
		fn(id, name)
	}
}

var (
	SAMPLE_SELECT_STRING, SAMPLE_ORDER_STRING string
	LOT_SELECT_STRING, LOT_ORDER_STRING       string
	DB_Select_samples                         *sql.Stmt
)

func DBinit(db *sql.DB) {
	DB.Check_db(db, true)
	DB.DBinit(db)

	SAMPLE_SELECT_STRING = `
	with lots as (
		select
			true internal, product_lot_id, lot_id, product_id, product_name_internal,product_customer_id, product_moniker_id
	from bs.product_lot
		join bs.product_line using (product_id)
union
		select
			false, inbound_lot_id, lot_id, inbound_product_id+1000,inbound_product_name, null product_customer_id, 8 product_moniker_id
		from bs.inbound_relabel
		join bs.inbound_lot using (inbound_lot_id)
		join bs.inbound_product using (inbound_product_id)
	)


	select
		internal,
		product_moniker_name,
		product_name_internal,
		product_name_customer,
		lot_name,
		sample_point,
		qc_sample_storage_name,
		time_stamp,
		ph ,
		specific_gravity ,
		density ,
		string_test ,
		viscosity
	from bs.qc_samples
		join bs.lot_list using (lot_id)
		join lots using (lot_id)
		join bs.product_moniker using (product_moniker_id)
		left join bs.product_customer_line using (product_customer_id, product_id)
		left join bs.product_sample_points using (sample_point_id)
		left join bs.qc_sample_storage_list using (qc_sample_storage_id)
	`
	SAMPLE_ORDER_STRING = `
	order by time_stamp desc
	`
	LOT_SELECT_STRING = `
	select
	product_lot_id, lot_name
	from (
		select
			product_lot_id, lot_id, product_id, product_customer_id, product_moniker_id
		from bs.product_lot
		join bs.product_line using (product_id)

		union
		select
			inbound_lot_id, lot_id, inbound_product_id, null product_customer_id, 8 product_moniker_id
		from bs.inbound_relabel
		join bs.inbound_lot using (inbound_lot_id)
		join bs.inbound_product using (inbound_product_id)

	)
	join bs.lot_list using (lot_id)
	join bs.product_moniker using (product_moniker_id)
	`
	LOT_ORDER_STRING = `
	order by lot_name desc
	`

	DB_Select_samples = DB.PrepareOrElse(db, SAMPLE_SELECT_STRING+SAMPLE_ORDER_STRING)
}
