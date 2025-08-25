package responses

import (
	"github.com/goravel/framework/contracts/validation"
)

type ErrorResponse struct {
	Errors interface{} `json:"errors"`
}

func NewErrorResponse(err interface{}) ErrorResponse {
	var errVal interface{}

	switch e := err.(type) {
	case error:
		errVal = e.Error()
	case string:
		errVal = e
	case validation.Errors:
		errs := e.All()

		for _, v := range errs {
			for _, vv := range v {
				errVal = vv
				break
			}
			break
		}

	case *validation.Errors:
		if e != nil {
			errVal = (*e).All()
		} else {
			errVal = "unknown error"
		}
	default:
		errVal = e
	}

	return ErrorResponse{
		Errors: errVal,
	}
}
