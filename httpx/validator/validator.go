package validator

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
)

type Error struct {
	problems map[string]string
}

func (e *Error) Error() string {

	var message string
	if e.problems != nil {
		for k, v := range e.problems {
			message = message + fmt.Sprintf("\n%s: %s", k, v)
		}
	}

	return message
}

// Validator 验证器接口
type Validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Valid(ctx context.Context) (err error)
}

func DecodeValid[T Validator](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
