package compare

import (
	"fmt"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type schema struct {
	rootTypes  map[ast.Operation]*ast.OperationTypeDefinition
	directives map[string]*ast.DirectiveDefinition
	types      map[string]*ast.Definition
}

func newSchema() *schema {
	return &schema{
		rootTypes:  make(map[ast.Operation]*ast.OperationTypeDefinition),
		directives: make(map[string]*ast.DirectiveDefinition),
		types:      make(map[string]*ast.Definition),
	}
}

func (s *schema) parse(sdl string) error {
	doc, err := parser.ParseSchema(&ast.Source{
		Input: sdl,
	})
	if err != nil {
		return fmt.Errorf("unable to parse schema: %v", err)
	}

	for _, def := range doc.Schema {
		if err := s.processSchemaDefinition(def); err != nil {
			return err
		}
	}

	for _, def := range doc.Directives {
		if err := s.processDirectiveDefinition(def); err != nil {
			return err
		}
	}

	for _, def := range doc.Definitions {
		if err := s.processTypeDefinition(def); err != nil {
			return err
		}
	}

	// TODO Extensions

	return nil
}

func (s *schema) processSchemaDefinition(def *ast.SchemaDefinition) error {
	for _, opDef := range def.OperationTypes {
		if _, ok := s.rootTypes[opDef.Operation]; ok {
			return fmt.Errorf("root type '%s' already exists", opDef.Operation)
		}
		s.rootTypes[opDef.Operation] = opDef
	}

	return nil
}

func (s *schema) processDirectiveDefinition(def *ast.DirectiveDefinition) error {
	if _, ok := s.directives[def.Name]; ok {
		return fmt.Errorf("directive '%s' already exists", def.Name)
	}
	s.directives[def.Name] = def

	return nil
}

func (s *schema) processTypeDefinition(def *ast.Definition) error {
	if _, ok := s.types[def.Name]; ok {
		return fmt.Errorf("%v type '%s' already exists", def.Kind, def.Name)
	}
	s.types[def.Name] = def

	return nil
}
