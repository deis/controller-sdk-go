package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

const usersFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "id": 1,
            "last_login": "2014-10-19T22:01:00.601Z",
            "is_superuser": true,
            "username": "test",
            "first_name": "test",
            "last_name": "testerson",
            "email": "test@example.com",
            "is_staff": true,
            "is_active": true,
            "date_joined": "2014-10-19T22:01:00.601Z",
            "groups": [],
            "user_permissions": []
        }
    ]
}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == "/v2/users/" && req.Method == "GET" {
		res.Write([]byte(usersFixture))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestUsersList(t *testing.T) {
	t.Parallel()

	expected := api.Users{
		{
			ID:          1,
			LastLogin:   "2014-10-19T22:01:00.601Z",
			IsSuperuser: true,
			Username:    "test",
			FirstName:   "test",
			LastName:    "testerson",
			Email:       "test@example.com",
			IsStaff:     true,
			IsActive:    true,
			DateJoined:  "2014-10-19T22:01:00.601Z",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(deis, 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
