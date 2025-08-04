package main

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/samuel-jimenez/qc_data_entry/DB"
)

var (
	CONTAINER_SAMPLE  = 1
	CONTAINER_TOTE    = 2
	CONTAINER_RAILCAR = 3
)

// ON UPDATE CASCADE
//        ON DELETE CASCADE

func dbinit(db *sql.DB) {

	// 	sqlStmt := `
	// PRAGMA foreign_keys = ON;
	// create schema bs;
	// create table bs.product_line (product_id integer not null primary key, product_name_internal text unique);
	// create table bs.product_lot (product_lot_id integer not null primary key, lot_name text, product_id references product_line, unique (lot_name,product_id));
	// create table bs.qc_samples (qc_id integer not null primary key, product_lot_id references product_lot, sample_point text, time_stamp integer, specific_gravity real,  ph real,   string_test real,   viscosity real);
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
	product_lot_id integer not null,
	lot_id not null,
	product_id not null,
	product_customer_id,
	recipe_id,
primary key (product_lot_id),
foreign key (lot_id) references lot_list,
foreign key (product_id) references product_line,
foreign key (product_customer_id) references product_customer_line,
foreign key (recipe_id) references recipe_list,
unique (lot_id,product_id)
);


create table bs.lot_list (
	lot_id integer not null,
	lot_name text not null,
unique (lot_name),
primary key (lot_id));



create table bs.product_sample_points (
	sample_point_id integer not null,
	sample_point text,
	primary key (sample_point_id),
	unique (sample_point)
);

create table bs.qc_samples (
	qc_id integer not null,
	product_lot_id integer not null,
	sample_point_id integer,
	time_stamp integer,
	ph real,
	specific_gravity real,
	string_test real,
	viscosity real,
foreign key (product_lot_id) references product_lot,
foreign key (sample_point_id) references product_sample_points,
primary key (qc_id)
);



create table bs.container_types (
	container_type_id integer not null,
	container_type_name text not null,
unique (container_type_name),
primary key (container_type_id)
);


create table bs.product_types (
	product_type_id integer not null,
	product_type_name text not null,
	container_type_id integer not null,
primary key (product_type_id),
foreign key (container_type_id) references container_types,
unique (product_type_name)
);



create table bs.product_appearance (
	product_appearance_id integer not null,
	product_appearance_text text not null,
unique (product_appearance_text),
primary key (product_appearance_id)
);

create table bs.product_ranges_measured (
	range_id 		integer not null,
	product_id 		not null,
	product_type_id 	integer not null,
	product_appearance_id 	integer,
	ph_min			real,
	ph_target		real,
	ph_max			real,
	specific_gravity_min	real,
	specific_gravity_target	real,
	specific_gravity_max	real,
	density_min		real,
	density_target		real,
	density_max		real,
	string_test_min		real,
	string_test_target	real,
	string_test_max		real,
	viscosity_min		real,
	viscosity_target	real,
	viscosity_max		real,
primary key (range_id),
foreign key (product_id) references product_line,
foreign key (product_type_id) references product_types,
foreign key (product_appearance_id) references product_appearance,
unique (product_id)
);

create table bs.product_ranges_published (
	qc_range_id 		integer not null,
	product_id 		not null,
	product_appearance_id 	integer not null,
	ph_min			real,
	ph_target		real,
	ph_max			real,
	specific_gravity_min	real,
	specific_gravity_target	real,
	specific_gravity_max	real,
	density_min		real,
	density_target		real,
	density_max		real,
	string_test_min		real,
	string_test_target	real,
	string_test_max		real,
	viscosity_min		real,
	viscosity_target	real,
	viscosity_max		real,
primary key (qc_range_id),
foreign key (product_id) references product_line,
foreign key (product_appearance_id) references product_appearance,
unique (product_id)
);




create table bs.container_list (
	container_id integer not null,
	container_name text not null,
	container_type_id not null default 3,
foreign key (container_type_id) references container_types,
unique (container_name),
primary key (container_id));



create table bs.inbound_provider_list (
	inbound_provider_id integer not null,
	inbound_provider_name text not null,
unique (inbound_provider_name),
primary key (inbound_provider_id));




create table bs.inbound_product (
	inbound_product_id integer not null,
	inbound_product_name text,
	unique (inbound_product_name),
	primary key (inbound_product_id));

create table bs.inbound_lot (
	inbound_lot_id integer not null,
	inbound_lot_name text,
	inbound_product_id,
	inbound_provider_id,
	container_id not null,
	status_id not null default 1,
unique (inbound_lot_name),
foreign key (inbound_product_id) references inbound_product,
foreign key (inbound_provider_id) references inbound_provider_list,
foreign key (container_id) references container_list,
foreign key (status_id) references status_list,
primary key (inbound_lot_id));



