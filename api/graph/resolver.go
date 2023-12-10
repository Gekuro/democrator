package graph

import (
	"log"
	"time"

	"github.com/Gekuro/democrator/api/graph/model"
	"github.com/Gekuro/democrator/api/store"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const STREAMER_TICK_RATE_MS = 200;

type Resolver struct{
	DB *gorm.DB
	Streamer *PollStreamer
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{DB: db, Streamer: NewPollStreamer(db)}
}

type PollStreamer struct{
	PollChannels map[string][]chan []*model.Option
	DB *gorm.DB
}

func (s *PollStreamer) Start() {
	go func(){
		for {
			time.Sleep(time.Millisecond * STREAMER_TICK_RATE_MS)
			
			for pollId, chs := range s.PollChannels {
				if len(chs) == 0 {
					continue
				}

				var poll store.Poll
				err := s.DB.First(&poll).
					Where("id = ?", pollId).
					Error
				
				if err != nil {
					log.Printf("error reading streamer-enabled poll: %s", err)
					continue 
				}

				var dbOpts []store.Option
				var opts []*model.Option

				result := s.DB.Find(&dbOpts, "poll_id = ?", pollId)
				if result.Error != nil || result.RowsAffected < 2 {
					continue
				}

				for i := range dbOpts {
					opts = append(opts, &model.Option{Name: dbOpts[i].Name, Votes: dbOpts[i].Votes})
				}

				for _, ch := range chs {
					ch <- opts
				}
			}
		}
	}()
}

func (s *PollStreamer) Subscribe(pollID string) <-chan []*model.Option {
	ch := make(chan []*model.Option)

	if _, ok := s.PollChannels[pollID]; !ok {
		s.PollChannels[pollID] = make([]chan []*model.Option, 0)
	}
	s.PollChannels[pollID] = append(s.PollChannels[pollID], ch)

	return ch
}

func NewPollStreamer(db *gorm.DB) *PollStreamer {
	s := &PollStreamer{DB: db, PollChannels: make(map[string][]chan []*model.Option)}
	s.Start()
	return s
}
