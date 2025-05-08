package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var qc_db *sql.DB
var db_select_product_id, db_insert_product, db_select_product_info,
	db_select_product_details, db_upsert_product_details,
	db_select_product_customer_id, db_upsert_product_customer, db_select_product_customer_info,
	db_select_lot_id, db_insert_lot, db_select_lot_info,
	db_insert_measurement *sql.Stmt
var err error

func PrepareOrElse(db *sql.DB, sqlStatement string) *sql.Stmt {
	preparedStatement, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		panic(err)

	}
	return preparedStatement
}

func dbinit(db *sql.DB) {

	// 	sqlStmt := `
	// PRAGMA foreign_keys = ON;
	// create schema bs;
	// create table bs.product_line (product_id integer not null primary key, product_name_internal text unique);
	// create table bs.product_lot (lot_id integer not null primary key, lot_name text, product_id references product_line, unique (lot_name,product_id));
	// create table bs.qc_samples (qc_id integer not null primary key, lot_id references product_lot, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real);
	// `
	sqlStmt := `
PRAGMA foreign_keys = ON;

create table bs.product_moniker (
	product_moniker_id integer not null,
	product_moniker_name text unique not null,
	primary key (product_moniker_id));

create table bs.product_line (
	product_id integer not null,
	product_name_internal text unique not null,
	product_moniker_id not null,
	foreign key (product_moniker_id) references product_moniker,
	primary key (product_id));

create table bs.product_customer_line (
	product_customer_id integer not null,
	product_id not null unique,
	product_name_customer text unique,
	foreign key (product_id) references product_line,
	primary key (product_customer_id));


create table bs.product_lot (
	lot_id integer not null,
	lot_name text,
	product_id not null,
	primary key (lot_id),
	foreign key (product_id) references product_line,
	unique (lot_name,product_id));

create table bs.qc_samples (
	qc_id integer not null,
	lot_id integer not null,
	sample_point text,
	time_stamp integer,
	specific_gravity real,
	ph real,
	string_test real,
	viscosity real,
	primary key (qc_id),
	foreign key (lot_id) references product_lot);



create table bs.product_types (
	product_type_id integer not null,
	product_type_name text,
	primary key (product_type_id));

create table bs.product_ranges (
	range_id integer not null,
	product_id not null,
	product_type_id integer not null,
	specific_gravity_min real,
	specific_gravity_max real,
	specific_gravity_target real,
	ph_min real,
	ph_max real,
	ph_target real,
	string_test_min real,
	string_test_max real,
	string_test_target real,
	viscosity_min real,
	viscosity_max real,
	viscosity_target real,
	primary key (range_id),
	foreign key (product_id) references product_line,
	foreign key (product_type_id) references product_types,
	unique (product_id));

`
	// _min real,
	// _max real,
	// _target real,

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		// return
	}
	db_select_product_info = PrepareOrElse(db, `
	select product_id, product_name_internal, product_moniker_name
		from bs.product_line
		join bs.product_moniker using (product_moniker_id)
		order by product_moniker_name,product_name_internal

	`)

	db_select_product_id = PrepareOrElse(db, `
	select product_id
		from bs.product_line
		join bs.product_moniker using (product_moniker_id)
		where product_name_internal = ?
		and product_moniker_name = ?
		`)

	db_insert_product = PrepareOrElse(db, `
	insert into bs.product_line
		(product_name_internal, product_moniker_id)
		select ?, product_moniker_id
			from bs.product_moniker
			where product_moniker_name = ?
		returning product_id
		`)

	db_select_product_details = PrepareOrElse(db, `
	select
		product_name_customer,
		product_type_id,
		specific_gravity_min,
		specific_gravity_max,
		specific_gravity_target,
		ph_min,
		ph_max,
		ph_target,
		string_test_min,
		string_test_max,
		string_test_target,
		viscosity_min,
		viscosity_max,
		viscosity_target
	from bs.product_customer_line
	full join bs.product_ranges using (product_id)
	where product_id = ?`)

	db_upsert_product_details = PrepareOrElse(db, `
	insert into bs.product_ranges
		(product_id,
		product_type_id,

		specific_gravity_min,
		specific_gravity_max,
		specific_gravity_target,

		ph_min,
		ph_max,
		ph_target,

		string_test_min,
		string_test_max,
		string_test_target,

		viscosity_min,
		viscosity_max,
		viscosity_target)
		values (?,?,
	?,?,?,
	?,?,?,
	?,?,?,
	?,?,?)
	on conflict(product_id) do update set
		product_type_id=excluded.product_type_id,
		specific_gravity_min=excluded.specific_gravity_min,
		specific_gravity_max=excluded.specific_gravity_max,
		specific_gravity_target=excluded.specific_gravity_target,

		ph_min=excluded.ph_min,
		ph_max=excluded.ph_max,
		ph_target=excluded.ph_target,

		string_test_min=excluded.string_test_min,
		string_test_max=excluded.string_test_max,
		string_test_target=excluded.string_test_target,

		viscosity_min=excluded.viscosity_min,
		viscosity_max=excluded.viscosity_max,
		viscosity_target=excluded.viscosity_target

		returning range_id
		`)

	db_select_product_customer_info = PrepareOrElse(db, `
	select product_customer_id, product_name_customer
		from bs.product_customer_line
		where product_id = ?
		`)

	db_select_product_customer_id = PrepareOrElse(db, `
	select product_customer_id
		from bs.product_customer_line
		where product_name_customer = ? and product_id = ?
		`)

	db_upsert_product_customer = PrepareOrElse(db, `
	insert into bs.product_customer_line
		(product_name_customer,product_id)
		values (?,?)
	on conflict(product_id) do update set
		product_name_customer=excluded.product_name_customer
		returning product_customer_id
		`)

	db_select_lot_info = PrepareOrElse(db, `
	select lot_id, lot_name
		from bs.product_lot
		where product_id = ?
		`)

	db_select_lot_id = PrepareOrElse(db, `
	select lot_id
		from bs.product_lot
		where lot_name = ? and product_id = ?
		`)

	db_insert_lot = PrepareOrElse(db, `
	insert into bs.product_lot
		(lot_name,product_id)
		values (?,?)
		returning lot_id
		`)

	db_insert_measurement = PrepareOrElse(db, `
	insert into bs.qc_samples
		(lot_id, sample_point, time_stamp, specific_gravity, ph, string_test, viscosity)
		values (?, ?, ?, ?, ?, ?, ?)
		returning qc_id
		`)

}

