package schedule

import (
	"sort"
	"time"
)

// ListExpression is the struct used to create cron list expressions.
type ListExpression struct {
	values []int
}

// List is an expression used to iterate the provided list of int parameters (*inclusive*).
func List(values []int) *ListExpression {
	if len(values) < 1 {
		panic("schedule: invalid ListExpression expression")
	}
	sort.Ints(values)
	return &ListExpression{
		values: values,
	}
}

// ListMilliseconds uses the regular list logic, ensuring valid millisecond parameters.
func ListMilliseconds(values ...int) *ListExpression {
	for _, v := range values {
		validateMillisecond(v)
	}
	return List(values)
}

// ListSeconds uses the regular list logic, ensuring valid second parameters.
func ListSeconds(values ...int) *ListExpression {
	for _, v := range values {
		validateSecond(v)
	}
	return List(values)
}

// ListMinutes uses the regular list logic, ensuring valid minute parameters.
func ListMinutes(values ...int) *ListExpression {
	for _, v := range values {
		validateMinute(v)
	}
	return List(values)
}

// ListHours uses the regular list logic, ensuring valid hour parameters.
func ListHours(values ...int) *ListExpression {
	for _, v := range values {
		validateHour(v)
	}
	return List(values)
}

// ListDays uses the regular list logic, ensuring valid day parameters.
func ListDays(values ...int) *ListExpression {
	for _, v := range values {
		validateDay(v)
	}
	return List(values)
}

// ListWeekdays uses the regular list logic, ensuring valid time.Weekday parameters.
func ListWeekdays(values ...time.Weekday) *ListExpression {
	converted := make([]int, len(values))
	for i, v := range values {
		vI := int(v)
		validateWeekday(vI)
		converted[i] = vI
	}
	return List(converted)
}

// ListMonths uses the regular list logic, ensuring valid time.Month parameters.
func ListMonths(values ...time.Month) *ListExpression {
	converted := make([]int, len(values))
	for i, v := range values {
		vI := int(v)
		validateMonth(vI)
		converted[i] = vI
	}
	return List(converted)
}

// ListYears uses the regular list logic, ensuring valid year parameters.
func ListYears(values ...int) *ListExpression {
	for _, v := range values {
		validateYear(v)
	}
	return List(values)
}

// Next allows retrieval of the next value from this expression.
// Expressions are stateless, the determination of their next value is based on input.
// Given a valid expression value, the parameter inc is used to specify if it should be included in the output.
// Given the last value of the expression or above, the inc parameter is ignored.
// It returns the next value according to provided parameters and a boolean indicating if it is the last value.
func (exp *ListExpression) Next(from int, inc bool) (int, bool) {
	lastI := len(exp.values) - 1
	lastV := exp.values[lastI]
	if from < lastV {
		for i, v := range exp.values {
			if from < v || (inc && from == v) {
				return v, i == lastI
			}
		}
	}
	return lastV, true
}

// Contains verifies if the provided value belongs to this expression
func (exp *ListExpression) Contains(val int) bool {
	firstV := exp.values[0]
	lastV := exp.values[len(exp.values)-1]
	if val < firstV || val > lastV {
		return false
	}
	for _, v := range exp.values {
		if val == v {
			return true
		}
	}
	return false
}
