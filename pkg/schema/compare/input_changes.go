package compare

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

const (
	InputFieldAdded               = ChangeType("INPUT_FIELD_ADDED")
	InputFieldRemoved             = ChangeType("INPUT_FIELD_REMOVED")
	InputFieldDescriptionChanged  = ChangeType("INPUT_FIELD_DESCRIPTION_CHANGED")
	InputFieldDefaultValueChanged = ChangeType("INPUT_FIELD_DEFAULT_VALUE_CHANGED")
	InputFieldTypeChanged         = ChangeType("INPUT_FIELD_TYPE_CHANGED")
)

func inputFieldAdded(i *ast.Definition, y *ast.FieldDefinition) Change {
	c := Change{
		Type: InputFieldAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Input field '%s' was added to input object type '%s'", y.Name, i.Name),
		Path:    strings.Join([]string{i.Name, y.Name}, "."),
	}

	if y.Type.NonNull {
		c.Severity = ChangeSeverity{
			Level:  Breaking,
			Reason: "Adding a non-null field to an existing input type will cause existing queries that use this input type to error.",
		}
	}

	return c
}

func inputFieldRemoved(i *ast.Definition, x *ast.FieldDefinition) Change {
	return Change{
		Type: InputFieldRemoved,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Removing an input field will cause existing queries that use this input field to error",
		},
		Message: fmt.Sprintf("Input field '%s' was removed from input object type '%s'", x.Name, i.Name),
		Path:    strings.Join([]string{i.Name, x.Name}, "."),
	}
}

func inputFieldDescriptionChanged(i *ast.Definition, x, y *ast.FieldDefinition) Change {
	return Change{
		Type: InputFieldDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Input field '%s.%s' description changed from '%s' to '%s'", i.Name, x.Name, x.Description, y.Description),
		Path:    strings.Join([]string{i.Name, x.Name}, "."),
	}
}

func inputFieldDefaultValueChanged(i *ast.Definition, x, y *ast.FieldDefinition) Change {
	return Change{
		Type: InputFieldDefaultValueChanged,
		Severity: ChangeSeverity{
			Level:  Dangerous,
			Reason: "Changing the default value for an input field may change the runtime behaviour of a field if it was never provided.",
		},
		Message: fmt.Sprintf("Input field '%s.%s' default value changed from '%v' to '%v'", i.Name, x.Name, x.DefaultValue.String(), y.DefaultValue.String()),
		Path:    strings.Join([]string{i.Name, x.Name}, "."),
	}
}

func inputFieldTypeChanged(i *ast.Definition, x, y *ast.FieldDefinition) Change {
	c := Change{
		Type: InputFieldTypeChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Input field '%s.%s' changed type from '%s' to '%s'", i.Name, x.Name, x.Type.String(), y.Type.String()),
		Path:    strings.Join([]string{i.Name, x.Name}, "."),
	}

	if isBreakingTypeChange(x.Type, y.Type, true) {
		c.Severity.Level = Breaking
	}

	return c
}
