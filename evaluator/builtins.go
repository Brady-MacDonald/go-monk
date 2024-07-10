package evaluator

import "monkey/object"

var BuiltIns = map[string]*object.Builtin{
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

			return &object.Error{
				Message: "Unsupported arg type",
			}
		},
	},
}
