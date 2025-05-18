package evaluator

import "github.com/dr8co/monke/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					return arg.Elements[length-1]
				}
				return NULL

			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1)
					copy(newElements, arg.Elements[1:length])
					return &object.Array{Elements: newElements}
				}
				return NULL
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newElements := make([]object.Object, length+1)
				copy(newElements, arg.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}

			default:
				return newError("argument to `push` not supported, got %s", args[0].Type())

			}
		},
	},
}
