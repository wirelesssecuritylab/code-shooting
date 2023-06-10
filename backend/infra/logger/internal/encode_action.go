package internal

import (
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// Predefined field keys
const (
	_emptyFieldName            = ""
	_timeFieldName             = "T"
	_levelFieldName            = "L"
	_hostnameFieldName         = "H"
	_modNameFieldName          = "N"
	_msgFieldName              = "M"
	_extendsFieldName          = "E"
	_callerFieldName           = "C"
	_callerFileFieldName       = "C_F"
	_callerLineNumberFieldName = "C_L"
	_namespaceFieldName        = "NS"
)

var (
	encoderMutex  sync.Mutex
	encodeActions = map[string]func(item CustomFormatItem) encodeAction{
		_emptyFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.AppendString(item.DefaultValue)
			}

		},
		_timeFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.AppendTime(entry.Time)
			}
		},
		_levelFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.appendLevel(entry.Level)
			}
		},
		_hostnameFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.appendHostName()
			}
		},
		_modNameFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.appendLoggerName(entry.LoggerName)
			}
		},
		_msgFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.AppendString(entry.Message)
			}
		},
		_extendsFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
			}
		},
		_callerFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.appendCaller(entry.Caller)
			}
		},
		_callerFileFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.appendCallerFile(entry.Caller)
			}
		},
		_callerLineNumberFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.appendCallerLine(entry.Caller)
			}
		},
		_namespaceFieldName: func(item CustomFormatItem) encodeAction {
			return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
				encoder.appendNamespace()
			}
		},
	}
)

type encodeAction func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field)

func NewEncoderActions(format CustomFormat) []encodeAction {
	return encoderActionFactorySingleton.CreateEncoderActions(format)
}

func RegisterEncoderActions(encodes map[string]func() string) error {
	eas := make(map[string]func(item CustomFormatItem) encodeAction)
	for k, v := range encodes {
		eas[k] = createEncoderActionFunc(v)
	}

	errTip := ""
	for k, v := range eas {
		if err := checkEncoderAction(k, v); err != nil {
			errTip += err.Error() + "."
		}
	}
	if len(errTip) != 0 {
		return errors.New(errTip)
	}

	encoderMutex.Lock()
	defer encoderMutex.Unlock()
	for k, v := range eas {
		encodeActions[k] = v
	}
	return nil
}

func createEncoderActionFunc(v func() string) func(item CustomFormatItem) encodeAction {
	return func(item CustomFormatItem) encodeAction {
		return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
			value := v()
			encoder.AppendString(value)
		}
	}
}

func checkEncoderAction(key string, value func(item CustomFormatItem) encodeAction) error {
	if len(key) == 0 {
		return errors.New("encodeActions has empty key")
	}
	encoderMutex.Lock()
	defer encoderMutex.Unlock()
	if _, ok := encodeActions[key]; ok {
		return errors.Errorf("key: %s already exist", key)
	}
	return nil
}

type EncoderActionFactory interface {
	CreateEncoderActions(format CustomFormat) []encodeAction
}

var encoderActionFactorySingleton = &encoderActionFactoryImpl{}

type encoderActionFactoryImpl struct {
}

func (s *encoderActionFactoryImpl) CreateEncoderActions(format CustomFormat) []encodeAction {
	specifiedFields := NewStringSet()
	for _, item := range format {
		specifiedFields.Add(item.Field)
	}

	actions := make([]encodeAction, 0, len(format))
	for _, item := range format {
		actions = append(actions, s.buildEncodeAction(item, specifiedFields))
	}
	return actions
}

func (s *encoderActionFactoryImpl) buildEncodeAction(customFormatItem CustomFormatItem, specifiedFields StringSet) encodeAction {
	encoderMutex.Lock()
	defer encoderMutex.Unlock()
	defaultFields := NewStringSet()
	for k := range encodeActions {
		defaultFields.Add(k)
	}
	encodeActions["E"] = func(item CustomFormatItem) encodeAction {
		return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
			userSpecifiedFields := specifiedFields.Difference(defaultFields)
			encoder.addExtendedFields(fields, nil, func(name string) bool {
				return !userSpecifiedFields.Contains(name)
			})
		}
	}

	defaultBuilder := func(item CustomFormatItem) encodeAction {
		return func(encoder *plainEncoder, entry *zapcore.Entry, fields []zapcore.Field) {
			encoder.appendFieldValueOrDefault(item.Field, fields, item.DefaultValue)
		}
	}

	build, ok := encodeActions[customFormatItem.Field]
	if ok {
		return build(customFormatItem)
	}

	return defaultBuilder(customFormatItem)
}
