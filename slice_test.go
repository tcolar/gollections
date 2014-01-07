// History: Jan 04 14 tcolar Creation

package goon

import (
	//"log"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

func TestSlice(t *testing.T) {
	s := testSlice()
	// result target
	var result int

	convey.Convey("Get", t, func() {
		s.Get(2, &result)
		convey.So(result, convey.ShouldEqual, 3)
		s.Get(-2, &result)
		convey.So(result, convey.ShouldEqual, 9)
	})

	convey.Convey("Index", t, func() {
		convey.So(s.Index(7), convey.ShouldEqual, 3)
		convey.So(s.Index(1), convey.ShouldEqual, 0)
		convey.So(s.Index(999), convey.ShouldEqual, -1)
	})

	convey.Convey("Contains", t, func() {
		convey.So(s.Contains(7), convey.ShouldEqual, true)
		convey.So(s.Contains(999), convey.ShouldEqual, false)
		convey.So(s.ContainsAll(7, 1, 15), convey.ShouldEqual, true)
		convey.So(s.ContainsAll(7, 999), convey.ShouldEqual, false)
		convey.So(s.ContainsAny(99, 7, 15), convey.ShouldEqual, true)
		convey.So(s.ContainsAll(97, 98, -99), convey.ShouldEqual, false)
	})

	convey.Convey("First & Last", t, func() {
		s.First(&result)
		convey.So(result, convey.ShouldEqual, 1)
		s.Last(&result)
		convey.So(result, convey.ShouldEqual, 15)
	})

	convey.Convey("Len", t, func() {
		convey.So(s.Len(), convey.ShouldEqual, 6)
	})

	convey.Convey("Clear", t, func() {
		s.Clear()
		convey.So(s.Len(), convey.ShouldEqual, 0)
		s.Append(7)
		convey.So(s.Len(), convey.ShouldEqual, 1)
	})

	convey.Convey("Append", t, func() {
		s.Clear()
		s.Append(5)
		convey.So(s.Len(), convey.ShouldEqual, 1)
		s.AppendAll(10, 22, 33)
		convey.So(s.Len(), convey.ShouldEqual, 4)
		s2 := NewSlice()
		s2.AppendAll(12, 13)
		s.AppendSlice(s2)
		convey.So(s.Len(), convey.ShouldEqual, 6)
		s.Get(-1, &result)
		convey.So(result, convey.ShouldEqual, 13)
	})

	convey.Convey("Clone", t, func() {
		s.Clear()
		s.Append(1)
		s2 := s.Clone()
		convey.So(s.Len(), convey.ShouldEqual, s2.Len())
		s.Append(2)
		s2.AppendAll(3, 4)
		convey.So(s.Len(), convey.ShouldEqual, 2)
		convey.So(s2.Len(), convey.ShouldEqual, 3)
		s.Get(-1, &result)
		convey.So(result, convey.ShouldEqual, 2)
		s2.Get(-1, &result)
		convey.So(result, convey.ShouldEqual, 4)
	})
	convey.Convey("To", t, func() {
		var results []int
		s.Clear()
		s.Append(1)
		s.AppendAll(2, 3, 4)
		s.To(&results)
		convey.So(len(results), convey.ShouldEqual, 4)
		convey.So(results[0], convey.ShouldEqual, 1)
		convey.So(results[3], convey.ShouldEqual, 4)
		results[0]++ // this is an actual number now
		convey.So(results[0], convey.ShouldEqual, 2)
		s.ToRange(&results, 1, 2)
		convey.So(len(results), convey.ShouldEqual, 2)
		convey.So(results[0], convey.ShouldEqual, 2)
		convey.So(results[1], convey.ShouldEqual, 3)
		s.ToRange(&results, -3, -1)
		convey.So(len(results), convey.ShouldEqual, 3)
		convey.So(results[0], convey.ShouldEqual, 2)
		convey.So(results[2], convey.ShouldEqual, 4)
	})
}

// Test for methods that take functions
func TestSliceFuncs(t *testing.T) {
	s := testSlice()

	convey.Convey("All", t, func() {
		f1 := func(e interface{}) bool {
			return e.(int) >= 1
		}
		f2 := func(e interface{}) bool {
			return e.(int) > 5
		}
		convey.So(s.All(f1), convey.ShouldEqual, true)
		convey.So(s.All(f2), convey.ShouldEqual, false)
	})

	convey.Convey("Any", t, func() {
		f1 := func(e interface{}) bool {
			return e.(int) == 7
		}
		f2 := func(e interface{}) bool {
			return e.(int) == 22
		}
		convey.So(s.Any(f1), convey.ShouldEqual, true)
		convey.So(s.Any(f2), convey.ShouldEqual, false)
	})

	convey.Convey("Each", t, func() {
		s.Clear()
		s.AppendAll("D", "E", "A", "D", "B", "E", "E", "F")
		// Test each
		a := "" // Bound varaiable that will be used by .Each closure
		f := func(i int, e interface{}) bool {
			a = fmt.Sprintf("%s%d:%s ", a, i, e.(string))
			return false
		}
		s.Each(f)
		convey.So(a, convey.ShouldEqual, "0:D 1:E 2:A 3:D 4:B 5:E 6:E 7:F ")
		// reverse each
		a = ""
		s.Eachr(f)
		convey.So(a, convey.ShouldEqual, "7:F 6:E 5:E 4:B 3:D 2:A 1:E 0:D ")
		// Test a range
		a = ""
		s.EachRange(2, -2, f)
		convey.So(a, convey.ShouldEqual, "2:A 3:D 4:B 5:E 6:E ")
		// Reversed range
		a = ""
		s.EachRange(6, 4, f)
		convey.So(a, convey.ShouldEqual, "6:E 5:E 4:B ")
		// Test stop
		f2 := func(i int, e interface{}) bool {
			a = fmt.Sprintf("%s%d:%s ", a, i, e.(string))
			return e.(string) == "B" // stop on B
		}
		a = ""
		s.Eachr(f2)
		convey.So(a, convey.ShouldEqual, "7:F 6:E 5:E 4:B ")
	})

}

func TestSliceSearch(t *testing.T) {
	s := testSlice()

	convey.Convey("Sort & Search", t, func() {
		// Example of sorting the slice
		s.Compare = func(a, b interface{}) int {
			ai := a.(int)
			bi := b.(int)
			if ai == bi {
				return 0
			}
			if ai > bi {
				return 1
			}
			return -1
		}

		// Using standard sorting
		sort.Sort(s)
		last := -9999
		var result int
		for i := 0; i != s.Len(); i++ {
			s.Get(i, &result)
			convey.So(result, convey.ShouldBeGreaterThanOrEqualTo, last)
			last = result
		}

		// Example of using standard search on sorted data
		i := sort.Search(s.Len(), func(i int) bool {
			var v int
			s.Get(i, &v)
			return v >= 2 // looking for first index of "2"
		})
		convey.So(i, convey.ShouldEqual, 1) // should be the second element
	})
}

func BenchmarkSlice(b *testing.B) {
	s := NewSlice()
	var result int
	for i := 0; i < b.N; i++ {
		s.Append(i)
		s.Last(&result)
	}
}

func BenchmarkSliceTo(b *testing.B) {
	s := NewSlice()
	var results []int
	for i := 0; i < b.N; i++ {
		s.Append(7)
		s.To(&results) // slow for large slices
	}
}

func testSlice() *Slice {
	s := NewSlice()
	s.AppendAll(1, 2, 3)
	s.Append(7)
	s.Append(9)
	s.Append(15)
	return s
}
