package api

import (
	"testing"
)

func TestUserString(t *testing.T) {
	user := User{
		ID:          1,
		Username:    "bacongobbler",
		Email:       "matthewf@deis.com",
		FirstName:   "Matthew",
		LastName:    "Fisher",
		LastLogin:   "Yesterday",
		IsSuperuser: true,
		IsStaff:     true,
		IsActive:    true,
		DateJoined:  "Yesterday",
	}

	expected := `ID: 1
Username: bacongobbler
Email: matthewf@deis.com
First Name: Matthew
Last Name: Fisher
Last Login: Yesterday
Is Superuser: true
Is Staff: true
Is Active: true
Date Joined: Yesterday`

	if user.String() != expected {
		t.Errorf("Got:\n\n%s\n\nExpected:\n\n%s", user.String(), expected)
	}
}
