package DB

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/util"
	"github.com/samuel-jimenez/windigo"
)

var (
	err error

	DB_VERSION = "0.0.4"

	// recipe_list
	DB_Select_product_recipe, DB_Insert_product_recipe,
	// recipe_components
	DB_Select_recipe_components_id, DB_Insert_recipe_component, DB_Update_recipe_component, DB_Delete_recipe_component, DB_Select_recipe_components_count,
	// component_types
	DB_Select_name_component_types, DB_Select_all_component_types, DB_Insert_component_types,
	DB_Select_component_type_product, DB_Insert_internal_product_component_type, DB_Insert_inbound_product_component_type,
	DB_Select_inbound_product_component_type_id,
	// inbound_product
	DB_Select_inbound_product_name, DB_Insert_inbound_product,
	// container_list
	DB_Select_container_id, DB_Select_container_all, DB_Insert_container, DB_Update_container_type,
	// inbound_provider_list
	DB_Select_inbound_provider_id, DB_Select_inbound_provider_all, DB_Insert_inbound_provider,
	// inbound_status_list
	DB_Select_all_inbound_status_list, DB_Select_name_inbound_status_list,
	// internal_status_list
	DB_Select_all_internal_status_list, DB_Select_name_internal_status_list,
	// inbound_lot
	DB_Select_inbound_lot_status, DB_Select_inbound_lot_all, DB_Insert_inbound_lot, DB_Update_inbound_lot_status,
	DB_Select_name_inbound_lot_status, DB_Select_inbound_lot_recipe, DB_Select_inbound_lot_components,
	// inbound_relabel
	DB_Insert_inbound_relabel, DB_Select_inbound_relabel_all,
	// component_list
	DB_Select_inbound_blend_component, DB_Insert_inbound_blend_component,
	DB_Select_internal_blend_component, DB_Insert_internal_blend_component,

	// blend_components
	DB_Insert_Product_blend,
	// lot_list
	DB_Insert_lot, DB_Select_lot, DB_Select_blend_lot, DB_Select_lot_list_all,
	DB_Select_lot_list_name, DB_Select_lot_list_for_name_status, DB_Update_lot_list__status, DB_Update_lot_list__component_status,
	DB_Select_product_lot_list_name,
	DB_Select_product_lot_list_sources,
	// product_lot
	db_select_id_lot, DB_Insert_product_lot,
	DB_Select_product_lot_all, DB_Select_product_lot_product,
	DB_Insert_blend_lot, DB_Update_lot_recipe,
	DB_Update_lot_customer,
	DB_Select_product_lot_components,
	// product_line
	db_select_product_id, db_insert_product,
	DB_Select_product_info,
	// product_customer_line
	DB_Select_product_customer_id, DB_Select_product_customer_info,
	db_select_product_customer, db_insert_product_customer,
	// bs.product_sample_points
	DB_Select_all_sample_points, DB_Select_product_sample_points,
	DB_Insel_sample_point,
	// bs.qc_tester_list
	DB_Select_all_qc_tester, DB_Insel_qc_tester,
	// bs.qc_samples
	DB_insert_measurement,
	DB_Update_qc_samples_storage,
	// bs.product_sample_storage
	DB_Select_product_sample_storage_capacity, DB_Select_gen_product_sample_storage,
	DB_Update_product_sample_storage_qc_sample, DB_Update_product_sample_storage_capacity,
	// bs.qc_sample_storage_list
	DB_Insert_sample_storage,
	// bs.product_appearance
	DB_Insert_appearance,
	DB_Select_product_details,
	DB_Upsert_product_details, DB_Upsert_product_type,
	DB_Select_product_coa_details, DB_Upsert_product_coa_details *sql.Stmt

	INVALID_ID     int64 = 0
	DEFAULT_LOT_ID int64 = 1
)

func PrepareOrElse(db *sql.DB, sqlStatement string) *sql.Stmt {
	preparedStatement, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		panic(err)

	}
	return preparedStatement
}

