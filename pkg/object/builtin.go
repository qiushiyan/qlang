package object

type BuiltinFunction func(env *Env, args ...Object) Object
type Builtin struct {
	RequiredParametersNum int
	Fn                    BuiltinFunction
}

func (b *Builtin) Inspect() string  { return "builtin" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
