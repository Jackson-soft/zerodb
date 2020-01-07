package server

import (
	"fmt"
	"testing"
)

func TestCreatToken(t *testing.T) {
	j := NewJWT()
	token, err := j.GenerateToken("admin", "admin")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(token)
}
