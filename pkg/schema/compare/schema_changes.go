package compare

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
)

const (
	SchemaQueryTypeChanged        = ChangeType("SCHEMA_QUERY_TYPE_CHANGED")
	SchemaMutationTypeChanged     = ChangeType("SCHEMA_MUTATION_TYPE_CHANGED")
	SchemaMutationTypeRemoved     = ChangeType("SCHEMA_MUTATION_TYPE_REMOVED")
	SchemaSubscriptionTypeChanged = ChangeType("SCHEMA_SUBSCRIPTION_TYPE_CHANGED")
	SchemaSubscriptionTypeRemoved = ChangeType("SCHEMA_SUBSCRIPTION_TYPE_REMOVED")
)

func schemaQueryTypeChanged(x, y *ast.OperationTypeDefinition) Change {
	return Change{
		Type: SchemaQueryTypeChanged,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Schema query type has changed from '%s' to '%s'.", x.Type, y.Type),
		Path:    x.Type,
	}
}

func schemaMutationTypeChanged(x, y *ast.OperationTypeDefinition) Change {
	return Change{
		Type: SchemaMutationTypeChanged,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Schema mutation type has changed from '%s' to '%s'.", x.Type, y.Type),
		Path:    x.Type,
	}
}

func schemaMutationTypeRemoved(x *ast.OperationTypeDefinition) Change {
	return Change{
		Type: SchemaMutationTypeRemoved,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Schema mutation type was removed."),
		Path:    x.Type,
	}
}

func schemaSubscriptionTypeChanged(x, y *ast.OperationTypeDefinition) Change {
	return Change{
		Type: SchemaSubscriptionTypeChanged,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Schema subscription type has changed from '%s' to '%s'.", x.Type, y.Type),
		Path:    x.Type,
	}
}

func schemaSubscriptionTypeRemoved(x *ast.OperationTypeDefinition) Change {
	return Change{
		Type: SchemaSubscriptionTypeRemoved,
		Severity: ChangeSeverity{
			Level: Breaking,
		},
		Message: fmt.Sprintf("Schema subscription type was removed."),
		Path:    x.Type,
	}
}
