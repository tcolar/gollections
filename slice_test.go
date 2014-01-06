// History: Jan 04 14 tcolar Creation

package goon

import (
	"github.com/smartystreets/goconvey/convey"
	//"log"
	"testing"
)

func TestSlice(t *testing.T) {
	// test slice
	s := NewSlice()
	s.AppendAll([]int{1, 2, 3})
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

	convey.Convey("Len", t, func() {
		convey.So(s.Len(), convey.ShouldEqual, 6)
	})

	convey.Convey("First & Last", t, func() {
		s.First(&result)
		convey.So(result, convey.ShouldEqual, 1)
		s.Last(&result)
		convey.So(result, convey.ShouldEqual, 15)
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