func select_product_name_customer(product_id int64) string {
	var (
		product_customer_id   int64
		customer_name         string
		product_name_customer string
	)
	product_name_customer = ""

	if db_select_product_customer_info.QueryRow(product_id).Scan(&product_customer_id, &customer_name) == nil {
		product_name_customer = customer_name

	}
	return product_name_customer

}

func DerefOrEmpty[T any](val *T) T {
	if val == nil {
		var empty T
		return empty
	}
	return *val
}

func ValidOr[T any](val *T, default_val T) T {
	if val == nil {
		return default_val
	}
	return *val

}

func OrNil[T comparable](val *T, default_val T) *T {
	if val == nil || *val == default_val {
		return nil
	}
	return val

}

/*
func OrNil[T any](val T, default_val T) *T {
	if val == default_val {
		return nil
	}
	return val

}

func OrNil_[T any](val *T, default_val T) *T {
	if val == nil {
		return nil
	}
	return val

}*/

func insel_product_id(product_name_full string) int64 {
	// product_id, err := db_select_product.Exec(product_name_internal)
	// product_id := db_select_product.QueryRow(product_name_internal)
	var product_id int64

	// v := strings.SplitN(s, " ", 2)

	// before, after, found := strings.Cut(s, sep)
	product_moniker_name, product_name_internal, _ := strings.Cut(product_name_full, " ")

	if db_select_product_id.QueryRow(product_name_internal, product_moniker_name).Scan(&product_id) != nil {
		//no rows
		result, err := db_insert_product.Exec(product_name_internal, product_moniker_name)
		if err != nil {
			log.Printf("%q: %s\n", err, "insel_product_id")
			return -1
		}
		product_id, err = result.LastInsertId()
		if err != nil {
			log.Printf("%q: %s\n", err, "insel_product_id")
			return -2
		}
	}
	return product_id
}

func insel_lot_id(lot_name string, product_id int64) int64 {
	// lot_id, err := db_select_lot.Exec(lot_name)
	// lot_id := db_select_lot.QueryRow(lot_name)
	var lot_id int64
	if db_select_lot_id.QueryRow(lot_name, product_id).Scan(&lot_id) != nil {
		//no rows
		result, err := db_insert_lot.Exec(lot_name, product_id)
		if err != nil {
			log.Printf("%q: %s\n", err, "insel_lot_id")
			return -1
		}
		lot_id, err = result.LastInsertId()
		if err != nil {
			log.Printf("%q: %s\n", err, "insel_lot_id")
			return -2
		}
	}
	return lot_id
}

func upsert_product_name_customer(product_name_customer string, product_id int64) int64 {
	var product_customer_id int64
	result, err := db_upsert_product_customer.Exec(product_name_customer, product_id)
	if err != nil {
		log.Printf("%q: %s\n", err, "insel_lot_id")
		return -1
	}
	product_customer_id, err = result.LastInsertId()
	if err != nil {
		log.Printf("%q: %s\n", err, "insel_lot_id")
		return -2
	}
	return product_customer_id
}
