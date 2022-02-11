package admin

import (
	"fmt"
	"io"
	"lan-chat/httpErrors"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	http.ServeFile(w, r, "./templates/admin.html")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		admin(w, r)
	case http.MethodPost:
		uploadMovie(w, r)
	}
}

func uploadMovie(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 32) // maximum upload size of 4.2 GB
	if err != nil {
		log.Println(err)
		httpErrors.UnProcessableEntry(w)
		return
	}
	var videoName = r.FormValue("movie-name")
	videoTyp := r.FormValue("type")

	if videoName == "" {
		log.Println("Empty Movie Name")
		httpErrors.UnProcessableEntry(w)
		return
	}

	var (
		epNo  int
		_epNo string
	)
	if r.FormValue("ep-no") == "" {
		epNo = 1
	} else {
		epNo, err = strconv.Atoi(r.FormValue("ep-no"))
		if err != nil || epNo < 0 {
			log.Println(err)
			httpErrors.UnProcessableEntry(w)
			return
		}
	}
	_epNo = strconv.Itoa(epNo)

	var fileName string
	switch videoTyp {
	case "movie":
		fileName = "movie"
	case "series":
		fileName = "ep-" + _epNo
	case "ova":
		fileName = "ova"
	default:
		{
			log.Println("Invalid Video Type")
			httpErrors.UnProcessableEntry(w)
			return
		}
	}

	videoFile, _, videoErr := r.FormFile("movie")
	imageFile, _, imageErr := r.FormFile("movie-image")
	subFile, _, subErr := r.FormFile("movie-sub")

	if (videoErr != nil) || (imageErr != nil) {
		log.Println("Err: ", videoErr, imageErr)
		httpErrors.UnProcessableEntry(w)
		return
	}
	_, err = os.Stat(fmt.Sprintf("videos/%s/%s/%s", videoTyp, videoName, fileName))

	var path string // path at which files are going to be saved
	if os.IsNotExist(err) {
		switch videoTyp {
		case "series":
			path = fmt.Sprintf("videos/%s/%s/%s", videoTyp, videoName, fileName) // create path in format videos/series/kochikame/ep1
		default:
			path = fmt.Sprintf("videos/%s/%s", videoTyp, videoName)
		}
		err = os.MkdirAll(path, 0755) // create directory in format videoTyp/videoName
		if err != nil {
			httpErrors.InternalServerError(w)
			return
		}

		var (
			videoDestPath = fmt.Sprintf("%s/%s.mp4", path, fileName)
			imageDestPath = fmt.Sprintf("%s/%s.png", path, fileName)
			subDestPath   = ""
		)
		if (saveFile(videoFile, videoDestPath) != nil) || (saveFile(imageFile, imageDestPath) != nil) {
			httpErrors.InternalServerError(w)
			return
		}
		if subErr == nil {
			subDestPath = fmt.Sprintf("%s/%s.vtt", path, fileName)
			if saveFile(subFile, subDestPath) != nil {
				httpErrors.InternalServerError(w)
				return
			}
		}

		video := Video{Name: videoName, Typ: videoTyp, EpNo: epNo,
			VideoLink: strings.Replace(videoDestPath, "videos", "file/video", 1),
			ImageLink: strings.Replace(imageDestPath, "videos", "file/poster", 1),
			SubLink:   strings.Replace(subDestPath, "videos", "file/sub", 1)}

		err = video.save()
		if err != nil {
			log.Println("Err while saving Video in Database", err)
			httpErrors.InternalServerError(w)
			return
		}
		w.WriteHeader(201)
	} else {
		httpErrors.StatusConflict(w)
	}
}

func saveFile(file multipart.File, path string) error {
	defer file.Close()
	Dst, DstErr := os.Create(path)

	if DstErr != nil {
		return DstErr
	}
	if _, err := io.Copy(Dst, file); err != nil {
		log.Println("Err while copying video", err)
		return err
	}
	return nil

}
