package utils

import (
	"strings"

	"github.com/labstack/gommon/log"
)

func ConvertOutKeyToInner(out string) (inner string) {
	key := strings.ReplaceAll(out, ".", "#")
	inner = strings.ReplaceAll(key, "##", ".")
	log.Debug(" ConvertOutKeyToInner ", inner)
	return
}

func ConvertInnerKeyToOut(inner string) (out string) {
	out = strings.ReplaceAll(inner, "#", ".")
	log.Debug(" ConvertInnerKeyToOut ", out)
	return
}
