"Generic" Go(Golang) collections
=================================

The collections provides many powerful methods, some taking functions, making for productive programming.

Here is a small example covering only a few of the features:

```go
    s := NewSlice()
    s.Append("_")
    s.AppendAll("A", "B", "Z", "J")
    var val string
    s.Get(2, &val)             // Extract the element at index 2 into val (int)
    s.Get(-2, &val)            // Extract the second to last element
    log.Print(s.ContainsAny("K", "Z")) // true
    s.Clear()
    s.AppendAll(1,2,3,4,5,6)
    // Example of calculating the sum using Reduce()
    sum := s.Reduce(0, func(reduction interface{}, i int, elem interface{}) interface{} {
      return reduction.(int) + elem.(int)
    })
    s.Reverse()                       // Reverse (in place)
    log.Print(s.Join(","))            // "6,5,4,3,2,1"
    s.Pop(&val)                       // Pop the last element (1)
```

**What does it do**

It provides "generic" Slice and Map elements for go.
They are feature full and mostly modeled against the Fantom [List](http://fantom.org/doc/sys/List.html) and Map implementations.

**Docs & Examples**

Gollections have some detailed Godocs:

[http://godoc.org/github.com/tcolar/gollections](http://godoc.org/github.com/tcolar/gollections)

You will find some examples here :

[http://godoc.org/github.com/tcolar/gollections#example-Slice](http://godoc.org/github.com/tcolar/gollections#example-Slice)

Even more details can be found in the detailed unit tests:

[https://github.com/tcolar/gollections/blob/master/slice_test.go](https://github.com/tcolar/gollections/blob/master/slice_test.go)

Installing
----------
go get github.com/tcolar/gollections

