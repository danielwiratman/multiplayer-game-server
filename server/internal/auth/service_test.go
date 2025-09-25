package auth

import "testing"

func TestRegister(t *testing.T) {
	svc := NewService()

	id, err := svc.Register("testuser", "testpass")
	if err != nil {
		t.Fatal(err)
	}
	if id == 0 {
		t.Fatal("id is 0")
	}
}
