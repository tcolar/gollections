// History: Jan 04 14 tcolar Creation

package goon

import (
	"github.com/smartystreets/goconvey/convey"
	"log"
	"sort"
	"testing"
)

func TestSlice(t *testing.T) {
	// test slice
	s := NewSlice()
	s.AppendAll(1, 2, 3)
	s.Append(7)
	s.Append(9)
	s.Append(15)

	// result target
	var result int

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
		f3 := func(e interface{}) bool {
			return e.(int) == 7
		}
		f4 := func(e interface{}) bool {
			return e.(int) == 22
		}
		convey.So(s.Any(f3), convey.ShouldEqual, true)
		convey.So(s.Any(f4), convey.ShouldEqual, false)
	})

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
	})

	convey.Convey("Len", t, func() {
		convey.So(s.Len(), convey.ShouldEqual, 6)
	})

	convey.Convey("First & Last", t, func() {
		s.First(&result)
		convey.So(result, convey.ShouldEqual, 1)
		s.Last(&result)
		convey.So(result, convey.ShouldEqual, 15)
	})

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
		sort.Sort(s)
		last := -9999
		for i := 0; i != s.Len(); i++ {
			s.Get(i, &result)
			convey.So(result, convey.ShouldBeGreaterThanOrEqualTo, last)
			last = result
		}
		// Example of using standard search on sorted data
		i := sort.Search(s.Len(), func(i int) bool {
			var v int
			s.Get(i, &v)
			log.Print(v)
			return v >= 2 // looking for first index of "2"
		})
		convey.So(i, convey.ShouldEqual, 1) // should be the second element
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

}

func BenchmarkSlice(b *testing.B) {
	s := NewSlice()
	var result int
	for i := 0; i < b.N; i++ {
		s.Append(7)
		s.Last(&result)
	}
}
