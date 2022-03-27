package hq

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type ReceiveWorker struct {
	Redis      *redis.Client
	Queue      string
	PoolIdle   time.Duration
	PoolActive time.Duration
	Handler    interface{ Handle(message []byte) error }
	NumWorkers uint
	Batch      uint
	stop       chan bool
}

func (r *ReceiveWorker) Work() {
	activeWorker := make(map[uint]bool, r.NumWorkers)
	done := make(chan uint, r.NumWorkers)
	r.stop = make(chan bool)
	t := time.NewTicker(r.PoolIdle)
	for {
		select {
		case <-r.stop:
			return
		case id := <-done:
			activeWorker[id] = false
		case <-t.C:
			for i := uint(0); i < r.NumWorkers; i++ {
				if activeWorker[i] {
					continue
				}
				activeWorker[i] = true
				go r.work(i, done)
			}
		}
	}
}

func (r *ReceiveWorker) work(id uint, done chan uint) {
	t := time.NewTicker(r.PoolActive)
	for {
		select {
		case <-r.stop:
			return
		case <-t.C:
			batch, err := r.Redis.LPopCount(context.Background(), r.Queue, int(r.Batch)).Result()
			if len(batch) == 0 || err != nil {
				if err != nil && !errors.Is(err, redis.Nil) {
					log.Printf("hq: %s: redis error: %s", r.Queue, err)
				}
				done <- id
				return
			}
			for _, m := range batch {
				if err := r.Handler.Handle([]byte(m)); err != nil {
					log.Printf("hq: %s: error: %s", r.Queue, err)
				}
			}
		}
	}
}

func (r *ReceiveWorker) Stop() { close(r.stop) }

type Sender struct {
	Redis  *redis.Client
	Queue  string
	MaxLen uint
	TTL    time.Duration
}

func (s Sender) Send(message []byte) error {
	ctx := context.Background()
	pipe := s.Redis.Pipeline()
	pipe.RPush(ctx, s.Queue, message)
	pipe.LTrim(ctx, s.Queue, 0, int64(s.MaxLen))
	pipe.Expire(ctx, s.Queue, s.TTL)
	_, err := pipe.Exec(ctx)
	return err
}
