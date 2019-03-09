package auth

import (
	"strconv"
	"testing"
)

func TestAdmin_IsAdminFalse(t *testing.T) {
	adminer := NewAdmin("admin", "admin")
	err := adminer.LoadSecret()
	if err != nil {
		panic(err)
		return
	}
	isAdmin, err := adminer.IsAdmin()
	println("isAdmin:" + strconv.FormatBool(isAdmin))
	if err != nil {
		panic(err)
		return
	}
}

func TestAdmin_IsAdminTrue(t *testing.T) {
	adminer := NewAdmin("admin", "1f2d1e2e67df")
	err := adminer.LoadSecret()
	if err != nil {
		panic(err)
		return
	}
	isAdmin, err := adminer.IsAdmin()
	println("isAdmin:" + strconv.FormatBool(isAdmin))
	if err != nil {
		panic(err)
		return
	}
}

func TestAdmin_LoadSecret(t *testing.T) {
	adminer := NewAdmin("admin", "admin")
	err := adminer.LoadSecret()
	if err != nil {
		panic(err)
		return
	}
}
