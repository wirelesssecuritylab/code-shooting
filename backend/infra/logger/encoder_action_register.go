package logger

import (
	"code-shooting/infra/logger/internal"

	"github.com/pkg/errors"
)

func RegisterEncoderActions(EncoderActions map[string]func() string) error {
	if len(EncoderActions) == 0 {
		return errors.New("EncoderActions is empty")
	}
	return internal.RegisterEncoderActions(EncoderActions)
}
