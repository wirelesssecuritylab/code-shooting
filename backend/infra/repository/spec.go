package repository

import (
	"fmt"
	"strings"

	"code-shooting/infra/logger"

	"gorm.io/gorm/schema"

	"code-shooting/domain/entity/spec"
)

func andSpecToQuery(s *spec.AndSpec) (string, []interface{}) {
	var qrys []string
	var args []interface{}
	for _, v := range s.Specs {
		q, a := specToQuery(v)
		if q != "" {
			qrys = append(qrys, "("+q+")")
			args = append(args, a...)
		}
	}
	return strings.Join(qrys, " AND "), args
}

func fieldSpecToQuery(s *spec.FieldSpec) (string, []interface{}) {
	if s.Operator == spec.AnyEq {
		return fmt.Sprintf("? %v (%v)", s.Operator, columnName(s.Field)),
			[]interface{}{s.Value}
	}
	return fmt.Sprintf("%v %v ?", columnName(s.Field), s.Operator),
		[]interface{}{s.Value}
}

func specToQuery(sp spec.Spec) (string, []interface{}) {
	switch s := sp.(type) {
	case *spec.AndSpec:
		return andSpecToQuery(s)
	case *spec.FieldSpec:
		return fieldSpecToQuery(s)
	}
	logger.Warnf("unknown spec type: %T", sp)
	return "", nil
}

func columnName(f spec.Field) string {
	return schema.NamingStrategy{}.ColumnName("", string(f))
}