func Check_db(db *sql.DB, showWindowp bool) {

	var found_database_version_major,
		found_database_version_minor,
		found_database_version_revision int

	sqlStmt := `
		select database_version_major, database_version_minor, database_version_revision
		from bs.database_info
		`
	err = db.QueryRow(sqlStmt).Scan(&found_database_version_major,
		&found_database_version_minor,
		&found_database_version_revision)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		panic(err)
	}

	found_db_version := fmt.Sprintf("%d.%d.%d", found_database_version_major,
		found_database_version_minor,
		found_database_version_revision)

	if DB_VERSION != found_db_version {
		message := fmt.Sprintf("Database version mismatch: Required: %s, found: %s", DB_VERSION, found_db_version)
		err = errors.New(message)
		if showWindowp {
			windigo.Error(nil, message)
		}
		log.Printf("%q\n", err)
		panic(err)
	}
}

func DBinit(db *sql.DB) {
	// recipe_list
	DB_Select_product_recipe = PrepareOrElse(db, `
	select recipe_id
		from bs.recipe_list
		where product_id = ?
	`)

	DB_Insert_product_recipe = PrepareOrElse(db, `
	insert into bs.recipe_list
		(product_id)
		values (?)
	returning recipe_id
	`)

	// recipe_components
	DB_Select_recipe_components_id = PrepareOrElse(db, `
	select recipe_components_id, component_type_id, component_type_name, component_type_amount, component_add_order
		from bs.recipe_components
		join bs.component_types
		using (component_type_id)
		where recipe_id = ?
		order by component_add_order
	`)

	DB_Insert_recipe_component = PrepareOrElse(db, `
	insert into bs.recipe_components
		(recipe_id,component_type_id,component_type_amount,component_add_order)
		values (?,?,?,?)
	returning recipe_components_id
	`)
	DB_Update_recipe_component = PrepareOrElse(db, `
	update  bs.recipe_components
	set
		component_type_id=?,
		component_type_amount=?,
		component_add_order=?
	where recipe_components_id = ?
	`)

	DB_Delete_recipe_component = PrepareOrElse(db, `
	delete from bs.recipe_components
	where recipe_components_id = ?
	`)

	DB_Select_recipe_components_count = PrepareOrElse(db, `
	select max(component_add_order)
		from bs.recipe_components
		where recipe_id = ?
	`)

	// component_types
	DB_Select_name_component_types = PrepareOrElse(db, `
	select component_type_id
		from bs.component_types
		where component_type_name = ?
	`)

	DB_Select_all_component_types = PrepareOrElse(db, `
	select component_type_id, component_type_name
		from bs.component_types
		order by component_type_name
	`)

	DB_Insert_component_types = PrepareOrElse(db, `
	insert into bs.component_types
		(component_type_name)
		values (?)
	returning component_type_id
	`)

	//TODO
	// 	DB_Select_component_type_product = PrepareOrElse(db, `
	// select component_type_id, component_type_name, product_id, inbound_product_id
	// from bs.component_types
	// left join bs.component_type_product_internal using (component_type_id)
	// left join bs.component_type_product_inbound using (component_type_id)
	// where component_type_id = ?
	// `)

	DB_Select_component_type_product = PrepareOrElse(db, `
select
false, product_lot_id, product_name_internal, lot_name, ''

from bs.component_type_product_internal
join bs.product_line using (product_id)
join bs.product_lot using (product_id)
join bs.lot_list using (lot_id)
join bs.internal_status_list using (internal_status_id)

where component_type_id = ?1
and internal_status_name = ?2

union

select
true, inbound_lot_id,inbound_product_name, inbound_lot_name, container_name

from bs.component_type_product_inbound
join bs.inbound_product using (inbound_product_id)
join bs.inbound_lot using (inbound_product_id)
join bs.container_list using (container_id)
join bs.inbound_status_list using (inbound_status_id)

where component_type_id = ?1
and inbound_status_name != ?3

`)

	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where blend_components.product_lot_id =?

	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.lot_list using (lot_id)
	// join bs.product_lot using (product_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where lot_name =?

	DB_Insert_internal_product_component_type = PrepareOrElse(db, `
	insert into bs.component_type_product_internal
		(component_type_id, product_id)
		values (?,?)
	returning component_type_product_internal_id
	`)

	DB_Insert_inbound_product_component_type = PrepareOrElse(db, `
	insert into bs.component_type_product_inbound
		(component_type_id, inbound_product_id)
		values (?,?)
	returning component_type_product_inbound_id
	`)
	DB_Select_inbound_product_component_type_id = PrepareOrElse(db, `
select inbound_product_id, inbound_product_name
from bs.component_type_product_inbound
left join bs.inbound_product using (inbound_product_id)
where component_type_id = ?1
		`)
	// inbound_product
	DB_Select_inbound_product_name = PrepareOrElse(db, `
select inbound_product_id
from bs.inbound_product
where inbound_product_name = ?
	`)
	DB_Insert_inbound_product = PrepareOrElse(db, `
	insert into bs.inbound_product
		(inbound_product_name)
		values (?)
		returning inbound_product_id
		`)

	// container_list
	DB_Select_container_id = PrepareOrElse(db, `
	select container_id
		from bs.container_list
		where container_name = ?
	`)

	DB_Select_container_all = PrepareOrElse(db, `
	select container_id, container_name
		from bs.container_list
		order by container_name
	`)

	DB_Insert_container = PrepareOrElse(db, `
	insert into bs.container_list
		(container_name)
		values (?)
		returning container_id
		`)

	DB_Update_container_type = PrepareOrElse(db, `
update
bs.container_list

set
container_type_id=?2

where
container_id = ?1
		`)

	// inbound_provider_list
	DB_Select_inbound_provider_id = PrepareOrElse(db, `
select inbound_provider_id
from bs.inbound_provider_list
where inbound_provider_name = ?
	`)

	DB_Select_inbound_provider_all = PrepareOrElse(db, `
	select inbound_provider_id, inbound_provider_name
		from bs.inbound_provider_list
		order by inbound_provider_name
	`)

	DB_Insert_inbound_provider = PrepareOrElse(db, `
	insert into bs.inbound_provider_list
		(inbound_provider_name)
		values (?)
		returning inbound_provider_id
		`)
	// inbound_status_list
	DB_Select_all_inbound_status_list = PrepareOrElse(db, `
select
	inbound_status_id, inbound_status_name
from bs.inbound_status_list
order by inbound_status_id
	`)
	DB_Select_name_inbound_status_list = PrepareOrElse(db, `
select
	inbound_status_id
from bs.inbound_status_list
where inbound_status_name = ?
	`)

	// internal_status_list
	DB_Select_all_internal_status_list = PrepareOrElse(db, `
select
	internal_status_id, internal_status_name
from bs.internal_status_list
order by internal_status_id
	`)
	DB_Select_name_internal_status_list = PrepareOrElse(db, `
select
	internal_status_id
from bs.internal_status_list
where internal_status_name = ?
	`)

	// inbound_lot
	DB_Select_inbound_lot_status = PrepareOrElse(db, `
select
	inbound_lot_id, inbound_lot_name, inbound_product_id, inbound_product_name, inbound_provider_id, inbound_provider_name, container_id, container_name, inbound_status_id, inbound_status_name
from bs.inbound_lot
join bs.inbound_product 		using (inbound_product_id)
join bs.inbound_provider_list 		using (inbound_provider_id)
join bs.container_list 			using (container_id)
join bs.inbound_status_list 			using (inbound_status_id)
		where inbound_status_name = ?
order by inbound_lot_name
	`)

	DB_Select_inbound_lot_all = PrepareOrElse(db, `
select
	inbound_lot_id, inbound_lot_name, inbound_product_id, inbound_product_name, inbound_provider_id, inbound_provider_name, container_id, container_name, inbound_status_id, inbound_status_name
from bs.inbound_lot
join bs.inbound_product 		using (inbound_product_id)
join bs.inbound_provider_list 		using (inbound_provider_id)
join bs.container_list 			using (container_id)
join bs.inbound_status_list 			using (inbound_status_id)
order by inbound_lot_name
	`)

	DB_Insert_inbound_lot = PrepareOrElse(db, `
insert into bs.inbound_lot
(inbound_lot_name,inbound_product_id,inbound_provider_id,container_id)
values (?,?,?,?)
returning inbound_lot_id
		`)
	DB_Update_inbound_lot_status = PrepareOrElse(db, `
update bs.inbound_lot
set
	inbound_status_id=?2
where inbound_lot_id=?1
`)

	DB_Select_name_inbound_lot_status = PrepareOrElse(db, `
select
	inbound_lot_name
from bs.inbound_lot
join bs.inbound_status_list 			using (inbound_status_id)
where inbound_status_name = ?
	`)

	DB_Select_inbound_lot_recipe = PrepareOrElse(db, `
select
	recipe_id, product_id, recipe_components_id, component_type_id, component_type_amount, component_add_order
from bs.inbound_lot
join bs.component_type_product_inbound 	using (inbound_product_id)
join bs.recipe_components 		using (component_type_id)
join bs.recipe_list			using (recipe_id)
where inbound_lot_id =  ?
order by recipe_id
	`)

	DB_Select_inbound_lot_components = PrepareOrElse(db, `
select
	inbound_lot_id, recipe_components_id, component_type_id, component_type_amount, component_add_order
from bs.inbound_lot
join bs.inbound_status_list 			using (inbound_status_id)
join bs.component_type_product_inbound 	using (inbound_product_id)
join bs.recipe_components 		using (component_type_id)
where recipe_id = ?1
	and component_type_id != ?2
	and inbound_status_name = ?3
order by component_add_order
	`)

	// inbound_relabel
	DB_Insert_inbound_relabel = PrepareOrElse(db, `
	insert into bs.inbound_relabel
		(lot_id,inbound_lot_id,container_id)
		values (?,?,?)
	returning inbound_relabel_id
	`)

	DB_Select_inbound_relabel_all = PrepareOrElse(db, `
	select
	inbound_relabel_id, lot_name
	from bs.inbound_relabel
	join bs.lot_list using (lot_id)
	`)
	// 			// TODO
	//
	// 		DB_Select_inbound_relabel = PrepareOrElse(db, `
	// 	select
	// 	lot_name inbound_product_name, inbound_lot_name, container_name
	// 	from bs.inbound_relabel
	// 	join bs.lot_list using (lot_id)
	//
	// join bs.inbound_lot 		using (inbound_lot_id)
	// join bs.inbound_lot 		using (inbound_lot_id)
	// join bs.inbound_product 		using (inbound_product_id)
	// join bs.inbound_provider_list 		using (inbound_provider_id)
	// join bs.container_list 			using (container_id)
	// 		(lot_id,inbound_lot_id,container_id)
	// 		values (?,?,?)
	// 	returning inbound_relabel_id
	// 	`)

	// component_list
	DB_Select_inbound_blend_component = PrepareOrElse(db, `
select component_id
from bs.component_list
where component_type_id = ?
and inbound_lot_id = ?
	`)
	DB_Insert_inbound_blend_component = PrepareOrElse(db, `
insert into bs.component_list
	(component_type_id,inbound_lot_id)
	values (?,?)
returning component_id
`)
	DB_Select_internal_blend_component = PrepareOrElse(db, `
select component_id
from bs.component_list
where component_type_id = ?
and product_lot_id = ?
	`)
	DB_Insert_internal_blend_component = PrepareOrElse(db, `
insert into bs.component_list
	(component_type_id,product_lot_id)
	values (?,?)
returning component_id
`)

	// blend_components
	DB_Insert_Product_blend = PrepareOrElse(db, `
	insert into bs.blend_components
		(product_lot_id, recipe_components_id, component_id, component_required_amount)
		values (?,?,?,?)
	returning blend_components_id
	`)

	// lot_list
	DB_Insert_lot = PrepareOrElse(db, `
insert into bs.lot_list
	(lot_name)
values (?)
returning lot_id
`)
	DB_Select_lot = PrepareOrElse(db, `
select lot_id
from bs.lot_list
where lot_name = ?
`)
	DB_Select_blend_lot = PrepareOrElse(db, `
select count(lot_name)
	from bs.lot_list
where lot_name like ?
`)
	DB_Select_lot_list_all = PrepareOrElse(db, `
	select
	lot_id, lot_name
	from bs.lot_list
	`)

	DB_Select_lot_list_name = PrepareOrElse(db, `
	select
	lot_id, lot_name
	from bs.lot_list
	where lot_name like ?
	`)

	DB_Select_lot_list_for_name_status = PrepareOrElse(db, `
	select
	lot_id, lot_name
	from bs.lot_list
	join bs.internal_status_list using (internal_status_id)
	where lot_name like ?
	and internal_status_name = ?
	`)

	DB_Update_lot_list__status = PrepareOrElse(db, `
update bs.lot_list
set
	internal_status_id=bs.internal_status_list.internal_status_id
from
bs.internal_status_list

where lot_id = ?
and internal_status_name = ?

		`)

	DB_Update_lot_list__component_status = PrepareOrElse(db, `
with upd_table as (
select

lot_id, internal_status_list.internal_status_id
from bs.lot_list
join bs.product_lot using (lot_id)
join bs.blend_components using (product_lot_id)
join bs.component_list using (component_id)
join bs.inbound_lot using (inbound_lot_id)
full join bs.internal_status_list
where inbound_lot_name = ?1
and internal_status_name = ?2

union

select

lot_id, internal_status_list.internal_status_id
from bs.lot_list
	join bs.inbound_relabel using (lot_id)
	join bs.inbound_lot using (inbound_lot_id)
full join bs.internal_status_list
where inbound_lot_name = ?1
and internal_status_name = ?2
	)

update bs.lot_list
set
internal_status_id = upd_table.internal_status_id
from
upd_table
where lot_list.lot_id = upd_table.lot_id
		`)

	// TODO blend012 tests
	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where blend_components.product_lot_id =?

	// 	TODO get sources

	// select inbound_product_name,inbound_lot_name,container_name from blend_components
	// join bs.component_list using (component_id)
	// join bs.inbound_lot using (inbound_lot_id)
	// join bs.lot_list using (lot_id)
	// join bs.product_lot using (product_lot_id)
	// join bs.container_list using (container_id)
	// join bs.inbound_product using (inbound_product_id)
	// where lot_name =?

	DB_Select_product_lot_list_name = PrepareOrElse(db, `
select lot_id,  format('%s %s',  product_moniker_name, product_name_internal), lot_name, product_name_customer
from bs.lot_list
join bs.product_lot using (lot_id)
join bs.product_line using (product_id)
join bs.product_moniker using (product_moniker_id)
left join bs.product_customer_line using (product_customer_id)
where lot_name = ?1

union

select lot_id, inbound_product_name, lot_name, null
from bs.lot_list
	join bs.inbound_relabel using (lot_id)
	join bs.inbound_lot using (inbound_lot_id)
	join bs.inbound_product using (inbound_product_id)
where lot_name = ?1
`)

	///TODO decide between:
	// 	create recipe
	//
	// 	create blend, assigning components, amounts
	//
	// 	capture actual amounts
	//
	// 	and:
	//
	// 		create recipe
	//
	// 	create blend, assigning amounts
	//
	// 		capture components, actual amounts

	// TODO sample_size := 500.
	DB_Select_product_lot_list_sources = PrepareOrElse(db, `
select inbound_product_name, inbound_lot_id, inbound_lot_name, container_name, component_required_amount, component_add_order, true

from bs.lot_list
join bs.product_lot using (lot_id)
join bs.blend_components using (product_lot_id)
join bs.component_list using (component_id)
join bs.recipe_components		using (recipe_components_id)
join bs.inbound_lot using (inbound_lot_id)
join bs.container_list using (container_id)
join bs.inbound_product using (inbound_product_id)

where lot_name = ?1

union

select inbound_product_name, inbound_lot_id, inbound_lot_name, container_name, 500, 0 component_add_order, true

from bs.lot_list
join bs.inbound_relabel using (lot_id)
join bs.inbound_lot using (inbound_lot_id)
join bs.container_list using (container_id)
join bs.inbound_product using (inbound_product_id)

where lot_name = ?1

		order by component_add_order
`)

	// product_lot
	db_select_id_lot = PrepareOrElse(db, `
select product_lot_id
from bs.product_lot
where lot_id = ? and product_id = ?
	`)

	DB_Select_product_lot_product = PrepareOrElse(db, `
select product_lot_id, lot_name
from bs.product_lot
join bs.lot_list using (lot_id)
where product_id = ?
`)
	DB_Select_product_lot_all = PrepareOrElse(db, `
select product_lot_id, lot_name
from bs.product_lot
join bs.lot_list using (lot_id)
order by lot_name
`)

	DB_Insert_product_lot = PrepareOrElse(db, `
	insert into bs.product_lot
		(lot_id,product_id)
		values (?,?)
		returning product_lot_id
	`)

	DB_Insert_blend_lot = PrepareOrElse(db, `
	insert into bs.product_lot
		(lot_id,product_id,product_customer_id,recipe_id)
		values (?,?,?,?)
	returning product_lot_id
	`)

	DB_Update_lot_recipe = PrepareOrElse(db, `
	update bs.product_lot
	set
		recipe_id=?
	where product_lot_id=?
		`)

	DB_Update_lot_customer = PrepareOrElse(db, `
	update bs.product_lot
	set
		product_customer_id=?
	where product_lot_id=?
		`)

	/// FIXME not actually used
	DB_Select_product_lot_components = PrepareOrElse(db, `
select inbound_product_name, inbound_lot_name, container_name from blend_components
join bs.component_list using (component_id)
join bs.inbound_lot using (inbound_lot_id)
join bs.container_list using (container_id)
join bs.inbound_product using (inbound_product_id)
where blend_components.product_lot_id =?
`)

	// product_line
	db_insert_product = PrepareOrElse(db, `
	insert into bs.product_line
		(product_name_internal, product_moniker_id)
		select ?, product_moniker_id
			from bs.product_moniker
			where product_moniker_name = ?
		returning product_id
		`)

	db_select_product_id = PrepareOrElse(db, `
	select product_id
		from bs.product_line
		join bs.product_moniker using (product_moniker_id)
		where product_name_internal = ?
		and product_moniker_name = ?
	`)

	DB_Select_product_info = PrepareOrElse(db, `
	select product_id, product_name_internal, product_moniker_name
		from bs.product_line
		join bs.product_moniker using (product_moniker_id)
	order by product_moniker_name, product_name_internal
	`)

	// product_customer_line
	DB_Select_product_customer_info = PrepareOrElse(db, `
	select product_customer_id, product_name_customer
		from bs.product_customer_line
		where product_id = ?
		`)

	DB_Select_product_customer_id = PrepareOrElse(db, `
	select product_customer_id
		from bs.product_customer_line
		where product_name_customer = ? and product_id = ?
	`)

	db_select_product_customer = PrepareOrElse(db, `
	select product_customer_id
		from bs.product_customer_line
		where product_name_customer = ? and product_id = ?
	`)

	db_insert_product_customer = PrepareOrElse(db, `
	insert into bs.product_customer_line
		(product_name_customer,product_id)
		values (?,?)
	returning product_customer_id
	`)

	// bs.product_sample_points
	DB_Insel_sample_point = PrepareOrElse(db, `
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

	DB_Select_all_sample_points = PrepareOrElse(db, `
	select sample_point_id, sample_point
		from bs.product_sample_points
		order by sample_point_id
	`)

	DB_Select_product_sample_points = PrepareOrElse(db, `
	select distinct sample_point_id, sample_point
		from bs.product_lot
		join bs.qc_samples using (lot_id)
		join bs.product_sample_points using (sample_point_id)
		where product_id = ?
		order by sample_point_id
	`)

	// bs.qc_tester_list
	DB_Select_all_qc_tester = PrepareOrElse(db, `
	select qc_tester_id, qc_tester_name
		from bs.qc_tester_list
		order by qc_tester_name
	`)

	DB_Insel_qc_tester = PrepareOrElse(db, `
	with val (qc_tester_name) as (
		values
			(?)
		),
		sel as (
			select qc_tester_name, qc_tester_id
			from val
			left join bs.qc_tester_list using (qc_tester_name)
		)
	insert into bs.qc_tester_list (qc_tester_name)
	select distinct qc_tester_name from sel where qc_tester_id is null
	returning qc_tester_id, qc_tester_name
	`)

	// bs.qc_samples
	DB_insert_measurement = PrepareOrElse(db, `
	with
		val (lot_id, sample_point, qc_tester_name, time_stamp,  ph, specific_gravity, string_test, viscosity) as (
			values
				(?, ?, ?, ?, ?, ?, ?, ?)
		),
		sel as (
			select lot_id, sample_point_id, qc_tester_id, time_stamp, ph, specific_gravity, string_test, viscosity
			from val
			left join bs.qc_tester_list using (qc_tester_name)
			left join bs.product_sample_points using (sample_point)
		)
	insert into bs.qc_samples (lot_id, sample_point_id, qc_tester_id, time_stamp, ph, specific_gravity, string_test, viscosity)
	select lot_id, sample_point_id, qc_tester_id, time_stamp, ph, specific_gravity, string_test, viscosity
		from   sel
	returning qc_id;
	`)

	DB_Update_qc_samples_storage = PrepareOrElse(db, `
update bs.qc_samples
	set
qc_sample_storage_id = ?2

where qc_id = ?1
	`)

	// bs.product_sample_storage
	DB_Select_product_sample_storage_capacity = PrepareOrElse(db, `
select

qc_sample_storage_id, qc_storage_capacity

from bs.product_sample_storage
join bs.product_line using (product_moniker_id)
where product_id = ?1
	`)

	DB_Select_gen_product_sample_storage = PrepareOrElse(db, `
			select

product_sample_storage_id, product_moniker_name, qc_sample_storage_offset, qc_sample_storage_name, min(time_stamp), max(time_stamp), retain_storage_duration

from bs.product_sample_storage
join bs.product_line using (product_moniker_id)
join bs.qc_sample_storage_list using (qc_sample_storage_id, product_moniker_id)
join bs.qc_samples using (qc_sample_storage_id)
join bs.product_moniker using (product_moniker_id)
where product_id = ?1
	`)

	DB_Update_product_sample_storage_qc_sample = PrepareOrElse(db, `
update bs.product_sample_storage
	set
qc_sample_storage_id = ?2,
qc_sample_storage_offset = ?3,
qc_storage_capacity = max_storage_capacity

where product_sample_storage_id = ?1
	`)

	DB_Update_product_sample_storage_capacity = PrepareOrElse(db, `
update bs.product_sample_storage
	set
qc_storage_capacity = qc_storage_capacity - ?2

where qc_sample_storage_id = ?1
	`)
	/*
	    *
	   			select

	   product_sample_storage_id, product_moniker_name, qc_sample_storage_offset, qc_sample_storage_name, min(time_stamp), max(time_stamp), retain_storage_duration, qc_storage_capacity

	   from product_sample_storage
	   join product_line using (product_moniker_id)
	   join qc_sample_storage_list using (qc_sample_storage_id, product_moniker_id)
	   left join qc_samples using (qc_sample_storage_id)
	   join product_moniker using (product_moniker_id)
	   group by product_moniker_id

	   *
	*/

	// bs.qc_sample_storage_list

	// DB_Select_all_sample_storage = PrepareOrElse(db, `
	// select qc_sample_storage_id, qc_sample_storage_name, product_moniker_id
	// 	from bs.qc_sample_storage_list
	// 	order by qc_sample_storage_name
	// `)
	DB_Insert_sample_storage = PrepareOrElse(db, `
	insert into bs.qc_sample_storage_list
		( qc_sample_storage_name, product_moniker_id )
		values (?, ?)
	returning qc_sample_storage_id
	`)

	// bs.product_appearance
	DB_Insert_appearance = PrepareOrElse(db, `
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

	DB_Select_product_details = PrepareOrElse(db, `
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
				left join bs.product_appearance using (product_appearance_id)
			where product_id = ?1),

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
				left join bs.product_appearance using (product_appearance_id)
			where product_id = ?1)
	select
		product_type_id,
		container_type_id,
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
	from measured
		full join published using (product_id)
		join bs.product_types using (product_type_id)
	where product_id = ?1
	`)

	DB_Select_product_coa_details = PrepareOrElse(db, `
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

	DB_Upsert_product_details = PrepareOrElse(db, `
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

	DB_Upsert_product_type = PrepareOrElse(db, `
	insert into bs.product_ranges_measured
		(product_id,
		product_type_id)
	values (?,?)
	on conflict(product_id) do update set

		product_type_id=excluded.product_type_id

		returning range_id
		`)
	DB_Upsert_product_coa_details = PrepareOrElse(db, `
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

}

func Select_Error(proc_name string, query *sql.Row, args ...any) error {
	err := query.Scan(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
	}
	return err
}

// TODO permit empty
// if err != sql.ErrNoRows { // no row? no problem!
// 		log.Printf("Error: [%s]: %q\n", proc_name, err)
// 	}
// 	return nil

func Select_Panic(proc_name string, query *sql.Row, args ...any) {
	err := query.Scan(args...)
	if err != nil {
		log.Printf("Critical hit: [%s]: %q\n", proc_name, err)
		panic(err)
	}
}

func Select_Panic_ErrorBox(proc_name string, query *sql.Row, args ...any) {
	err := query.Scan(args...)
	if err != nil {
		log.Printf("Critical hit: [%s]: %q\n", proc_name, err)
		// 			TODO: make windigo.Error to avoid
		// [printf] (default) non-constant format string in call to github.com/samuel-jimenez/windigo.Errorf
		windigo.Errorf(nil, "Something's gone wrong.")
		panic(err)
	}
}

func Insert(proc_name string, insert_statement *sql.Stmt, args ...any) int64 {
	var insert_id int64
	result, err := insert_statement.Exec(args...)
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return INVALID_ID
	}
	insert_id, err = result.LastInsertId()
	if err != nil {
		log.Printf("Err: [%s]: %q\n", proc_name, err)
		return INVALID_ID
	}
	return insert_id
}

