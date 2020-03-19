package julian

import (
	"fmt"
	"testing"
)

// TestJulianDay tests function JulianDay
// all test cases derived from the book <astronomical algorithms> 2nd edition by Jean Meeus
func TestJulianDay(t *testing.T) {
	cases := map[JulianDay]struct {
		Year        int
		Month       int
		Day         int
		offsetOfDay float64
	}{
		2451545.0:  {2000, 1, 1, 0.5},
		2451179.5:  {1999, 1, 1, 0},
		2446822.5:  {1987, 1, 27, 0},
		2446966.0:  {1987, 6, 19, 0.5},
		2447187.5:  {1988, 1, 27, 0},
		2447332.0:  {1988, 6, 19, 0.5},
		2436116.31: {1957, 10, 4, 0.81},
		2415020.5:  {1900, 1, 1, 0},
		2305447.5:  {1600, 1, 1, 0},
		2305812.5:  {1600, 12, 31, 0},
		2026871.8:  {837, 4, 10, 0.3},
		1842713.0:  {333, 1, 27, 0.5},
		1676496.5:  {-123, 12, 31, 0},
		1676497.5:  {-122, 1, 1, 0},
		1356001.0:  {-1000, 7, 12, 0.5},
		1355866.5:  {-1000, 2, 29, 0},
		1355671.4:  {-1001, 8, 17, 0.9},
		0.0:        {-4712, 1, 1, 0.5},
	}

	for expect, input := range cases {
		d, _ := NewDate(input.Year, input.Month, input.Day, input.offsetOfDay)
		v := d.JulianDay()

		if v != expect {
			t.Fatalf("expect Juliay day of %s to be %f, got %f",
				fmt.Sprintf("%d-%d-%f", input.Year, input.Month, float64(input.Day)+input.offsetOfDay),
				expect,
				v,
			)
		}
	}
}
