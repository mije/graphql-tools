package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareDirective(x, y *ast.DirectiveDefinition) {
	// Description
	if x.Description != y.Description {
		r.reportChange(directiveDescriptionChanged(x, y))
	}

	{ // Locations
		cast := func(v interface{}) ast.DirectiveLocation {
			return v.(ast.DirectiveLocation)
		}
		res := compareSlices(x.Locations, y.Locations, func(x, y interface{}) bool {
			return cast(x) == cast(y)
		})
		for _, loc := range res.added {
			r.reportChange(directiveLocationAdded(x, cast(loc)))
		}
		for _, loc := range res.removed {
			r.reportChange(directiveLocationRemoved(x, cast(loc)))
		}
	}

	{ // Arguments
		cast := func(v interface{}) *ast.ArgumentDefinition {
			return v.(*ast.ArgumentDefinition)
		}
		res := compareSlices(x.Arguments, y.Arguments, func(x, y interface{}) bool {
			return cast(x).Name == cast(y).Name
		})
		for _, def := range res.added {
			r.reportChange(directiveArgumentAdded(x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(directiveArgumentRemoved(x, cast(def)))
		}
		for _, pair := range res.common {
			r.compareDirectiveArgument(x, cast(pair.x), cast(pair.y))
		}
	}
}

func (r *Result) compareDirectiveArgument(d *ast.DirectiveDefinition, x, y *ast.ArgumentDefinition) {
	// Description
	if x.Description != y.Description {
		r.reportChange(directiveArgumentDescriptionChanged(d, x, y))
	}

	// Type
	if !typeEquals(x.Type, y.Type) {
		r.reportChange(directiveArgumentTypeChanged(d, x, y))
	}

	// Default value
	if !valueEquals(x.DefaultValue, y.DefaultValue) {
		r.reportChange(directiveArgumentDefaultValueChanged(d, x, y))
	}
}
