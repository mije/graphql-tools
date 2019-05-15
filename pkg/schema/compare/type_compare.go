package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareType(x, y *ast.Definition) {
	if x.Description != y.Description {
		r.reportChange(typeDescriptionChanged(x, y))
	}

	if x.Kind != y.Kind {
		r.reportChange(typeKindChanged(x, y))
		return
	}

	switch x.Kind {
	case ast.Scalar:
		r.compareScalar(x, y)
	case ast.Object:
		r.compareObject(x, y)
	case ast.Interface:
		r.compareInterface(x, y)
	case ast.Union:
		r.compareUnion(x, y)
	case ast.Enum:
		r.compareEnum(x, y)
	case ast.InputObject:
		r.compareInput(x, y)
	}
}
