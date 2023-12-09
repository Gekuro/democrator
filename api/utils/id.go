package utils

import (
	"fmt"
	"math/rand"

	"github.com/Gekuro/democrator/api/store"
	"gorm.io/gorm"
)

// GetRandomId func replaces X-es in the format to random captial letters or digits
const ID_FORMAT = "XXX-XXX-XXX"
const CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetRandomId() string {
	id := []byte(ID_FORMAT)
	for i := range id {
		if id[i] != 'X' {
			continue
		}
		id[i] = CHARS[rand.Intn(len(CHARS))]
	}

	return string(id)
}

func GetUnusedPollId(db *gorm.DB) (id string, err error) {
	for {
		id = GetRandomId()
		var exists bool;

		err = db.Model(&store.Poll{}).
			Select("COUNT(*) > 0").
			Where("id = ?", id).
			Find(&exists).
			Error
		
		if err != nil { // TODO log this
			return "", fmt.Errorf("error querying the database")
		}

		if !exists {
			break
		}
	}

	return
}