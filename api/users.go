package api

import (
	"fmt"
)

// User is the definition of the user object.
type User struct {
	ID          int    `json:"id"`
	LastLogin   string `json:"last_login"`
	IsSuperuser bool   `json:"is_superuser"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	IsStaff     bool   `json:"is_staff"`
	IsActive    bool   `json:"is_active"`
	DateJoined  string `json:"date_joined"`
}

func (u User) String() string {
	tpl := `ID: %d
Username: %s
Email: %s
First Name: %s
Last Name: %s
Last Login: %s
Is Superuser: %t
Is Staff: %t
Is Active: %t
Date Joined: %s`

	return fmt.Sprintf(tpl, u.ID, u.Username, u.Email, u.FirstName, u.LastName, u.LastLogin,
		u.IsSuperuser, u.IsStaff, u.IsActive, u.DateJoined)
}

type Users []User

func (u Users) Len() int           { return len(u) }
func (u Users) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u Users) Less(i, j int) bool { return u[i].Username < u[j].Username }
