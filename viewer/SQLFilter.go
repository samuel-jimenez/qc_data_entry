package viewer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samuel-jimenez/qc_data_entry/nullable"
)

//TODO can prob rewite this whole thing

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

func NewSQLFilterDiscrete(key string, set []string) *SQLFilterDiscrete {
	return &SQLFilterDiscrete{
		key,
		set,
	}
}

func (filter SQLFilterDiscrete) Key() string { return filter.key }
func (filter SQLFilterDiscrete) Get() string {
	if len(filter.set) == 0 {
		return ""
	}
	return fmt.Sprintf("%s in ('%s')",
		filter.key, strings.Join(filter.set, "','"))
}

func (filter *SQLFilterDiscrete) Set(set []string) {
	filter.set = set
}

/* SQLFilterIntDiscrete
 *
 */
type SQLFilterIntDiscrete struct {
	key string
	set []int
}

func NewSQLFilterIntDiscrete(key string, set []int) *SQLFilterIntDiscrete {
	return &SQLFilterIntDiscrete{
		key,
		set,
	}
}

func (filter SQLFilterIntDiscrete) Key() string { return filter.key }
func (filter SQLFilterIntDiscrete) Get() string {
	if len(filter.set) == 0 {
		return ""
	}
	mapped := make([]string, len(filter.set))

	for i, e := range filter.set {
		mapped[i] = strconv.Itoa(e)
	}

	return fmt.Sprintf("%s in (%s)",
		filter.key, strings.Join(mapped, ","))
}

func (filter *SQLFilterIntDiscrete) Set(set []int) {
	filter.set = set
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
	return filter.GetFMT(identity)
}

func identity(in string, endp bool) string {
	return in
}

// for values which are one point in a range, such as dates, always choosing the start effectively converts >= into >
// takes  func(string, bool) to specify start or end
func (filter SQLFilterContinuous) GetFMT(parse func(string, bool) string) string {
	var filter_out string
	if filter.range_lo.Valid {
		filter_out = fmt.Sprintf("%s >= %s",
			filter.key, parse(filter.range_lo.String, false))
		if filter.range_hi.Valid {
			filter_out = fmt.Sprintf("%s and ", filter_out)
		}
	}
	if filter.range_hi.Valid {
		filter_out = fmt.Sprintf("%s%s <= %s", filter_out, filter.key, parse(filter.range_hi.String, true))
	}
	return filter_out
}
