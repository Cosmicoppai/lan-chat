package admin

import (
	"errors"
	"fmt"
	"io"
	"lan-chat/httpErrors"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func uploadMovie(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 32) // maximum upload size of 4.2 GB
	if err != nil {
		log.Println(err)
		httpErrors.UnProcessableEntry(w)
		return
	}
	videoName, err := cleanUserInput(r.FormValue("movie-name"))
	if err != nil {
		httpErrors.UnProcessableEntry(w)
		return
	}
	videoTyp := r.FormValue("type")

	var (
		epNo  int
		_epNo string
	)

	var fileName string
	switch videoTyp {
	case "movie":
		fileName = "movie"
	case "series":
		fileName = "ep-" + _epNo
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

	var _path string // path at which files are going to be saved
	if os.IsNotExist(err) {
		switch videoTyp {
		case "series":
			_path = fmt.Sprintf("videos/%s/%s/%s", videoTyp, videoName, fileName) // create path in format videos/series/kochikame/ep1
		default:
			_path = fmt.Sprintf("videos/%s/%s", videoTyp, videoName)
		}
		err = os.MkdirAll(_path, 0755) // create directory in format videoTyp/videoName
		if err != nil {
			httpErrors.InternalServerError(w)
			return
		}

		var (
			videoDestPath = fmt.Sprintf("%s/%s.mp4", _path, fileName)
			imageDestPath = fmt.Sprintf("%s/%s.png", _path, fileName)
			subDestPath   = ""
		)
		if (saveFile(videoFile, videoDestPath) != nil) || (saveFile(imageFile, imageDestPath) != nil) {
			httpErrors.InternalServerError(w)
			return
		}
		if subErr == nil {
			subDestPath = fmt.Sprintf("%s/%s.vtt", _path, fileName)
			if saveFile(subFile, subDestPath) != nil {
				httpErrors.InternalServerError(w)
				return
			}
		}

		_video := Video{Name: videoName, Typ: videoTyp,
			VideoLink:  strings.Replace(videoDestPath, "videos", "file/video", 1),
			PosterLink: strings.Replace(imageDestPath, "videos", "file/poster", 1),
			SubLink:    strings.Replace(subDestPath, "videos", "file/sub", 1)}
		switch videoTyp {
		case "movie":
			_movie := Movie{Video: _video, Parts: 5}
			err = _movie.save()
		case "series":
			_series := Series{Video: _video, TotalEp: 10, EpAired: epNo}
			err = _series.save()
		case "ova":
			// err = _video.save()

		}

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

func saveFile(file multipart.File, path string) (err error) {
	Dst, DstErr := os.Create(path)

	if DstErr != nil {
		return DstErr
	}
	if _, err := io.Copy(Dst, file); err != nil {
		log.Println("Err while copying video", err)
		return err
	}
	if err = file.Close(); err != nil { // check for any errors while closing the file
		return err
	}
	return nil

}

func cleanUserInput(userInput string) (cleanInput string, err error) {
	cleanInput = strings.TrimPrefix(path.Clean("/"+userInput), "/")

	if cleanInput != "" {
		return cleanInput, nil
	}
	return cleanInput, errors.New("invalid Movie Name")

}
