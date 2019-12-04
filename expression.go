package schedule

// Expression is used as an empty interface, for readability purposes.
type Expression interface {
}

// IteratorExpression is used by expressions that represent a sequential set of values.
type IteratorExpression interface {
	Next(from int, inc bool) (int, bool)
	Contains(val int) bool
}
