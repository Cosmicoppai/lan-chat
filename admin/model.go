package admin

import (
	"fmt"
	"log"
)

type Video struct {
	Typ  string `json:"typ"` // type of video (movie/series/ova)
	Name string `json:"name"`

	// addresses in filepath
	PosterLink string `json:"imageLink"`
	VideoLink  string `json:"videoLink"`
	SubLink    string `json:"subLink"`
}

type Series struct {
	Video   Video
	TotalEp int `json:"totalEp"`
	EpAired int `json:"epAired"`
}

type Movie struct {
	Video Video
	Parts int `json:"parts"`
}

func (v *Video) save() string { // function to save data in postgres

	id := ""

	err := Db.QueryRow("INSERT INTO show (typ, name, posterLink, videolink, sublink) VALUES($1, $2, $3, $4, $5) RETURNING id",
		v.Typ, v.Name, v.PosterLink, v.VideoLink, v.SubLink).Scan(&id)

	if err != nil {
		log.Printf("Error in saving video: %s", err.Error())
	}

	return id
}

func (m *Movie) save() error {
	id := m.Video.save()
	if id == "" {
		return fmt.Errorf("invalid Id")
	}
	_, err := Db.Exec("INSERT INTO movie_data (id, parts) VALUES ($1, $2)", id, m.Parts)
	if err != nil {
		return err
	}
	return nil
}

func (s *Series) save() error {
	id := s.Video.save()

	if id == "" {
		return fmt.Errorf("inavlid Id")
	}

	_, err := Db.Exec("INSERT INTO series_data  (id, totalEp, epAired) VALUES($1, $2, $3)", id, s.TotalEp, s.EpAired)
	if err != nil {
		return err
	}
	return nil
}

func (v *Video) List() ([]Video, error) {
	rows, err := Db.Query("SELECT typ, name, posterLink, videoLink, subLink FROM show ORDER BY createdAT DESC")
	if err != nil {
		log.Println("Error while querying ListVideo: ", err)
		return []Video{}, err
	}
	defer rows.Close()
	var result []Video

	for rows.Next() {
		err = rows.Scan(&v.Typ, &v.Name, &v.PosterLink, &v.VideoLink, &v.SubLink)
		if err != nil {
			log.Println("Error while scanning the row: ", err)
			return []Video{}, nil
		}
		result = append(result, *v)
	}
	return result, nil

}
