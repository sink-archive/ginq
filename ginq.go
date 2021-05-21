package ginq
import "go/types"

// Queryable wraps a slice and some operations to perform to transform it to a result
type Queryable struct {
	Type       types.Type
	Count      int64
	Base       []interface{}
	Operations []Operation
}

// Operation provides a type to explain a step in the transformation from start to end in a Queryable
type Operation struct {
	ParentQueryable    *Queryable
	NeedsFullEnumerate bool
	CalcNext           func(*Operation) interface{}
	CalcFull           func(*Operation) Queryable
}

// EnumerateAll takes a Queryable and returns a slice with the final result of the transformation.
func (q *Queryable) EnumerateAll() []interface{} {
	working := *q
	for i := 0; i < len(working.Operations); i++ {
		currentOp := working.Operations[i]

		if currentOp.NeedsFullEnumerate {
			working = currentOp.CalcFull(&currentOp)
		} else {
			for j := 0; j < int(working.Count); j++ {

			}
		}
	}

	return working.Base
}