package schedule

import (
	"testing"
	"time"
)

type testDate struct {
	time.Time
}

func BenchmarkCronInstance_Next(b *testing.B) {
	crons := [...]*CronExpression{
		Cron().OnMinutes(23).OnHours(Between(0, 20).Every(2)),
		Cron().OnMinutes(5).OnHours(4).OnWeekdays(time.Sunday),
		Cron().OnMinutes(5).EveryDay().OnMonths(time.August),
		Cron().OnMinutes(15).OnHours(14).OnDays(1),
		Cron().OnHours(22).OnWeekdays(BetweenWeekdays(time.Monday, time.Friday)),
		Cron().OnHours(ListHours(0, 12)).OnMonths(BetweenMonths(time.February, time.December).Every(2)),
		Cron().OnDays(ListDays(1, 15)).OnWeekdays(time.Wednesday),
		Cron().OnMonths(time.January),
		Cron().EveryMonth(),
		Cron().OnWeekdays(time.Sunday),
		Cron().EveryDay(),
		Cron().EveryHour(),
		Cron().EveryMinute(),
		Cron().EverySecond(),
		Cron().EveryMillisecond(),
	}
	var crnI *CronInstance
	startDate := date().Time
	for n := 0; n < b.N; n++ {
		crnI = crons[n%len(crons)].NewInstance(startDate)
		if err := crnI.Next(); err != nil {
			b.Error(err.Error())
		}
	}
}

