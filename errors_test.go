package deis

import (
	"errors"
	"io"
	"net/http"
	"testing"
)

const failureMessage = "Expected error '%v', Got '%v'"

type mockReadCloser struct {
	msg    string
	closed bool
}

func (m *mockReadCloser) Close() error {
	m.closed = true
	return nil
}

func (m *mockReadCloser) Read(msg []byte) (int, error) {
	if m.closed == true {
		return 0, errors.New("You can't read on a closed ReadCloser")
	}

	if m.msg == "" {
		return 0, io.EOF
	}

	for i, b := range []byte(m.msg) {
		msg[i] = b
	}
	len := len(m.msg)
	m.msg = ""
	return len, nil
}

func readCloser(msg string) *mockReadCloser {
	return &mockReadCloser{msg: msg}
}

type errorTest struct {
	expected error
	res      *http.Response
}

func TestErrors(t *testing.T) {
	tests := []errorTest{
		errorTest{
			res: &http.Response{
				StatusCode: 200,
				Body:       readCloser(""),
			},
			expected: nil,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 201,
				Body:       readCloser(""),
			},
			expected: nil,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"username":["This field is required."]}`),
			},
			expected: ErrInvalidUsername,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"username":["Enter a valid username. This value may contain only letters, numbers and @/./+/-/_ characters."]}`),
			},
			expected: ErrInvalidUsername,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"username":["A user with that username already exists."]}`),
			},
			expected: ErrDuplicateUsername,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"password":["This field is required."]}`),
			},
			expected: ErrMissingPassword,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"non_field_errors":["Unable to log in with provided credentials."]}`),
			},
			expected: ErrLogin,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"id":["App name can only contain a-z (lowercase), 0-9 and hypens","Enter a valid \"slug\" consisting of letters, numbers, underscores or hyphens."]}`),
			},
			expected: ErrInvalidAppName,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"id":["App with this id already exists."]}`),
			},
			expected: ErrDuplicateApp,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"certificate": ["This field is required."]}`),
			},
			expected: ErrInvalidCertificate,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"certificate":["Could not load certificate: [('PEM routines', 'PEM_read_bio', 'no start line')]"]}`),
			},
			expected: ErrInvalidCertificate,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"key": ["This field is required."]}`),
			},
			expected: ErrMissingKey,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"public":["Key contains invalid base64 chars"]}`),
			},
			expected: ErrMissingKey,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"public":["This field is required."]}`),
			},
			expected: ErrMissingKey,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"name": ["This field is required."]}`),
			},
			expected: ErrInvalidName,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"name":["Can only contain a-z (lowercase), 0-9 and hypens"]}`),
			},
			expected: ErrInvalidName,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"detail":"Container type foo does not exist in application"}`),
			},
			expected: ErrPodNotFound,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"domain":["Hostname does not look valid."]}`),
			},
			expected: ErrInvalidDomain,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"image":["This field is required."]}`),
			},
			expected: ErrInvalidImage,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"detail":"version cannot be below 0"}`),
			},
			expected: ErrInvalidVersion,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"id":["This field is required."]}`),
			},
			expected: ErrMissingID,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"email":["Enter a valid email address."]}`),
			},
			expected: ErrInvalidEmail,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"detail":"No nodes matched the provided labels: foo=bar"}`),
			},
			expected: ErrTagNotFound,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 401,
				Body:       readCloser(""),
			},
			expected: ErrUnauthorized,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 403,
				Body:       readCloser(""),
			},
			expected: ErrForbidden,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 404,
				Body:       readCloser(""),
			},
			expected: ErrNotFound,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 405,
				Body:       readCloser(""),
			},
			expected: ErrMethodNotAllowed,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 409,
				Body:       readCloser(""),
			},
			expected: ErrConflict,
		},
		errorTest{
			res: &http.Response{
				StatusCode: 500,
				Body:       readCloser(""),
			},
			expected: ErrServerError,
		},
	}

	for _, check := range tests {
		actual := checkForErrors(check.res)

		if actual != check.expected {
			t.Errorf(failureMessage, check.expected, actual)
		}
	}
}
