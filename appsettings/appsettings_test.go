package appsettings

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

const appSettingsFixture string = `
{
    "owner": "test",
    "app": "example-go",
    "maintenance": true,
    "routable": true,
    "whitelist": ["1.2.3.4", "0.0.0.0/0"],
    "autoscale": {"cmd": {"min": 3, "max": 8, "cpu_percent": 40}},
    "label": {"git_repo": "https://github.com/deis/controller-sdk-go", "team" : "deis"},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const appSettingsUnsetFixture string = `
{
    "owner": "test",
    "app": "unset-test",
    "maintenance": true,
    "routable": true,
    "whitelist": ["1.2.3.4", "0.0.0.0/0"],
    "autoscale": {"cmd": {"min": 3, "max": 8, "cpu_percent": 40}},
    "label": {"git_repo": "https://github.com/deis/controller-sdk-go", "team" : "deis"},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const appSettingsSetExpected string = `{"maintenance":true,"routable":true,"whitelist":["1.2.3.4","0.0.0.0/0"],"autoscale":{"cmd":{"min":3,"max":8,"cpu_percent":40}},"label":{"git_repo":"https://github.com/deis/controller-sdk-go","team":"deis"}}`
const appSettingsUnsetExpected string = `{"maintenance":true,"routable":true,"whitelist":["1.2.3.4","0.0.0.0/0"],"autoscale":{"cmd":{"min":3,"max":8,"cpu_percent":40}},"label":{"git_repo":"https://github.com/deis/controller-sdk-go","team":"deis"}}`

var trueVar = true

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/settings/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appSettingsSetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appSettingsSetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(appSettingsFixture))
		return
	}

	if req.URL.Path == "/v2/apps/unset-test/settings/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appSettingsUnsetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appSettingsUnsetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(appSettingsUnsetFixture))
		return
	}

	if req.URL.Path == "/v2/apps/invalidjson-test/settings/" && req.Method == "POST" {
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`"maintenance": "test"`))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/settings/" && req.Method == "GET" {
		res.Write([]byte(appSettingsFixture))
		return
	}

	if req.URL.Path == "/v2/apps/invalidjson-test/settings/" && req.Method == "GET" {
		res.Write([]byte(`"maintenance": "test"`))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestAppSettingsSet(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.AppSettings{
		Owner:       "test",
		App:         "example-go",
		Routable:    api.NewRoutable(),
		Maintenance: &trueVar,
		Whitelist:   []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/deis/controller-sdk-go",
			"team":     "deis",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	appSettingsVars := api.AppSettings{
		Maintenance: &trueVar,
		Routable:    api.NewRoutable(),
		Whitelist:   []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/deis/controller-sdk-go",
			"team":     "deis",
		},
	}

	actual, err := Set(deis, "example-go", appSettingsVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppSettingsUnset(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.AppSettings{
		Owner:       "test",
		App:         "unset-test",
		Maintenance: &trueVar,
		Routable:    api.NewRoutable(),
		Whitelist:   []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/deis/controller-sdk-go",
			"team":     "deis",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	appSettingsVars := api.AppSettings{
		Maintenance: &trueVar,
		Routable:    api.NewRoutable(),
		Whitelist:   []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/deis/controller-sdk-go",
			"team":     "deis",
		},
	}

	actual, err := Set(deis, "unset-test", appSettingsVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppSettingsList(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.AppSettings{
		Owner:       "test",
		App:         "example-go",
		Maintenance: &trueVar,
		Routable:    api.NewRoutable(),
		Whitelist:   []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/deis/controller-sdk-go",
			"team":     "deis",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	actual, err := List(deis, "example-go")

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppSettingsInvalidJson(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	deis, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	_, err = List(deis, "invalidjson-test")
	expected := "json: cannot unmarshal string into Go value of type api.AppSettings"
	if err == nil || !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("Expected %v, Got %v", expected, err)
	}

	appSettingsVars := api.AppSettings{
		Maintenance: &trueVar,
	}
	_, err = Set(deis, "invalidjson-test", appSettingsVars)
	if err == nil || !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("Expected %v, Got %v", expected, err)
	}
}
