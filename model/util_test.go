package model

import (
	"testing"
)

func TestUtil0000_toSnake(t *testing.T) {
	s := toSnake("Moke")
	if s != "moke" {
		t.Errorf("result must be moke but %s", s)
		return
	}
	s2 := toSnake("MokeTable")
	if s2 != "moke_table" {
		t.Errorf("result must be moke_table but %s", s2)
		return
	}
}
