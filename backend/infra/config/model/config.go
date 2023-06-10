package model

import (
	"strings"
)

type ConfigItem struct {
	Value   interface{}
	Key     string
	Version string
}

type ConfigItemSlice []*ConfigItem

func (s ConfigItemSlice) Len() int {
	return len(s)
}

func (s ConfigItemSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ConfigItemSlice) Less(i, j int) bool {
	return strings.Compare(s[i].Key, s[j].Key) < 0
}
