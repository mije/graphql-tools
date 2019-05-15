package compare

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

const (
	DirectiveAdded                       = ChangeType("DIRECTIVE_ADDED")
	DirectiveRemoved                     = ChangeType("DIRECTIVE_REMOVED")
	DirectiveDescriptionChanged          = ChangeType("DIRECTIVE_DESCRIPTION_CHANGED")
	DirectiveLocationAdded               = ChangeType("DIRECTIVE_LOCATION_ADDED")
	DirectiveLocationRemoved             = ChangeType("DIRECTIVE_LOCATION_REMOVED")
	DirectiveArgumentAdded               = ChangeType("DIRECTIVE_ARGUMENT_ADDED")
	DirectiveArgumentRemoved             = ChangeType("DIRECTIVE_ARGUMENT_REMOVED")
	DirectiveArgumentDescriptionChanged  = ChangeType("DIRECTIVE_ARGUMENT_DESCRIPTION_CHANGED")
	DirectiveArgumentDefaultValueChanged = ChangeType("DIRECTIVE_ARGUMENT_DEFAULT_VALUE_CHANGED")
	DirectiveArgumentTypeChanged         = ChangeType("DIRECTIVE_ARGUMENT_TYPE_CHANGED")
)

func directiveAdded(y *ast.DirectiveDefinition) Change {
	return Change{
		Type: DirectiveAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Directive '%s' was added", y.Name),
		Path:    y.Name,
	}
}

func directiveRemoved(x *ast.DirectiveDefinition) Change {
	return Change{
		Type: DirectiveRemoved,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Directive '%s' was removed", x.Name),
		Path:    x.Name,
	}
}

func directiveDescriptionChanged(x, y *ast.DirectiveDefinition) Change {
	return Change{
		Type: DirectiveDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Directive '%s' description changed from '%s' to '%s'", x.Name, x.Description, y.Description),
		Path:    x.Name,
	}
}

func directiveLocationAdded(y *ast.DirectiveDefinition, loc ast.DirectiveLocation) Change {
	return Change{
		Type: DirectiveLocationAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Location '%s' was added to directive '%s'", loc, y.Name),
		Path:    y.Name,
	}
}

func directiveLocationRemoved(y *ast.DirectiveDefinition, loc ast.DirectiveLocation) Change {
	return Change{
		Type: DirectiveLocationRemoved,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Location '%s' was removed from directive '%s'", loc, y.Name),
		Path:    y.Name,
	}
}

func directiveArgumentAdded(y *ast.DirectiveDefinition, arg *ast.ArgumentDefinition) Change {
	c := Change{
		Type: DirectiveArgumentAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Argument '%s' was added to directive '%s'", arg.Name, y.Name),
		Path:    y.Name,
	}

	if arg.Type.NonNull {
		c.Severity.Level = Breaking
	}

	return c
}

func directiveArgumentRemoved(x *ast.DirectiveDefinition, arg *ast.ArgumentDefinition) Change {
	return Change{
		Type: DirectiveArgumentRemoved,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Argument '%s' was removed from directive '%s'", arg.Name, x.Name),
		Path:    strings.Join([]string{x.Name, arg.Name}, "."),
	}
}

func directiveArgumentDescriptionChanged(d *ast.DirectiveDefinition, x, y *ast.ArgumentDefinition) Change {
	return Change{
		Type: DirectiveArgumentDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Description for argument '%s' on directive '%s' changed from '%s' to '%s'", x.Name, d.Name, x.Description, y.Description),
		Path:    strings.Join([]string{d.Name, x.Name}, "."),
	}
}

func directiveArgumentDefaultValueChanged(def *ast.DirectiveDefinition, x, y *ast.ArgumentDefinition) Change {
	return Change{
		Type: DirectiveArgumentDefaultValueChanged,
		Severity: ChangeSeverity{
			Level: Breaking, //TODO Assess the value change, it may change the severity level.
		},
		Message: fmt.Sprintf("Default value for argument '%s' on directive '%s' changed from %v to %v", x.Name, def.Name, x.DefaultValue.String(), y.DefaultValue.String()),
		Path:    strings.Join([]string{def.Name, x.Name}, "."),
	}
}

func directiveArgumentTypeChanged(def *ast.DirectiveDefinition, x, y *ast.ArgumentDefinition) Change {
	return Change{
		Type: DirectiveArgumentTypeChanged,
		Severity: ChangeSeverity{
			Level: Breaking, //TODO Asses the type change, it may change the severity level.
		},
		Message: fmt.Sprintf("Type for argument '%s' on directive '%s' changed from '%s' to '%s'", x.Name, def.Name, x.Type.String(), y.Type.String()),
		Path:    strings.Join([]string{def.Name, x.Name}, "."),
	}
}
