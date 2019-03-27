package tools

import "testing"

func TestMapToString(t *testing.T) {
	s := make(map[string]string)
	s["ddd"] = "ddd"
	s["eee"] = "eee"
	s["fff"] = "fff"
	t.Log(MapToString(s))
}
