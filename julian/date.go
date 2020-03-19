package julian

import (
	"errors"
	"fmt"
	"time"
)

var (
	JulianCalendar    = 0
	GregorianCalendar = 1
)

// DayOffset 计算该时刻相对于一天开始的偏移
// 例:
//		00:00:00, 应返回0
//		04:48:00,即走过了17280秒,即86400秒的1/5, 所以应返回0.2
//		以此类推
func DayOffset(hour, minute, second int) float64 {
	seconds := float64(hour*3600 + 60*minute + second)
	return seconds / float64(86400)
}

// Date represents a date in Julian calendar or Gregorian calendar
type Date struct {
	Year        int
	Month       int
	Day         int
	OffsetOfDay float64
	Calendar    int
}

func NewDate(year, month, day int, offset ...float64) (*Date, error) {
	ymd := year*10000 + month*100 + day
	if ymd <= 15821004 {
		return newJulianDate(year, month, day, offset...)
	} else if ymd >= 15821015 {
		return newGregorianDate(year, month, day, offset...)
	} else {
		return nil, fmt.Errorf("the date %d %s,%d does not exist", day, time.Month(month).String()[:3], year)
	}
}

func NewDateFromTime(t *time.Time) (*Date, error) {
	year, month, day := t.Year(), int(t.Month()), t.Day()

	if year*10000+month*100+day < 15821015 {
		return nil, errors.New("not support dates before 15 Oct,1582")
	}

	return newGregorianDate(year, month, day, DayOffset(t.Hour(), t.Minute(), t.Second()))
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

func (date *Date) AddDays(days int) {
	daysSum := []int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
	dayNumber := daysSum[date.Month-1] + date.Day
	if isLeapYear(date.Year, date.Calendar) && date.Month > 3 {
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
			if isLeapYear(date.Year, date.Calendar) {
				leap = 1
			}
			if dayNumber > 365+leap {
				date.Year += 1
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
				date.Year -= 1
				if isLeapYear(date.Year-1, date.Calendar) {
					dayNumber += 366
				} else {
					dayNumber += 365
				}
			}
		}
	}

	leap := 0
	if isLeapYear(date.Year, date.Calendar) {
		leap = 1
	}
	monthDays := []int{31, 28 + leap, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	month := 1
	for ; monthDays[month-1] < dayNumber; month++ {
		dayNumber = dayNumber - monthDays[month-1]
	}
	date.Month = month
	date.Day = dayNumber
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
