package object

import "fmt"

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Hash() DictKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return DictKey{Type: b.Type(), Value: value}
}