func TestCronInstance_Next(t *testing.T) {
	var crnI *CronInstance
	var startDate = date().Time
	var expectedAt time.Time

	// At 00:00 on day-of-month 29 and on Sunday in February. (0 0 29 2 0)
	crnI = Cron().
		OnMonths(time.February).
		OnDays(29).
		OnWeekdays(time.Sunday).NewInstance(startDate)
	expectedAt = date().setYear(2032).setMonth(time.February).setDay(29).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setYear(2060).setMonth(time.February).setDay(29).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At minute 23 past every 2nd hour from 0 through 20. (23 0-20/2 * * *)
	crnI = Cron().
		OnMinutes(23).
		OnHours(Between(0, 20).Every(2)).NewInstance(startDate)
	expectedAt = date().setMinute(23).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMinute(23).setHour(20).Time
	if crnI.advanceX(t, 10) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMinute(23).setDay(2).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 04:05 on Sunday. (5 4 * * sun)
	crnI = Cron().
		OnMinutes(5).
		OnHours(4).
		OnWeekdays(time.Sunday).NewInstance(startDate)
	expectedAt = date().setMinute(5).setHour(4).setDay(6).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMinute(5).setHour(4).setDay(13).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMinute(5).setHour(4).setDay(3).setMonth(time.February).Time
	if crnI.advanceX(t, 3) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 00:05 in August. (5 0 * 8 *)
	crnI = Cron().
		OnMinutes(5).EveryDay().
		OnMonths(time.August).NewInstance(startDate)
	expectedAt = date().setMinute(5).setMonth(time.August).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMinute(5).setDay(10).setMonth(time.August).Time
	if crnI.advanceX(t, 9) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMinute(5).setMonth(time.August).setYear(2020).Time
	if crnI.advanceX(t, 22) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 14:15 on day-of-month 1. (15 14 1 * *)
	crnI = Cron().
		OnMinutes(15).
		OnHours(14).
		OnDays(1).NewInstance(startDate)
	expectedAt = date().setMinute(15).setHour(14).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMinute(15).setHour(14).setMonth(time.February).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 22:00 on every day-of-week from Monday through Friday. (0 22 * * 1-5)
	crnI = Cron().
		OnHours(22).
		OnWeekdays(BetweenWeekdays(time.Monday, time.Friday)).NewInstance(startDate)
	expectedAt = date().setHour(22).setDay(4).Time
	if crnI.advanceX(t, 4) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setHour(22).setDay(7).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At minute 0 past hour 0 and 12 on day-of-month 1 in every 2nd month. (0 0,12 1 */2 *)
	crnI = Cron().
		OnHours(ListHours(0, 12)).
		OnMonths(BetweenMonths(time.February, time.December).Every(2)).NewInstance(startDate)
	expectedAt = date().setMonth(time.February).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setHour(12).setMonth(time.February).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMonth(time.April).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMonth(time.February).setYear(2020).Time
	if crnI.advanceX(t, 10) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 00:00 on day-of-month 1 and 15 and on Wednesday. (0 0 1,15 * 3)
	crnI = Cron().
		OnDays(ListDays(1, 15)).
		OnWeekdays(time.Wednesday).NewInstance(startDate)
	expectedAt = date().setMonth(time.May).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(15).setMonth(time.May).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setYear(2020).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 00:00 on day-of-month 1 in January (@annually).
	crnI = Cron().OnMonths(time.January).NewInstance(startDate)
	expectedAt = date().setYear(2020).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setYear(2021).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 00:00 on day-of-month 1 (@monthly).
	crnI = Cron().EveryMonth().NewInstance(startDate)
	expectedAt = date().setMonth(time.February).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMonth(time.March).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	// Or
	crnI = Cron().OnDays(1).NewInstance(startDate)
	expectedAt = date().setMonth(time.February).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMonth(time.March).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 00:00 on Sunday (@weekly).
	crnI = Cron().OnWeekdays(time.Sunday).NewInstance(startDate)
	expectedAt = date().setDay(6).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(13).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At 00:00 (@daily).
	crnI = Cron().EveryDay().NewInstance(startDate)
	expectedAt = date().setDay(2).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(3).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(31).Time
	if crnI.advanceX(t, 28) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(1).setMonth(time.March).Time
	if crnI.advanceX(t, 29) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(29).setYear(2020).setMonth(time.February).Time
	if crnI.advanceX(t, 365) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	// Or
	crnI = Cron().OnHours(0).NewInstance(startDate)
	expectedAt = date().setDay(2).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(3).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(31).Time
	if crnI.advanceX(t, 28) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(1).setMonth(time.March).Time
	if crnI.advanceX(t, 29) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(29).setYear(2020).setMonth(time.February).Time
	if crnI.advanceX(t, 365) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// At minute 0 (@hourly).
	crnI = Cron().EveryHour().NewInstance(startDate)
	expectedAt = date().setHour(1).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setHour(2).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(2).Time
	if crnI.advanceX(t, 22) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	// Or
	crnI = Cron().OnMinutes(0).NewInstance(startDate)
	expectedAt = date().setHour(1).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setHour(2).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setDay(2).Time
	if crnI.advanceX(t, 22) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// additional tests for code coverage
	crnI = Cron().
		OnMilliseconds(BetweenMilliseconds(499, 501)).
		OnSeconds(BetweenSeconds(29, 31)).
		OnMinutes(BetweenMinutes(29, 31)).NewInstance(startDate)
	expectedAt = date().setMillisecond(501).setSecond(29).setMinute(29).Time
	if crnI.advanceX(t, 3) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMillisecond(499).setSecond(30).setMinute(29).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMillisecond(499).setSecond(29).setMinute(30).Time
	if crnI.advanceX(t, 6) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	crnI = Cron().OnMilliseconds(BetweenMilliseconds(0, 750).Every(250)).NewInstance(startDate)
	expectedAt = date().setMillisecond(250).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMillisecond(750).Time
	if crnI.advanceX(t, 2) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setSecond(1).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	crnI = Cron().
		EveryMillisecond().
		EverySecond().
		EveryMinute().NewInstance(startDate)
	expectedAt = date().setMillisecond(999).Time
	if crnI.advanceX(t, 999) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMillisecond(0).setSecond(1).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	expectedAt = date().setMillisecond(0).setSecond(0).setMinute(1).Time
	if crnI.advanceX(t, 59000) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}
	crnI = Cron().
		OnMilliseconds(ListMilliseconds(0, 500)).
		OnSeconds(ListSeconds(0, 30)).
		OnMinutes(ListMinutes(0, 30)).
		OnWeekdays(ListWeekdays(time.Monday, time.Friday)).EveryDay().
		OnMonths(ListMonths(time.February, time.December)).
		OnYears(ListYears(2020, 2024)).NewInstance(startDate)
	expectedAt = date().setYear(2020).setMonth(time.February).setDay(3).Time
	if crnI.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected CronExpression date returned.")
	}

	// Test outdated crnI
	crnI = Cron().
		OnYears(2018).
		OnMonths(time.February).
		OnDays(28).NewInstance(startDate)
	err := crnI.Next()
	if err == nil || err != CronOutdatedInvalidError {
		t.Error("Unexpected CronExpression behavior.")
	} else if err.Error() != "schedule: outdated or invalid CronExpression" {
		t.Error("Unexpected ErrorOutdatedInvalidCron message.")
	}

	// Test invalid crnI
	crnI = Cron().
		OnYears(2018).
		OnMonths(time.February).
		OnDays(29).NewInstance(startDate)
	if err := crnI.Next(); err != CronOutdatedInvalidError {
		t.Error("Unexpected CronExpression behavior.")
	}
	crnI = Cron().
		OnMonths(time.February).
		OnDays(30).NewInstance(startDate)
	if err := crnI.Next(); err != CronOutdatedInvalidError {
		t.Error("Unexpected CronExpression behavior.")
	}
}

