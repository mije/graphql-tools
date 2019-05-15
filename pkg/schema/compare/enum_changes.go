package compare

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

const (
	EnumValueAdded   = ChangeType("ENUM_VALUE_ADDED")
	EnumValueRemoved = ChangeType("ENUM_VALUE_REMOVED")
)

func enumValueAdded(e *ast.Definition, v *ast.EnumValueDefinition) Change {
	return Change{
		Type: EnumValueAdded,
		Severity: ChangeSeverity{
			Level:  Dangerous,
			Reason: "Adding an enum value may break existing clients that were not programming defensively against an added case when querying an enum.",
		},
		Message: fmt.Sprintf("Enum value '%s' was added to enum '%s'", v.Name, e.Name),
		Path:    strings.Join([]string{e.Name, v.Name}, "."),
	}
}

func enumValueRemoved(e *ast.Definition, v *ast.EnumValueDefinition) Change {
	return Change{
		Type: EnumValueRemoved,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Removing an enum value will cause existing queries that use this enum value to error.",
		},
		Message: fmt.Sprintf("Enum value '%s' was removed from enum '%s'", v.Name, e.Name),
		Path:    strings.Join([]string{e.Name, v.Name}, "."),
	}
}
