package parser

import (
	"testing"
)

func TestTag_0000_1entry(t *testing.T) {
	tag := newTag("`name:value`")
	v := tag.get("name")
	if v != "value" {
		t.Errorf("value must be value but %s", v)
	}
	v2 := tag.get("name2")
	if v2 != "" {
		t.Errorf("value must be empty string but %s", v2)
	}
}
