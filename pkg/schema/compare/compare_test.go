package compare

import (
	"fmt"
	"strings"
	"testing"
)

func TestSchemaCompare(t *testing.T) {
	testData := map[string][]struct {
		name string
		x, y string
		want Change
	}{
		"Schema": {
			{
				name: "Change of query type is a breaking change",
				x:    "schema { query: A } type A { a: String } type B { b: String }",
				y:    "schema { query: B } type A { a: String } type B { b: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     SchemaQueryTypeChanged,
				},
			},
			{
				name: "Change of query type is a breaking change",
				x:    "schema { mutation: A } type A { a: String } type B { b: String }",
				y:    "schema { mutation: B } type A { a: String } type B { b: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     SchemaMutationTypeChanged,
				},
			},
			{
				name: "Adding mutation type is a non-breaking change",
				x:    "schema { query: Q } type Q { q: String }",
				y:    "schema { query: Q mutation: M } type Q { q: String } type M { m: String }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     TypeAdded,
				},
			},

			{
				name: "Removing mutation type is a breaking change",
				x:    "schema { query: Q mutation: M } type Q { q: String } type M { m: String }",
				y:    "schema { query: Q } type Q { q: String } type M { m: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     SchemaMutationTypeRemoved,
				},
			},
			{
				name: "Change of subscription type is a breaking change",
				x:    "schema { subscription: A } type A { a: String } type B { b: String }",
				y:    "schema { subscription: B } type A { a: String } type B { b: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     SchemaSubscriptionTypeChanged,
				},
			},
			{
				name: "Adding subscription type is a non-breaking change",
				x:    "schema { query: Q } type Q { q: String }",
				y:    "schema { query: Q subscription: S } type Q { q: String } type S { s: String }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     TypeAdded,
				},
			},
			{
				name: "Removing subscription type is a breaking change",
				x:    "schema { query: Q subscription: S } type Q { q: String } type S { s: String }",
				y:    "schema { query: Q } type Q { q: String } type S { s: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     SchemaSubscriptionTypeRemoved,
				},
			},
		},
		"Directive": {
			// TODO Add test scenarios for directive changes.
		},
		"Type": {
			{
				name: "Adding type is a non-breaking change",
				x:    "type A { a: String }",
				y:    "type A { a: String } type B { b: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     TypeAdded,
				},
			},
			{
				name: "Removing type is a breaking change",
				x:    "type A { a: String } type B { b: Int }",
				y:    "type A { a: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     TypeRemoved,
				},
			},
			{
				name: "Changing type's kind is a breaking change",
				x:    "type A { a: String }",
				y:    "scalar A ",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     TypeKindChanged,
				},
			},
			{
				name: "Changing type's description is a non-breaking change",
				x:    "type A { a: String }",
				y:    ` "Type description" type A { a: String }`,
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     TypeDescriptionChanged,
				},
			},
		},
		"Enum": {
			{
				name: "Adding enum value may be a breaking change",
				x:    "enum E { X Y }",
				y:    "enum E { X Y Z }",
				want: Change{
					Severity: ChangeSeverity{Level: Dangerous},
					Type:     EnumValueAdded,
				},
			},
			{
				name: "Removing enum value is a breaking change",
				x:    "enum E { X Y Z }",
				y:    "enum E { X Y }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     EnumValueRemoved,
				},
			},
		},
		"Input": {
			{
				name: "Adding optional input field is a non-breaking change",
				x:    "input I { i: Int }",
				y:    "input I { i: Int f: Float }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InputFieldAdded,
				},
			},
			{
				name: "Adding mandatory input field is a breaking change",
				x:    "input I { i: Int }",
				y:    "input I { i: Int f: Float! }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InputFieldAdded,
				},
			},
			{
				name: "Removing optional input field is a breaking change",
				x:    "input I { i: Int f: Float }",
				y:    "input I { i: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InputFieldRemoved,
				},
			},
			{
				name: "Removing mandatory input field is a breaking change",
				x:    "input I { i: Int f: Float! }",
				y:    "input I { i: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InputFieldRemoved,
				},
			},
			{
				name: "Changing input field type is a breaking change",
				x:    "input I { i: Int }",
				y:    "input I { i: Float }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InputFieldTypeChanged,
				},
			},
			{
				name: "Making input field mandatory is a breaking change",
				x:    "input I { i: Int }",
				y:    "input I { i: Int! }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InputFieldTypeChanged,
				},
			},
			{
				name: "Making input field optional is a non-breaking change",
				x:    "input I { i: Int! }",
				y:    "input I { i: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InputFieldTypeChanged,
				},
			},
			{
				name: "Changing input field description is a non-breaking change",
				x:    "input I { i: Int }",
				y:    `input I { "Field description" i: Int }`,
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InputFieldDescriptionChanged,
				},
			},
			{
				name: "Changing input field's default value may be a breaking change",
				x:    "input I { i: Int = 100 }",
				y:    "input I { i: Int = 500 }",
				want: Change{
					Severity: ChangeSeverity{Level: Dangerous},
					Type:     InputFieldDefaultValueChanged,
				},
			},
		},
		"Interface": {
			{
				name: "Adding field is a non-breaking change",
				x:    "interface I { i: String }",
				y:    "interface I { i: String j: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InterfaceTypeFieldAdded,
				},
			},
			{
				name: "Removing field is a breaking change",
				x:    "interface I { i: String j: Int }",
				y:    "interface I { i: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldRemoved,
				},
			},
			{
				name: "Changing field type is a breaking change",
				x:    "interface I { i: String }",
				y:    "interface I { i: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldTypeChanged,
				},
			},
			{
				name: "Making field mandatory is a non-breaking change",
				x:    "interface I { i: String }",
				y:    "interface I { i: String! }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InterfaceTypeFieldTypeChanged,
				},
			},
			{
				name: "Making field optional is a breaking change",
				x:    "interface I { i: String! }",
				y:    "interface I { i: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldTypeChanged,
				},
			},
			{
				name: "Adding optional field argument is a non-breaking change",
				x:    "interface I { i: String }",
				y:    "interface I { i(x: Boolean): String }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InterfaceTypeFieldArgumentAdded,
				},
			},
			{
				name: "Adding mandatory field argument with a default value is a non-breaking change",
				x:    "interface i { i: String }",
				y:    "interface i { i(x: Boolean = true): String }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InterfaceTypeFieldArgumentAdded,
				},
			},
			{
				name: "Adding mandatory field argument is a breaking change",
				x:    "interface i { i: String }",
				y:    "interface i { i(x: Boolean!): String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldArgumentAdded,
				},
			},
			{
				name: "Removing optional field argument is a breaking change",
				x:    "interface I { i(x: Boolean): String }",
				y:    "interface I { i: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldArgumentRemoved,
				},
			},
			{
				name: "Removing mandatory field argument is a breaking change",
				x:    "interface I { i(x: Boolean!): String }",
				y:    "interface I { i: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldArgumentRemoved,
				},
			},
			{
				name: "Removing field argument with default value is a breaking change",
				x:    "interface I { i(x: Boolean = true): String }",
				y:    "interface I { i: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldArgumentRemoved,
				},
			},
			{
				name: "Changing field's description is a non-breaking change",
				x:    "interface I { i: String }",
				y:    `interface I { "A field description" i: String }`,
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InterfaceTypeFieldDescriptionChanged,
				},
			},
			{
				name: "Changing field's description is a non-breaking change",
				x:    "interface I { i(x: Boolean): String }",
				y:    `interface I { i( "A field argument description" x: Boolean): String }`,
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     InterfaceTypeFieldArgumentDescriptionChanged,
				},
			},
			{
				name: "Changing field's argument default value may be a breaking change",
				x:    "interface I { i(x: Boolean = true): String }",
				y:    "interface I { i(x: Boolean = false): String }",
				want: Change{
					Severity: ChangeSeverity{Level: Dangerous},
					Type:     InterfaceTypeFieldArgumentDefaultValueChanged,
				},
			},
			{
				name: "Changing field's argument type is a breaking change",
				x:    "interface I { i(s: String): String }",
				y:    "interface I { i(s: Int): String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     InterfaceTypeFieldArgumentTypeChanged,
				},
			},
		},
		"Object": {
			{
				name: "Adding field is a non-breaking change",
				x:    "type A { a: String }",
				y:    "type A { a: String b: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     ObjectTypeFieldAdded,
				},
			},
			{
				name: "Removing field is a breaking change",
				x:    "type A { a: String b: Int }",
				y:    "type A { b: Int}",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldRemoved,
				},
			},
			{
				name: "Changing field type a breaking change",
				x:    "type A { a: String }",
				y:    "type A { a: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldTypeChanged,
				},
			},
			{
				name: "Making field mandatory is a non-breaking change",
				x:    "type A { a: String }",
				y:    "type A { a: String! }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     ObjectTypeFieldTypeChanged,
				},
			},
			{
				name: "Making field optional is a breaking change",
				x:    "type A { a: String! }",
				y:    "type A { a: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldTypeChanged,
				},
			},
			{
				name: "Adding optional field argument is a non-breaking change",
				x:    "type A { a: String }",
				y:    "type A { a(x: Boolean): String }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     ObjectTypeFieldArgumentAdded,
				},
			},
			{
				name: "Adding mandatory field argument with a default value is a non-breaking change",
				x:    "type A { a: String }",
				y:    "type A { a(x: Boolean = true): String }",
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     ObjectTypeFieldArgumentAdded,
				},
			},
			{
				name: "Adding mandatory field argument is a breaking change",
				x:    "type A { a: String }",
				y:    "type A { a(x: Boolean!): String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldArgumentAdded,
				},
			},
			{
				name: "Removing optional field argument is a breaking change",
				x:    "type A { a(x: Boolean): String }",
				y:    "type A { a: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldArgumentRemoved,
				},
			},
			{
				name: "Removing mandatory field argument is a breaking change",
				x:    "type A { a(x: Boolean!): String }",
				y:    "type A { a: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldArgumentRemoved,
				},
			},
			{
				name: "Removing field argument with default value is a breaking change",
				x:    "type A { a(x: Boolean = true): String }",
				y:    "type A { a: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldArgumentRemoved,
				},
			},
			{
				name: "Adding interface may be a breaking change",
				x:    "interface I { a: String } interface J { b: String } type A implements I { a: String b: String }",
				y:    "interface I { a: String } interface J { b: String } type A implements I & J { a: String b: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Dangerous},
					Type:     ObjectTypeInterfaceAdded,
				},
			},
			{
				name: "Removing interface is a breaking change",
				x:    "interface I { a: String } interface J { b: String } type A implements I & J { a: String b: String }",
				y:    "interface I { a: String } interface J { b: String } type A implements I { a: String b: String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeInterfaceRemoved,
				},
			},
			{
				name: "Changing field's description is a non-breaking change",
				x:    "type A { a: String }",
				y:    `type A { "A field description" a: String }`,
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     ObjectTypeFieldDescriptionChanged,
				},
			},
			{
				name: "Changing field's description is a non-breaking change",
				x:    "type A { a(x: Boolean): String }",
				y:    `type A { a( "A field argument description" x: Boolean): String }`,
				want: Change{
					Severity: ChangeSeverity{Level: NonBreaking},
					Type:     ObjectTypeFieldArgumentDescriptionChanged,
				},
			},
			{
				name: "Changing field's argument default value may be a breaking change",
				x:    "type A { a(x: Boolean = true): String }",
				y:    "type A { a(x: Boolean = false): String }",
				want: Change{
					Severity: ChangeSeverity{Level: Dangerous},
					Type:     ObjectTypeFieldArgumentDefaultValueChanged,
				},
			},
			{
				name: "Changing field's argument type is a breaking change",
				x:    "type A { a(s: String): String }",
				y:    "type A { a(s: Int): String }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     ObjectTypeFieldArgumentTypeChanged,
				},
			},
		},
		"Scalar": {
			// TODO Add test scenarios for scalar changes.
		},
		"Union": {
			{
				name: "Adding union member may be a breaking change",
				x:    "union U = A type A { a: String } type B { b: Int }",
				y:    "union U = A | B type A { a: String } type B { b: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: Dangerous},
					Type:     UnionMemberAdded,
				},
			},
			{
				name: "Removing union member is a breaking change",
				x:    "union U = A | B type A { a: String } type B { b: Int }",
				y:    "union U = A type A { a: String } type B { b: Int }",
				want: Change{
					Severity: ChangeSeverity{Level: Breaking},
					Type:     UnionMemberRemoved,
				},
			},
		},
	}

	for category, scenarios := range testData {
		for _, s := range scenarios {
			name := fmt.Sprintf("%s/%s", category, s.name)
			t.Run(name, func(t *testing.T) {
				xr, yr := strings.NewReader(s.x), strings.NewReader(s.y)
				res, err := Schema(xr, yr)
				if err != nil {
					t.Fatalf("unable to process schema: %v", err)
				}
				changes := res.Changes()
				if l := len(changes); l == 0 {
					t.Fatal("no changes")
				} else if l > 1 {
					t.Fatal("too many changes")
				}
				have := changes[0]
				if s.want.Type != have.Type {
					t.Errorf("invalid change type: want %q, have %q", s.want.Type, have.Type)
				}
				if s.want.Severity.Level != have.Severity.Level {
					t.Errorf("invalid severity level: want %q, have %q", s.want.Severity.Level, have.Severity.Level)
				}
			})
		}
	}
}
