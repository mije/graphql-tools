package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareInput(x, y *ast.Definition) {
	{ // Fields
		cast := func(v interface{}) *ast.FieldDefinition {
			return v.(*ast.FieldDefinition)
		}
		res := compareSlices(x.Fields, y.Fields, func(x, y interface{}) bool {
			return cast(x).Name == cast(y).Name
		})
		for _, def := range res.added {
			r.reportChange(inputFieldAdded(x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(inputFieldRemoved(x, cast(def)))
		}
		for _, pair := range res.common {
			r.compareInputField(x, cast(pair.x), cast(pair.y))
		}
	}

	// TODO Directives
}

func (r *Result) compareInputField(i *ast.Definition, x, y *ast.FieldDefinition) {
	// Description
	if x.Description != y.Description {
		r.reportChange(inputFieldDescriptionChanged(i, x, y))
	}

	// Type
	if !typeEquals(x.Type, y.Type) {
		r.reportChange(inputFieldTypeChanged(i, x, y))
	}

	// Default value
	if !valueEquals(x.DefaultValue, y.DefaultValue) {
		r.reportChange(inputFieldDefaultValueChanged(i, x, y))
	}

	// TODO Directives
}
