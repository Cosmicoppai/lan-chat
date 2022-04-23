package httpErrors

import (
	"net/http"
)

func NotFound(w http.ResponseWriter) {
	http.Error(w, "Nothing here stranger, Go Back!..", http.StatusNotFound)
}

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
