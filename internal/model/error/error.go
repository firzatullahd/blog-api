package customerror

import (
	"fmt"
	"net/http"
)

var (
	ErrEmailExists      = fmt.Errorf("email already exists")
	ErrValidation       = fmt.Errorf("validation error")
	ErrNotFound         = fmt.Errorf("not found")
	ErrWrongPass        = fmt.Errorf("wrong password")
	ErrUnknown          = fmt.Errorf("error unknown")
	ErrUnauthorized     = fmt.Errorf("unauthorized")
	ErrForbidden        = fmt.Errorf("forbidden")
	ErrSexEdit          = fmt.Errorf("sex is edited when cat is already requested to match")
	ErrDeleteForbidden  = fmt.Errorf("forbidden to delete match")
	ErrAlreadyProcessed = fmt.Errorf("match is already approved / rejected")
	ErrAlreadyMatch     = fmt.Errorf("cat already matched")
	ErrSameOwner        = fmt.Errorf("cat have same owner")
	ErrNotOwned         = fmt.Errorf("cat not owned by user")
)

var mapErrorCode = map[error]int{
	ErrEmailExists:      http.StatusConflict,
	ErrValidation:       http.StatusBadRequest,
	ErrNotFound:         http.StatusNotFound,
	ErrWrongPass:        http.StatusBadRequest,
	ErrUnknown:          http.StatusInternalServerError,
	ErrUnauthorized:     http.StatusUnauthorized,
	ErrSexEdit:          http.StatusBadRequest,
	ErrDeleteForbidden:  http.StatusForbidden,
	ErrAlreadyProcessed: http.StatusBadRequest,
	ErrAlreadyMatch:     http.StatusBadRequest,
	ErrSameOwner:        http.StatusBadRequest,
	ErrNotOwned:         http.StatusNotFound,
	ErrForbidden:        http.StatusForbidden,
}

func ParseError(err error) (code int, errMsg string) {
	if err != nil {
		if errCode, ok := mapErrorCode[err]; ok {
			code = errCode
			return code, err.Error()
		} else {
			return http.StatusInternalServerError, ErrUnknown.Error()
		}
	}
	return http.StatusOK, ""
}
