package compare

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/ast"
)

const (
	ObjectTypeInterfaceAdded                   = ChangeType("OBJECT_TYPE_INTERFACE_ADDED")
	ObjectTypeInterfaceRemoved                 = ChangeType("OBJECT_TYPE_INTERFACE_REMOVED")
	ObjectTypeFieldAdded                       = ChangeType("OBJECT_TYPE_FIELD_ADDED")
	ObjectTypeFieldRemoved                     = ChangeType("OBJECT_TYPE_FIELD_REMOVED")
	ObjectTypeFieldDescriptionChanged          = ChangeType("OBJECT_TYPE_FIELD_DESCRIPTION_CHANGED")
	ObjectTypeFieldTypeChanged                 = ChangeType("OBJECT_TYPE_FIELD_TYPE_CHANGED")
	ObjectTypeFieldArgumentAdded               = ChangeType("OBJECT_TYPE_FIELD_ARGUMENT_ADDED")
	ObjectTypeFieldArgumentRemoved             = ChangeType("OBJECT_TYPE_FIELD_ARGUMENT_REMOVED")
	ObjectTypeFieldArgumentDescriptionChanged  = ChangeType("OBJECT_TYPE_FIELD_ARGUMENT_DESCRIPTION_CHANGED")
	ObjectTypeFieldArgumentDefaultValueChanged = ChangeType("OBJECT_TYPE_FIELD_ARGUMENT_DEFAULT_VALUE_CHANGED")
	ObjectTypeFieldArgumentTypeChanged         = ChangeType("OBJECT_TYPE_FIELD_ARGUMENT_TYPE_CHANGED")
)

func objectTypeInterfaceAdded(o *ast.Definition, inf string) Change {
	return Change{
		Type: ObjectTypeInterfaceAdded,
		Severity: ChangeSeverity{
			Level:  Dangerous,
			Reason: "Adding an interface to an object type may break existing clients that were not programming defensively against a new possible type.",
		},
		Message: fmt.Sprintf("'%s' object type implements interface '%s'", o.Name, inf),
		Path:    o.Name,
	}
}

func objectTypeInterfaceRemoved(o *ast.Definition, inf string) Change {
	return Change{
		Type: ObjectTypeInterfaceRemoved,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Removing an interface from an object type can cause existing queries that use this in a fragment spread to error.",
		},
		Message: fmt.Sprintf("'%s' object type no longer implements interface '%s'", o.Name, inf),
		Path:    o.Name,
	}
}

func objectFieldAdded(o *ast.Definition, y *ast.FieldDefinition) Change {
	return Change{
		Type: ObjectTypeFieldAdded,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Field '%s' was added to type '%s'", y.Name, o.Name),
		Path:    strings.Join([]string{o.Name, y.Name}, "."),
	}
}

func objectFieldRemoved(o *ast.Definition, x *ast.FieldDefinition) Change {
	return Change{
		Type: ObjectTypeFieldRemoved,
		Severity: ChangeSeverity{
			Level: Breaking,
			//TODO Check the directives for deprecation
			//if f.directives.hasDeprecation() Change {
			//	Reason: "Removing a deprecated field is a breaking change. Before removing it, you may want to look at the field's usage to see the impact of removing the field.'%s'
			//} else{
			//	Reason: "Removing a field is a breaking change. It is preferable to deprecate the field before removing it."
			//}
		},
		Message: fmt.Sprintf("Field '%s' was removed from type '%s'", x.Name, o.Name),
		Path:    strings.Join([]string{o.Name, x.Name}, "."),
	}
}

func objectFieldDescriptionChanged(o *ast.Definition, x, y *ast.FieldDefinition) Change {
	return Change{
		Type: ObjectTypeFieldDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Field '%s.%s' description changed from '%s' to '%s'", o.Name, x.Name, x.Description, y.Description),
		Path:    strings.Join([]string{o.Name, x.Name}, "."),
	}
}

func objectFieldTypeChanged(o *ast.Definition, x, y *ast.FieldDefinition) Change {
	c := Change{
		Type: ObjectTypeFieldTypeChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Field '%s.%s' changed type from '%s' to '%s'", x.Name, o.Name, x.Type.String(), y.Type.String()),
		Path:    strings.Join([]string{o.Name, x.Name}, "."),
	}

	if isBreakingTypeChange(x.Type, y.Type, false) {
		c.Severity.Level = Breaking
	}

	return c
}

func objectFieldArgumentAdded(o *ast.Definition, f *ast.FieldDefinition, y *ast.ArgumentDefinition) Change {
	c := Change{
		Type: ObjectTypeFieldArgumentAdded,
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

func objectFieldArgumentRemoved(o *ast.Definition, f *ast.FieldDefinition, x *ast.ArgumentDefinition) Change {
	return Change{
		Type: ObjectTypeFieldArgumentRemoved,
		Severity: ChangeSeverity{
			Level:  Breaking,
			Reason: "Removing a field argument is a breaking change because it will cause existing queries that use this argument to error.",
		},
		Message: fmt.Sprintf("Argument '%s' was removed from field '%s.%s'", x.Name, o.Name, f.Name),
		Path:    strings.Join([]string{o.Name, f.Name, x.Name}, "."),
	}
}

func objectFieldArgumentDescriptionChanged(o *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) Change {
	return Change{
		Type: ObjectTypeFieldArgumentDescriptionChanged,
		Severity: ChangeSeverity{
			Level: NonBreaking,
		},
		Message: fmt.Sprintf("Description for argument '%s' on field '%s.%s' changed from '%s' to '%s'", x.Name, o.Name, f.Name, x.Description, y.Description),
		Path:    strings.Join([]string{o.Name, f.Name, x.Name}, "."),
	}
}

func objectFieldArgumentDefaultValueChanged(o *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) Change {
	c := Change{
		Type: ObjectTypeFieldArgumentDefaultValueChanged,
		Severity: ChangeSeverity{
			Level:  Dangerous,
			Reason: "Changing the default value for an argument may change the runtime behaviour of a field if it was never provided.",
		},
		Path: strings.Join([]string{o.Name, x.Name}, "."),
	}

	if x.DefaultValue == nil {
		c.Message = fmt.Sprintf("Default value '%s' was added to argument '%s' on field '%s.%s'", y.DefaultValue.String(), y.Name, o.Name, f.Name)
	} else {
		c.Message = fmt.Sprintf("Default value for argument '%s' on field '%s.%s' changed from '%v' to '%v'", y.Name, o.Name, f.Name, x.DefaultValue.String(), y.DefaultValue.String())
	}

	return c
}

func objectFieldArgumentTypeChanged(o *ast.Definition, f *ast.FieldDefinition, x, y *ast.ArgumentDefinition) Change {
	c := Change{
		Type: ObjectTypeFieldArgumentTypeChanged,
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
