package schedule

import (
	"sync/atomic"
	"time"
)

// CronExpression is used to represent the complete cron expression.
// It is used to create new CronInstances.
type CronExpression struct {
	milliseconds Expression
	seconds      Expression
	minutes      Expression
	hours        Expression
	days         Expression
	weekdays     Expression
	months       Expression
	years        Expression

	initialized *uint32
}

// Cron creates and returns a reference to a new CronExpression.
func Cron() *CronExpression {
	return &CronExpression{
		initialized: new(uint32),
	}
}

//------Public Functions------//

// NewInstance creates and returns a reference to a new CronInstance for the referenced CronExpression.
func (crn *CronExpression) NewInstance(from time.Time) *CronInstance {
	crn.initialize()
	return &CronInstance{
		crn:       crn,
		following: from,
		location:  from.Location(),

		ms:  from.Nanosecond() / int(time.Millisecond),
		s:   from.Second(),
		min: from.Minute(),
		h:   from.Hour(),
		d:   from.Day(),
		am:  NewAttunedMonth(int(from.Month()), from.Year()),
	}
}

//------Expression------//

// EveryMillisecond sets this expression to return a date for every millisecond.
func (crn *CronExpression) EveryMillisecond() *CronExpression {
	return crn.OnMilliseconds(Between(0, 999))
}

// OnMilliseconds sets this expression to return a date on the provided millisecond.
func (crn *CronExpression) OnMilliseconds(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateMillisecond(exp)
	}
	crn.milliseconds = exp
	crn.reset()
	return crn
}

// EverySecond sets this expression to return a date for every second.
func (crn *CronExpression) EverySecond() *CronExpression {
	return crn.OnSeconds(Between(0, 59))
}

// OnSeconds sets this expression to return a date on the provided seconds.
func (crn *CronExpression) OnSeconds(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateSecond(exp)
	}
	crn.seconds = exp
	crn.handleMillisecond()
	crn.reset()
	return crn
}

// EveryMinute sets this expression to return a date for every minute.
func (crn *CronExpression) EveryMinute() *CronExpression {
	return crn.OnMinutes(Between(0, 59))
}

// OnMinutes sets this expression to return a date on the provided minutes.
func (crn *CronExpression) OnMinutes(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateMinute(exp)
	}
	crn.minutes = exp
	crn.handleSecond()
	crn.reset()
	return crn
}

// EveryHour sets this expression to return a date for every hour.
func (crn *CronExpression) EveryHour() *CronExpression {
	return crn.OnHours(Between(0, 23))
}

// OnHours sets this expression to return a date on the provided hours.
func (crn *CronExpression) OnHours(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateHour(exp)
	}
	crn.hours = exp
	crn.handleMinute()
	crn.reset()
	return crn
}

// EveryDay sets this expression to return a date for every day.
func (crn *CronExpression) EveryDay() *CronExpression {
	return crn.OnDays(Between(1, 31))
}

// OnDays sets this expression to return a date on the provided days.
func (crn *CronExpression) OnDays(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateDay(exp)
	}
	crn.days = exp
	crn.handleHour()
	crn.reset()
	return crn
}

// OnWeekdays sets this expression to return a date on the provided weekdays.
func (crn *CronExpression) OnWeekdays(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateWeekday(exp)
	}
	crn.weekdays = exp
	crn.handleHour()
	crn.reset()
	return crn
}

// EveryMonth sets this expression to return a date for every month.
func (crn *CronExpression) EveryMonth() *CronExpression {
	return crn.OnMonths(BetweenMonths(time.January, time.December))
}

// OnMonths sets this expression to return a date on the provided months.
func (crn *CronExpression) OnMonths(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateMonth(exp)
	}
	crn.months = exp
	crn.handleDay()
	crn.reset()
	return crn
}

// OnYears sets this expression to return a date on the provided years.
func (crn *CronExpression) OnYears(exp Expression) *CronExpression {
	if exp, iOf := exp.(int); iOf {
		validateYear(exp)
	}
	crn.years = exp
	crn.handleMonth()
	crn.reset()
	return crn
}

//------Initialization------//

func (crn *CronExpression) initialize() {
	if atomic.CompareAndSwapUint32(crn.initialized, 0, 1) {
		crn.ensureSeconds()
		crn.ensureMinutes()
		crn.ensureHours()
		crn.ensureDays()
		crn.ensureWeekdays()
		crn.ensureMonths()
		crn.ensureYears()
	}
}

func (crn *CronExpression) reset() {
	atomic.CompareAndSwapUint32(crn.initialized, 1, 0)
}

func (crn *CronExpression) ensureSeconds() {
	if crn.seconds == nil {
		crn.seconds = BetweenSeconds(0, 59)
	}
}

func (crn *CronExpression) ensureMinutes() {
	if crn.minutes == nil {
		crn.minutes = BetweenMinutes(0, 59)
	}
}

func (crn *CronExpression) ensureHours() {
	if crn.hours == nil {
		crn.hours = BetweenHours(0, 23)
	}
}

func (crn *CronExpression) ensureDays() {
	if crn.days == nil {
		crn.days = BetweenDays(1, 31)
	}
}

func (crn *CronExpression) ensureWeekdays() {
	if crn.weekdays == nil {
		crn.weekdays = BetweenWeekdays(time.Sunday, time.Saturday)
	}
}

func (crn *CronExpression) ensureMonths() {
	if crn.months == nil {
		crn.months = BetweenMonths(time.January, time.December)
	}
}

func (crn *CronExpression) ensureYears() {
	if crn.years == nil {
		crn.years = BetweenYears(1970, 200000000)
	}
}

//------Utils------//

func (crn *CronExpression) handleMillisecond() {
	if crn.milliseconds == nil {
		crn.milliseconds = 0
	}
}

func (crn *CronExpression) handleSecond() {
	if crn.seconds == nil {
		crn.seconds = 0
	}
	crn.handleMillisecond()
}

func (crn *CronExpression) handleMinute() {
	if crn.minutes == nil {
		crn.minutes = 0
	}
	crn.handleSecond()
}

func (crn *CronExpression) handleHour() {
	if crn.hours == nil {
		crn.hours = 0
	}
	crn.handleMinute()
}

func (crn *CronExpression) handleDay() {
	if crn.days == nil {
		crn.days = 1
	}
	crn.handleHour()
}

func (crn *CronExpression) handleMonth() {
	if crn.months == nil {
		crn.months = 1
	}
	crn.handleDay()
}

func next(exp Expression, from int, inc bool) (int, bool) {
	switch exp := exp.(type) {
	case IteratorExpression:
		return exp.Next(from, inc)
	case time.Month:
		return int(exp), true
	}
	return exp.(int), true
}

func (crn *CronExpression) containsWeekday(wd time.Weekday) bool {
	if weekdays, iOf := crn.weekdays.(IteratorExpression); iOf {
		return weekdays.Contains(int(wd))
	}
	return crn.weekdays.(time.Weekday) == wd
}
