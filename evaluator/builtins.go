package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Invalid number of args, Got=%d, expected=1", len(args))
			}

			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(obj.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(obj.Value))}
			}

			return newError("Unsupported arg type to len(): Got=%s", args[0].Type())
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Invalid number of args, Got%d, want=1", len(args))
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				newError("first() only supports Array's, Got=%s", args[0].Type())
			}

			if len(arr.Value) == 0 {
				return NULL
			}

			return arr.Value[0]
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Invalid number of args, Got%d, want=1", len(args))
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				newError("last() only supports Array's, Got=%s", args[0].Type())
			}

			if len(arr.Value) == 0 {
				return NULL
			}

			return arr.Value[len(arr.Value)-1]
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Invalid number of args, Got%d, want=1", len(args))
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				newError("rest() only supports Array's, Got=%s", args[0].Type())
			}

			if len(arr.Value) == 0 {
				return NULL
			}

			return &object.Array{
				Value: arr.Value[1:len(arr.Value)],
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, val := range args {
				fmt.Println(val.Inspect())
			}

			return NULL
		},
	},
}
