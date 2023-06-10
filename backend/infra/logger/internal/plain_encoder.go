package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap/buffer"
)

const (
	_hex         = "0123456789abcdef"
	_nullLiteral = "null"
)

var bufferPool = buffer.NewPool()
var namespace string

var _plainEncoderPool = sync.Pool{New: func() interface{} {
	return &plainEncoder{}
}}

func getPlainEncoder() *plainEncoder {
	return _plainEncoderPool.Get().(*plainEncoder)
}

func putPlainEncoder(enc *plainEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	enc.inJsonValueScene = false
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_plainEncoderPool.Put(enc)
}

func SetNamespace(ns string) {
	namespace = ns
}

func getNamespace() string {
	return namespace
}

type plainEncoder struct {
	*zapcore.EncoderConfig
	buf              *buffer.Buffer
	spaced           bool // include spaces after colons and commas
	openNamespaces   int
	inJsonValueScene bool // is json or console(plain) scene, in json value is "value" in other is value
	// don't manual modify this field, but use wrapInJsonValueScene method replace it
	// for encoding generic values by reflection
	ConsoleSeparator           string
	reflectBuf                 *buffer.Buffer
	reflectEnc                 *json.Encoder
	userSpecifiedEncodeActions []encodeAction
}

func NewPlainEncoder(cfg zapcore.EncoderConfig, format CustomFormat) zapcore.Encoder {
	return &plainEncoder{
		EncoderConfig:              &cfg,
		buf:                        bufferPool.Get(),
		spaced:                     true,
		inJsonValueScene:           false,
		ConsoleSeparator:           "\t",
		userSpecifiedEncodeActions: NewEncoderActions(format),
	}
}

func PlainCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}

	enc.AppendString("[")
	index := strings.LastIndexByte(caller.File, '/')
	enc.AppendString(caller.File[index+1:])
	enc.AppendString("]")
	enc.AppendString("[")
	enc.AppendInt(caller.Line)
	enc.AppendString("]")
}

func (enc *plainEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *plainEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *plainEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *plainEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.wrapOfJsonValueScene(func() {
		enc.AppendByteString(val)
	})
}

func (enc *plainEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *plainEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *plainEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.AppendDuration(val)
	enc.buf.AppendByte('"')
}

func (enc *plainEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.wrapOfJsonValueScene(func() {
		enc.AppendFloat64(val)
	})
}

func (enc *plainEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.wrapOfJsonValueScene(func() {
		enc.AppendInt64(val)
	})
}

func (enc *plainEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferPool.Get()
		enc.reflectEnc = json.NewEncoder(enc.reflectBuf)

		// For consistency with our custom JSON encoder.
		enc.reflectEnc.SetEscapeHTML(false)
	} else {
		enc.reflectBuf.Reset()
	}
}

// Only invoke the standard JSON encoder if there is actually something to
// encode; otherwise write JSON null literal directly.
func (enc *plainEncoder) encodeReflected(obj interface{}) ([]byte, error) {
	if obj == nil {
		return []byte(_nullLiteral), nil
	}
	enc.resetReflectBuf()
	if err := enc.reflectEnc.Encode(obj); err != nil {
		return nil, err
	}
	enc.reflectBuf.TrimNewline()
	return enc.reflectBuf.Bytes(), nil
}

func (enc *plainEncoder) AddReflected(key string, obj interface{}) error {
	valueBytes, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *plainEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *plainEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.wrapOfJsonValueScene(func() {
		enc.AppendString(val)
	})
}

func (enc *plainEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.AppendTime(val)
	enc.buf.AppendByte('"')
}

func (enc *plainEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.wrapOfJsonValueScene(func() {
		enc.AppendUint64(val)
	})
}

func (enc *plainEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := enc.wrapWithReturnOfJsonValueScene(func() error {
		return arr.MarshalLogArray(enc)
	})
	enc.buf.AppendByte(']')
	return err
}

