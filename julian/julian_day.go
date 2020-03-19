package julian

import (
	"math"
)

// JuliayDay
type JuliayDay float64

func (jd JuliayDay) MJD() MJD {
	return MJD(jd - 240000.5)
}

// MJD Modified Julian Day
type MJD float64

func (date *Date) JulianDay() JuliayDay {
	year := float64(date.Year)
	month := float64(date.Month)
	if date.Month < 3 {
		year -= 1
		month += 12
	}
	a := math.Floor(year / 100)
	b := math.Floor(2 - a + math.Floor(a/4))
	d := float64(date.Day) + date.OffsetOfDay
	if date.Calendar == JulianCalendar {
		b = 0
	}

	jd := math.Floor(365.25*(year+4716)) +
		math.Floor(30.6001*(month+1)) +
		d + b - 1524.5
	return JuliayDay(jd)
}
