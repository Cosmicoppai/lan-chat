package httpErrors

import (
	"net/http"
)

func NotFound(w http.ResponseWriter, error ...string) {
	_error := "Nothing here stranger, Go Back!.."
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusNotFound)
}

func MethodNotAllowed(w http.ResponseWriter, error ...string) {
	_error := "Invalid Method"
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusMethodNotAllowed)
}

func InternalServerError(w http.ResponseWriter, error ...string) {
	_error := "Internal Server Error"
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusInternalServerError)
}

func UnProcessableEntry(w http.ResponseWriter, error ...string) {
	_error := "One or more fields are invalid"
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusUnprocessableEntity)
}

func StatusConflict(w http.ResponseWriter, error ...string) {
	_error := "Resource Already Exist"
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusConflict)
}

func BadRequest(w http.ResponseWriter, error ...string) {
	_error := "Bad Request"
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusBadRequest)
}

func UnAuthorized(w http.ResponseWriter, error ...string) {
	_error := "Provide valid authorization"
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter, error ...string) {
	_error := "insufficient privilege"
	if len(error) > 0 {
		_error = error[0]
	}
	http.Error(w, _error, http.StatusForbidden)
}
