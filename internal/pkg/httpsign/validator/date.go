package validator

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const maxTimeGap = 30 * time.Second

func newPublicError(msg string) *gin.Error {
	return &gin.Error{
		Err:  errors.New(msg),
		Type: gin.ErrorTypePublic,
	}
}

var ErrDateNotInRange = newPublicError("date submit is not in acceptable range")

type DateValidator struct {
	TimeGap time.Duration
}

func NewDateValidator() *DateValidator {
	return &DateValidator{
		TimeGap: maxTimeGap,
	}
}

func (v *DateValidator) Validator(r *http.Request) error {
	t, err := http.ParseTime(r.Header.Get("date"))
	if err != nil {
		return newPublicError(fmt.Sprintf("Could not parse date header. Error: %s", err.Error()))
	}

	serverTime := time.Now()
	start := serverTime.Add(-v.TimeGap)
	stop := serverTime.Add(v.TimeGap)
	if t.Before(start) || t.After(stop) {
		return ErrDateNotInRange
	}

	return nil
}
