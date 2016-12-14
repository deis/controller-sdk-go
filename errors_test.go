package deis

import (
	"errors"
	"fmt"
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
	if m.closed {
		return 0, errors.New("You can't read on a closed ReadCloser")
	}

	if m.msg == "" {
		return 0, io.EOF
	}

	copy(msg, []byte(m.msg))

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
		{
			res: &http.Response{
				StatusCode: 200,
				Body:       readCloser(""),
			},
			expected: nil,
		},
		{
			res: &http.Response{
				StatusCode: 201,
				Body:       readCloser(""),
			},
			expected: nil,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"username":["This field may not be blank."]}`),
			},
			expected: ErrInvalidUsername,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"username":["Enter a valid username. This value may contain only letters, numbers and @/./+/-/_ characters."]}`),
			},
			expected: ErrInvalidUsername,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"username":["A user with that username already exists."]}`),
			},
			expected: ErrDuplicateUsername,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"password":["This field may not be blank."]}`),
			},
			expected: ErrMissingPassword,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"non_field_errors":["Unable to log in with provided credentials."]}`),
			},
			expected: ErrLogin,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"id":["App name can only contain a-z (lowercase), 0-9 and hyphens","Enter a valid \"slug\" consisting of letters, numbers, underscores or hyphens."]}`),
			},
			expected: ErrInvalidAppName,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"id":["Application with this id already exists."]}`),
			},
			expected: ErrDuplicateApp,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"certificate": ["This field may not be blank."]}`),
			},
			expected: ErrInvalidCertificate,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"certificate":["Could not load certificate: [('PEM routines', 'PEM_read_bio', 'no start line')]"]}`),
			},
			expected: ErrInvalidCertificate,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"key": ["This field may not be blank."]}`),
			},
			expected: ErrMissingKey,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"key": ["Public Key is already in use"]}`),
			},
			expected: ErrDuplicateKey,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"public":["Key contains invalid base64 chars"]}`),
			},
			expected: ErrMissingKey,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"public":["This field may not be blank."]}`),
			},
			expected: ErrMissingKey,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"name": ["This field may not be blank."]}`),
			},
			expected: ErrInvalidName,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"name":["Can only contain a-z (lowercase), 0-9 and hyphens"]}`),
			},
			expected: ErrInvalidName,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"detail":"Container type foo does not exist in application"}`),
			},
			expected: ErrPodNotFound,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"domain":["Hostname does not look valid."]}`),
			},
			expected: ErrInvalidDomain,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"domain":["Domain is already in use by another application"]}`),
			},
			expected: ErrDuplicateDomain,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"image":["This field may not be blank."]}`),
			},
			expected: ErrInvalidImage,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"detail":"version cannot be below 0"}`),
			},
			expected: ErrInvalidVersion,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"id":["This field may not be blank."]}`),
			},
			expected: ErrMissingID,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"email":["Enter a valid email address."]}`),
			},
			expected: ErrInvalidEmail,
		},
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"detail":"No nodes matched the provided labels: foo=bar"}`),
			},
			expected: ErrTagNotFound,
		},
		{
			res: &http.Response{
				StatusCode: 401,
				Body:       readCloser(""),
			},
			expected: ErrUnauthorized,
		},
		{
			res: &http.Response{
				StatusCode: 403,
				Body:       readCloser(""),
			},
			expected: ErrForbidden,
		},
		{
			res: &http.Response{
				StatusCode: 404,
				Body:       readCloser(""),
			},
			expected: ErrNotFound{"Not Found"},
		},
		{
			res: &http.Response{
				StatusCode: 404,
				Body:       readCloser("App not found"),
			},
			expected: ErrNotFound{"App not found"},
		},
		{
			res: &http.Response{
				StatusCode: 405,
				Body:       readCloser(""),
			},
			expected: ErrMethodNotAllowed,
		},
		{
			res: &http.Response{
				StatusCode: 409,
				Body:       readCloser(`{"detail":"foo still has applications assigned. Delete or transfer ownership"}`),
			},
			expected: ErrCancellationFailed,
		},
		{
			res: &http.Response{
				StatusCode: 422,
				Body:       readCloser(`{"detail":"test does not exist under values"}`),
			},
			expected: ErrUnprocessable{"test does not exist under values"},
		},
		{
			res: &http.Response{
				StatusCode: 500,
				Body:       readCloser(""),
			},
			expected: ErrServerError,
		},
		// ensure unknown errors at least look pretty
		{
			res: &http.Response{
				StatusCode: 400,
				Body:       readCloser(`{"detail":"unknown error\nnewline"}`),
			},
			expected: errors.New(`Unknown Error (400): {"detail":"unknown error
newline"}`),
		},
	}

	for _, check := range tests {
		actual := checkForErrors(check.res)

		// specifically check error output rather than value comparison
		if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", check.expected) {
			t.Errorf(failureMessage, check.expected, actual)
		}
	}
}
