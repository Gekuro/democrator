package store

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var cfg = &gorm.Config{}

type Poll struct {
	gorm.Model
	Id string
	ExpiresAt time.Time
}

type Option struct {
	gorm.Model
	PollId string
	Name string
	Votes int
}

func NewStore() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), cfg)
	if err != nil {
		return nil, err
	}
	
	db.AutoMigrate(&Poll{}, &Option{})

	return db, nil
}