func TestCronExpression_OnMillisecondsPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid millisecond value")
	Cron().OnMilliseconds(-1)
}

func TestCronExpression_OnSecondsPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid second value")
	Cron().OnSeconds(-1)
}

func TestCronExpression_OnMinutesPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid minute value")
	Cron().OnMinutes(-1)
}

func TestCronExpression_OnHoursPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid hour value")
	Cron().OnHours(-1)
}

func TestCronExpression_OnDaysPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid day value")
	Cron().OnDays(-1)
}

func TestCronExpression_OnWeekdaysPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid weekday value")
	Cron().OnWeekdays(-1)
}

func TestCronExpression_OnMonthsPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid month value")
	Cron().OnMonths(-1)
}

func TestCronExpression_OnYearsPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid year value")
	Cron().OnYears(-1)
}

func TestBetween_Contains(t *testing.T) {
	exp := Between(0, 10).Every(3)
	for _, v := range []int{0, 3, 6, 9, 10} {
		if !exp.Contains(v) {
			t.Errorf("Expected Between contains %d.", v)
		}
	}
	for _, v := range []int{-1, 1, 2, 4, 5, 7, 8} {
		if exp.Contains(v) {
			t.Errorf("Unexpected Between contains %d.", v)
		}
	}
}

func TestBetweenPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid BetweenExpression expression")
	Between(1, 0)
}

func TestBetween_EveryPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid step value")
	Between(0, 1).Every(0)
}

func TestList_Contains(t *testing.T) {
	exp := List([]int{0, 3, 6, 9, 10})
	for _, v := range []int{0, 3, 6, 9, 10} {
		if !exp.Contains(v) {
			t.Errorf("Expected List contains %d.", v)
		}
	}
	for _, v := range []int{-1, 1, 2, 4, 5, 7, 8} {
		if exp.Contains(v) {
			t.Errorf("Unexpected List contains %d.", v)
		}
	}
}

func TestListPanic(t *testing.T) {
	defer ensurePanic(t, "schedule: invalid ListExpression expression")
	List([]int{})
}

func TestMonthYear_MonthLastDay(t *testing.T) {
	my := NewAttunedMonth(2, 2020)
	if my.MonthLastDay() != 29 {
		t.Errorf("Unexpected NewAttunedMonth behavior.")
	}
}

func TestAt(t *testing.T) {
	dt1 := date().Time
	dt2 := date().setDay(2).Time

	sch := At(dt1, dt2)
	if sch.at[0] != dt1 || sch.at[1] != dt2 {
		t.Error("Unexpected Schedule behavior.")
	}
}

func TestIn(t *testing.T) {
	dt := time.Now()
	sch := In(time.Second, time.Millisecond)
	dt1 := dt.Add(time.Second)
	dt2 := dt1.Add(time.Millisecond)
	if toMilliseconds(sch.at[0]) != toMilliseconds(dt1) || toMilliseconds(sch.at[1]) != toMilliseconds(dt2) {
		t.Error("Unexpected Schedule behavior.")
	}
}

