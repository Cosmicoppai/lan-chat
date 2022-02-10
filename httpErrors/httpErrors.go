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

func UnProcessableEntry(w http.ResponseWriter) {
	http.Error(w, "One or more fields are invalid", http.StatusUnprocessableEntity)
}

func StatusConflict(w http.ResponseWriter) {
	http.Error(w, "Resource Already Exist", http.StatusConflict)
}
