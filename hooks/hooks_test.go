package hooks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/deis/controller-sdk-go"
	"github.com/deis/controller-sdk-go/api"
)

const keyFixture string = `
{
		"username": "foo",
		"apps": ["test", "testing"]
}`

const configFixture string = `
{
	"owner": "test",
	"app": "example-go",
	"values": {
		"TEST": "testing",
		"FOO": "bar"
	},
	"memory": {
		"web": "1G"
	},
	"cpu": {
		"web": "1000"
	},
	"tags": {
		"test": "tests"
	},
	"registry": {
		"username": "bob"
	},
	"created": "2014-01-01T00:00:00UTC",
	"updated": "2014-01-01T00:00:00UTC",
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const buildFixture = `
{
	"release": {
		"version": 2
	}
}
`

const (
	testingClientFingerprint = `78:b9:21:20:1a:ed:e6:10:05:35:47:da:d4:1f:b6:73`
	configHookExpected       = `{"receive_user":"test","receive_repo":"example-go"}`
	buildHookExpected        = `{"sha":"abc123","receive_user":"test","receive_repo":"example-go","image":"test:abc123","procfile":{"web":"./run"},"dockerfile":"true"}`
	pushHookExpected         = `{"sha":"abc123","receive_user":"test","receive_repo":"example-go","fingerprint":"testing","ssh_connection":"1234","ssh_original_command":"foo"}`
)

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == fmt.Sprintf("/v2/hooks/key/%s", testingClientFingerprint) && req.Method == "GET" {
		res.Write([]byte(keyFixture))
		return
	}

	if req.URL.Path == "/v2/hooks/config/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != configHookExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", configHookExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(configFixture))
		return
	}

	if req.URL.Path == "/v2/hooks/build/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != buildHookExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", buildHookExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(buildFixture))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestUserFromKey(t *testing.T) {
	t.Parallel()

	expected := api.UserApps{
		Username: "foo",
		Apps:     []string{"test", "testing"},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := UserFromKey(deis, testingClientFingerprint)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestConfigHook(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Config{
		Owner: "test",
		App:   "example-go",
		Values: map[string]interface{}{
			"TEST": "testing",
			"FOO":  "bar",
		},
		Memory: map[string]interface{}{
			"web": "1G",
		},
		CPU: map[string]interface{}{
			"web": "1000",
		},
		Tags: map[string]interface{}{
			"test": "tests",
		},
		Registry: map[string]interface{}{
			"username": "bob",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	actual, err := GetAppConfig(deis, "test", "example-go")

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestBuildHook(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := 2

	actual, err := CreateBuild(deis, "test", "example-go", "test:abc123", "abc123", map[string]string{"web": "./run"}, true)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
