# GINQ
**<u>G</u>o <u>IN</u>tegrated <u>Q</u>uery**

[![`GPL-3.0-or-later`](https://img.shields.io/badge/license-GPL--3.0--or--later-blue)](https://github.com/yellowsink/ginq/blob/master/LICENSE.md)

GINQ ports the amazing LINQ features of .NET to Golang.
Note that this library only deals with the `.Where()`, `.Select()`, and similar extension methods.
The awesome `(from x in list select x.y where x.z)`
query syntax would require special language support that a library cannot provide.

Read about LINQ [here!](https://docs.microsoft.com/en-us/dotnet/standard/linq)

## Some notes and drawbacks
To get this to work properly like the C# extension methods,
I'd need a way to take a generic type and use that for a receiver.
Instead, the functions have no receiver.

Instead of chaining left-to-right:
```go
x.Where(/* where stuff */).Select(/* select stuff */)
```
They chain outwards from the middle right-to-left:
```go
Select(/* select stuff */, Where(/* where stuff */, x))
```
### Functions that weren't implemented
- `.OfType()` - Go type system, and being OOP, and I don't know how I'd even approach this.
- `.ThenBy()` - This function provides a fallback for `OrderBy`, so that items that match the same when compared are sorted a different way.
- `.LongCount()` - Go's `len()` returns an `int` and not an `int64`.

## List of GINQ functions
### `Select[Tin, Tout](operation func(Tin) Tout, slice []Tin) []Tout`
`Select` takes a list of `Tin` and performs `operation` on every item, returning the resultant list of `Tout`.

### `Skip[T](amount int, slice []T) []T`
This function returns every element in the list, except the first `amount` items.

### `SkipLast[T](amount int, slice []T) []T`
This function returns every element in the list, except the last `amount` items.

### `SkipWhile[T](check func(T) bool, slice []T) []T`
Every element in the list will be ignored until `check` returns false, at which point every following element is returned.

### `Take[T](amount int, slice []T) []T`
This function returns the first `amount` items of the list.

### `TakeLast[T](amount int, slice []T) []T`
This function returns the last `amount` items of the list.

### `TakeWhile[T](check func(T) bool, slice []T) []T`
Every element in the list is returned until `check` returns false, at which point all following items are ignored.

### `Join[Tfirst, Tsecond, Tkey, Tout](...) []Tout`
The actual function signature with args is the following:
```go
Join[Tfirst, Tsecond, Tkey, Tout](
	firstSelector func(Tfirst) Tkey,
	secondSelector func(Tsecond) Tkey,
	resultFunc func(Tfirst, Tsecond) Tout,
	firstSlice []Tfirst,
	secondSlice []Tsecond) []Tout
```

This function takes two slices (of types `Tfirst` and `Tsecond`), and using the `firstSelector` and `secondSelector`
functions, they are processed into keys that can match up elements between the slices.
These will be type `Tkey`.

When two items with identical keys are found, both items are passed to `resultFunc` to generate the final value.
This final value of type `Tout` is added to the list that is returned.

Each item can have multiple matches - this doesn't take the first match for each item, it finds *all* of them.

### `GroupJoin[Tfirst, Tsecond, Tkey, Tout](...) []Tout`
The actual function signature with args is the following:
```go
Join[Tfirst, Tsecond, Tkey, Tout](
	firstSelector func(Tfirst) Tkey,
	secondSelector func(Tsecond) Tkey,
	resultFunc func(Tfirst, []Tsecond) Tout,
	firstSlice []Tfirst,
	secondSlice []Tsecond) []Tout
```

This function takes two slices (of types `Tfirst` and `Tsecond`), and using the `firstSelector` and `secondSelector`
functions, they are processed into keys that can match up elements between the slices.
These will be type `Tkey`.

For each item in `firstSlice`, it will find all matches in `secondSlice` and pass them to `resultFunc`,
which should process a `Tfirst` and a list of `Tsecond` to produce a final value of `Tout`.

All of these `Tout` values are returned as one list.

### `GroupBy[Tsource, Tkey, Telement, Tout] []Tout`
The actual function signature with args is the following:
```go
GroupBy[Tsource, Tkey, Telement, Tout](
	keySelector func (Tsource) Tkey,
	elementSelector func (Tsource) Telement,
	resultSelector func (Tkey, []Telement) Tout,
	slice []Tsource) []Tout
```

`GroupBy` takes a slice of `Tsource` and groups it into multiple groups based on a function.

Each item is passed through `keySelector` to get a key of type `Tkey`.
All items with the same key are put into the same group.

Each item is then passed through `elementSelector` to get a type of `Telement`.
This processes the item into the form that will make up the group items.

Finally, for each group, `resultSelector` is given the key of the group `TKey` and all items in the group `[]Telement`,
and should return a `Tout` object to be added to the list which is then returned.

### `Where[T](check func(T) bool, slice []T) []T`
This passes every item in the list to `check`, and return a list in which only includes values where `check` was true.

### TODO: document orderby, reverse, and all aggregation methods