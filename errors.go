package schedule

// ErrorOutdatedInvalidCron is used to represent a cron that has entered an invalid state or is outdated.
type ErrorOutdatedInvalidCron string

// CronOutdatedInvalidError is a constant equivalent of the ErrorOutdatedInvalidCron error.
const CronOutdatedInvalidError = ErrorOutdatedInvalidCron("schedule: outdated or invalid CronExpression")

// Error produces a string message of this error.
func (e ErrorOutdatedInvalidCron) Error() string {
	return string(e)
}

// ErrorOutdated is used to represent a schedule that became outdated.
type ErrorOutdated string

// OutdatedError is a constant equivalent of the ErrorOutdated error.
const OutdatedError = ErrorOutdated("schedule: outdated or invalid Schedule")

// Error produces a string message of this error.
func (e ErrorOutdated) Error() string {
	return string(e)
}
