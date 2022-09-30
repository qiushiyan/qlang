package eval

import (
	"github.com/qiushiyan/qlang/pkg/ast"
	"github.com/qiushiyan/qlang/pkg/object"
)

func evalIndexExpression(node *ast.IndexExpression, env *object.Env) object.Object {
	left := Eval(node.Left, env)
	if object.IsError(left) {
		return left
	}
	index := Eval(node.Index, env)
	if object.IsError(index) {
		return index
	}
	switch left.(type) {
	case object.IVector:
		if index.Type() == object.NUMBER_OBJ || index.Type() == object.RANGE_OBJ {
			return evalVectorIndexExpression(left, index)
		} else {
			return object.NewError("index must be a number or range")
		}
	case *object.Dict:
		switch index.(type) {
		case object.Hashable:
			return evalDictIndexExpression(left, index)
		default:
			return object.NewError("index must be a hashable type")
		}
	default:
		return object.NewError("invliad index operation for type %s", left.Type())
	}
}

func evalVectorIndexExpression(v object.Object, index object.Object) object.Object {
	switch v.Type() {
	case object.VECTOR_OBJ:
		vv := v.(*object.Vector)
		if start, end, valid := getIndexBounds(index, vv.Length()); valid {
			if end == 0 {
				return vv.Elements[start]
			} else {
				return &object.Vector{BaseVector: object.BaseVector{Elements: vv.Elements[start:end]}}
			}
		} else {
			return object.NewError("index out of bounds for vector")
		}
	case object.NUMERIC_VECTOR_OBJ:
		nv := v.(*object.NumericVector)
		if start, end, valid := getIndexBounds(index, nv.Length()); valid {
			if end == 0 {
				return nv.Elements[start]
			} else {
				return &object.NumericVector{BaseVector: object.BaseVector{Elements: nv.Elements[start:end]}}
			}
		} else {
			return object.NewError("index out of bounds for numeric vector")
		}
	case object.CHARACTER_VECTOR_OBJ:
		cv := v.(*object.CharacterVector)
		if start, end, valid := getIndexBounds(index, cv.Length()); valid {
			if end == 0 {
				return cv.Elements[start]
			} else {
				return &object.CharacterVector{BaseVector: object.BaseVector{Elements: cv.Elements[start:end]}}
			}
		} else {
			return object.NewError("index out of bounds for character vector")
		}
	case object.LOGICAL_VECTOR_OBJ:
		lv := v.(*object.LogicalVector)
		if start, end, valid := getIndexBounds(index, lv.Length()); valid {
			if end == 0 {
				return lv.Elements[start]
			} else {
				return &object.LogicalVector{BaseVector: object.BaseVector{Elements: lv.Elements[start:end]}}
			}
		} else {
			return object.NewError("index out of bounds for logical vector")
		}
	default:
		return object.NewError("invliad index operation for type %s", v.Type())
	}
}

// start, end, valid
func getIndexBounds(index object.Object, length int) (int, int, bool) {
	switch index.(type) {
	case *object.Number:
		idx := int(index.(*object.Number).Value) - 1
		if idx < 0 || idx >= length {
			return 0, 0, false
		}
		return idx, 0, true
	case *object.Range:
		var start, end int
		rangeObj := index.(*object.Range)
		if rangeObj.Start == -1 {
			start = 0
		} else {
			start = rangeObj.Start - 1
		}
		if rangeObj.End == -1 {
			end = length
		} else {
			end = rangeObj.End
		}
		if indexOutofBounds(start, end, length) || start > end {
			return 0, 0, false
		}

		return start, end, true
	default:
		return 0, 0, false
	}
}

func indexOutofBounds(start, end, length int) bool {
	return start < 0 || end > length
}

// index dict
func evalDictIndexExpression(d object.Object, index object.Object) object.Object {
	// index is already verified as Hashable in evalIndexExpression
	dict := d.(*object.Dict)
	key := index.(object.Hashable).Hash()
	if pair, ok := dict.Pairs[key]; ok {
		return pair.Value
	} else {
		return NULL
	}
}
