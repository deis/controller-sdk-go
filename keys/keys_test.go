package keys

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	deis "github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

const keysFixture string = `
{
		"created": "2014-01-01T00:00:00UTC",
		"id": "test@example.com",
		"owner": "test",
		"public": "ssh-rsa abc test@example.com",
		"updated": "2014-01-01T00:00:00UTC",
		"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`

const keysListFixture string = `
{
		"count": 1,
		"next": null,
		"previous": null,
		"results": [
				{
						"created": "2014-01-01T00:00:00UTC",
						"id": "test@example.com",
						"owner": "test",
						"public": "ssh-rsa abc test@example.com",
						"updated": "2014-01-01T00:00:00UTC",
						"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
				}
		]
}`

const keyCreateExpected string = `{"id":"test@example.com","public":"ssh-rsa abc test@example.com"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == "/v2/keys/" && req.Method == "GET" {
		res.Write([]byte(keysListFixture))
		return
	}

	if req.URL.Path == "/v2/keys/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != keyCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", keyCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(keysFixture))
		return
	}

	if req.URL.Path == "/v2/keys/test@example.com" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestKeysList(t *testing.T) {
	t.Parallel()

	expected := api.Keys{
		{
			Created: "2014-01-01T00:00:00UTC",
			ID:      "test@example.com",
			Owner:   "test",
			Public:  "ssh-rsa abc test@example.com",
			Updated: "2014-01-01T00:00:00UTC",
			UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
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

func TestKeyCreate(t *testing.T) {
	t.Parallel()

	expected := api.Key{
		Created: "2014-01-01T00:00:00UTC",
		ID:      "test@example.com",
		Owner:   "test",
		Public:  "ssh-rsa abc test@example.com",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := New(deis, "test@example.com", "ssh-rsa abc test@example.com")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestKeysDestroy(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(deis, "test@example.com"); err != nil {
		t.Fatal(err)
	}
}
