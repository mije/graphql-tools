package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareObject(x, y *ast.Definition) {
	{ // Interfaces
		cast := func(v interface{}) string {
			return v.(string)
		}
		res := compareSlices(x.Interfaces, y.Interfaces, func(x, y interface{}) bool {
			return cast(x) == cast(y)
		})
		for _, def := range res.added {
			r.reportChange(objectTypeInterfaceAdded(y, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(objectTypeInterfaceRemoved(x, cast(def)))
		}
	}

	{ // Fields
		cast := func(v interface{}) *ast.FieldDefinition {
			return v.(*ast.FieldDefinition)
		}
		res := compareSlices(x.Fields, y.Fields, func(x, y interface{}) bool {
			return cast(x).Name == cast(y).Name
		})
		for _, def := range res.added {
			r.reportChange(objectFieldAdded(x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(objectFieldRemoved(x, cast(def)))
		}
		for _, pair := range res.common {
			r.compareObjectField(x, cast(pair.x), cast(pair.y))
		}
	}

	// TODO Directives
}

func (r *Result) compareObjectField(o *ast.Definition, x, y *ast.FieldDefinition) {
	// Description
	if x.Description != y.Description {
		r.reportChange(objectFieldDescriptionChanged(o, x, y))
	}

	// Type
	if !typeEquals(x.Type, y.Type) {
		r.reportChange(objectFieldTypeChanged(o, x, y))
	}

	{ // Arguments
		cast := func(v interface{}) *ast.ArgumentDefinition {
			return v.(*ast.ArgumentDefinition)
		}
		res := compareSlices(x.Arguments, y.Arguments, func(x, y interface{}) bool {
			return cast(x).Name == cast(y).Name
		})
		for _, def := range res.added {
			r.reportChange(objectFieldArgumentAdded(o, x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(objectFieldArgumentRemoved(o, x, cast(def)))
		}
		for _, pair := range res.common {
			r.compareObjectFieldArgument(o, x, cast(pair.x), cast(pair.y))
		}
	}

	// TODO Directives
}

func (r *Result) compareObjectFieldArgument(o *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) {
	// Description
	if x.Description != y.Description {
		r.reportChange(objectFieldArgumentDescriptionChanged(o, f, x, y))
	}

	// Type
	if !typeEquals(x.Type, y.Type) {
		r.reportChange(objectFieldArgumentTypeChanged(o, f, x, y))
	}

	// Default Value
	if !valueEquals(x.DefaultValue, y.DefaultValue) {
		r.reportChange(objectFieldArgumentDefaultValueChanged(o, f, x, y))
	}

	// TODO Directives
}
