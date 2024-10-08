package validator

import (
	"bytes"
	"context"
	"github.com/goccy/go-json"
	"github.com/matryer/is"
	"net/http"
	"testing"
)

type requestTest struct {
	Name string
	Age  int
}

func (r *requestTest) Valid(ctx context.Context) (problems map[string]string) {

	problems = map[string]string{}
	if len(r.Name) == 0 {
		problems["required"] = "name is required"
	}

	if r.Age == 0 {
		problems["required"] = "age is required"
	}

	return problems
}

func TestValidator(t *testing.T) {

	_is := is.New(t)

	body := &struct {
		Name string
		Age  int
	}{
		Name: "",
		Age:  0,
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(body)
	_is.NoErr(err) // json.NewEncoder

	req, err := http.NewRequest(http.MethodPost, "/", payloadBuf)
	_is.NoErr(err)

	valid, m, err := DecodeValid[*requestTest](req)
	_is.NoErr(err)

	_is.Equal(valid, body)

	_is.Equal(len(m), 0)

}
