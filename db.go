package main

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/qc_data_entry/DB"
)

var (
	db_select_product_customer_id, db_select_product_customer_info,
	db_select_lot_info,
	db_select_sample_points *sql.Stmt
	err error

	CONTAINER_TOTE    = 1
	CONTAINER_RAILCAR = 2
)

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

create table bs.database_info (
	database_id integer not null,
	database_version_major integer not null,
	database_version_minor integer not null,
	database_version_revision integer not null,
	check (database_id = 0),
	primary key (database_id));

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
	product_id not null,
	product_name_customer text unique,
	foreign key (product_id) references product_line,
	primary key (product_customer_id));


create table bs.product_lot (
	lot_id integer not null,
	lot_name text,
	product_id not null,
	product_customer_id,
	primary key (lot_id),
	foreign key (product_id) references product_line,
	foreign key (product_customer_id) references product_customer_line,
	unique (lot_name,product_id));




create table bs.product_sample_points (
	sample_point_id integer not null,
	sample_point text,
	primary key (sample_point_id),
	unique(sample_point));

create table bs.qc_samples (
	qc_id integer not null,
	lot_id integer not null,
	sample_point_id integer,
	time_stamp integer,
	ph real,
	specific_gravity real,
	string_test real,
	viscosity real,
	primary key (qc_id),
	foreign key (lot_id) references product_lot),
	foreign key (sample_point_id) references product_sample_points;



create table bs.container_types (
	container_type_id integer not null,
	container_type_name text not null,
	primary key (container_type_id),
	unique(container_type_name));


create table bs.product_types (
	product_type_id integer not null,
	product_type_name text not null,
	container_type_id integer not null,
	primary key (product_type_id),
	foreign key (container_type_id) references container_types,
	unique(product_type_name));



create table bs.product_appearance (
	product_appearance_id integer not null,
	product_appearance_text text not null,
	primary key (product_appearance_id),
	unique(product_appearance_text));

create table bs.product_ranges_measured (
	range_id integer not null,
	product_id not null,
	product_type_id integer not null,
	product_appearance_id integer,
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

insert into bs.container_types
	(container_type_name)
	values ('Tote'),('Railcar');



`
	// TODO add ranges table
	// _min real,
	// _max real,
	// _target real,

	db.Exec(sqlStmt)
	// _, err = db.Exec(sqlStmt)
	if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	// return
	// }

	DB.Check_db(db)
	DB.DBinit(db)



	db_select_product_customer_info = DB.PrepareOrElse(db, `
	select product_customer_id, product_name_customer
		from bs.product_customer_line
		where product_id = ?
		`)

	db_select_product_customer_id = DB.PrepareOrElse(db, `
	select product_customer_id
		from bs.product_customer_line
		where product_name_customer = ? and product_id = ?
	`)

	db_select_lot_info = DB.PrepareOrElse(db, `
	select lot_id, lot_name
		from bs.product_lot
		where product_id = ?
	`)

	db_select_sample_points = DB.PrepareOrElse(db, `
	select distinct sample_point_id, sample_point
		from bs.product_lot
		join bs.qc_samples using (lot_id)
		join bs.product_sample_points using (sample_point_id)
		where product_id = ?
		order by sample_point_id
	`)

}
