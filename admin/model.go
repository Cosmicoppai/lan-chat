package admin

import (
	"log"
	"time"
)

type Video struct {
	Typ  string `json:"typ"` // type of video (movie/series/ova)
	Name string `json:"name"`
	EpNo int    `json:"epNo"`

	// addresses in filepath
	ImageLink string    `json:"imageLink"`
	VideoLink string    `json:"videoLink"`
	SubLink   string    `json:"subLink"`
	CreatedAt time.Time `json:"createdAt"`
}

func (v *Video) save() error { // function to save data in postgres

	_, err := Db.Exec("INSERT INTO videos (typ, name, epno, imagelink, videolink, sublink) VALUES($1, $2, $3, $4, $5, $6)",
		v.Typ, v.Name, v.EpNo, v.ImageLink, v.VideoLink, v.SubLink)

	if err != nil {
		return err
	}
	return nil
}

func (v *Video) List() ([]Video, error) {
	rows, err := Db.Query("SELECT typ, name, epNo, imageLink, videoLink, subLink, createdAt from videos ORDER BY createdAT DESC")
	if err != nil {
		log.Println("Error while querying ListVideo: ", err)
		return []Video{}, err
	}
	defer rows.Close()
	var result []Video

	for rows.Next() {
		err = rows.Scan(&v.Typ, &v.Name, &v.EpNo, &v.ImageLink, &v.VideoLink, &v.SubLink, &v.CreatedAt)
		if err != nil {
			log.Println("Error while scanning the row: ", err)
			return []Video{}, nil
		}
		result = append(result, *v)
	}
	return result, nil

}
