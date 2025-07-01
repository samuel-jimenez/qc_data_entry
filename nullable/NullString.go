package nullable

import (
	"database/sql"
	"strings"
)

/* NullString
 *
 */
type NullString struct {
	sql.NullString
}

func (a_n NullString) Compare(b_n NullString) int {
	var a, b string
	if a_n.Valid {
		a = a_n.String
	}
	if b_n.Valid {
		b = b_n.String
	}
	return strings.Compare(a, b)
}
