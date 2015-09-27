package parser

import (
	"strings"
)

type tag map[string]string

func newTag(s string) tag {
	t := tag{}
	s = s[1 : len(s)-1]
	entries := strings.SplitN(s, ",", -1)
	for _, entry := range entries {
		keyValue := strings.SplitN(entry, ":", -1)
		if len(keyValue) != 2 {
			continue
		}
		t[keyValue[0]] = keyValue[1]
	}
	return t
}

func (t tag) get(key string) string {
	if value, ok := t[key]; ok {
		return value
	} else {
		return ""
	}
}
