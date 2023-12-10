package graph

import (
	"github.com/Gekuro/democrator/api/store"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	DB *gorm.DB
	PubSub *store.PollPubSub
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{DB: db, PubSub: store.NewPollPubSub(db)}
}

