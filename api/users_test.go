package api

import (
	"sort"
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

func TestUsersSorted(t *testing.T) {
	users := Users{
		{1, "", false, "Zulu", "", "", "", false, false, ""},
		{2, "", false, "Beta", "", "", "", false, false, ""},
		{3, "", false, "Gamma", "", "", "", false, false, ""},
		{4, "", false, "Alpha", "", "", "", false, false, ""},
	}

	sort.Sort(users)
	expectedUsernames := []string{"Alpha", "Beta", "Gamma", "Zulu"}

	for i, user := range users {
		if expectedUsernames[i] != user.Username {
			t.Errorf("Expected users to be sorted %v, Got %v at index %v", expectedUsernames[i], user.Username, i)
		}
	}
}