func TestSchedule_Next(t *testing.T) {
	dt1 := date().Time
	dt2 := date().setYear(2020).setMonth(2).setDay(28).Time
	sch := At(dt1, dt2)
	sch.Cron(Cron().EveryDay())
	if !sch.Following().IsZero() {
		t.Error("Unexpected Schedule date returned.")
	}
	expectedAt := date().Time
	if sch.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected Schedule date returned.")
	}
	expectedAt = date().setYear(2020).setMonth(2).setDay(28).Time
	if sch.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected Schedule date returned.")
	}
	expectedAt = date().setYear(2020).setMonth(2).setDay(29).Time
	if sch.advanceX(t, 1) != expectedAt {
		t.Error("Unexpected Schedule date returned.")
	}
	expectedAt = date().setYear(2020).setMonth(4).Time
	if sch.advanceX(t, 32) != expectedAt {
		t.Error("Unexpected Schedule date returned.")
	}

	sch = In(time.Hour, time.Hour*6)
	sch.Cron(Cron().EveryHour())
	expectedAt = time.Now().Add(time.Hour)
	if toMilliseconds(sch.advanceX(t, 1)) != toMilliseconds(expectedAt) {
		t.Error("Unexpected Schedule date returned.")
	}
	expectedAt = expectedAt.Add(time.Hour * 6)
	if toMilliseconds(sch.advanceX(t, 1)) != toMilliseconds(expectedAt) {
		t.Error("Unexpected Schedule date returned.")
	}

	expectedAt = expectedAt.Add(time.Hour)
	expectedAt = date().overwrite(expectedAt).setMinute(0).setSecond(0).setMillisecond(0).Time
	if toMilliseconds(sch.advanceX(t, 1)) != toMilliseconds(expectedAt) {
		t.Error("Unexpected Schedule date returned.")
	}
	expectedAt = expectedAt.Add(time.Hour * 1337)
	if toMilliseconds(sch.advanceX(t, 1337)) != toMilliseconds(expectedAt) {
		t.Error("Unexpected Schedule date returned.")
	}

	sch = In(time.Hour)
	expectedAt = time.Now().Add(time.Hour)
	if toMilliseconds(sch.advanceX(t, 1)) != toMilliseconds(expectedAt) {
		t.Error("Unexpected Schedule date returned.")
	}

	err := sch.Next()
	if err == nil || err != OutdatedError {
		t.Error("Unexpected Schedule behavior.")
	} else if err.Error() != "schedule: outdated or invalid Schedule" {
		t.Error("Unexpected ErrorOutdated message.")
	}
}

func TestAtPanicAtLeastOneTimeRequired(t *testing.T) {
	defer ensurePanic(t, "schedule: at least one time must be provided")
	At()
}

func TestAtPanicInvalidTimeOrder(t *testing.T) {
	defer ensurePanic(t, "schedule: time order provided is invalid")
	At(date().setMillisecond(1).Time, date().Time)
}

func TestInPanicAtLeastOneDurationRequired(t *testing.T) {
	defer ensurePanic(t, "schedule: at least one duration must be provided")
	In()
}

func date() *testDate {
	return &testDate{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func (crnI *CronInstance) advanceX(t *testing.T, times int) time.Time {
	for x := 0; x < times; x++ {
		if err := crnI.Next(); err != nil {
			t.Error(err.Error())
			break
		}
	}
	return crnI.Following()
}

func (sch *Schedule) advanceX(t *testing.T, times int) time.Time {
	for x := 0; x < times; x++ {
		if err := sch.Next(); err != nil {
			t.Error(err.Error())
			break
		}
	}
	return sch.Following()
}

func (date *testDate) setYear(y int) *testDate {
	date.Time = time.Date(y, date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
	return date
}

func (date *testDate) setMonth(mon time.Month) *testDate {
	date.Time = time.Date(date.Year(), mon, date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
	return date
}

func (date *testDate) setDay(d int) *testDate {
	date.Time = time.Date(date.Year(), date.Month(), d, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
	return date
}

func (date *testDate) setHour(h int) *testDate {
	date.Time = time.Date(date.Year(), date.Month(), date.Day(), h, date.Minute(), date.Second(), date.Nanosecond(), date.Location())
	return date
}

func (date *testDate) setMinute(min int) *testDate {
	date.Time = time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), min, date.Second(), date.Nanosecond(), date.Location())
	return date
}

func (date *testDate) setSecond(s int) *testDate {
	date.Time = time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), s, date.Nanosecond(), date.Location())
	return date
}

func (date *testDate) setMillisecond(ms int) *testDate {
	date.Time = time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), ms*int(time.Millisecond), date.Location())
	return date
}

func (date *testDate) overwrite(t time.Time) *testDate {
	date.Time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	return date
}

func toMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func ensurePanic(t *testing.T, p string) {
	if r := recover(); r != p {
		t.Errorf("Expected panic: %s", p)
	}
}
