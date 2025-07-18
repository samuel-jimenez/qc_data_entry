package nullable

import (
	"database/sql"
)

/* NullInt64
 * Nullable Int64 that overrides sql.NullInt64
 *
 */
type NullInt64 struct {
	sql.NullInt64
}

func NullInt64Default() NullInt64 {
	return NullInt64{}
}

func NewNullInt64(val int64) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: val, Valid: true}}
}
