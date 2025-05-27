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
	db_insert_appearance,
	db_select_product_details, db_upsert_product_details,
	db_select_product_coa_details, db_upsert_product_coa_details,
	db_select_product_customer_id, db_upsert_product_customer, db_select_product_customer_info,
	db_select_lot_id, db_insert_lot, db_select_lot_info,
	db_insert_sample_point,
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

// ON UPDATE CASCADE
//        ON DELETE CASCADE

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


create table bs.product_sample_points (
	sample_point_id integer not null,
	sample_point text,
	primary key (sample_point_id),
	unique(sample_point));

create table bs.qc_samples (
	qc_id integer not null,
	lot_id integer not null,
	sample_point_id integer not null,
	time_stamp integer,
	ph real,
	specific_gravity real,
	string_test real,
	viscosity real,
	primary key (qc_id),
	foreign key (lot_id) references product_lot),
	foreign key (sample_point_id) references product_sample_points;



create table bs.product_types (
	product_type_id integer not null,
	product_type_name text,
	primary key (product_type_id),
	unique(product_type_name));


create table bs.product_appearance (
	product_appearance_id integer not null,
	product_appearance_text text,
	primary key (product_appearance_id),
	unique(product_appearance_text));

create table bs.product_ranges_measured (
	range_id integer not null,
	product_id not null,
	product_type_id integer not null,
	product_appearance_id integer not null,
	ph_min real,
	ph_target real,
	ph_max real,
	specific_gravity_min real,
	specific_gravity_target real,
	specific_gravity_max real,
	density_min real,
	density_target real,
	density_max real,
	string_test_min real,
	string_test_target real,
	string_test_max real,
	viscosity_min real,
	viscosity_target real,
	viscosity_max real,
	primary key (range_id),
	foreign key (product_id) references product_line,
	foreign key (product_type_id) references product_types,
	foreign key (product_appearance_id) references product_appearance,
	unique (product_id));

create table bs.product_ranges_published (
	qc_range_id integer not null,
	product_id not null,
	product_appearance_id integer not null,
	ph_min real,
	ph_target real,
	ph_max real,
	specific_gravity_min real,
	specific_gravity_target real,
	specific_gravity_max real,
	density_min real,
	density_target real,
	density_max real,
	string_test_min real,
	string_test_target real,
	string_test_max real,
	viscosity_min real,
	viscosity_target real,
	viscosity_max real,
	primary key (qc_range_id),
	foreign key (product_id) references product_line,
	foreign key (visual_id) references product_appearance,
	unique (product_id));



