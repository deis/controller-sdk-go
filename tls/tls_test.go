package tls

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

const (
	tlsDisabledFixture string = `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "foo",
	"owner": "test",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": false
}`
	tlsEnabledFixture string = `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "foo",
	"owner": "test",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": true
}`
	tlsEnableExpected  string = `{"https_enforced":true}`
	tlsDisableExpected string = `{"https_enforced":false}`
)

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "GET" {
		res.Write([]byte(tlsDisabledFixture))
		return
	}

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == tlsEnableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsEnabledFixture))
			return
		} else if string(body) == tlsDisableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsDisableExpected))
			return
		} else {
			fmt.Printf("Expected '%s' or '%s', Got '%s'\n",
				tlsEnableExpected,
				tlsDisableExpected,
				body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

type badJSONFakeHTTPServer struct{}

func (badJSONFakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", deis.APIVersion)

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "GET" {
		res.Write([]byte(tlsDisabledFixture))
		return
	}

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == tlsEnableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsEnabledFixture + "blarg"))
			return
		} else if string(body) == tlsDisableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsDisableExpected + "blarg"))
			return
		} else {
			fmt.Printf("Expected '%s' or '%s', Got '%s'\n",
				tlsEnableExpected,
				tlsDisableExpected,
				body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestTLSInfo(t *testing.T) {
	t.Parallel()

	expected := api.TLS{
		Created:       "2016-08-22T17:40:16Z",
		Updated:       "2016-08-22T17:40:16Z",
		App:           "foo",
		Owner:         "test",
		UUID:          "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
		HTTPSEnforced: new(bool),
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Info(dClient, "foo")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}

	// now test with bad JSON in response, expecting command to return an error
	badHandler := badJSONFakeHTTPServer{}
	badServer := httptest.NewServer(badHandler)
	defer badServer.Close()

	dClient, err = deis.New(false, badServer.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = Info(dClient, "foo"); err != nil {
		t.Errorf("Expected Info() with poorly JSON response to fail")
	}
}

func TestTLSEnable(t *testing.T) {
	t.Parallel()

	b := true
	expected := api.TLS{
		Created:       "2016-08-22T17:40:16Z",
		Updated:       "2016-08-22T17:40:16Z",
		App:           "foo",
		Owner:         "test",
		UUID:          "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
		HTTPSEnforced: &b,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Enable(dClient, "foo")

	if err != nil {
		t.Fatal(err)
	}

	if expected.String() != actual.String() {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}

	// now test with bad JSON in response, expecting command to return an error
	badHandler := badJSONFakeHTTPServer{}
	badServer := httptest.NewServer(badHandler)
	defer badServer.Close()

	dClient, err = deis.New(false, badServer.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = Enable(dClient, "foo"); err != nil {
		t.Errorf("Expected Enable() with poorly JSON response to fail")
	}
}

func TestTLSDisable(t *testing.T) {
	t.Parallel()

	b := false
	expected := api.TLS{
		Created:       "2016-08-22T17:40:16Z",
		Updated:       "2016-08-22T17:40:16Z",
		App:           "foo",
		Owner:         "test",
		UUID:          "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
		HTTPSEnforced: &b,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := deis.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Disable(dClient, "foo")

	if err != nil {
		t.Fatal(err)
	}

	if expected.String() != actual.String() {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}

	// now test with bad JSON in response, expecting command to return an error
	badHandler := badJSONFakeHTTPServer{}
	badServer := httptest.NewServer(badHandler)
	defer badServer.Close()

	dClient, err = deis.New(false, badServer.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = Disable(dClient, "foo"); err != nil {
		t.Errorf("Expected Disable() with poorly JSON response to fail")
	}
}
