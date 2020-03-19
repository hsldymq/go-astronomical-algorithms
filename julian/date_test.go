package julian

import (
	"fmt"
	"testing"
)

func TestDayOffset(t *testing.T) {
	cases := []struct {
		Hour   int
		Minute int
		Second int
		Expect float64
	}{
		{Hour: 0, Minute: 0, Second: 0, Expect: 0},
		{Hour: 4, Minute: 48, Second: 0, Expect: 0.2},
		{Hour: 9, Minute: 36, Second: 0, Expect: 0.4},
		{Hour: 12, Minute: 0, Second: 0, Expect: 0.5},
	}

	for _, each := range cases {
		actual := DayOffset(each.Hour, each.Minute, each.Second)
		if actual != each.Expect {
			t.Fatalf("%02d:%02d:%02d, expect %f, got %f", each.Hour, each.Minute, each.Second, each.Expect, actual)
		}
	}
}

func TestIsLeapYear(t *testing.T) {
	cases := []struct {
		Year     int
		Calendar int
		Expect   bool
		ErrText  string
	}{
		{2004, GregorianCalendar, true, "year 2004 on Gregorian calendar is a leap year"},
		{2004, JulianCalendar, true, "year 2004 on Julian calendar is a leap year"},
		{1900, GregorianCalendar, false, "year 1900 on Gregorian calendar is a leap year"},
		{1900, JulianCalendar, true, "year 1900 on Julian calendar is a leap year"},
		{2000, GregorianCalendar, true, "year 2000 on Gregorian calendar is a leap year"},
		{2000, JulianCalendar, true, "year 2000 on Julian calendar is a leap year"},
	}

	for _, each := range cases {
		if isLeapYear(each.Year, each.Calendar) != each.Expect {
			t.Fatal(each.ErrText)
		}
	}
}

func TestNewDate(t *testing.T) {
	cases := []struct {
		Year       int
		Month      int
		Day        int
		Offset     float64
		ExpectDate *Date
	}{
		{1985, 1, 1, 0.5, &Date{Year: 1985, Month: 1, Day: 1, OffsetOfDay: 0.5, Calendar: GregorianCalendar}},
		{2000, 2, 29, 0, &Date{Year: 2000, Month: 2, Day: 29, OffsetOfDay: 0, Calendar: GregorianCalendar}},
		{1582, 10, 15, 0, &Date{Year: 1582, Month: 10, Day: 15, OffsetOfDay: 0, Calendar: GregorianCalendar}},
		{1582, 10, 5, 0, &Date{Year: 1582, Month: 10, Day: 15, OffsetOfDay: 0, Calendar: GregorianCalendar}},
		{1582, 10, 6, 0, &Date{Year: 1582, Month: 10, Day: 16, OffsetOfDay: 0, Calendar: GregorianCalendar}},
		{1582, 10, 14, 0, &Date{Year: 1582, Month: 10, Day: 24, OffsetOfDay: 0, Calendar: GregorianCalendar}},
		{1582, 10, 4, 0, &Date{Year: 1582, Month: 10, Day: 4, OffsetOfDay: 0, Calendar: JulianCalendar}},
		{1000, 2, 29, 0, &Date{Year: 1000, Month: 2, Day: 29, OffsetOfDay: 0, Calendar: JulianCalendar}},
		{0, 1, 1, 0, &Date{Year: 0, Month: 1, Day: 1, OffsetOfDay: 0, Calendar: JulianCalendar}},
		{-4712, 1, 1, 0.5, &Date{Year: -4712, Month: 1, Day: 1, OffsetOfDay: 0.5, Calendar: JulianCalendar}},
		{1999, 2, 29, 0, nil},
		{1211, 2, 29, 0, nil},
		{1000, 1, 32, 0, nil},
		{1000, 0, 20, 0, nil},
		{1000, 13, 20, 0, nil},
		{1000, 1, 1, -1, nil},
		{1000, 1, 1, 1, nil},
	}

	for _, each := range cases {
		d, err := NewDate(each.Year, each.Month, each.Day, each.Offset)
		if err != nil {
			if each.ExpectDate != nil {
				t.Fatalf("year:%d, month:%d, day:%d, offset:%f, expect:%+v, got error:%s", each.Year, each.Month, each.Day, each.Offset, each.ExpectDate, err.Error())
			}
		} else {
			if each.ExpectDate == nil {
				t.Fatalf("year:%d, month:%d, day:%d, offset:%f, expect an error, got %+v", each.Year, each.Month, each.Day, each.Offset, d)
			} else if d.Year != each.ExpectDate.Year ||
				d.Month != each.ExpectDate.Month ||
				d.Day != each.ExpectDate.Day ||
				d.OffsetOfDay != each.ExpectDate.OffsetOfDay {
				t.Fatalf("year:%d, month:%d, day:%d, offset:%f, expect:%+v, got:%+v", each.Year, each.Month, each.Day, each.Offset, each.ExpectDate, d)
			}
		}
	}
}

func TestDate_AddDays(t *testing.T) {
	cases := []struct {
		Year        int
		Month       int
		Day         int
		Calendar    int
		AddDays     int
		ExpectYear  int
		ExpectMonth int
		ExpectDay   int
	}{
		{2000, 1, 1, GregorianCalendar, 60, 2000, 3, 1},
		{1999, 2, 28, GregorianCalendar, 365, 2000, 2, 28},
		{1999, 3, 2, GregorianCalendar, 366, 2000, 3, 2},
		{1601, 1, 1, GregorianCalendar, 146097, 2001, 1, 1},
		{1200, 2, 28, JulianCalendar, 155, 1200, 8, 1},
		{1199, 3, 2, JulianCalendar, 366, 1200, 3, 2},
		{1001, 1, 1, JulianCalendar, 146100, 1401, 1, 1},
		{2001, 1, 1, GregorianCalendar, -146097, 1601, 1, 1},
		{2000, 1, 1, GregorianCalendar, -2451545, -4713, 11, 24},
		{1401, 1, 1, JulianCalendar, -146100, 1001, 1, 1},
		{837, 4, 10, JulianCalendar, -2026871, -4712, 1, 1},
		{2000, 8, 1, GregorianCalendar, -155, 2000, 2, 28},
		{1200, 8, 1, JulianCalendar, -155, 1200, 2, 28},
	}

	for _, each := range cases {
		d := &Date{
			Year:     each.Year,
			Month:    each.Month,
			Day:      each.Day,
			Calendar: each.Calendar,
		}
		d.AddDays(each.AddDays)
		if d.Year != each.ExpectYear ||
			d.Month != each.ExpectMonth ||
			d.Day != each.ExpectDay {
			calendar := "Gregorian"
			if d.Calendar == JulianCalendar {
				calendar = "Julian"
			}

			t.Fatalf("add %d days to %s calendar date %s, expect %s, got %s",
				each.AddDays,
				calendar,
				fmt.Sprintf("%d-%d-%d", each.Year, each.Month, each.Day),
				fmt.Sprintf("%d-%d-%d", each.ExpectYear, each.ExpectMonth, each.ExpectDay),
				fmt.Sprintf("%d-%d-%d", d.Year, d.Month, d.Day),
			)
		}
	}
}