func Exec_Error(proc_name string, statement *sql.Stmt, args ...any) error {
	_, err := statement.Exec(args...)
	util.LogError(proc_name, err)
	return err
}
func Update(proc_name string, update_statement *sql.Stmt, args ...any) error {
	return Exec_Error(proc_name, update_statement, args...)
}
func Delete(proc_name string, delete_statement *sql.Stmt, args ...any) error {
	return Exec_Error(proc_name, delete_statement, args...)
}

func Insel(proc_name string, insert_statement, select_statement *sql.Stmt, args ...any) int64 {
	var insel_id int64
	if select_statement.QueryRow(args...).Scan(&insel_id) != nil {
		//no rows
		insel_id = Insert(proc_name, insert_statement, args...)
	}
	return insel_id
}

func Insel_product_id(product_name_full string) int64 {

	product_moniker_name, product_name_internal, _ := strings.Cut(product_name_full, " ")

	return Insel("Insel_product_id", db_insert_product, db_select_product_id, product_name_internal, product_moniker_name)
}

func Insel_lot_id(lot_name string) int64 {
	return Insel("Insel_lot_id", DB_Insert_lot, DB_Select_lot, lot_name)
}

func Insel_product_lot_id(Lot_id, product_id int64) int64 {
	return Insel("Insel_product_lot_id", DB_Insert_product_lot, db_select_id_lot, Lot_id, product_id)
}

func Insel_product_name_customer(product_name_customer string, product_id int64) int64 {
	return Insel("Insel_product_name_customer", db_insert_product_customer, db_select_product_customer, product_name_customer, product_id)
}

//TODO add variant with if err != sql.ErrNoRow

func Forall(proc_name string, start_fn func(), row_fn func(row *sql.Rows), select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		row_fn(rows)
	}
}

// TODO count?
func Forall_err(proc_name string, start_fn func(), row_fn func(row *sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Printf("error: [%s]: %q\n", proc_name, err)
		}
	}
}

func Forall_exit(proc_name string, start_fn func(), row_fn func(row *sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n", proc_name, err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Fatalf("Crit: [%s]: %q\n", proc_name, err)
			panic("Forall_exit")
		}
	}
}

/*
// TODO func Forall(calling_fn_name string, start_fn func(), row_fn func(*sql.Rows) err, select_statement *sql.Stmt, args ...any) {
func Forall(proc_name string, start_fn func(), row_fn func(*sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: [%s]: %q\n",  proc_name,  err)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Printf("error: [%s]: %q\n",  proc_name,  err)
			return
		}
	}
}*/
