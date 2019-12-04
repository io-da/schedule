package schedule

func validateMillisecond(t int) {
	if t < 0 || t > 999 {
		panic("schedule: invalid millisecond value")
	}
}

func validateSecond(t int) {
	if t < 0 || t > 59 {
		panic("schedule: invalid second value")
	}
}

func validateMinute(t int) {
	if t < 0 || t > 59 {
		panic("schedule: invalid minute value")
	}
}

func validateHour(t int) {
	if t < 0 || t > 23 {
		panic("schedule: invalid hour value")
	}
}

func validateDay(t int) {
	if t < 1 || t > 31 {
		panic("schedule: invalid day value")
	}
}

func validateWeekday(t int) {
	if t < 0 || t > 6 {
		panic("schedule: invalid weekday value")
	}
}

func validateMonth(t int) {
	if t < 1 || t > 12 {
		panic("schedule: invalid month value")
	}
}

func validateYear(t int) {
	if t < 1970 || t > 200000000 {
		panic("schedule: invalid year value")
	}
}