create table bs.status_list (
	status_id integer not null,
	status_name text not null,
	primary key (status_id),
	unique (status_name));


create table bs.inbound_relabel (
	inbound_relabel_id integer not null,
	lot_id not null,
	inbound_lot_id not null,
	container_id not null,
foreign key (lot_id) references lot_list,
foreign key (inbound_lot_id) references inbound_lot,
foreign key (container_id) references container_list,
primary key (inbound_relabel_id),
unique (lot_id)
);






create table bs.component_types (
	component_type_id integer not null,
	component_type_name text not null,
	primary key (component_type_id),
	unique (component_type_name));





create table bs.component_type_product_internal (
	component_type_product_internal_id integer not null,
	component_type_id not null,
	product_id not null,
	foreign key (component_type_id) references component_types,
	foreign key (product_id) references product_line,
	unique (component_type_id,product_id),
	primary key (component_type_product_internal_id));


create table bs.component_type_product_inbound (
	component_type_product_inbound_id integer not null,
	component_type_id not null,
	inbound_product_id not null,
	foreign key (component_type_id) references component_types,
	foreign key (inbound_product_id) references inbound_product,
	unique (component_type_id,inbound_product_id),
	primary key (component_type_product_inbound_id));







create table bs.component_list (
	component_id integer not null,
	component_type_id not null,
	inbound_lot_id,
	product_lot_id,
unique (component_type_id,inbound_lot_id),
unique (component_type_id,product_lot_id),
foreign key (component_type_id) references component_types,
foreign key (inbound_lot_id) references inbound_lot,
foreign key (product_lot_id) references product_lot,
primary key (component_id));




create table bs.recipe_list (
	recipe_id integer not null,
	product_id not null,
	foreign key (product_id) references product_line,
	primary key (recipe_id));

create table bs.recipe_components (
	recipe_components_id integer not null,
	recipe_id integer not null,
	component_type_id not null,
	component_type_amount real,
	component_add_order not null,
	unique (recipe_id,component_add_order),
	foreign key (recipe_id) references recipe_list,
	foreign key (component_type_id) references component_types,
	primary key (recipe_components_id));




create table bs.blend_components (
	blend_components_id integer not null,
	product_lot_id not null,
	recipe_id not null,
	component_id not null,
	component_type_amount real,
	foreign key (product_lot_id) references product_lot,
	foreign key (component_id) references component_list,
	foreign key (recipe_id) references recipe_list,
	primary key (blend_components_id));





insert into bs.container_types
	(container_type_name)
	values ('Sample'), ('Tote'), ('Railcar');


insert into bs.container_list
	(container_name, container_type_id)
	values ('SAMPLE',1), ('TOTE',2);

insert into bs.status_list
	(status_name)
values
	('AVAILABLE'),
	('TESTED'),
	('UNAVAILABLE');

insert into bs.lot_list
	(lot_name)
values
	('');



`
	//TODO

	// select lot_name,container_name,inbound_lot_name from product_lot
	//
	// join blend_components using (product_lot_id)
	// join component_list using (component_id)
	// join inbound_lot using (inbound_lot_id)
	// join container_list using (container_id)
	//
	//  where lot_name like "BSQL%"

	// 	//TODO
	//
	// select produced_product_lot.lot_name,component_product_lot.lot_name,* from product_lot produced_product_lot
	//
	// join blend_components using (product_lot_id)
	// join component_list using (component_id)
	// join product_lot component_product_lot on (component_list.product_lot_id = component_product_lot.product_lot_id)
	//
	//
	//  where produced_product_lot.lot_name like "BSQL%"
	//TODO

	// select produced_product_lot.lot_name,container_name,inbound_lot_name,* from product_lot produced_product_lot
	//
	// --
	// join blend_components using (product_lot_id)
	// join component_list using (component_id)
	// left join product_lot component_product_lot on (component_list.product_lot_id = component_product_lot.product_lot_id)
	// left join inbound_lot using (inbound_lot_id)
	// left join container_list using (container_id)
	//
	//  where produced_product_lot.lot_name like "BSQL%"

	//TODO

	/*

	   create table bs.sample_container_list (
	   	sample_container_id integer not null,
	   	sample_container_name text not null,
	   	primary key (sample_container_id),
	   	unique (sample_container_name));*/

	//TODO lot number authority
	// https://www.jujens.eu/posts/en/2021/Apr/08/sequence-reset-every-day/

	// TODO add ranges table
	// _min real,
	// _max real,
	// _target real,

	db.Exec(sqlStmt)
	// _, err = db.Exec (sqlStmt)
	// if err != nil {
	// 	log.Printf ("[%s]: %q\n",  sqlStmt,  err)
	// 	// return
	// }

	DB.Check_db(db)
	DB.DBinit(db)

}
