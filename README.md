"Generic" Go(Golang) collections
=================================

The collections provides many powerful methods, some taking functions, making for productive programming.

Here is a small example covering only a few of the features:

    s := NewSlice()
    s.Append("_")
    s.AppendAll("A", "B", "Z", "J")
    var val string
    s.Get(2, &val)
    log.Print(s.ContainsAny("K", "Z"))
    s.Clear()
    s.AppendAll(1,2,3,4,5,6)
    val := s.Reduce(0, func(reduction interface{}, i int, elem interface{}) interface{} {
      return reduction.(int) + elem.(int)
    })

**What does it do**

It provides "generic" Slice and Map elements for go.
They are feature full and mostly modeled against the Fantom [List](http://fantom.org/doc/sys/List.html) and Map implementations.

**Docs & Examples**

Gollections has some detailed Godocs:

[http://godoc.org/github.com/tcolar/gollections](http://godoc.org/github.com/tcolar/gollections)

Also even more details are available as unit tests:

[http://godoc.org/github.com/tcolar/gollections#example-Slice](http://godoc.org/github.com/tcolar/gollections#example-Slice)

Even more details can be found in the detailed unit tests:

[https://github.com/tcolar/gollections/blob/master/slice_test.go](https://github.com/tcolar/gollections/blob/master/slice_test.go)

Installing
----------
go get github.com/tcolar/gollections

