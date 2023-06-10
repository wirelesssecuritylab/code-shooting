package spec

import "time"

type Spec interface {
}

type AndSpec struct {
	Specs []Spec
}

func NewAndSpec(specs ...Spec) *AndSpec {
	return &AndSpec{Specs: specs}
}

func (s *AndSpec) Add(specs ...Spec) {
	s.Specs = append(s.Specs, specs...)
}

func (s *AndSpec) AddIf(c bool, spec ...Spec) {
	if c {
		s.Add(spec...)
	}
}

type Field string

const (
	Institute  Field = "Institute"
	Center     Field = "Center"
	Department Field = "Department"
	Team       Field = "Team"
	ImportTime Field = "ImportTime"

	CustomLabel    Field = "CustomLabelInfo"
	ExtendedLabels Field = "ExtendedLabel"
)

type Operator string

const (
	Equal Operator = "="
	Since Operator = ">="
	Until Operator = "<="
	AnyEq Operator = "= ANY"
)

type FieldSpec struct {
	Field    Field
	Operator Operator
	Value    interface{}
}

func (s Field) Equal(value interface{}) Spec {
	return &FieldSpec{Field: s, Operator: Equal, Value: value}
}

func (s Field) Since(value time.Time) Spec {
	return &FieldSpec{Field: s, Operator: Since, Value: value}
}

func (s Field) Until(value time.Time) Spec {
	return &FieldSpec{Field: s, Operator: Until, Value: value}
}

func (s Field) AnyEq(value interface{}) Spec {
	return &FieldSpec{Field: s, Operator: AnyEq, Value: value}
}