func (enc *plainEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *plainEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *plainEncoder) AppendByteString(val []byte) {
	if enc.inJsonValueScene {
		enc.addElementSeparator()
	}
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

func (enc *plainEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *plainEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	if e := enc.EncodeDuration; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		zapcore.StringDurationEncoder(val, enc)
	}
}

func (enc *plainEncoder) AppendInt64(val int64) {
	if enc.inJsonValueScene {
		enc.addElementSeparator()
	}
	enc.buf.AppendInt(val)
}

func (enc *plainEncoder) AppendReflected(val interface{}) error {
	valueBytes, err := enc.encodeReflected(val)
	if err != nil {
		return err
	}
	enc.addElementSeparator()
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *plainEncoder) AppendString(val string) {
	if enc.inJsonValueScene {
		enc.addElementSeparator()
		enc.buf.AppendByte('"')
		enc.safeAddString(val)
		enc.buf.AppendByte('"')
	} else {
		enc.buf.AppendString(val)
	}
}

func (enc *plainEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.buf.AppendTime(time, layout)
}

func (enc *plainEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	if e := enc.EncodeTime; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		zapcore.RFC3339TimeEncoder(val, enc)
	}
}

func (enc *plainEncoder) AppendUint64(val uint64) {
	if enc.inJsonValueScene {
		enc.addElementSeparator()
	}
	enc.buf.AppendUint(val)
}

func (enc *plainEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *plainEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *plainEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *plainEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *plainEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *plainEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

// custom append and add method
func (enc *plainEncoder) appendLevel(level zapcore.Level) {
	encodeLevel := enc.EncodeLevel
	if encodeLevel == nil {
		encodeLevel = zapcore.CapitalLevelEncoder
	}
	encodeLevel(level, enc)
}

func (enc *plainEncoder) appendHostName() {
	hostName, err := os.Hostname()
	if err != nil {
		enc.AppendString("unknown")
	} else {
		enc.AppendString(hostName)
	}
}

func (enc *plainEncoder) appendLoggerName(loggerName string) {
	nameEncoder := enc.EncodeName
	if nameEncoder == nil {
		nameEncoder = zapcore.FullNameEncoder
	}

	cur := enc.buf.Len()
	nameEncoder(loggerName, enc)

	if cur == enc.buf.Len() {
		enc.AppendString(loggerName)
	}
}

func (enc *plainEncoder) appendCaller(caller zapcore.EntryCaller) {
	if caller.Defined {
		encodeCaller := enc.EncodeCaller
		if encodeCaller == nil {
			encodeCaller = PlainCallerEncoder
		}
		encodeCaller(caller, enc)
	}
}

func (enc *plainEncoder) appendCallerFile(caller zapcore.EntryCaller) {
	if caller.Defined {
		index := strings.LastIndexByte(caller.File, '/')
		enc.AppendString(caller.File[index+1:])
	}
}

func (enc *plainEncoder) appendCallerLine(caller zapcore.EntryCaller) {
	if caller.Defined {
		enc.AppendInt(caller.Line)
	}
}

func (enc *plainEncoder) appendNamespace() {
	namespace := getNamespace()
	enc.AppendString(namespace)
}

func (enc *plainEncoder) appendLineEnding() {
	if enc.LineEnding != "" {
		enc.buf.AppendString(enc.LineEnding)
	} else {
		enc.buf.AppendString(zapcore.DefaultLineEnding)
	}
}

func (enc *plainEncoder) wrapAddConsoleItem(addItem func()) {
	cur := enc.buf.Len()
	addItem()
	if cur != enc.buf.Len() {
		enc.addConsoleSeparator()
	}
}

func (enc *plainEncoder) appendFieldValueOrDefault(fieldName string, fields []zapcore.Field, defaultValue string) {
	field := fieldArray(fields).find(fieldName)
	if field != nil {
		enc.appendFieldValue(*field)
	} else {
		enc.AppendString(defaultValue)
	}
}

func (enc *plainEncoder) addExtendedFields(fields []zapcore.Field, parentEncoder *plainEncoder, fieldPredicate func(string) bool) {
	var extendedFields []zapcore.Field

	for _, field := range fields {
		if fieldPredicate(field.Key) {
			extendedFields = append(extendedFields, field)
		}
	}

	if len(extendedFields) == 0 && (parentEncoder == nil || parentEncoder.buf.Len() == 0) {
		return
	}

	enc.buf.AppendByte('{')
	if parentEncoder != nil {
		enc.buf.Write(parentEncoder.buf.Bytes()) // #nosec
	}

	for _, field := range extendedFields {
		field.AddTo(enc)
	}
	enc.closeOpenNamespaces()
	enc.buf.AppendByte('}')
}

func (enc *plainEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes()) // #nosec
	return clone
}

