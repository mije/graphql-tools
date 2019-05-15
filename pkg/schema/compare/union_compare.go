package compare

import (
	"github.com/vektah/gqlparser/ast"
)

func (r *Result) compareUnion(x, y *ast.Definition) {
	{ // Types
		cast := func(v interface{}) string {
			return v.(string)
		}
		res := compareSlices(x.Types, y.Types, func(x, y interface{}) bool {
			return cast(x) == cast(y)
		})
		for _, def := range res.added {
			r.reportChange(unionMemberAdded(x, cast(def)))
		}
		for _, def := range res.removed {
			r.reportChange(unionMemberRemoved(x, cast(def)))
		}
	}

	// TODO Directives
}
