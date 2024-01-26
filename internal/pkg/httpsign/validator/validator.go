package validator

import "net/http"

type Validator interface {
	Validator(r *http.Request) error
}
