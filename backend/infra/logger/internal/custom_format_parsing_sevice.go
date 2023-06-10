package internal

import "github.com/pkg/errors"

func _getFieldBeginToken() []byte {
	return []byte{'$', '{'}
}

func _getFieldBeginTokenLen() int {
	return len(_getFieldBeginToken())
}

func _getFieldEndToken() []byte {
	return []byte{'}'}
}

func _getFieldEndTokenLen() int {
	return len(_getFieldEndToken())
}

func _getDefaultToken() []byte {
	return []byte{':'}
}

func _getDefaultTokenLen() int {
	return len(_getDefaultToken())
}

type CustomFormatParsingService interface {
	Parse(format string) (CustomFormat, error)
}

var customFormatParsingServiceSingleton CustomFormatParsingService = &customFormatParsingServiceImpl{}

func GetCustomFormatParsingService() CustomFormatParsingService {
	return customFormatParsingServiceSingleton
}

type customFormatParsingServiceImpl struct {
}

func (s *customFormatParsingServiceImpl) Parse(format string) (CustomFormat, error) {
	f := bytes([]byte(format))
	var items []CustomFormatItem

	startPos := 0
	for startPos < len(f) {
		parsedItems, nextPos, err := s.parse(f, startPos)
		if err != nil {
			return nil, errors.Wrap(err, "parse format")
		}
		items = append(items, parsedItems...)
		startPos = nextPos
	}

	return CustomFormat(items), nil
}

func (s *customFormatParsingServiceImpl) parse(f bytes, startPos int) ([]CustomFormatItem, int, error) {
	fieldPos := f.findFromPos(_getFieldBeginToken(), startPos)
	if fieldPos == -1 {
		return []CustomFormatItem{{Field: "", DefaultValue: string(f[startPos:])}}, len(f), nil
	}

	fieldEndPos := f.findFromPos(_getFieldEndToken(), fieldPos+_getFieldBeginTokenLen())
	if fieldEndPos == -1 {
		return nil, 0, errors.Errorf("without field ending token } in %s", string(f[fieldPos:]))
	}

	var items []CustomFormatItem

	if fieldPos != startPos {
		items = append(items, CustomFormatItem{Field: "", DefaultValue: string(f[startPos:fieldPos])})
	}

	item, err := s.buildFieldItem(f, fieldPos, fieldEndPos)
	if err != nil {
		return nil, 0, errors.Wrap(err, "build custom field item")
	}

	items = append(items, item)

	return items, fieldEndPos + _getFieldEndTokenLen(), nil
}

func (s *customFormatParsingServiceImpl) buildFieldItem(f bytes, fieldPos, fieldEndPos int) (CustomFormatItem, error) {
	nextFieldPos := f.findFromPos(_getFieldBeginToken(), fieldPos+_getFieldBeginTokenLen())
	if nextFieldPos != -1 && nextFieldPos+_getFieldBeginTokenLen() <= fieldEndPos {
		return CustomFormatItem{}, errors.Errorf("without field ending token in %s", string(f[fieldPos:]))
	}

	defaultValuePos := f.findInRange(_getDefaultToken(), fieldPos+_getFieldBeginTokenLen(), fieldEndPos)

	fieldName := ""
	defaultValue := ""
	if defaultValuePos == -1 {
		fieldName = string(f[fieldPos+_getFieldBeginTokenLen() : fieldEndPos])
	} else {
		fieldName = string(f[fieldPos+_getFieldBeginTokenLen() : defaultValuePos])
		defaultValue = string(f[defaultValuePos+_getDefaultTokenLen() : fieldEndPos])
	}

	if len(fieldName) == 0 {
		return CustomFormatItem{}, errors.Errorf("without field name in %s", string(f[fieldPos:]))
	}

	return CustomFormatItem{Field: fieldName, DefaultValue: defaultValue}, nil
}

type bytes []byte

func (s *bytes) findFromPos(bs []byte, pos int) int {
	return s.findInRange(bs, pos, len(*s))
}

func (s *bytes) findInRange(bs []byte, startPos, endPos int) int {
	if startPos+len(bs) > len(*s) {
		return -1
	}

	if len(bs) == 0 {
		return startPos
	}

	for i := startPos; i <= len(*s)-len(bs) && i < endPos; i++ {
		matched := true
		for j := 0; j < len(bs); j++ {
			if bs[j] != (*s)[i+j] {
				matched = false
				break
			}
		}
		if matched {
			return i
		}
	}
	return -1
}
