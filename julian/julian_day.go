package julian

import (
	"math"
)

// J2000 is the julian day for 1 Jan,2000 12:00:00
const J2000 JulianDay = 2451545.0

// JulianDay julian day
type JulianDay float64

func (jd JulianDay) MJD() MJD {
	return MJD(jd - 2400000.5)
}

// MJD Modified Julian Day
type MJD float64

func (mjd MJD) JulianDay() JulianDay {
	return JulianDay(mjd + 2400000.5)
}

func (date *Date) JulianDay() JulianDay {
	year, month := float64(date.Year), float64(date.Month)
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
	return JulianDay(jd)
}
