package movieHandler

import (
	"encoding/json"
	"lan-chat/admin"
	"log"
)

func serialize(videos []admin.Video) ([]byte, error) {
	serializedResult, err := json.Marshal(videos)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return serializedResult, nil

}
