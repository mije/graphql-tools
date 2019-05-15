package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareInterface(x, y *ast.Definition) {
	{ // Fields
		cast := func(v interface{}) *ast.FieldDefinition {
			return v.(*ast.FieldDefinition)
		}
		res := compareSlices(x.Fields, y.Fields, func(x, y interface{}) bool {
			return cast(x).Name == cast(y).Name
		})
		for _, def := range res.added {
			r.reportChange(interfaceFieldAdded(x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(interfaceFieldRemoved(x, cast(def)))
		}
		for _, pair := range res.common {
			r.compareInterfaceField(x, cast(pair.x), cast(pair.y))
		}
	}

	// TODO Directives
}

func (r *Result) compareInterfaceField(i *ast.Definition, x, y *ast.FieldDefinition) {
	// Description
	if x.Description != y.Description {
		r.reportChange(interfaceFieldDescriptionChanged(i, x, y))
	}

	// Type
	if !typeEquals(x.Type, y.Type) {
		r.reportChange(interfaceFieldTypeChanged(i, x, y))
	}

	{ // Arguments
		cast := func(v interface{}) *ast.ArgumentDefinition {
			return v.(*ast.ArgumentDefinition)
		}
		res := compareSlices(x.Arguments, y.Arguments, func(x, y interface{}) bool {
			return cast(x).Name == cast(y).Name
		})
		for _, def := range res.added {
			r.reportChange(interfaceFieldArgumentAdded(i, x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(interfaceFieldArgumentRemoved(i, x, cast(def)))
		}
		for _, pair := range res.common {
			r.compareInterfaceFieldArgument(i, x, cast(pair.x), cast(pair.y))
		}
	}

	// TODO Directives
}

func (r *Result) compareInterfaceFieldArgument(i *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) {
	// Description
	if x.Description != y.Description {
		r.reportChange(interfaceFieldArgumentDescriptionChanged(i, f, x, y))
	}

	// Type
	if !typeEquals(x.Type, y.Type) {
		r.reportChange(interfaceFieldArgumentTypeChanged(i, f, x, y))
	}

	// Default value
	if !valueEquals(x.DefaultValue, y.DefaultValue) {
		r.reportChange(interfaceFieldArgumentDefaultValueChanged(i, f, x, y))
	}
}
