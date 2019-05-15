package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareEnum(x, y *ast.Definition) {
	{ // Values
		cast := func(v interface{}) *ast.EnumValueDefinition {
			return v.(*ast.EnumValueDefinition)
		}
		res := compareSlices(x.EnumValues, y.EnumValues, func(x, y interface{}) bool {
			return cast(x).Name == cast(y).Name
		})
		for _, def := range res.added {
			r.reportChange(enumValueAdded(x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(enumValueRemoved(x, cast(def)))
		}
		for _, pair := range res.common {
			r.compareEnumValue(x, cast(pair.x), cast(pair.y))
		}
	}

	// TODO Directives
}

func (r *Result) compareEnumValue(e *ast.Definition, x, y *ast.EnumValueDefinition) {
	// TODO Directives
}
