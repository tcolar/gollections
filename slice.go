// History: Jan 04 14 tcolar Creation

package goon

import (
	"reflect"
)

// Custom "Generic" (Sorta) slice
// Note: Satisfies sort.Interface so can use sort, search etc...
type Slice struct {
	// A comparator function. Returns 0 if a==b; -1 if a < b; 1 if a>b
	// **NOT INITIALIZED**, must be set for sorting to work
	Compare func(a, b interface{}) int

	// Returns whether two items are equal
	// Default imlementation uses reflect.DeepEqual
	Equals func(a, b interface{}) bool

	// internal slice that hold the objects
	slice []interface{}
	// value of pointer to slice
	sliceValPtr reflect.Value
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

// Append a goon.Slice to this slice
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

// Set value of ptr to this slice first element
func (s *Slice) First(ptr interface{}) {
	s.Get(0, ptr)
}

// Set value of ptr to slice[idx]
// If idx is negative then idx element from the end -> slice[len(slice)+idx]
// ie Get(-1) would return the last element
func (s *Slice) Get(idx int, ptr interface{}) {
	if idx < 0 {
		idx = s.Len() + idx
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
// Also used as impl of sort.Interface
func (s *Slice) Len() int {
	return len(s.slice)
}

// Check if a < b (used as impl of sort.Interface)
// S.Compare must be defined !
func (s Slice) Less(a, b int) bool {
	if s.Compare == nil {
		panic("Slice.Compare function was not implemented !")
	}
	return s.Compare(s.slice[a], s.slice[b]) == -1
}

// Returns pointer to the underlying slice ([]interface{})
func (s *Slice) Slice() *[]interface{} {
	return &s.slice
}

// Swap 2 elements (used as impl of sort.Interface)
func (s *Slice) Swap(a, b int) {
	s.slice[a], s.slice[b] = s.slice[b], s.slice[a]
}

/*func (s *Slice) Test(ptr interface{}) {
  obj := reflect.ValueOf(ptr).Elem()
  t := reflect.SliceOf(reflect.TypeOf(5))
  l := reflect.MakeSlice(t, 2, 2)
  // then fill it
  obj.Set(reflect.Indirect(l))
}*/
