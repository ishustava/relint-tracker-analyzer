package worktime

import (
	"time"
)

const (
	startHour = 8
	endHour = 18
	workingHoursInDay = 10
)
var local *time.Location

func init() {
	var err error
	local, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
}

func Duration(start, end time.Time) (time.Duration, error) {
	if !isWorkday(start) {
		return 0, new(StartOnWeekend)
	}

	if !isWorkday(end) {
		return 0, new(EndOnWeekend)
	}

	if start.After(end) {
		return 0, new(StartAfterEnd)
	}

	if start.Year() == end.Year() && start.YearDay() == end.YearDay() {
		return end.Sub(start), nil
	}

	calendarDays := calendarDays(start, end)
	fullWorkingDays := 0
	for day := 1; day < calendarDays; day++ {
		daysSinceStart := time.Duration(day) * 24 * time.Hour
		if isWorkday(start.Add(daysSinceStart)) {
			fullWorkingDays++
		}
	}

	workOnFirstDay := time.Duration(0)
	if start.Before(endOfWorkday(start)) {
		workOnFirstDay = endOfWorkday(start).Sub(start)
	}
	return workOnFirstDay + time.Duration(fullWorkingDays) * time.Duration(workingHoursInDay) * time.Hour + end.Sub(startOfWorkday(end)), nil
}

func isWorkday(day time.Time) bool {
	return day.Weekday() != time.Saturday && day.Weekday() != time.Sunday
}

func startOfWorkday(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), startHour, 0,0,0, local)
}

func endOfWorkday(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), endHour,0,0,0, local)
}

func calendarDays(start, end time.Time) int {
	return int(startOfWorkday(end).Sub(startOfWorkday(start)).Hours()/24)
}
