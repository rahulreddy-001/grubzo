package dto

import (
	"encoding/json"
	"testing"
)

func TestCreateUserAcceptsLowercaseJSONKeys(t *testing.T) {
	var got CreateUser
	if err := json.Unmarshal([]byte(`{"email":"rohan@gmail.com","password":"Rohan@001","name":"Rohan"}`), &got); err != nil {
		t.Fatal(err)
	}

	if got.Email != "rohan@gmail.com" {
		t.Fatalf("email = %q", got.Email)
	}
	if got.Password != "Rohan@001" {
		t.Fatalf("password = %q", got.Password)
	}
	if got.Name != "Rohan" {
		t.Fatalf("name = %q", got.Name)
	}
}
