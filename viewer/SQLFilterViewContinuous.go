package viewer

import (
	"strconv"
	"time"

	"github.com/samuel-jimenez/qc_data_entry/GUI"
	"github.com/samuel-jimenez/windigo"
)

/* SQLFilterViewContinuous
 *
 */
type SQLFilterViewContinuous struct {
	SQLFilterView

	field_panel          *windigo.AutoPanel
	min_field, max_field NullStringView
}

func NewSQLFilterViewContinuous(parent windigo.Controller, key, label string) *SQLFilterViewContinuous {
	view := new(SQLFilterViewContinuous)
	view.key = key

	view.AutoPanel = windigo.NewAutoPanel(parent)

	view.label = SQLFilterViewHeader_from_new(view, label)

	view.field_panel = windigo.NewAutoPanel(view.AutoPanel)
	view.label.SetHidePanel(view.field_panel)

	view.field_panel.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
	view.min_field = NewNullStringView(view.field_panel)
	view.max_field = NewNullStringView(view.field_panel)
	view.field_panel.Dock(view.min_field, windigo.Left)
	view.field_panel.Dock(view.max_field, windigo.Left)
	view.AutoPanel.Dock(view.label, windigo.Top)
	view.AutoPanel.Dock(view.field_panel, windigo.Top)

	return view
}

// TODO replace these with Get()  string
func (view *SQLFilterViewContinuous) Get() string {
	return SQLFilterContinuous{view.key,
		view.min_field.Get(),
		view.max_field.Get()}.Get()
}

func (view *SQLFilterViewContinuous) RefreshSize() {
	view.SQLFilterView.RefreshSize()
	view.field_panel.SetSize(GUI.OFF_AXIS, GUI.EDIT_FIELD_HEIGHT)
}

func (view *SQLFilterViewContinuous) SetFont(font *windigo.Font) {
	view.SQLFilterView.SetFont(font)
	view.field_panel.SetFont(font)
	view.min_field.SetFont(font)
	view.max_field.SetFont(font)
}

func (view *SQLFilterViewContinuous) Clear() {
	view.SQLFilterView.Clear()
	view.min_field.Clear()
	view.max_field.Clear()
}

/* SQLFilterViewContinuousTime
 *
 */
type SQLFilterViewContinuousTime struct {
	SQLFilterViewContinuous
}

func SQLFilterViewContinuousTime_from_new(parent windigo.Controller, key, label string) *SQLFilterViewContinuousTime {

	view := new(SQLFilterViewContinuousTime)
	view.SQLFilterViewContinuous = *NewSQLFilterViewContinuous(parent, key, label)

	return view
}

func wrapTime(Ttime time.Time) string {
	// [compiler] (IncompatibleAssign) cannot use Ttime.UnixNano() (value of type int64) as  value in argument to strconv.Itoa
	return strconv.Itoa(int(Ttime.UTC().UnixNano()))
	// return Ttime.Format(time.DateOnly)
}
func parseTime(Time_string string, endp bool) string {
	fmt_need_year := []string{"1/2", "1-2", "Jan 2", "2Jan", "2 Jan"}
	fmt_need_year_mth := []string{"1", "Jan"}
	fmt_full := []string{time.DateOnly, "1/2/6", "1/2/2006", "Jan 2 2006", "Jan 2, 2006", "2 Jan 2006"}
	fmt_full_year := []string{"2006"}

	for _, fmt := range fmt_full {
		if Ttime, err := time.Parse(fmt, Time_string); err == nil {
			if endp {
				Ttime = Ttime.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
			}
			return wrapTime(Ttime)
		}
	}

	for _, fmt := range fmt_need_year_mth {
		if Ttime, err := time.Parse(fmt, Time_string); err == nil {
			Ttime = Ttime.AddDate(time.Now().Year(), 0, 0)
			if endp {
				Ttime = Ttime.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)
			}
			return wrapTime(Ttime)
		}
	}

	for _, fmt := range fmt_need_year {
		if Ttime, err := time.Parse(fmt, Time_string); err == nil {
			Ttime = Ttime.AddDate(time.Now().Year(), 0, 0)
			if endp {
				Ttime = Ttime.AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
			}
			return wrapTime(Ttime)
		}
	}
	for _, fmt := range fmt_full_year {
		if Ttime, err := time.Parse(fmt, Time_string); err == nil {
			if endp {
				Ttime = Ttime.AddDate(1, 0, 0).Add(-1 * time.Nanosecond)
			}
			return wrapTime(Ttime)
		}
	}
	return "0"
}

func (view *SQLFilterViewContinuousTime) Get() string {
	return SQLFilterContinuous{view.key,
		view.min_field.Get(),
		view.max_field.Get()}.GetFMT(parseTime)
}
