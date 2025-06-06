package DB

import (
	"database/sql"
	"log"
)

func PrepareOrElse(db *sql.DB, sqlStatement string) *sql.Stmt {
	preparedStatement, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		panic(err)

	}
	return preparedStatement
}
