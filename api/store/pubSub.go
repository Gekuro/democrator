package store

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Gekuro/democrator/api/graph/model"
	"gorm.io/gorm"
)

type PollPubSub struct {
	DB           *gorm.DB
	PollChannels map[string][]chan []*model.Option
	tickRate	 time.Duration
	mu           sync.Mutex
	wg			 sync.WaitGroup
}

func (s *PollPubSub) Start() {
	go func() {
		for {
			time.Sleep(s.tickRate)

			s.mu.Lock()
			for pollId, chs := range s.PollChannels {
				if len(chs) == 0 {
					continue
				}

				var poll Poll
				err := s.DB.First(&poll).
					Where("id = ?", pollId).
					Error

				if err != nil {
					log.Printf("error reading streamer-enabled poll: %s", err)
					continue
				}

				dbOpts := make([]Option, 0)
				opts := make([]*model.Option, 0)

				result := s.DB.Find(&dbOpts, "poll_id = ?", pollId)
				if result.Error != nil || result.RowsAffected < 2 {
					continue
				}

				for i := range dbOpts {
					opts = append(opts, &model.Option{Name: dbOpts[i].Name, Votes: dbOpts[i].Votes})
				}

				for _, ch := range chs {
					s.wg.Add(1)
					go func(c chan []*model.Option) {
						defer s.wg.Done()

						select {
						case c <- opts:
						case <-time.After(20 * time.Millisecond):
							// TODO stuck channels should be removed
						}
					}(ch)
				}
				s.wg.Wait()
			}
			s.mu.Unlock()
		}
	}()
}

func (s *PollPubSub) Subscribe(pollID string) <-chan []*model.Option {
	ch := make(chan []*model.Option)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.PollChannels[pollID]; !ok {
		s.PollChannels[pollID] = make([]chan []*model.Option, 0)
	}
	s.PollChannels[pollID] = append(s.PollChannels[pollID], ch)

	return ch
}

func NewPollPubSub(db *gorm.DB) *PollPubSub {
	var rate time.Duration
	ms, err := strconv.Atoi(os.Getenv("STREAMER_TICK_RATE_MS"))
	if err != nil {
		log.Printf("error parsing STREAMER_TICK_RATE_MS environment variable, will set the tick rate to default 200 ms")
		rate = time.Millisecond * 300
	}else{
		rate = time.Millisecond * time.Duration(ms)
	}

	s := &PollPubSub{
		DB: db, 
		PollChannels: make(map[string][]chan []*model.Option), 
		tickRate: rate, 
		mu: sync.Mutex{},
		wg: sync.WaitGroup{},
	}
	s.Start()
	return s
}