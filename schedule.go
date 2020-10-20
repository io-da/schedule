package schedule

import "time"

// Schedule is the struct used to represent a set of retrievable time.Time structs.
type Schedule struct {
	at   []time.Time
	crn  *CronExpression
	crnI *CronInstance

	followingIndex int
}

// At creates a new schedule that produces the dates provided.
func At(at ...time.Time) *Schedule {
	if len(at) == 0 {
		panic("schedule: at least one time must be provided")
	}

	sch := &Schedule{
		at:             make([]time.Time, len(at)),
		followingIndex: -1,
	}
	currentTime := time.Time{}
	for i, t := range at {
		if t.Before(currentTime) || t.Equal(currentTime) {
			panic("schedule: time order provided is invalid")
		}
		sch.at[i] = t
		currentTime = t
	}

	return sch
}

// In creates a new schedule that produces dates based on provided durations.
// The durations are added sequentially.
// Example: In(time.Second, time.Minute):
// 		date = time.Now().Add(time.Second);
//		date = date.Add(time.Minute).
func In(in ...time.Duration) *Schedule {
	if len(in) == 0 {
		panic("schedule: at least one duration must be provided")
	}

	sch := &Schedule{
		at:             make([]time.Time, len(in)),
		followingIndex: -1,
	}

	currentTime := time.Now()
	for i, d := range in {
		currentTime = currentTime.Add(d)
		sch.at[i] = currentTime
	}

	return sch
}

// As creates a new schedule that produces dates based on the provided CronExpression.
// Example: As(Cron().EveryDay()):
// 		date = 00:00:00 of the following day;
// 		...
func As(crn *CronExpression) *Schedule {
	return &Schedule{
		crn:            crn,
		crnI:           crn.NewInstance(time.Now()),
	}
}

// AddCron is used to setup a CronExpression that starts operating after the scheduled times pass.
// Example: In(time.Hour * 24 * 7).AddCron(Cron().EveryDay()):
// 		date = time.Now().Add(time.Hour * 24 * 7);
//		date = 00:00:00 of the following day;
//		...
func (sch *Schedule) AddCron(crn *CronExpression) {
	sch.crn = crn
}

// Next is used to determine the following date to be produced.
func (sch *Schedule) Next() error {
	if sch.followingIndex < len(sch.at)-1 {
		sch.followingIndex++
		return nil
	}
	if sch.crn == nil {
		return OutdatedError
	}
	if sch.crnI == nil {
		sch.crnI = sch.crn.NewInstance(sch.at[sch.followingIndex])
		sch.followingIndex++
	}
	return sch.crnI.Next()
}

// Following returns the determined following date.
func (sch *Schedule) Following() time.Time {
	if sch.followingIndex < 0 {
		return time.Time{}
	}
	if sch.followingIndex < len(sch.at) {
		return sch.at[sch.followingIndex]
	}
	return sch.crnI.Following()
}
