package movieHandler

import (
	"lan-chat/admin"
	"lan-chat/httpErrors"
	"net/http"
	"os"
	"strings"
)

func StaticPages(w http.ResponseWriter, r *http.Request) { // To serve static pages
	fileLocation := strings.TrimPrefix(r.URL.Path, "/static/")
	http.ServeFile(w, r, "./templates/"+fileLocation)
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpErrors.MethodNotAllowed(w)
		return
	}
	w.Header().Set("content-type", "text/html")
	fn := r.URL.Path

	switch fn {
	case "/":
		http.ServeFile(w, r, "./templates/index.html")
	case "/favicon.ico":
		http.Redirect(w, r, "/static"+fn, 301)
	default:
		file := fn + ".html"
		http.ServeFile(w, r, "./templates/"+file)
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpErrors.MethodNotAllowed(w)
		return
	}
	moviePath := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/file/"), "/") // no need to clean url, http package cleans it by default
	params := strings.Split(moviePath, "/")
	if len(params) < 3 {
		httpErrors.NotFound(w)
		return
	}
	switch params[0] {
	case "video":
		{
			moviePath = "videos" + strings.TrimPrefix(moviePath, "video")
			if isExist(moviePath) {
				w.Header().Set("content-type", "video/mp4")
				w.Header().Set("accept-ranges", "bytes")
				http.ServeFile(w, r, moviePath)
			} else {
				httpErrors.NotFound(w)
			}
		}
	case "poster":
		{
			var posterPath = "videos" + strings.TrimPrefix(moviePath, "poster")
			if isExist(posterPath) {
				w.Header().Set("content-type", "images/png")
				http.ServeFile(w, r, posterPath)
			} else {
				httpErrors.NotFound(w)
			}

		}
	case "sub":
		{
			var subPath = "videos" + strings.TrimPrefix(moviePath, "sub")
			if isExist(subPath) {
				w.Header().Set("accept-ranges", "bytes")
				http.ServeFile(w, r, subPath)
			} else {
				httpErrors.NotFound(w)
			}
		}
	default:
		httpErrors.NotFound(w)

	}
}

func ListVideos(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpErrors.MethodNotAllowed(w)
		return
	}
	movies := admin.Video{}
	resp, err := movies.List()
	if err != nil {
		httpErrors.InternalServerError(w)
		return
	}
	jsonResp, err := serialize(resp)
	if err != nil {
		httpErrors.InternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
