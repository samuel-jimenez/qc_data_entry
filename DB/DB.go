package DB

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	err error

	DB_VERSION = "0.0.2"

	DB_Select_product_recipe, DB_Insert_product_recipe,
	DB_Select_recipe_components, DB_Insert_recipe_component,
	DB_Select_name_component_types, DB_Select_all_component_types, DB_Insert_component_types,
	db_select_product_id, db_insert_product,
	DB_Select_product_info,
	db_select_lot_id, db_insert_lot,
	DB_Select_lot_all, DB_Select_lot_info,
	DB_Select_product_customer_id, DB_Select_product_customer_info,
	db_select_product_customer, db_insert_product_customer,
	DB_Update_lot_customer,
	DB_Select_all_sample_points, DB_Select_product_sample_points,
	DB_insert_sample_point,
	DB_insert_measurement,
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

func Check_db(db *sql.DB) {

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
		err = errors.New(fmt.Sprintf("Database version mismatch: Required: %s, found: %s", DB_VERSION, found_db_version))
		log.Printf("%q\n", err)
		panic(err)
	}
}

func DBinit(db *sql.DB) {

	DB_Select_product_recipe = PrepareOrElse(db, `
	select recipe_list_id
		from bs.recipe_list
		where product_id = ?
	`)

	DB_Insert_product_recipe = PrepareOrElse(db, `
	insert into bs.recipe_list
		(product_id)
		values (?)
	returning recipe_list_id
	`)

	DB_Select_recipe_components = PrepareOrElse(db, `
	select component_type_id, component_type_name, component_type_amount, component_add_order
		from bs.recipe_components
		join bs.component_types
		using (component_type_id)
		where recipe_list_id = ?
		order by component_add_order
	`)

	DB_Insert_recipe_component = PrepareOrElse(db, `
	insert into bs.recipe_components
		(recipe_list_id,component_type_id,component_type_amount,component_add_order)
		values (?,?,?,?)
	returning recipe_components_id
	`)

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
	order by product_moniker_name,product_name_internal
	`)

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

	db_select_lot_id = PrepareOrElse(db, `
	select lot_id
		from bs.product_lot
		where lot_name = ? and product_id = ?
	`)

	DB_Select_lot_info = PrepareOrElse(db, `
	select lot_id, lot_name
		from bs.product_lot
		where product_id = ?
	`)
	DB_Select_lot_all = PrepareOrElse(db, `
	select lot_id, lot_name
		from bs.product_lot
		order by lot_name
	`)

	db_insert_lot = PrepareOrElse(db, `
	insert into bs.product_lot
		(lot_name,product_id)
		values (?,?)
		returning lot_id
	`)

	DB_insert_sample_point = PrepareOrElse(db, `
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

	DB_insert_measurement = PrepareOrElse(db, `
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

	DB_Update_lot_customer = PrepareOrElse(db, `
	update bs.product_lot
	set
		product_customer_id=?
	where lot_id=?
		`)

}

func insert(insert_statement *sql.Stmt, proc_name string, args ...any) int64 {
	var insert_id int64
	result, err := insert_statement.Exec(args...)
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return INVALID_ID
	}
	insert_id, err = result.LastInsertId()
	if err != nil {
		log.Printf("%q: %s\n", err, proc_name)
		return -2
	}
	return insert_id
}

func Insel(insert_statement, select_statement *sql.Stmt, proc_name string, args ...any) int64 {
	var insel_id int64
	if select_statement.QueryRow(args...).Scan(&insel_id) != nil {
		//no rows
		insel_id = insert(insert_statement, proc_name, args...)
	}
	return insel_id
}

func Insel_product_id(product_name_full string) int64 {

	product_moniker_name, product_name_internal, _ := strings.Cut(product_name_full, " ")

	return Insel(db_insert_product, db_select_product_id, "Debug: insel_product_id", product_name_internal, product_moniker_name)
}

func Insel_lot_id(lot_name string, product_id int64) int64 {
	return Insel(db_insert_lot, db_select_lot_id, "Debug: insel_lot_id", lot_name, product_id)
}

func Insel_product_name_customer(product_name_customer string, product_id int64) int64 {
	return Insel(db_insert_product_customer, db_select_product_customer, "Debug: Insel_product_name_customer", product_name_customer, product_id)
}

func Forall(proc_name string, start_fn func(), row_fn func(*sql.Rows), select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: %q: %s\n", err, proc_name)
		return
	}
	start_fn()
	for rows.Next() {
		row_fn(rows)
	}
}

/*
// TODO func Forall(calling_fn_name string, start_fn func(), row_fn func(*sql.Rows) err, select_statement *sql.Stmt, args ...any) {
func Forall(proc_name string, start_fn func(), row_fn func(*sql.Rows) error, select_statement *sql.Stmt, args ...any) {
	rows, err := select_statement.Query(args...)
	if err != nil {
		log.Printf("error: %q: %s\n", err, proc_name)
		return
	}
	start_fn()
	for rows.Next() {
		if err = row_fn(rows); err != nil {
			log.Printf("error: %q: %s\n", err, proc_name)
			return
		}
	}
}*/
