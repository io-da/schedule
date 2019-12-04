package schedule

import "time"

const secondsPerDay int64 = 24 * 60 * 60
const firstWeekDay = secondsPerDay * 4
const secondsBeforeUnix = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay

var daysPerMonth = [...]int{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

// AttunedMonth is an upgrade to the regular time.Month attuned to the year provided.
// It is self aware of it's last day and week days.
type AttunedMonth struct {
	y         int
	mon       int
	unixYear  int64
	unixMonth int64

	isLeap  bool
	lastDay int
}

// NewAttunedMonth returns a reference to a new AttunedMonth.
func NewAttunedMonth(mon int, y int) *AttunedMonth {
	my := &AttunedMonth{}
	my.UpdateMonthYear(mon, y)
	return my
}

// UpdateMonth updates the month and recalculates the necessary components.
func (my *AttunedMonth) UpdateMonth(mon int) {
	if my.mon != mon {
		my.mon = mon
		my.determineUnixMonth()
		my.determineLastDay()
	}
}

// UpdateMonthYear updates both month and year and recalculates the necessary components.
func (my *AttunedMonth) UpdateMonthYear(mon int, y int) {
	if my.y != y || my.mon != mon {
		my.mon = mon
		my.y = y
		my.determineUnixYear()
		my.determineIsLeap()
		my.determineUnixMonth()
		my.determineLastDay()
	}
}

// Month returns the time.Month type representation of the month.
func (my *AttunedMonth) Month() time.Month {
	return time.Month(my.mon)
}

// Year returns the int value of the year.
func (my *AttunedMonth) Year() int {
	return my.y
}

// Contains verifies if this month contains the provided day.
func (my *AttunedMonth) Contains(d int) bool {
	return d >= 1 && d <= my.lastDay
}

// IsMonthLastDay verifies if the provided day is the last day of this month.
func (my *AttunedMonth) IsMonthLastDay(d int) bool {
	return d == my.lastDay
}

// MonthLastDay returns the last day of this month
func (my *AttunedMonth) MonthLastDay() int {
	return my.lastDay
}

// WeekDay returns the time.Weekday type representation of the provided day.
func (my *AttunedMonth) WeekDay(d int) time.Weekday {
	return time.Weekday((my.unixDay(d) + firstWeekDay) / secondsPerDay % 7)
}

func (my *AttunedMonth) unixDay(d int) int64 {
	return my.unixMonth + int64(d-1)*secondsPerDay
}

func (my *AttunedMonth) determineUnixYear() {
	y := my.y - 1
	my.unixYear = (int64(y*365+y/4-y/100+y/400) * secondsPerDay) - secondsBeforeUnix
}

func (my *AttunedMonth) determineIsLeap() {
	my.isLeap = my.y%4 == 0 && (my.y%100 != 0 || my.y%400 == 0)
}

func (my *AttunedMonth) determineUnixMonth() {
	d := int64(daysPerMonth[my.mon-1])
	if my.isLeap && my.mon > 2 {
		d++
	}
	my.unixMonth = my.unixYear + d*secondsPerDay
}

func (my *AttunedMonth) determineLastDay() {
	if my.isLeap && my.mon == 2 {
		my.lastDay = 29
		return
	}
	my.lastDay = daysPerMonth[my.mon] - daysPerMonth[my.mon-1]
}
