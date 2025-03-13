package server

import (
	"errors"
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/commons"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func errAbortWithValidatorPayload(c *gin.Context, err error) {
	var errs validator.ValidationErrors
	errors.As(err, &errs)

	errorMap := make(map[string]string)
	for _, v := range errs {
		key, txt := validatorErrorToText(&v)
		errorMap[key] = txt
	}

	resErr := newServiceError(illegalArgument, fmt.Errorf("validation error: %v (%w)", commons.ValuesString(errorMap), err))
	c.Header(api.HeaderContentType, api.HeaderContentTypeApplicationJson)
	_ = c.AbortWithError(http.StatusBadRequest, resErr)
	return
}

func errToHttpStatus(err error) int {
	var e *serviceError
	switch {
	case errors.As(err, &e):
		if e.Status == illegalArgument {
			return http.StatusBadRequest
		} else if e.Status == unauthorized {
			return http.StatusUnauthorized
		} else if e.Status == forbidden {
			return http.StatusForbidden
		} else if e.Status == notFound {
			return http.StatusNotFound
		} else if e.Status == methodNotAllowed {
			return http.StatusMethodNotAllowed
		} else if e.Status == conflict {
			return http.StatusConflict
		} else if e.Status == general {
			return http.StatusInternalServerError
		}
	default:
		return http.StatusInternalServerError
	}

	return -1
}

func errCodeToStr(err error) string {
	var e *serviceError
	ok := errors.As(err, &e)

	if ok {
		return string(e.Status)
	}

	return string(general)
}

func validatorErrorToText(e *validator.FieldError) (string, string) {
	x := *e

	switch x.Tag() {
	case "required":
		return x.Field(), fmt.Sprintf("%s is required", x.Field())
	case "max":
		return x.Field(), fmt.Sprintf("%s cannot be longer than %s", x.Field(), x.Param())
	case "min":
		return x.Field(), fmt.Sprintf("%s must be longer than %s", x.Field(), x.Param())
	case "len":
		return x.Field(), fmt.Sprintf("%s must be %s characters long", x.Field(), x.Param())
	}
	return x.Field(), fmt.Sprintf("%s is not valid", x.Field())
}
