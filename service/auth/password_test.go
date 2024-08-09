package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}
	if hash == "" {
		t.Error("expected hash to be not empty")
	}
	if hash == "password" {
		t.Error("expected hash to be different from password")
	}
}

func TestComparePassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}
	if !ComparePasswords(hash, []byte("password")) {
		t.Error("expeceted password to match hash")
	}
	if ComparePasswords(hash, []byte("notpassword")) {
		t.Error("expected password to match hash")
	}
}
