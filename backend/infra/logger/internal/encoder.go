package internal

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Encoder string

const (
	JsonEncoder  Encoder = "json"
	PlainEncoder Encoder = "plain"
)

func (s *Encoder) Check() error {
	validEncoders := []Encoder{JsonEncoder, PlainEncoder}
	for _, encoder := range validEncoders {
		if strings.ToLower(string(*s)) == strings.ToLower(string(encoder)) {
			return nil
		}
	}
	return errors.Errorf("unsupported encoder: %s", string(*s))
}

func (s *Encoder) Equal(e Encoder) bool {
	if err := s.Check(); err != nil {
		return false
	}
	return strings.ToLower(string(*s)) == strings.ToLower(string(e))
}

func (s *Encoder) String() string {
	if err := s.Check(); err != nil {
		return fmt.Sprintf("Encoder(%s)", *s)
	}
	return strings.ToLower(string(*s))
}
