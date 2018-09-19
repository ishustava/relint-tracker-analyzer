package worktime_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/ishustava/relint-tracker-analyzer/worktime"
	"time"
	. "github.com/onsi/gomega"
)

const layout = "Jan 2 2006 3:04 PM"

var _ = Describe("Time", func() {
	var sanFranciscoLocation *time.Location

	BeforeEach(func() {
		var err error
		sanFranciscoLocation, err = time.LoadLocation("America/Los_Angeles")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns duration when start time and end time are within working hours of the same day", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 14 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 12:00 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(2*time.Hour))
	})

	It("returns duration when start time and end time are on consecutive working days of the same week", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 13 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 12:00 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(12 * time.Hour))
	})

	It("returns duration when start time and end time are on consecutive working days of the same week but less than 24 hours apart", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 13 2018 11:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(9 * time.Hour))
	})

	It("returns duration when start time and end time are on consecutive working days of the same week but more than 24 hours apart", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 13 2018 11:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 12:00 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(11 * time.Hour))
	})

	It("returns duration when start time and end time are two days apart in the same week but less than 48 hours apart", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 12 2018 11:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(19 * time.Hour))
	})

	It("returns duration when start time and end time are two days apart in the same week but more than 48 hours apart", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 12 2018 11:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 12:00 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(21 * time.Hour))
	})

	It("returns duration when start time and end time are on non-consecutive working days of the same week", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 12 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 12:00 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(22*time.Hour))
	})

	It("returns duration when start time and end time are separated by one weekend", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 7 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 10 2018 12:00 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(12*time.Hour))
	})

	It("returns duration when start and end time are on the same day in PT, but not UTC", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 14 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 6:00 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(8*time.Hour))
	})

	It("returns duration when end time is past the end of the work day", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 14 2018 10:00 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 7:32 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(9 * time.Hour + 32 * time.Minute))
	})

	It("returns duration when start is before start of work day", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 14 2018 7:32 AM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 2:15 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(6 * time.Hour + 43 * time.Minute))
	})

	It("returns duration when start and end time are after the end of a work day", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 14 2018 6:15 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 7:32 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
		Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(1 * time.Hour + 17 * time.Minute))
	})

	It("returns duration when start is after the end of the work day and end time is during work hours another day", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 12 2018 6:15 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 13 2018 4:32 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		duration, err := worktime.Duration(startTime, endTime)
	 	Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(8 * time.Hour + 32 * time.Minute))
	})

	It("returns duration when start is after the end of the work day and end time is after work hours another day", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 12 2018 6:15 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 14 2018 7:32 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

	 	duration, err := worktime.Duration(startTime, endTime)
	 	Expect(err).NotTo(HaveOccurred())
		Expect(duration).To(Equal(21 * time.Hour + 32 * time.Minute))
	})

	It("errors when start is a weekend", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 15 2018 2:15 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 17 2018 12:32 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		_, err = worktime.Duration(startTime, endTime)
		Expect(err).To(MatchError("start time is on a weekend"))
	})

	It("errors when end is a weekend", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 12 2018 2:15 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 15 2018 12:32 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		_, err = worktime.Duration(startTime, endTime)
		Expect(err).To(MatchError("end time is on a weekend"))
	})

	It("errors when start is after end", func() {
		startTime, err := time.ParseInLocation(layout, "Sep 14 2018 7:32 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())
		endTime, err := time.ParseInLocation(layout, "Sep 12 2018 6:15 PM", sanFranciscoLocation)
		Expect(err).NotTo(HaveOccurred())

		_, err = worktime.Duration(startTime, endTime)
		Expect(err).To(MatchError("start time is after end time"))
	})
})
