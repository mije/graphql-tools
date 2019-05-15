package compare

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

const (
	TypeAdded              = ChangeType("TYPE_ADDED")
	TypeRemoved            = ChangeType("TYPE_REMOVED")
	TypeKindChanged        = ChangeType("TYPE_KIND_CHANGED")
	TypeDescriptionChanged = ChangeType("TYPE_DESCRIPTION_CHANGED")
)

func typeAdded(y *ast.Definition) Change {
	return Change{
		Type: TypeAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Type '%s' was added", y.Name),
		Path:    strings.Join([]string{y.Name}, "."),
	}
}

func typeRemoved(x *ast.Definition) Change {
	return Change{
		Type: TypeRemoved,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Removing a type is a breaking change. It is preferable to deprecate and remove all references to this type first.",
		},
		Message: fmt.Sprintf("Type '%s' was removed", x.Name),
		Path:    x.Name,
	}
}

func typeKindChanged(x, y *ast.Definition) Change {
	return Change{
		Type: TypeKindChanged,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Changing the kind of a type is a breaking change because it can cause existing queries to error.reportChange( For example, turning an object type to a scalar type would break queries that define a selection set for this type.",
		},
		Message: fmt.Sprintf("'%s' kind changed from '%s' to '%s'", x.Name, x.Kind, y.Kind),
		Path:    x.Name,
	}
}

func typeDescriptionChanged(x, y *ast.Definition) Change {
	return Change{
		Type: TypeDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Description on type '%s' has changed from '%s' to '%s'", x.Name, x.Description, y.Description),
		Path:    x.Name,
	}
}