func (enc *plainEncoder) clone() *plainEncoder {
	clone := getPlainEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferPool.Get()
	clone.ConsoleSeparator = enc.ConsoleSeparator
	clone.userSpecifiedEncodeActions = enc.userSpecifiedEncodeActions
	return clone
}

func (enc *plainEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if len(enc.userSpecifiedEncodeActions) > 0 {
		return enc.customFormatEncodeEntry(ent, fields)
	}

	final := enc.clone()

	final.AppendTime(ent.Time)
	final.addConsoleSeparator()

	final.appendLevel(ent.Level)
	final.addConsoleSeparator()

	final.appendHostName()
	final.addConsoleSeparator()

	final.appendLoggerName(ent.LoggerName)
	final.addConsoleSeparator()

	final.AppendString(ent.Message)
	final.addConsoleSeparator()

	final.wrapAddConsoleItem(func() {
		final.addExtendedFields(fields, enc, func(s string) bool {
			return true
		})
	})

	final.appendCaller(ent.Caller)

	final.appendLineEnding()

	ret := final.buf
	putPlainEncoder(final)
	return ret, nil
}

func (enc *plainEncoder) customFormatEncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()

	for _, action := range enc.userSpecifiedEncodeActions {
		action(final, &ent, fields)
	}

	ret := final.buf
	putPlainEncoder(final)
	return ret, nil
}

func (enc *plainEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
}

func (enc *plainEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(key)
	enc.buf.AppendByte('"')
	enc.buf.AppendByte(':')
	if enc.spaced {
		enc.buf.AppendByte(' ')
	}
}

func (enc *plainEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

func (enc *plainEncoder) addConsoleSeparator() {
	if enc.ConsoleSeparator == "" {
		enc.buf.AppendByte('\t')
	} else {
		enc.buf.AppendString(enc.ConsoleSeparator)
	}
}

func (enc *plainEncoder) appendFloat(val float64, bitSize int) {
	if enc.inJsonValueScene {
		enc.addElementSeparator()
	}
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *plainEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *plainEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size]) // #nosec
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *plainEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *plainEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

// only for string scene
func (enc *plainEncoder) appendFieldValue(field zapcore.Field) {
	switch field.Type {
	case zapcore.StringType:
		enc.AppendString(field.String)
	case zapcore.Int64Type:
		enc.AppendInt64(int64(field.Integer))
	case zapcore.DurationType:
		enc.AppendDuration(time.Duration(field.Integer))
	default:
		enc.AppendString(fmt.Sprintf("unsupported field type: %v", field))
	}
}

func (enc *plainEncoder) wrapOfJsonValueScene(valueFunc func()) {
	cur := enc.inJsonValueScene
	defer func() {
		enc.inJsonValueScene = cur
	}()
	enc.inJsonValueScene = true
	valueFunc()
}

func (enc *plainEncoder) wrapWithReturnOfJsonValueScene(valueFunc func() error) error {
	cur := enc.inJsonValueScene
	defer func() {
		enc.inJsonValueScene = cur
	}()
	enc.inJsonValueScene = true
	return valueFunc()
}

// extended filed slice
type fieldArray []zapcore.Field

func (s fieldArray) find(key string) *zapcore.Field {
	for _, field := range s {
		if key == field.Key {
			value := field
			return &value
		}
	}
	return nil
}
