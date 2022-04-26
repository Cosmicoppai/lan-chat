package suggestions

import (
	"fmt"
	"lan-chat/audio"
	"lan-chat/httpErrors"
	"lan-chat/logger"
	"net/http"
	"os"
)

const suggestionFilePath = "suggestions/suggestion.txt"

// FormHandler to accept form-data
func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		_ = r.ParseForm() // parse the form
		movieName := r.Form.Get("movie_name")
		date := r.Form.Get("date")
		msg := r.Form.Get("msg")
		done := make(chan bool) // channel to receive completion of audio.Notify() function
		if movieName != "" && date != "" {
			message := fmt.Sprintf("%s has been requested on %s. [ msg:- %s]", movieName, date, msg)
			go audio.Notify(done)
			err := writeMessage(message)
			if err != nil {
				logger.ErrorLog.Println(err)
				httpErrors.InternalServerError(w)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Form Received! Thank you for your suggestion"))
			<-done // when we receive something return
			return
		}
		httpErrors.BadRequest(w)
		return
	}
	httpErrors.MethodNotAllowed(w)
	return
}

func writeMessage(msg string) error {
	f, err := os.OpenFile(suggestionFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if _, err = f.WriteString(msg + "\n"); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	return nil
}
