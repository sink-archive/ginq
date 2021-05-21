package ginq

import "go/types"

// Queryable wraps a slice and some operations to perform to transform it to a result
type Queryable struct {
	Type           types.Type
	Length         int64
	Base           []interface{}
	Operations     []Operation
	EnumerateIndex int64
}

// Operation provides a type to explain a step in the transformation from start to end in a Queryable
type Operation struct {
	ParentQueryable    *Queryable
	NeedsFullEnumerate bool
	CalcOneOrAll       func(op *Operation, previous *interface{}) *OperationResult
	NewType            types.Type
}

// OperationResult is a type to hold the results of Operation.CalcOneOrAll
type OperationResult struct {
	IsFullyEnumerated bool
	FullyEnumerated   []interface{}
	SingleEnumerated  interface{}
	IsFinished        bool
}

// FindNextValue gets the next value of the Queryable, and makes any necessary changes to the Queryable in the process.
func (q *Queryable) FindNextValue() interface{} {
	working := q.Base[q.EnumerateIndex]
	for i, o := range q.Operations {
		if o.NeedsFullEnumerate {
			// fully enumerate all previous steps up to here, AND this step.
			for j := 0; j <= i; i++ {
				*q = *q.calculateFirstStep()
			}
			continue // we're done now.
		}

		calculated := o.CalcOneOrAll(&o, &working)
		working = calculated.SingleEnumerated
	}
	q.EnumerateIndex++

	return working
}

// calculateFirstStep fully enumerates the first operation and returns a new Queryable from that
func (q *Queryable) calculateFirstStep() *Queryable {
	op := q.Operations[0]
	return &Queryable{
		Type:       op.NewType,
		Length:     q.Length,
		Base:       op.getAllResults(),
		Operations: q.Operations[1:],
	}
}

// getAllResults fully enumerates an operation
func (o *Operation) getAllResults() []interface{} {
	if o.NeedsFullEnumerate {
		calculated := o.CalcOneOrAll(o, &o.ParentQueryable.Base[o.ParentQueryable.EnumerateIndex])
		if calculated.IsFullyEnumerated {
			return calculated.FullyEnumerated
		}
		return []interface{}{}
	}

	var results []interface{}
	for true {
		calculated := o.CalcOneOrAll(o, &o.ParentQueryable.Base[o.ParentQueryable.EnumerateIndex])
		if calculated.IsFinished {
			break
		}
		results = append(results, calculated.SingleEnumerated)
	}

	return results
}

// EnumerateAll takes a Queryable and returns a slice with the final result of the transformation.
func (q *Queryable) EnumerateAll() []interface{} {
}
