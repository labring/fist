package terminal

import (
	"encoding/json"
	"testing"
)

func TestQuery_Query(t *testing.T) {
	query := &ListQuery{
		CookieUserName: "admin",
		TerminalID:     "",
		Namespace:      "",
	}
	got, _ := query.Query()
	ss, _ := json.Marshal(got)
	t.Log(string(ss))
}
