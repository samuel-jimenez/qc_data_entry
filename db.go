package main

import (
	"database/sql"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var DB_PATH = "C:/Users/QC Lab/Documents/golang/qc_data_entry/qc.sqlite3"

var qc_db *sql.DB
var db_get_product, db_insert_product,
	db_get_lot, db_insert_lot *sql.Stmt
var err error

func dbinit(db *sql.DB) {

	// 	sqlStmt := `
	// PRAGMA foreign_keys = ON;
	// create schema bs;
	// create table bs.product_line (product_id integer not null primary key, product_name text unique);
	// create table bs.product_lot (lot_id integer not null primary key, lot_name text, product_id references product_line, unique (lot_name,product_id));
	// create table bs.qc_samples (qc_id integer not null primary key, lot_id references product_lot, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real);
	// `
	sqlStmt := `
PRAGMA foreign_keys = ON;
create table bs.product_line (
	product_id integer not null,
	product_name text unique,
	primary key (product_id));

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

`

	// 	sqlStmt := `
	// PRAGMA foreign_keys = ON;
	// create table product_lot (lot_id integer not null primary key, lot_name text, product_id references product);
	// create table qc_samples (qc_id integer not null primary key, lot_id references product_lot, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real);
	// `

	/*
	   	sqlStmt := `
	   drop table qc_samples
	   create table qc_samples (qc_id integer not null primary key, batch_id references product_batch, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real, );
	   `*/
	// foreign key(trackartist) references product(product_id)
	// references product
	// recipe
	//qc_values

	// db.Exec(sqlStmt)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		// return
	}

	get_product_statement := `select product_id from bs.product_line where product_name = ?`
	db_get_product, err = db.Prepare(get_product_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, get_product_statement)
		return
	}

	insert_product_statement := `insert into bs.product_line (product_name) values (?) returning product_id`
	db_insert_product, err = db.Prepare(insert_product_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, insert_product_statement)
		return
	}

	get_lot_statement := `select lot_id from bs.product_lot join product_line using (product_id) where lot_name = ? and product_name = ?`
	db_get_lot, err = db.Prepare(get_lot_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, get_lot_statement)
		return
	}

	// insert_lot_statement := `insert into bs.product_lot (lot_name,product_id) values (?,?) returning lot_id`
	insert_lot_statement := `insert into bs.product_lot (lot_name,product_id)
	select ? as lot_name, product_id from bs.product_line where product_name = ?
	returning lot_id`
	db_insert_lot, err = db.Prepare(insert_lot_statement)
	if err != nil {
		log.Printf("%q: %s\n", err, insert_lot_statement)
		return
	}
}


func get_init_product_id(product_name string) int64 {
	// product_id, err := db_get_product.Exec(product_name)
	// product_id := db_get_product.QueryRow(product_name)
	var product_id int64
	if db_get_product.QueryRow(product_name).Scan(&product_id) != nil {
		//no rows
		result, err := db_insert_product.Exec(product_name)
		if err != nil {
			log.Printf("%q: %s\n", err, "get_init_product_id")
			return -1
		}
		product_id, err = result.LastInsertId()
		if err != nil {
			log.Printf("%q: %s\n", err, "get_init_product_id")
			return -2
		}
	}
	return product_id
}

func get_init_lot_id(lot_name string, product_name string) int64 {
	// lot_id, err := db_get_lot.Exec(lot_name)
	// lot_id := db_get_lot.QueryRow(lot_name)
	var lot_id int64
	if db_get_lot.QueryRow(lot_name, product_name).Scan(&lot_id) != nil {
		//no rows
		result, err := db_insert_lot.Exec(lot_name, product_name)
		if err != nil {
			log.Printf("%q: %s\n", err, "get_init_lot_id")
			return -1
		}
		lot_id, err = result.LastInsertId()
		if err != nil {
			log.Printf("%q: %s\n", err, "get_init_lot_id")
			return -2
		}
	}
	return lot_id
}
