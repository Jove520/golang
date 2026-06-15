package utils

import "testing"

func TestPasswordHash(t *testing.T) {
	salt, err := GenSalt(12)
	if err != nil {
		t.Fatal(err)
	}

	hashed := GenPassword("123456", salt)

	if !CheckPassword("123456", salt, hashed) {
		t.Fatal("expected password check to pass")
	}

	if CheckPassword("wrong", salt, hashed) {
		t.Fatal("expected wrong password check to fail")
	}
}
