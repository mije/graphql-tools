package compare

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Compare compares two GraphQL schemas and returns a set of detected changes.
// Each change has a severity and reason assigned to be able to further evaluate its impact.
// Schemas must be encoded using SDL, no other form is accepted.
func Schema(x, y io.Reader) (*Result, error) {
	parseSchema := func(name string, r io.Reader) (*schema, error) {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("unable to read schema '%s': %v", name, err)
		}
		s := newSchema()
		if err := s.parse(string(b)); err != nil {
			return nil, fmt.Errorf("unable to parse schema '%s': %v", name, err)
		}
		return s, nil
	}

	sx, err := parseSchema("x", x)
	if err != nil {
		return nil, err
	}
	sy, err := parseSchema("y", y)
	if err != nil {
		return nil, err
	}

	r := new(Result)
	r.compareSchema(sx, sy)
	return r, nil
}

// Result stores the detected changes.
type Result struct {
	breaking    []Change
	dangerous   []Change
	nonBreaking []Change
}

func (r *Result) reportChange(c Change) {
	switch l := c.Severity.Level; l {
	case Breaking:
		r.breaking = append(r.breaking, c)
	case Dangerous:
		r.dangerous = append(r.dangerous, c)
	case NonBreaking:
		r.nonBreaking = append(r.nonBreaking, c)
	default:
		panic(fmt.Errorf("invalid severity level: %s", l))
	}
}

// Breaking returns list of changes which are not backward compatible.
func (r Result) Breaking() []Change {
	return r.breaking
}

// Dangerous returns list of changes which in some cases may be considered breaking.
func (r Result) Dangerous() []Change {
	return r.dangerous
}

// NonBreaking returns list of changes which are backward compatible.
func (r Result) NonBreaking() []Change {
	return r.nonBreaking
}

// Changes returns list of all changes.
func (r Result) Changes() []Change {
	var changes []Change
	changes = append(changes, r.breaking...)
	changes = append(changes, r.dangerous...)
	changes = append(changes, r.nonBreaking...)
	return changes
}

// Change materialize a schema modification.
type Change struct {

	// Severity of the change
	Severity ChangeSeverity

	// Type of the change
	Type ChangeType

	// Message provides human-readable explanation of the change
	Message string

	// Path to the changed item
	Path string
}

// ChangeSeverity defined how serious a change is.
type ChangeSeverity struct {

	// Level indicates backward compatibility
	Level ChangeSeverityLevel

	// Reason provides human-readable explanation of why the severity is chosen
	Reason string
}

// Level indicates backward compatibility.
type ChangeSeverityLevel string

const (
	// Breaking causes existing clients to error
	Breaking = ChangeSeverityLevel("BREAKING")

	// Dangerous may in some cases cause existing clients to error
	Dangerous = ChangeSeverityLevel("DANGEROUS")

	// NonBreaking keeps schema backward compatible
	NonBreaking = ChangeSeverityLevel("NON_BREAKING")
)

// ChangeType allows categorization of changes.
type ChangeType string
