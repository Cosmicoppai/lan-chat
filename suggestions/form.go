package suggestions

import (
	"fmt"
	"lan-chat/httpErrors"
	"net/http"
	"os"
)

const suggestionFilePath = "suggestions/suggestion.txt"

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		_ = r.ParseForm() // parse the form
		movieName := r.Form.Get("movie_name")
		date := r.Form.Get("date")
		msg := r.Form.Get("msg")
		if movieName != "" && date != "" {
			message := fmt.Sprintf("%s has been requested on %s. <br> msg:- %s", movieName, date, msg)
			err := writeMessage(message)
			if err != nil {
				httpErrors.InternalServerError(w)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Form Received! Thank you for your suggestion"))
			return
		}
		httpErrors.BadRequest(w)
		return
	}
	httpErrors.MethodNotAllowed(w)
	return
}

func writeMessage(msg string) error {
	f, err := os.OpenFile(suggestionFilePath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	if _, err = f.WriteString(msg); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}
