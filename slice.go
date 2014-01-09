// History: Jan 04 14 Thibaut Colar Creation

package gollections

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

// Custom "Generic" (Sorta) slice
// Note: Satisfies sort.Interface so can use sort, search as long as Compare is
// implemented
type Slice struct {

	// internal slice that hold the items
	slice []interface{}
	// value of pointer to slice
	sliceValPtr reflect.Value

	// Returns whether two items are equal
	// Default imlementation uses reflect.DeepEqual (==)
	Equals func(a, b interface{}) bool

	// Optional comparator function, must return 0 if a==b; -1 if a < b; 1 if a>b
	// **Nil by default**
	// **MUST** be defined for sorting to work.
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
func (s *Slice) Append(elem interface{}) {
	s.slice = append(s.slice, elem)
}

// Append several values
func (s *Slice) AppendAll(elems ...interface{}) {
	s.slice = append(s.slice, elems...)
}

// Append another Slice to this slice
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

// Clone part of this slice into a new Slice
// From and To are both inclusive
func (s *Slice) CloneRange(from, to int) *Slice {
	var err error
	if from, err = s.handleIndex(from); err != nil {
		panic(err.Error())
	}
	if to, err = s.handleIndex(to); err != nil {
		panic(err.Error())
	}
	clone := NewSlice()
	clone.slice = append(clone.slice, s.slice[from:to+1]...)
	return clone
}

// Does the slice contain the given element (by equality)
// Note, this uses simple iteration, use sort methods if meeding more performance
func (s *Slice) Contains(elem interface{}) bool {
	return s.IndexOf(elem) != -1
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
// From and To are both inclusive
// if from is < to it will be called in reversed order
// If the function returns true (stop), iteration will stop
func (s *Slice) EachRange(from, to int, f func(int, interface{}) (stop bool)) {
	var err error
	if from, err = s.handleIndex(from); err != nil {
		panic(err.Error())
	}
	if to, err = s.handleIndex(to); err != nil {
		panic(err.Error())
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

// Fill(append to) the slice with 'count' times the 'elem' value
func (s *Slice) Fill(elem interface{}, count int) {
	for i := 0; i != count; i++ {
		s.Append(elem)
	}
}

// Apply a function to find an element in the slice (iteratively)
// Returns the index if found, or -1 if no matches.
// The function is expected to return true when the index is found.
func (s *Slice) Find(f func(int, interface{}) (found bool)) (index int) {
	for i, e := range s.slice {
		if f(i, e) {
			return i
		}
	}
	return -1
}

// Apply a function to find all element in the slice for which the function returns true
// Returns a new Slice made of the matches.
func (s *Slice) FindAll(f func(int, interface{}) (found bool)) *Slice {
	results := NewSlice()
	for i, e := range s.slice {
		if f(i, e) {
			results.slice = append(results.slice, e)
		}
	}
	return results
}

// Set value of ptr to this slice first element
// Return an error if slice is empty
func (s *Slice) First(ptr interface{}) {
	if _, err := s.handleIndex(0); err != nil {
		panic(err.Error())
	}
	s.Get(0, ptr)
}

// Set value of ptr to slice[idx]
// If idx is negative then idx element from the end -> slice[len(slice)+idx]
// ie Get(-1) would return the last element
func (s *Slice) Get(idx int, ptr interface{}) {
	var err error
	if idx, err = s.handleIndex(idx); err != nil {
		panic(err.Error())
	}
	obj := reflect.ValueOf(ptr).Elem()
	obj.Set(reflect.Indirect(s.sliceValPtr).Index(idx).Elem())
}

// Return the (lowest) index of given element (using Equals() method)
// Return -1 if the lement is part of the slice
// Note, this uses simple iteration, use sort methods if meeding more performance
func (s *Slice) IndexOf(elem interface{}) int {
	for i, e := range s.slice {
		if s.Equals(e, elem) {
			return i
		}
	}
	return -1
}

// Insert the element before index idx
// Can use negative index
func (s *Slice) Insert(idx int, elem interface{}) {
	var err error
	if idx, err = s.handleIndex(idx); err != nil {
		panic(err.Error())
	}
	s.slice = append(s.slice, 0)
	copy(s.slice[idx+1:], s.slice[idx:])
	s.slice[idx] = elem
}

// Insert All the element before index idx
// Can use negative index
func (s *Slice) InsertAll(idx int, elems ...interface{}) {
	var err error
	if idx, err = s.handleIndex(idx); err != nil {
		panic(err.Error())
	}
	// Expand the slice by elems size
	s.slice = append(s.slice, make([]interface{}, len(elems))...)
	// Shift "in place" elements to the right of index to the right
	copy(s.slice[idx+len(elems):], s.slice[idx:])
	// fill in the space with the elements to be inserted
	copy(s.slice[idx:], elems)
}

// Insert All the element of the slice before index idx
// Can use negative index
func (s *Slice) InsertSlice(idx int, slice *Slice) {
	s.InsertAll(idx, slice.slice...)
}

// Is this slice empty
func (s *Slice) IsEmpty() bool {
	return len(s.slice) == 0
}

// Create a string by jining all the elements with the given seprator
// Note: Use fmt.Sprintf("%v", e) to get each element as a string
func (s *Slice) Join(sep string) string {
	var buf bytes.Buffer
	for i, e := range s.slice {
		if i != 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(fmt.Sprintf("%v", e))
	}
	return buf.String()
}

// Set value of ptr to this slice last element
// Will panic if slice is empty
func (s *Slice) Last(ptr interface{}) {
	if _, err := s.handleIndex(-1); err != nil {
		panic(err.Error())
	}
	s.Get(-1, ptr)
}

// Length of this slice
// Also used for impl of sort.Interface
func (s *Slice) Len() int {
	return len(s.slice)
}

// Check if element at index a < b (used as impl of sort.Interface)
// S.Compare must be defined !
func (s *Slice) Less(a, b int) bool {
	if s.Compare == nil {
		panic("Slice.Compare function was not implemented !")
	}
	var err error
	if a, err = s.handleIndex(a); err != nil {
		panic(err.Error())
	}
	if b, err = s.handleIndex(b); err != nil {
		panic(err.Error())
	}
	return s.Compare(s.slice[a], s.slice[b]) == -1
}

// Set ptr to the minimum value in the slice (pamic if slice is empty)
// NOTE: Compare function **MUST** be implemented
// This uses simple iteration (0n time) and does not modify the slice
// Alternatively use sort.Sort(slice).Get(0, ptr) when performance is needed
func (s *Slice) Min(ptr interface{}) {
	if s.IsEmpty() {
		panic("Can't find Min of empty slice !")
	}
	minIdx := 0
	for i := 1; i < len(s.slice); i++ {
		if s.Less(i, minIdx) {
			minIdx = i
		}
	}
	s.Get(minIdx, ptr)
}

// Set ptr to the maximum value in the slice (pamic if slice is empty)
// NOTE: Compare function **MUST** be implemented
// This uses simple iteration (0n time) and does not modify the slice
// Alternatively use sort.Sort(slice).Get(-1, ptr) when performance is needed
func (s *Slice) Max(ptr interface{}) {
	if s.IsEmpty() {
		panic("Can't find Max of empty slice !")
	}
	maxIdx := 0
	for i := 1; i < len(s.slice); i++ {
		if s.Less(maxIdx, i) {
			maxIdx = i
		}
	}
	s.Get(maxIdx, ptr)
}

// Set ptr to the last element
// Will panic if slice is empty
func (s *Slice) Peek(ptr interface{}) {
	s.Last(ptr)
}

// Pop (return & remove) and set ptr to the last element
// Will panic if slice is empty
func (s *Slice) Pop(ptr interface{}) {
	s.Last(ptr)
	// remove last elem of slice
	s.slice = s.slice[:len(s.slice)-1]
}

// Push an elem at the end of the slice (same as Append)
func (s *Slice) Push(elem interface{}) {
	s.Append(elem)
}

// Reduce is used to iterate through every item in the list to reduce the list
// into a single value called the reduction.
// The initial value (startVal) of the reduction is passed in as the init parameter
// then passed to the closure along with each item (which returns the updated reduction)
// See Tests / Examples in the test file for more info
func (s *Slice) Reduce(startVal interface{}, f func(reduction interface{}, index int, elem interface{}) interface{}) interface{} {
	reduction := startVal
	for i, e := range s.slice {
		reduction = f(reduction, i, e)
	}
	return reduction
}

// Remove the element at the given index
func (s *Slice) RemoveAt(idx int) {
	copy(s.slice[idx:], s.slice[idx+1:]) // shift elements past index to the left
	s.slice = s.slice[:len(s.slice)-1]   // lose last element
}

// Remove the first element found by value equality (found by IndexFrom method)
func (s *Slice) RemoveElem(elem interface{}) {
	idx := s.IndexOf(elem)
	if idx >= 0 {
		s.RemoveAt(idx)
	}
}

// Remove all elements by value equality (using Equals function)
func (s *Slice) RemoveElems(elem interface{}) {
	s.RemoveFunc(func(idx int, e interface{}) bool {
		return s.Equals(elem, e)
	})
}

// Remove the elements that match the function (where the function return true)
func (s *Slice) RemoveFunc(f func(idx int, elem interface{}) bool) {
	for i := 0; i < len(s.slice); i++ {
		if f(i, s.slice[i]) {
			s.RemoveAt(i)
			i--
		}
	}
}

// Remove the elements within the given index range
func (s *Slice) RemoveRange(from, to int) {
	var err error
	if from, err = s.handleIndex(from); err != nil {
		panic(err.Error())
	}
	if to, err = s.handleIndex(to); err != nil {
		panic(err.Error())
	}
	copy(s.slice[from:], s.slice[to:])         // shift elements
	s.slice = s.slice[:len(s.slice)-(to-from)] // lose last elements
}

// Reverse the slice in place (first element becomes last etc...)
func (s *Slice) Reverse() {
	start := 0
	end := len(s.slice) - 1
	for end > start { // Otherwise 0 or 1 element left, nothing to swap
		s.slice[start], s.slice[end] = s.slice[end], s.slice[start]
		end--
		start++
	}
}

// Set the element at the given index
func (s *Slice) Set(idx int, elem interface{}) {
	s.slice[idx] = elem
}

// Returns pointer to the raw underlying slice ([]interface{})
func (s *Slice) Slice() *[]interface{} {
	return &s.slice
}

// impl String interface
func (s *Slice) String() string {
	return fmt.Sprintf("Slice[%d] %v", len(s.slice), s.slice)
}

// Swap 2 elements (used as impl of sort.Interface)
// Return an error if the indexes are out of bounds
func (s *Slice) Swap(a, b int) {
	var err error
	if a, err = s.handleIndex(a); err != nil {
		panic(err.Error())
	}
	if b, err = s.handleIndex(b); err != nil {
		panic(err.Error())
	}

	s.slice[a], s.slice[b] = s.slice[b], s.slice[a]
}

// Export our "generic" slice to a typed slice (say []int)
// Ptr needs to be a pointer to a slice
// Note that it can't be a simple cast and instead the data needs to be copied
// so it's definitely a VERY costly operation.
func (s *Slice) To(ptr interface{}) {
	s.ToRange(0, len(s.slice)-1, ptr)
}

// Same as To() but only get a subset(range) of the slice
// From and To are both inclusive
// Note that from and to can use negative index to indicate "from the end"
func (s *Slice) ToRange(from, to int, ptr interface{}) {
	var err error
	if from, err = s.handleIndex(from); err != nil {
		panic(err.Error())
	}
	if to, err = s.handleIndex(to); err != nil {
		panic(err.Error())
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

// Validate the index is in the slice bounds
// Also turm negative indexes into index from the end of the slice (-1 = last)
func (s *Slice) handleIndex(idx int) (int, error) {
	if idx < 0 {
		idx = len(s.slice) + idx
	}
	if idx >= len(s.slice) || idx < 0 {
		return idx, errors.New(fmt.Sprintf("Invalid slice index: %d", idx))
	}
	return idx, nil
}
