// History: Jan 04 14 tcolar Creation

package gollections

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"sort"
	"strings"
	"testing"
)

// #################### EXAMPLES ##############################################

// Some usage examples for gollection.Slice
//This does not demonstarte all the methods, see godoc and tests for more details
func ExampleSlice() {
	s := NewSlice()                 // Create a new slice
	s.Append("_")                   // add something to it
	s.AppendAll("A", "B", "Z", "J") // add several more things
	log.Print(s)                    // Slice[4] [_ A B Z J]
	i := s.Index("A")               // Find the index of the (first) element equal to "A" (1)
	log.Print(i)                    // 1

	var val string                  // We will get an element of the slice into this strongly typed var
	s.Get(2, &val)                  // set 'val' to the value of slice element at index 2
	log.Print(strings.ToLower(val)) // b
	s.Last(&val)                    // Get the last element
	log.Print(val)                  // J
	s.Get(-2, &val)                 // Get "secnd to last" element
	log.Print(val)                  // Z

	log.Print(s.ContainsAny("K", "Z")) // Does s contain either K or Z ? -> true

	log.Print(s.Join("|")) // "_|A|B|Z|J"

	// Using Each() closure to create a string of the elements joined by '-'
	val = ""
	s.Each(func(i int, e interface{}) bool {
		val = fmt.Sprintf("%s-%s", val, e.(string))
		return false // No "stop" condition is this closure
	})
	log.Print(val) // _-A-B-Z-J

	// More complex Each() form, iterating over a range with a stop condition
	val = ""
	s.EachRange(1, -2, func(i int, e interface{}) bool { // skip first and last elements
		val = fmt.Sprintf("%s-%s", val, e.(string))
		return e == "B" // But stop if we encountered a B
	})
	log.Print(val) // -A-B (We iterated from 'A' to 'Z' but stopped iteratin after 'B')

	// Example: using Any() to see if at least one element satisfies a condition
	any := s.Any(func(e interface{}) bool {
		str := e.(string) // we are working on strings, so doing an assertion
		// Is the string the same in upper and lower case ?
		return strings.ToLower(str) == str
	})
	log.Print(any) // true because '_' is the same in upper and lower case

	// Copying some of the slice content back into a strongly typed slice
	// Note that it's a costly operation as all elements have to be copied individually
	var raw []string
	s.ToRange(1, -2, &raw) // retrieving all but first and last element
	log.Print(raw)         // [A B Z]  ("standard" string slice)

	// Using findAll function to create a new list
	found := s.FindAll(func(i int, e interface{}) bool {
		return e.(string) < "X"
	})
	log.Print(found) // Slice[3] [A B J]

	// sort / search -> see TestSliceSearch

}

func TestSliceExample(t *testing.T) {
	ExampleSlice()
}

// #################### TESTS #################################################

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
		convey.So(s.Contains(7), convey.ShouldBeTrue)
		convey.So(s.Contains(999), convey.ShouldBeFalse)
		convey.So(s.ContainsAll(7, 1, 15), convey.ShouldBeTrue)
		convey.So(s.ContainsAll(7, 999), convey.ShouldBeFalse)
		convey.So(s.ContainsAny(99, 7, 15), convey.ShouldBeTrue)
		convey.So(s.ContainsAny(97, 98, -99), convey.ShouldBeFalse)
	})

	convey.Convey("First & Last", t, func() {
		s.First(&result)
		convey.So(result, convey.ShouldEqual, 1)
		s.Last(&result)
		convey.So(result, convey.ShouldEqual, 15)
	})

	convey.Convey("Join", t, func() {
		convey.So(s.Join("|"), convey.ShouldEqual, "1|2|3|7|9|15")
	})

	convey.Convey("Len", t, func() {
		convey.So(s.Len(), convey.ShouldEqual, 6)
	})

	convey.Convey("Clear", t, func() {
		s.Clear()
		convey.So(s.IsEmpty(), convey.ShouldBeTrue)
		convey.So(s.Len(), convey.ShouldEqual, 0)
		s.Append(7)
		convey.So(s.IsEmpty(), convey.ShouldBeFalse)
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
		s2 := s.Clone() // 1
		convey.So(s.Len(), convey.ShouldEqual, s2.Len())
		s.Append(2)        // 1 2
		s2.AppendAll(3, 4) // 1 3 4
		convey.So(s.Len(), convey.ShouldEqual, 2)
		convey.So(s2.Len(), convey.ShouldEqual, 3)
		s.Get(-1, &result)
		convey.So(result, convey.ShouldEqual, 2)
		s2.Get(-1, &result)
		convey.So(result, convey.ShouldEqual, 4)
		s3 := s2.CloneRange(1, -1) // 3 4
		convey.So(s3.Len(), convey.ShouldEqual, 2)
		s3.Get(0, &result)
		convey.So(result, convey.ShouldEqual, 3)
		s3.Get(1, &result)
		convey.So(result, convey.ShouldEqual, 4)
	})

	convey.Convey("Fill", t, func() {
		s.Clear()
		s.Append("X")
		s.Fill("A", 5)
		convey.So(s.Join(""), convey.ShouldEqual, "XAAAAA")
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
		s.ToRange(1, 2, &results)
		convey.So(len(results), convey.ShouldEqual, 2)
		convey.So(results[0], convey.ShouldEqual, 2)
		convey.So(results[1], convey.ShouldEqual, 3)
		s.ToRange(-3, -1, &results)
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
		convey.So(s.All(f1), convey.ShouldBeTrue)
		convey.So(s.All(f2), convey.ShouldBeFalse)
	})

	convey.Convey("Any", t, func() {
		f1 := func(e interface{}) bool {
			return e.(int) == 7
		}
		f2 := func(e interface{}) bool {
			return e.(int) == 22
		}
		convey.So(s.Any(f1), convey.ShouldBeTrue)
		convey.So(s.Any(f2), convey.ShouldBeFalse)
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
		a = ""
		s.Eachr(func(i int, e interface{}) bool {
			a = fmt.Sprintf("%s%d:%s ", a, i, e.(string))
			return e == "B" // stop on B
		})
		convey.So(a, convey.ShouldEqual, "7:F 6:E 5:E 4:B ")
	})

	convey.Convey("Find", t, func() {
		s.Clear()
		s.AppendAll("D", "E", "A", "D", "B", "E", "E", "F")
		i := s.Find(func(i int, e interface{}) bool {
			return e == "E"
		})
		convey.So(i, convey.ShouldEqual, 1)
		i = s.Find(func(i int, e interface{}) bool {
			return e == "Z"
		})
		convey.So(i, convey.ShouldEqual, -1)

		es := s.FindAll(func(i int, e interface{}) bool {
			return e == "E"
		})
		convey.So(es.Join(""), convey.ShouldEqual, "EEE")
	})

	var result string
	convey.Convey("Stack", t, func() {
		s.Clear()
		convey.So(func() { s.Peek(&result) }, convey.ShouldPanic)
		s.Push("A")
		s.Push("B")
		convey.So(s.Join(""), convey.ShouldEqual, "AB")
		s.Peek(&result)
		convey.So(result, convey.ShouldEqual, "B")
		convey.So(s.Join(""), convey.ShouldEqual, "AB")
		s.Pop(&result)
		convey.So(result, convey.ShouldEqual, "B")
		convey.So(s.Join(""), convey.ShouldEqual, "A")
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

// #################### BENCHMARKS ############################################

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

// #################### TESTS DATA ############################################

func testSlice() *Slice {
	s := NewSlice()
	s.AppendAll(1, 2, 3)
	s.Append(7)
	s.Append(9)
	s.Append(15)
	return s
}
