package utils

import (
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	secret := "test-secret"

	token, err := GenToken(1001, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatal(err)
	}

	if claims.UserId != 1001 {
		t.Fatalf("expected user id 1001, got %d", claims.UserId)
	}
}

func TestJWTExpired(t *testing.T) {
	secret := "test-secret"

	token, err := GenToken(1001, secret, -time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ParseToken(token, secret); err == nil {
		t.Fatal("expected expired token error")
	}
}
