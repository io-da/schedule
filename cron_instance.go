package schedule

import (
	"time"
)

// CronInstance is created from CronExpression's and is used to hold state and determine next cron date according to it's expression.
type CronInstance struct {
	crn       *CronExpression
	following time.Time
	location  *time.Location

	ms    int
	s     int
	min   int
	h     int
	d     int
	lastD bool
	am    *AttunedMonth

	err error
}

// Next uses it's following date to determine the next valid cron date according to it's expression.
// Each subsequent execution advances the instance's following date.
func (crnI *CronInstance) Next() error {
	crnI.nextMs()
	if crnI.err != nil {
		return crnI.err
	}
	d := time.Date(crnI.am.Year(), crnI.am.Month(), crnI.d, crnI.h, crnI.min, crnI.s, crnI.ms*int(time.Millisecond), crnI.location)
	if !crnI.following.IsZero() && d.Before(crnI.following) {
		return CronOutdatedInvalidError
	}
	crnI.following = d
	return nil
}

// Following returns the following valid cron date determined by the Next function without modifying its state.
func (crnI *CronInstance) Following() time.Time {
	return crnI.following
}

func (crnI *CronInstance) nextMs() {
	fromMs := crnI.ms
	crnI.ms, _ = next(crnI.crn.milliseconds, fromMs, false)
	if crnI.nextS(crnI.ms > fromMs) {
		crnI.ms, _ = next(crnI.crn.milliseconds, 0, true)
	}
}

func (crnI *CronInstance) nextS(inc bool) bool {
	fromS := crnI.s
	crnI.s, _ = next(crnI.crn.seconds, fromS, inc)
	if crnI.nextMin(crnI.s > fromS || inc && crnI.s == fromS) {
		crnI.s, _ = next(crnI.crn.seconds, 0, true)
		return true
	}
	return crnI.s != fromS
}

func (crnI *CronInstance) nextMin(inc bool) bool {
	fromMin := crnI.min
	crnI.min, _ = next(crnI.crn.minutes, fromMin, inc)
	if crnI.nextH(crnI.min > fromMin || inc && crnI.min == fromMin) {
		crnI.min, _ = next(crnI.crn.minutes, 0, true)
		return true
	}
	return crnI.min != fromMin
}

func (crnI *CronInstance) nextH(inc bool) bool {
	fromH := crnI.h
	crnI.h, _ = next(crnI.crn.hours, fromH, inc)
	if crnI.nextD(crnI.h > fromH || inc && crnI.h == fromH) {
		crnI.h, _ = next(crnI.crn.hours, 0, true)
		return true
	}
	return crnI.h != fromH
}

func (crnI *CronInstance) nextD(inc bool) bool {
	fromD, fromY, last, invalid, reset := crnI.d, crnI.am.y, false, true, false
	if reset = crnI.nextMonY(true); reset {
		crnI.d, inc = 1, true
	}
	for invalid {
		crnI.d, last, invalid = crnI.nextDay(crnI.d, inc)
		inc = false
		if invalid && last {
			if reset = crnI.nextMonY(false); !reset || crnI.am.y-fromY > 400 {
				crnI.err = CronOutdatedInvalidError
				return false
			}
			crnI.d, inc = 1, true
		}
	}

	return reset || crnI.d != fromD
}

func (crnI *CronInstance) nextMonY(inc bool) bool {
	fromMon, fromY := crnI.am.mon, crnI.am.y
	nextMon, _ := next(crnI.crn.months, fromMon, inc)
	nextY, _ := next(crnI.crn.years, fromY, nextMon > fromMon || inc && nextMon == fromMon)
	if nextY != fromY {
		nextMon, _ := next(crnI.crn.months, 1, true)
		crnI.am.UpdateMonthYear(nextMon, nextY)
		return true
	}
	crnI.am.UpdateMonth(nextMon)
	return nextMon != fromMon
}

func (crnI *CronInstance) nextDay(fromD int, inc bool) (d int, last bool, invalid bool) {
	d, last = next(crnI.crn.days, fromD, inc)
	last = last || crnI.am.IsMonthLastDay(d)
	invalid = d < fromD || !inc && d == fromD || !crnI.am.Contains(d) || !crnI.crn.containsWeekday(crnI.am.WeekDay(d))
	return d, last, invalid
}
