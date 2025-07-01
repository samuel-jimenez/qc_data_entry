package viewer

import (
	"fmt"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

/* SQLFilter
 *
 */
type SQLFilter interface {
	Key() string
	Get() string
}

// /* SQLFilterBase
//  *
//  */
// type SQLFilterBase struct {
// 	key string
// }
//
// func (filter SQLFilterBase) Key() string { return filter.key }
// SQLFilterBase

/* SQLFilterDiscrete
 *
 */
type SQLFilterDiscrete struct {
	key string
	set []string
}

func (filter SQLFilterDiscrete) Key() string { return filter.key }
func (filter SQLFilterDiscrete) Get() string {
	if len(filter.set) == 0 {
		return ""
	}
	return fmt.Sprintf("%s in ('%s')",
		filter.key, strings.Join(filter.set, "','"))
}

/* SQLFilterContinuous
 *
 */
type SQLFilterContinuous struct {
	key                string
	range_lo, range_hi nullable.NullString
}

func (filter SQLFilterContinuous) Key() string { return filter.key }
func (filter SQLFilterContinuous) Get() string {
	var filter_out string
	if filter.range_lo.Valid {
		filter_out = fmt.Sprintf("%s >= %s",
			filter.key, filter.range_lo.String)
		if filter.range_hi.Valid {
			filter_out = fmt.Sprintf("%s and ", filter_out)
		}
	}
	if filter.range_hi.Valid {
		filter_out = fmt.Sprintf("%s%s <= %s", filter_out, filter.key, filter.range_hi.String)
	}
	return filter_out
}
