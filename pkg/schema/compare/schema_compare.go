package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareSchema(x, y *schema) {
	r.compareRootTypes(x.rootTypes, y.rootTypes)
	r.compareDirectives(x.directives, y.directives)
	r.compareTypes(x.types, y.types)
	//TODO extensions
}

func (r *Result) compareRootTypes(x, y map[ast.Operation]*ast.OperationTypeDefinition) {
	if xq, ok := x[ast.Query]; ok {
		if yq, ok := y[ast.Query]; ok {
			if xq.Type != yq.Type {
				r.reportChange(schemaQueryTypeChanged(xq, yq))
			}
		}
	}

	if xm, ok := x[ast.Mutation]; ok {
		if ym, ok := y[ast.Mutation]; ok {
			if xm.Type != ym.Type {
				r.reportChange(schemaMutationTypeChanged(xm, ym))
			}
		} else {
			r.reportChange(schemaMutationTypeRemoved(xm))
		}
	}

	if xs, ok := x[ast.Subscription]; ok {
		if ys, ok := y[ast.Subscription]; ok {
			if xs.Type != ys.Type {
				r.reportChange(schemaSubscriptionTypeChanged(xs, ys))
			}
		} else {
			r.reportChange(schemaSubscriptionTypeRemoved(xs))
		}
	}
}

func (r *Result) compareDirectives(x, y map[string]*ast.DirectiveDefinition) {
	cast := func(v interface{}) *ast.DirectiveDefinition {
		return v.(*ast.DirectiveDefinition)
	}

	res := compareMaps(x, y)
	for _, def := range res.added {
		r.reportChange(directiveAdded(cast(def)))
	}
	for _, def := range res.removed {
		r.reportChange(directiveRemoved(cast(def)))
	}
	for _, pair := range res.common {
		r.compareDirective(cast(pair.x), cast(pair.y))
	}
}

func (r *Result) compareTypes(x, y map[string]*ast.Definition) {
	cast := func(v interface{}) *ast.Definition {
		return v.(*ast.Definition)
	}

	res := compareMaps(x, y)
	for _, def := range res.added {
		r.reportChange(typeAdded(cast(def)))
	}
	for _, def := range res.removed {
		r.reportChange(typeRemoved(cast(def)))
	}
	for _, pair := range res.common {
		r.compareType(cast(pair.x), cast(pair.y))
	}
}
