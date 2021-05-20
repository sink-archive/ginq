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
### Projection
#### `Select[Tin, Tout](operation func(Tin) Tout, slice []Tin) []Tout`
`Select` takes a list of `Tin` and performs `operation` on every item, returning the resultant list of `Tout`.

#### `SelectMany[Tin, Tout](operation func(Tin) []Tout, slice []Tin) []Tout`
Like `Select` except `operation` returns multiple items. All the lists returned are concatenated into one.

### Partitioning
#### `Skip[T](amount int, slice []T) []T`
This function returns every element in the list, except the first `amount` items.

#### `SkipLast[T](amount int, slice []T) []T`
This function returns every element in the list, except the last `amount` items.

#### `SkipWhile[T](check func(T) bool, slice []T) []T`
Every element in the list will be ignored until `check` returns false, at which point every following element is returned.

#### `Take[T](amount int, slice []T) []T`
This function returns the first `amount` items of the list.

#### `TakeLast[T](amount int, slice []T) []T`
This function returns the last `amount` items of the list.

#### `TakeWhile[T](check func(T) bool, slice []T) []T`
Every element in the list is returned until `check` returns false, at which point all following items are ignored.

### Joining
#### `Join[Tfirst, Tsecond, Tkey, Tout](...) []Tout`
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

#### `GroupJoin[Tfirst, Tsecond, Tkey, Tout](...) []Tout`
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

### Grouping
#### `GroupBy[Tsource, Tkey, Telement, Tout] []Tout`
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

### Filtering
#### `Where[T](check func(T) bool, slice []T) []T`
This passes every item in the list to `check`, and return a list in which only includes values where `check` was true.

### Sorting
#### `OrderBy[T](isFirstLess func(T, T) bool, slice []T) []T`
Orders a slice using a function that determines if the first item is less or not.

#### `OrderByNumKey[T](keySelector func(T) float64, slice []T) []T`
Orders a slice where a key to sort by can be selected as a float64.

#### `OrderByDescending[T](...) []T`
See above.

#### `OrderByNumKeyDescending[T](...) []T`
See above.

#### `Reverse[T](slice []T) []T`
Reverses a slice.

### Aggregation
#### `Average(slice []float64) float64`
This returns the mean average of the floats passed.

#### `Count[T](slice []T) int64`
This returns the amount of items in the slice.

#### `MaxInt(slice []int64) int64`
This returns the largest int in the slice.

#### `MaxFloat(slice []float64) float64`
See above.

#### `MinInt(slice []int64) int64`
This returns the smallest int in the slice.

#### `MinFloat(slice []float64) float64`
See above.

#### `SumInts(slice []int64) int64`
This returns the sum of all ints in the slice.

#### `SumFloats(slice []float64) float64`
See above.

#### `Aggregate[T](base T, operation func (working, next T) T, slice []T) T`
This takes the item `base` and for each item in the slice,
passes `working`, and `next` - which is the current slice item,
into `operation`, and `working` for the next item becomes the output of that.
Once done the final value of `working` is returned.

Use this to do things like appending all strings in a list or something IDK.

### Set
#### `Distinct[T](slice []T) []T`
This returns the slice with duplicates removed.

#### `Except[T](first, second []T) []T`
This returns every item in `first` that is not in `second`.

#### `Intersect[T](first, second []T []T`
This returns every item that exists in both `first` and `second`.

#### `Union[T](first, second []T) []T`
This returns both lists concatenated, but with duplicates removed.

#### `Concat[T](first, second []T) []T`
This concatenates two slices together and returns the result.

#### `Zip[T1, T2, Tout](zipper func(T1, T2) Tout, first []T1, second []T2) []Tout`
Zip combines two lists together.
It does this by running through both lists simultaneously and passing the values at the same index as each to `zipper`.
The results of this are collected until one of the lists runs out. These results are then returned.

### Quantifier
#### `All[T](validate func(T) bool, slice []T) bool`
If ALL the items in the slice pass `validate`, then true, else false.

#### `AnyItems[T](slice []T) bool`
If the slice contains any items.

#### `Any[T](validate func(T) bool, slice []T) bool`
If the slice contains any items that pass `validate`.

### Generation
#### `DefaultIfEmpty[T](default_ T, slice []T) []T`
Returns the slice, or if it's empty, returns `default_` instead.

#### `Empty[T]() []T`
Gets an empty slice of `T`.

#### `Range(start, end int64) []int64`
Returns a slice with every value from start to end, BOTH INCLUSIVE.

#### `Repeat[T](element T, count int64) []T`
Returns a slice with `element` repeated `count` times.

### Element
#### `First[T](slice []T) T`
Get the first item in the slice.

#### `FirstOrDefault[T](default_ T, slice []T) T`
Get the first item in the slice, or a default.

#### `Last[T](slice []T) T`
See above.

#### `LastOrDefault[T](default_ T, slice []T) T`
See above.

#### `Single[T](slice []T) (T, error)`
Get the only item in the list, or if there is more or none throw an error.

#### `SingleOrDefault[T](default_ T, slice []T) T`
Single, but instead of an error return `default_`.

### Equality
#### `SequenceEqual[T](first, second []T) bool`
Checks if the sequences are equal.

### Conversions and things
#### `ToMap[Tkey, Tvalye](keySelector func(Tvalue) Tkey, slice []Tvalue) map[Tkey]Tvalue`
For each item in the list, calls `keySelector` and uses that as the key for a map, with the item as the value.
Duplicate items override the previous copies.