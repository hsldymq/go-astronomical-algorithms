package astro_algo

import (
	"errors"
	"fmt"
)

var (
	JulianCalendar    = 0
	GregorianCalendar = 1
)

// JulianDate represents a Julian calendar date
type Date struct {
	Year        int
	Month       int
	Day         int
	OffsetOfDay float64
	Calendar    int
}

func NewDate(year, month, day int, offset ...float64) (*Date, error) {
	if year > 1582 {
		return newGregorianDate(year, month, day, offset...)
	} else if year < 1582 {
		return newJulianDate(year, month, day, offset...)
	} else if month < 10 || (month == 10 && day <= 4) {
		return newJulianDate(year, month, day, offset...)
	} else if month == 10 && day > 4 && day < 15 {
		return newGregorianDate(year, month, day+10, offset...)
	} else {
		return newGregorianDate(year, month, day, offset...)
	}
}

// newJulianDate creates a Date on Julian Calendar
// offset should be in [0, 1), default 0
func newJulianDate(year, month, day int, offset ...float64) (*Date, error) {
	offsetOfDay := 0.0
	if len(offset) > 0 {
		offsetOfDay = offset[0]
	}

	if err := validateDate(year, month, day, offsetOfDay, JulianCalendar); err != nil {
		return nil, err
	}

	return &Date{
		Year:        year,
		Month:       month,
		Day:         day,
		OffsetOfDay: offsetOfDay,
		Calendar:    JulianCalendar,
	}, nil
}

// newGregorianDate creates a Date on Gregorian calendar
// offset should be in [0, 1), default 0
func newGregorianDate(year, month, day int, offset ...float64) (*Date, error) {
	offsetOfDay := 0.0
	if len(offset) > 0 {
		offsetOfDay = offset[0]
	}

	if err := validateDate(year, month, day, offsetOfDay, GregorianCalendar); err != nil {
		return nil, err
	}

	return &Date{
		Year:        year,
		Month:       month,
		Day:         day,
		OffsetOfDay: offsetOfDay,
		Calendar:    GregorianCalendar,
	}, nil
}

func (d *Date) AddDays(days int) {
	daysSum := []int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
	dayNumber := daysSum[d.Month-1] + d.Day
	if isLeapYear(d.Year, d.Calendar) && d.Month > 3 {
		dayNumber += 1
	}

	if days > 0 {
		for days > 0 {
			if days >= 365 {
				days -= 365
				dayNumber += 365
			} else {
				dayNumber += days
				days = 0
			}

			leap := 0
			if isLeapYear(d.Year, d.Calendar) {
				leap = 1
			}
			if dayNumber > 365+leap {
				d.Year += 1
				dayNumber = dayNumber - (365 + leap)
			}
		}
	} else {
		for days < 0 {
			if days <= -365 {
				days += 365
				dayNumber -= 365
			} else {
				dayNumber += days
				days = 0
			}

			if dayNumber <= 0 {
				d.Year -= 1
				if isLeapYear(d.Year-1, d.Calendar) {
					dayNumber += 366
				} else {
					dayNumber += 365
				}
			}
		}
	}

	leap := 0
	if isLeapYear(d.Year, d.Calendar) {
		leap = 1
	}
	monthDays := []int{31, 28 + leap, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	month := 1
	for ; monthDays[month-1] < dayNumber; month++ {
		dayNumber = dayNumber - monthDays[month-1]
	}
	d.Month = month
	d.Day = dayNumber
}

func validateDate(year, month, day int, offset float64, calendar int) error {
	if offset < 0 || offset >= 1 {
		return errors.New("value of offset should be in [0,1)")
	}

	if month > 12 || month < 1 {
		return fmt.Errorf("invalid month: %d", month)
	}

	if day > 31 || day < 1 {
		return fmt.Errorf("invalid day: %d", day)
	}

	leap := 0
	if isLeapYear(year, calendar) {
		leap = 1
	}

	monthDays := []int{31, 28 + leap, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if day < 1 || day > monthDays[month-1] {
		return fmt.Errorf("invalid day: %d", day)
	}

	return nil
}

func isLeapYear(year, calendar int) bool {
	if year%100 == 0 {
		if calendar == JulianCalendar || year%400 == 0 {
			return true
		}
	} else if year%4 == 0 {
		return true
	}

	return false
}
