package schedule

import "time"

// BetweenExpression is the struct used to create cron between expressions.
type BetweenExpression struct {
	x    int
	y    int
	step int
}

// Between is an expression that generates integers between the provided parameters (*inclusive*).
func Between(x int, y int) *BetweenExpression {
	if x > y {
		panic("schedule: invalid BetweenExpression expression")
	}

	return &BetweenExpression{
		x:    x,
		y:    y,
		step: 1,
	}
}

// BetweenMilliseconds uses the regular between logic, ensuring valid millisecond parameters.
func BetweenMilliseconds(x int, y int) *BetweenExpression {
	validateMillisecond(x)
	validateMillisecond(y)
	return Between(x, y)
}

// BetweenSeconds uses the regular between logic, ensuring valid second parameters.
func BetweenSeconds(x int, y int) *BetweenExpression {
	validateSecond(x)
	validateSecond(y)
	return Between(x, y)
}

// BetweenMinutes uses the regular between logic, ensuring valid minute parameters.
func BetweenMinutes(x int, y int) *BetweenExpression {
	validateMinute(x)
	validateMinute(y)
	return Between(x, y)
}

// BetweenHours uses the regular between logic, ensuring valid hour parameters.
func BetweenHours(x int, y int) *BetweenExpression {
	validateHour(x)
	validateHour(y)
	return Between(x, y)
}

// BetweenDays uses the regular between logic, ensuring valid day parameters.
func BetweenDays(x int, y int) *BetweenExpression {
	validateDay(x)
	validateDay(y)
	return Between(x, y)
}

// BetweenWeekdays uses the regular between logic, ensuring valid time.Weekday parameters.
func BetweenWeekdays(x time.Weekday, y time.Weekday) *BetweenExpression {
	xI := int(x)
	yI := int(y)
	validateWeekday(xI)
	validateWeekday(yI)
	return Between(xI, yI)
}

// BetweenMonths uses the regular between logic, ensuring valid time.Month parameters.
func BetweenMonths(x time.Month, y time.Month) *BetweenExpression {
	xI := int(x)
	yI := int(y)
	validateMonth(xI)
	validateMonth(yI)
	return Between(xI, yI)
}

// BetweenYears uses the regular between logic, ensuring valid year parameters.
func BetweenYears(x int, y int) *BetweenExpression {
	validateYear(x)
	validateYear(y)
	return Between(x, y)
}

// Every allows optional specification of the stepping used for the between logic.
// It's important to understand the behavior of the expression when step > 1. It may produce some unexpected values.
// Example: Between(0,10).Every(3)
//		- Next(-1, true || false)  = 0
//		- Next(0, true)   = 0
//		- Next(0, false)  = 3
//		- Next(1, true || false)  = 3
//		- Next(3, false)  = 6
//		- Next(6, false)  = 9
//		- Next(9, false)  = 10
//		- Next(10, true || false) = 10
func (exp *BetweenExpression) Every(s int) *BetweenExpression {
	if s < 1 {
		panic("schedule: invalid step value")
	}
	exp.step = s
	return exp
}

// Next allows retrieval of the next value from this expression.
// Expressions are stateless, the determination of their next value is based on input.
// Given a valid expression value, the parameter inc is used to specify if it should be included in the output.
// Given the last value of the expression or above, the inc parameter is ignored.
// It returns the next value according to provided parameters and a boolean indicating if it is the last value.
func (exp *BetweenExpression) Next(from int, inc bool) (int, bool) {
	if from < exp.x || inc && from == exp.x {
		return exp.x, false
	}

	if from >= exp.y {
		return exp.y, true
	}

	diff := exp.step - (from-exp.x)%exp.step
	if inc && diff == exp.step {
		diff = 0
	}

	next := from + diff
	if next >= exp.y {
		return exp.y, true
	}
	return next, false
}

// Contains verifies if the provided value belongs to this expression.
func (exp *BetweenExpression) Contains(val int) bool {
	if val < exp.x || val > exp.y {
		return false
	}

	if val == exp.x || val == exp.y {
		return true
	}

	return (exp.step - (val-exp.x)%exp.step) == exp.step
}
