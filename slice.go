// History: Jan 04 14 Thibaut Colar Creation

package gollections

import (
	"fmt"
	"log"
	"reflect"
)

// Custom "Generic" (Sorta) slice
// Note: Satisfies sort.Interface so can use sort, search etc...
type Slice struct {

	// internal slice that hold the objects
	slice []interface{}
	// value of pointer to slice
	sliceValPtr reflect.Value

	// Returns whether two items are equal
	// Default imlementation uses reflect.DeepEqual
	Equals func(a, b interface{}) bool

	// Optional comparator function, must return 0 if a==b; -1 if a < b; 1 if a>b
	// **Nil by default**, MUST be defined for sorting to work.
	Compare func(a, b interface{}) int
}

// Initialize a new empty slice
func NewSlice() *Slice {
	s := &Slice{}
	s.sliceValPtr = reflect.ValueOf(&s.slice)
	s.Equals = func(a, b interface{}) bool { return reflect.DeepEqual(a, b) }
	return s
}

// Return true if f returns true for all of the items in the list.
func (s *Slice) All(f func(interface{}) bool) bool {
	for _, e := range s.slice {
		if !f(e) {
			return false
		}
	}
	return true
}

// Return true if c returns true for any(at least 1) of the items in the list
func (s *Slice) Any(f func(interface{}) bool) bool {
	for _, e := range s.slice {
		if f(e) {
			return true
		}
	}
	return false
}

// Append a single value
func (s *Slice) Append(more interface{}) {
	s.slice = append(s.slice, more)
}

// Append several values
func (s *Slice) AppendAll(more ...interface{}) {
	s.slice = append(s.slice, more...)
}

// Append another goon.Slice to this slice
func (s *Slice) AppendSlice(slice *Slice) {
	s.slice = append(s.slice, slice.slice...)
}

// Current slice capacity
func (s *Slice) Cap() int {
	return cap(s.slice)
}

// Clear (empty) the list
func (s *Slice) Clear() {
	// Note: A nil slice in go is valid and can then be used just as if empty
	s.slice = nil
}

// Create and return a clone of this slice
func (s *Slice) Clone() *Slice {
	clone := NewSlice()
	clone.slice = append(clone.slice, s.slice)
	return clone
}

// Does the slice contain the given element (by equality)
// Note, this uses simple iteration, use sort methods if meeding more performance
func (s *Slice) Contains(elem interface{}) bool {
	return s.Index(elem) != -1
}

// Does the slice contain all the given values
func (s *Slice) ContainsAll(elems ...interface{}) bool {
	for _, elem := range elems {
		if !s.Contains(elem) {
			return false
		}
	}
	return true
}

// Does the slice contain at least one of the given values
func (s *Slice) ContainsAny(elems ...interface{}) bool {
	for _, elem := range elems {
		if s.Contains(elem) {
			return true
		}
	}
	return false
}

// Apply the function to the whole slice (in order)
// If the function returns true (stop), iteration will stop
func (s *Slice) Each(f func(int, interface{}) (stop bool)) {
	s.EachRange(0, len(s.slice)-1, f)
}

// Apply the function to the slice range
// if from is < to it will be called in reversed order
// If the function returns true (stop), iteration will stop
func (s *Slice) EachRange(from, to int, f func(int, interface{}) (stop bool)) {
	l := len(s.slice)
	// Deal with negative indexes
	if from < 0 {
		from = l + from
	}
	if to < 0 {
		to = l + to
	}
	// Figure if we are to step forward or backwars
	step := 1
	steps := to - from
	if from > to {
		step = -1
		steps = -steps
	}
	var stop bool
	// Process the each
	for i := 0; i != steps+1; i++ {
		stop = f(from, s.slice[from])
		if stop {
			break
		}
		from += step
	}
}

// Apply the function to the whole slice (reverse order)
// If the function returns true (stop), iteration will stop
func (s *Slice) Eachr(f func(int, interface{}) (stop bool)) {
	s.EachRange(len(s.slice)-1, 0, f)
}

// Set value of ptr to this slice first element
func (s *Slice) First(ptr interface{}) {
	s.Get(0, ptr)
}

// Set value of ptr to slice[idx]
// If idx is negative then idx element from the end -> slice[len(slice)+idx]
// ie Get(-1) would return the last element
func (s *Slice) Get(idx int, ptr interface{}) {
	if idx < 0 {
		idx = len(s.slice) + idx
	}
	obj := reflect.ValueOf(ptr).Elem()
	obj.Set(reflect.Indirect(s.sliceValPtr).Index(idx).Elem())
}

// Return the (lowest) index of given element (using Equals() method)
// Return -1 if the lement is part of the slice
// Note, this uses simple iteration, use sort methods if meeding more performance
func (s *Slice) Index(val interface{}) int {
	for i, e := range s.slice {
		if s.Equals(e, val) {
			return i
		}
	}
	return -1
}

// Set value of ptr to this slice last element
func (s *Slice) Last(ptr interface{}) {
	s.Get(-1, ptr)
}

// Length of this slice
// Also used for impl of sort.Interface
func (s *Slice) Len() int {
	return len(s.slice)
}

// Check if a < b (used as impl of sort.Interface)
// S.Compare must be defined !
func (s *Slice) Less(a, b int) bool {
	if s.Compare == nil {
		panic("Slice.Compare function was not implemented !")
	}
	return s.Compare(s.slice[a], s.slice[b]) == -1
}

// Returns pointer to the raw underlying slice ([]interface{})
func (s *Slice) Slice() *[]interface{} {
	return &s.slice
}

// impl String interface
func (s *Slice) String() string {
	return fmt.Sprintf("Slice[%d] %v", len(s.slice), s.slice)
}

// Export our "generic" slice to a typed slice (say []int)
// Ptr needs to be a pointer to a slice
// Note that it can't be a simple cast and instead the data needs to be copied
// so it's definitely a VERY costly operation.
func (s *Slice) To(ptr interface{}) {
	s.ToRange(0, len(s.slice)-1, ptr)
}

// Same as To() but only get a subset(range) of the slice
// Note that from and to can use negative index to indicate "from the end"
func (s *Slice) ToRange(from, to int, ptr interface{}) {
	l := len(s.slice)
	if from < 0 {
		from = l + from
	}
	if to < 0 {
		to = l + to
	}
	if to < from || from < 0 || to > l-1 {
		log.Fatalf("ToSub: Indexes(%d:%d) out of range(0:%d)", from, to, l-1)
	}

	// Value of the pointer to the target
	obj := reflect.Indirect(reflect.ValueOf(ptr))
	// We can't just convert from interface{} to whatever the target is (diff memory layout),
	// so we need to create a New slice of the proper type and copy the values individually
	t := reflect.TypeOf(ptr).Elem()
	slice := reflect.MakeSlice(t, to-from+1, to-from+1)
	// Copying the data, val is an adressable Pointer of the actual target type
	val := reflect.Indirect(reflect.New(t.Elem()))
	for i := from; i <= to; i++ {
		v := reflect.ValueOf(s.slice[i])
		val.Set(v)
		slice.Index(i - from).Set(v)
	}
	// Ok now assign our slice to the target pointer
	obj.Set(slice)
}

// Swap 2 elements (used as impl of sort.Interface)
func (s *Slice) Swap(a, b int) {
	s.slice[a], s.slice[b] = s.slice[b], s.slice[a]
}
