package compare

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
)

const (
	UnionMemberRemoved = ChangeType("UNION_MEMBER_REMOVED")
	UnionMemberAdded   = ChangeType("UNION_MEMBER_ADDED")
)

func unionMemberAdded(u *ast.Definition, member string) Change {
	return Change{
		Type: UnionMemberAdded,
		Severity: ChangeSeverity{
			Level:  Dangerous,
			Reason: "Adding a possible type to Unions may break existing clients that were not programming defensively against a new possible type.",
		},
		Message: fmt.Sprintf("Union member '%s' was added to union type '%s'", member, u.Name),
		Path:    u.Name,
	}
}

func unionMemberRemoved(u *ast.Definition, member string) Change {
	return Change{
		Type: UnionMemberRemoved,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Removing a union member from a union can cause existing queries that use this union member in a fragment spread to error.",
		},
		Message: fmt.Sprintf("Union member '%s' was removed from union type '%s'", member, u.Name),
		Path:    u.Name,
	}
}
