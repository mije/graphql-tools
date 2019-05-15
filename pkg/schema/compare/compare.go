package compare

import (
	"fmt"
	"io"
	"io/ioutil"
)

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

type Result struct {
	changes []Change
}

func (r *Result) reportChange(c Change) {
	r.changes = append(r.changes, c)
}

func (r Result) Changes() []Change {
	return r.changes
}

type Change struct {
	Severity ChangeSeverity
	Type     ChangeType
	Message  string
	Path     string
}

type ChangeSeverity struct {
	Level  ChangeSeverityLevel
	Reason string
}

type ChangeSeverityLevel string

const (
	Breaking    = ChangeSeverityLevel("BREAKING")
	NonBreaking = ChangeSeverityLevel("NON_BREAKING")
	Dangerous   = ChangeSeverityLevel("DANGEROUS")
)

type ChangeType string