`
	// TODO

	// _min real,
	// _max real,
	// _target real,

	db.Exec(sqlStmt)
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	// return
	// }

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

	db_insert_appearance = PrepareOrElse(db, `
	with val (product_appearance_text) as (
		values
			(?)
		),
		sel as (
			select product_appearance_text, product_appearance_id
			from val
			left join bs.product_appearance using (product_appearance_text)
		)
	insert into bs.product_appearance (product_appearance_text)
	select distinct product_appearance_text from sel where product_appearance_id is null and product_appearance_text is not null
	returning product_appearance_id, product_appearance_text
	`)

	db_select_product_details = PrepareOrElse(db, `
	with
		measured as (
			select
				product_id,
				product_type_id,
				product_appearance_text,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from bs.product_ranges_measured
			join bs.product_appearance using (product_appearance_id)),

		published as (
			select
				product_id,
				product_appearance_text,

				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max

			from bs.product_ranges_published
			join bs.product_appearance using (product_appearance_id))
	select
		product_name_customer,
		product_type_id,
		coalesce(measured.product_appearance_text, published.product_appearance_text) as product_appearance_text,

		coalesce(measured.ph_min, published.ph_min) as ph_min,
		coalesce(measured.ph_target, published.ph_target) as ph_target,
		coalesce(measured.ph_max, published.ph_max) as ph_max,

		coalesce(measured.specific_gravity_min, published.specific_gravity_min) as specific_gravity_min,
		coalesce(measured.specific_gravity_target, published.specific_gravity_target) as specific_gravity_target,
		coalesce(measured.specific_gravity_max, published.specific_gravity_max) as specific_gravity_max,

		coalesce(measured.density_min, published.density_min) as density_min,
		coalesce(measured.density_target, published.density_target) as density_target,
		coalesce(measured.density_max, published.density_max) as density_max,

		coalesce(measured.string_test_min, published.string_test_min) as string_test_min,
		coalesce(measured.string_test_target, published.string_test_target) as string_test_target,
		coalesce(measured.string_test_max, published.string_test_max) as string_test_max,

		coalesce(measured.viscosity_min, published.viscosity_min) as viscosity_min,
		coalesce(measured.viscosity_target, published.viscosity_target) as viscosity_target,
		coalesce(measured.viscosity_max, published.viscosity_max) as viscosity_max
	from bs.product_customer_line
		full join measured using (product_id)
		full join published using (product_id)
	where product_id = ?
	`)

	db_select_product_coa_details = PrepareOrElse(db, `
	select
		product_appearance_text,

		ph_min,
		ph_target,
		ph_max,

		specific_gravity_min,
		specific_gravity_target,
		specific_gravity_max,

		density_min,
		density_target,
		density_max,

		string_test_min,
		string_test_target,
		string_test_max,

		viscosity_min,
		viscosity_target,
		viscosity_max

	from bs.product_ranges_published
	join bs.product_appearance using (product_appearance_id)
	where product_id = ?
	`)

	db_upsert_product_details = PrepareOrElse(db, `
	with
		val
			(product_id,
			product_type_id,
			product_appearance_text,

			ph_min,
			ph_target,
			ph_max,

			specific_gravity_min,
			specific_gravity_target,
			specific_gravity_max,

			density_min,
			density_target,
			density_max,

			string_test_min,
			string_test_target,
			string_test_max,

			viscosity_min,
			viscosity_target,
			viscosity_max)
		as (
			values (
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?)
		),
		sel as (
			select
				product_id,
				product_type_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from val
			left join bs.product_appearance using (product_appearance_text)
		)
	insert into bs.product_ranges_measured
		(product_id,
		product_type_id,
		product_appearance_id,


		ph_min,
		ph_target,
		ph_max,

		specific_gravity_min,
		specific_gravity_target,
		specific_gravity_max,

		density_min,
		density_target,
		density_max,

		string_test_min,
		string_test_target,
		string_test_max,

		viscosity_min,
		viscosity_target,
		viscosity_max)
	select
				product_id,
				product_type_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from sel
			where true
	on conflict(product_id) do update set

		product_type_id=excluded.product_type_id,
		product_appearance_id=excluded.product_appearance_id,

		ph_min=excluded.ph_min,
		ph_target=excluded.ph_target,
		ph_max=excluded.ph_max,

		specific_gravity_min=excluded.specific_gravity_min,
		specific_gravity_target=excluded.specific_gravity_target,
		specific_gravity_max=excluded.specific_gravity_max,

		density_min=excluded.density_min,
		density_target=excluded.density_target,
		density_max=excluded.density_max,

		string_test_min=excluded.string_test_min,
		string_test_target=excluded.string_test_target,
		string_test_max=excluded.string_test_max,

		viscosity_min=excluded.viscosity_min,
		viscosity_target=excluded.viscosity_target,
		viscosity_max=excluded.viscosity_max

		returning range_id
		`)
	db_upsert_product_coa_details = PrepareOrElse(db, `
	with
		val
			(product_id,
			product_type_id,
			product_appearance_text,

			ph_min,
			ph_target,
			ph_max,

			specific_gravity_min,
			specific_gravity_target,
			specific_gravity_max,

			density_min,
			density_target,
			density_max,

			string_test_min,
			string_test_target,
			string_test_max,

			viscosity_min,
			viscosity_target,
			viscosity_max)
		as (
			values (
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?,
				?,?,?)
		),
		sel as (
			select
				product_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from val
			left join bs.product_appearance using (product_appearance_text)
		)
	insert into bs.product_ranges_published
		(product_id,
		product_appearance_id,


		ph_min,
		ph_target,
		ph_max,

		specific_gravity_min,
		specific_gravity_target,
		specific_gravity_max,

		density_min,
		density_target,
		density_max,

		string_test_min,
		string_test_target,
		string_test_max,

		viscosity_min,
		viscosity_target,
		viscosity_max)
	select
				product_id,
				product_appearance_id,


				ph_min,
				ph_target,
				ph_max,

				specific_gravity_min,
				specific_gravity_target,
				specific_gravity_max,

				density_min,
				density_target,
				density_max,

				string_test_min,
				string_test_target,
				string_test_max,

				viscosity_min,
				viscosity_target,
				viscosity_max
			from sel
			where true
	on conflict(product_id) do update set

		product_appearance_id=excluded.product_appearance_id,

		ph_min=excluded.ph_min,
		ph_target=excluded.ph_target,
		ph_max=excluded.ph_max,

		specific_gravity_min=excluded.specific_gravity_min,
		specific_gravity_target=excluded.specific_gravity_target,
		specific_gravity_max=excluded.specific_gravity_max,

		density_min=excluded.density_min,
		density_target=excluded.density_target,
		density_max=excluded.density_max,

		string_test_min=excluded.string_test_min,
		string_test_target=excluded.string_test_target,
		string_test_max=excluded.string_test_max,

		viscosity_min=excluded.viscosity_min,
		viscosity_target=excluded.viscosity_target,
		viscosity_max=excluded.viscosity_max

		returning qc_range_id
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

	db_insert_sample_point = PrepareOrElse(db, `
	with val (sample_point) as (
		values
			(?)
		),
		sel as (
			select sample_point, sample_point_id
			from val
			left join bs.product_sample_points using (sample_point)
		)
	insert into bs.product_sample_points (sample_point)
	select distinct sample_point from sel where sample_point_id is null
	returning sample_point_id, sample_point
	`)

	db_insert_measurement = PrepareOrElse(db, `
	with
		val (lot_id, sample_point, time_stamp, ph, specific_gravity, string_test, viscosity) as (
			values
				(?, ?, ?, ?, ?, ?, ?)
		),
		sel as (
			select lot_id, sample_point_id, sample_point, time_stamp, ph, specific_gravity, string_test, viscosity
			from val
			left join bs.product_sample_points using (sample_point)
		)
	insert into bs.qc_samples (lot_id, sample_point_id, time_stamp, ph, specific_gravity, string_test, viscosity)
	select lot_id, sample_point_id, time_stamp, ph, specific_gravity, string_test, viscosity
	from   sel
	returning qc_id;
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
