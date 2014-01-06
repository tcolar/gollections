// History: Jan 04 14 tcolar Creation

package goon

import (
	"reflect"
)

type Slice struct {
	// internal slice that hold the objects
	slice []interface{}
	// value of pointer to slice
	sliceValPtr reflect.Value
}

// Create a new empty slice
func NewSlice() *Slice {
	s := &Slice{}
	s.slice = []interface{}{}
	s.sliceValPtr = reflect.ValueOf(&s.slice)
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

// Append several values (slice of)
func (s *Slice) AppendAll(more interface{}) {
	v := reflect.ValueOf(more)
	// assert v.Kind is Slice ?
	for i := 0; i != v.Len(); i++ {
		s.Append(v.Index(i).Interface())
	}
}

// Concatenate this slice with vaues of another slice
// TODO
//func (s *Slice) Concat(slice interface{}) {
//}

// Set value of ptr to this slice first element
func (s *Slice) First(ptr interface{}) {
	obj := reflect.ValueOf(ptr).Elem()
	obj.Set(reflect.Indirect(s.sliceValPtr).Index(0).Elem())
}

// Set value of ptr to this slice last element
func (s *Slice) Last(ptr interface{}) {
	obj := reflect.ValueOf(ptr).Elem()
	obj.Set(reflect.Indirect(s.sliceValPtr).Index(len(s.slice) - 1).Elem())
}

// Length of this slice
func (s *Slice) Len() int {
	return len(s.slice)
}

func (s *Slice) Test(ptr interface{}) {
	obj := reflect.ValueOf(ptr).Elem()
	t := reflect.SliceOf(reflect.TypeOf(5))
	l := reflect.MakeSlice(t, 2, 2)
	// then fill it
	obj.Set(reflect.Indirect(l))
}
