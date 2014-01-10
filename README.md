"Generic" Go(Golang) collections
=================================

The collections provides many powerful methods, some taking functions, making for productive programming.

Here is a small example covering only a few of the features:

```go
    s := NewSlice()
    s.Append("_")
    s.AppendAll("A", "B", "Z", "J")
    var val string
    s.Get(2, &val)             // Extract the element at index 2 (Z) into val (int)
    s.Get(-3, &val)            // Extract the second to last element (B) into val
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

Why ?
----
I'm actually not much of a generics lover, I never liked the way they where implemented in Java for example.

I've used Fantom a lot and while it has no generics either it provides very powerful collections that make you rarely miss them.

On the other Hand Go has neither generics nor collections with a lot of features, so this is an attempt to fill that gap.

How does it work
----------------

The custom collections rely on slice of "generic" elements ([]interface{} in Go).
Obviously that means that we lose some type safety, however it is mitigated by
the fact that you can retrieve elements into a srongly typed variable pointer.

For example
```Go
    s := NewSlice()
    s.AppendAll(5,6,7)
    var myInt int
    s.Get(1, &myint) // now myInt is a strongly typed int with the value 6
```

A benefit of this "trick" is that we do regain some type safety since we are getting the
value back into a strongly typed variable(int) that the compiler can watch for us from then on.

**Performance**

Overall the performance is actually better than expected.
Getting values from the generic slice into a type variable as an extra cost due to the use of reflection,
however so far benchmarking indicates it's not unreasonable. (More becnhmarking TBD)

I did put extra attention trying to make all the slice operations as efficient as I could.
Most operations are done in place unless otherwise noted and try not to allocate any unnecessay space.

One operation that is very costly is To() which "exports" the slice contents into a strongly typed slice
(native go slice), that requires the use of reflection and copy of each elements one at a time.
So it's best to not use it at all or only use it as the very last step once all operations are completed.





