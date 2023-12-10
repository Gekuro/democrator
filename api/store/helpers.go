package store

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

func IsPollExpired(db *gorm.DB, id string) (bool, error) {
	var poll Poll
	result := db.Where("id = ?", id).Find(&poll)

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("poll doesn't exist")
	}
	if poll.ExpiresAt.Before(time.Now()) {
		return true, nil
	}
	return false, nil
}