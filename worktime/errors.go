package worktime

type StartOnWeekend struct {}
type EndOnWeekend struct {}
type StartAfterEnd struct {}

func (StartOnWeekend) Error() string {
	return "start time is on a weekend"
}

func (EndOnWeekend) Error() string {
	return "end time is on a weekend"
}

func (StartAfterEnd) Error() string {
	return "start time is after end time"
}
