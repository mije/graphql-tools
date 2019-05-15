package compare

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

const (
	InterfaceTypeFieldAdded                       = ChangeType("INTERFACE_TYPE_FIELD_ADDED")
	InterfaceTypeFieldRemoved                     = ChangeType("INTERFACE_TYPE_FIELD_REMOVED")
	InterfaceTypeFieldDescriptionChanged          = ChangeType("INTERFACE_TYPE_FIELD_DESCRIPTION_CHANGED")
	InterfaceTypeFieldTypeChanged                 = ChangeType("INTERFACE_TYPE_FIELD_TYPE_CHANGED")
	InterfaceTypeFieldArgumentAdded               = ChangeType("INTERFACE_TYPE_FIELD_ARGUMENT_ADDED")
	InterfaceTypeFieldArgumentRemoved             = ChangeType("INTERFACE_TYPE_FIELD_ARGUMENT_REMOVED")
	InterfaceTypeFieldArgumentDescriptionChanged  = ChangeType("INTERFACE_TYPE_FIELD_ARGUMENT_DESCRIPTION_CHANGED")
	InterfaceTypeFieldArgumentDefaultValueChanged = ChangeType("INTERFACE_TYPE_FIELD_ARGUMENT_DEFAULT_VALUE_CHANGED")
	InterfaceTypeFieldArgumentTypeChanged         = ChangeType("INTERFACE_TYPE_FIELD_ARGUMENT_TYPE_CHANGED")
)

func interfaceFieldAdded(o *ast.Definition, y *ast.FieldDefinition) Change {
	return Change{
		Type: InterfaceTypeFieldAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Field '%s' was added to interface '%s'", y.Name, o.Name),
		Path:    strings.Join([]string{o.Name, y.Name}, "."),
	}
}

func interfaceFieldRemoved(o *ast.Definition, x *ast.FieldDefinition) Change {
	return Change{
		Type: InterfaceTypeFieldRemoved,
		Severity: ChangeSeverity{
			Level: Breaking,
			//TODO Check the directives for deprecation
			//if f.directives.hasDeprecation() Change {
			//	Reason: "Removing a deprecated field is a breaking change. Before removing it, you may want to look at the field's usage to see the impact of removing the field."
			//} else{
			//	Reason: "Removing a field is a breaking change. It is preferable to deprecate the field before removing it."
			//}
		},
		Message: fmt.Sprintf("Field '%s' was removed from interface '%s'", x.Name, o.Name),
		Path:    strings.Join([]string{o.Name, x.Name}, "."),
	}
}

func interfaceFieldDescriptionChanged(o *ast.Definition, x, y *ast.FieldDefinition) Change {
	return Change{
		Type: InterfaceTypeFieldDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Field %s.%s description changed from '%s' to '%s'", o.Name, x.Name, x.Description, y.Description),
		Path:    strings.Join([]string{o.Name, x.Name}, "."),
	}
}

func interfaceFieldTypeChanged(o *ast.Definition, x, y *ast.FieldDefinition) Change {
	c := Change{
		Type: InterfaceTypeFieldTypeChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Field %s.%s changed type from '%s' to '%s'", x.Name, o.Name, x.Type.String(), y.Type.String()),
		Path:    strings.Join([]string{o.Name, x.Name}, "."),
	}

	if isBreakingTypeChange(x.Type, y.Type, false) {
		c.Severity.Level = Breaking
	}

	return c
}

func interfaceFieldArgumentAdded(o *ast.Definition, f *ast.FieldDefinition, y *ast.ArgumentDefinition) Change {
	c := Change{
		Type: InterfaceTypeFieldArgumentAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Argument '%s' was added to field '%s.%s'", y.Name, o.Name, f.Name),
		Path:    strings.Join([]string{o.Name, f.Name, y.Name}, "."),
	}

	if y.Type.NonNull {
		c.Severity = ChangeSeverity{
			Level:  Breaking,
			Reason: "Adding a required argument to an existing field is a breaking change because it will cause existing uses of this field to error.",
		}
	}

	return c
}

func interfaceFieldArgumentRemoved(o *ast.Definition, f *ast.FieldDefinition, x *ast.ArgumentDefinition) Change {
	return Change{
		Type: InterfaceTypeFieldArgumentRemoved,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Removing a field argument is a breaking change because it will cause existing queries that use this argument to error.",
		},
		Message: fmt.Sprintf("Argument '%s' was removed from field '%s.%s'", x.Name, o.Name, f.Name),
		Path:    strings.Join([]string{o.Name, f.Name, x.Name}, "."),
	}
}

func interfaceFieldArgumentDescriptionChanged(o *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) Change {
	return Change{
		Type: InterfaceTypeFieldArgumentDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Description for argument '%s' on field '%s.%s' changed from '%s' to '%s'", x.Name, o.Name, f.Name, x.Description, y.Description),
		Path:    strings.Join([]string{o.Name, f.Name, x.Name}, "."),
	}
}

func interfaceFieldArgumentDefaultValueChanged(o *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) Change {
	c := Change{
		Type: InterfaceTypeFieldArgumentDefaultValueChanged,
		Severity: ChangeSeverity{
			Level:  Dangerous,
			Reason: "Changing the default value for an argument may change the runtime behaviour of a field if it was never provided.",
		},
		Path: strings.Join([]string{o.Name, x.Name}, "."),
	}

	if x.DefaultValue == nil {
		c.Message = fmt.Sprintf("Default value '%v' was added to argument '%s' on field '%s.%s'", y.DefaultValue.String(), y.Name, o.Name, f.Name)
	} else {
		c.Message = fmt.Sprintf("Default value for argument '%s' on field '%s.%s' changed from '%v' to '%v'", y.Name, o.Name, f.Name, x.DefaultValue.String(), y.DefaultValue.String())
	}

	return c
}

func interfaceFieldArgumentTypeChanged(o *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) Change {
	c := Change{
		Type: InterfaceTypeFieldArgumentTypeChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Type for argument '%s' on field '%s.%s' changed from '%s' to '%s'", x.Name, o.Name, f.Name, x.Type.String(), y.Type.String()),
		Path:    strings.Join([]string{o.Name, f.Name, x.Name}, "."),
	}

	if isBreakingTypeChange(x.Type, y.Type, false) {
		c.Severity.Level = Breaking
	}

	return c
}
