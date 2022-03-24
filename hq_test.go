package hq_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nikolaydubina/hq"
)

type CountHandler struct {
	vs   map[string]int
	lock sync.RWMutex
}

func (c *CountHandler) Handle(message []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.vs[string(message)]++
	return nil
}

// To test
// 1. start redis server
// 2. REDIS_HOST=localhost REDIS_PORT=6379 go test -coverprofile hq.cover ./...
func TestSendAndReceive(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
	})

	sender := hq.Sender{
		Redis:  rdb,
		Queue:  "my-queue",
		MaxLen: 10,
		TTL:    time.Minute,
	}

	handler := CountHandler{
		vs:   map[string]int{},
		lock: sync.RWMutex{},
	}

	worker := hq.ReceiveWorker{
		Redis:      rdb,
		Queue:      "my-queue",
		PoolIdle:   time.Microsecond * 50,
		PoolActive: time.Microsecond * 5,
		Handler:    &handler,
		NumWorkers: 10,
	}

	go worker.Work()

	counters := map[string]int{
		"a":                     100,
		"b":                     5,
		"c":                     17,
		"some-very-long-string": 1111,
	}

	for k, v := range counters {
		for i := 0; i < v; i++ {
			sender.Send([]byte(k))
		}
	}

	time.Sleep(time.Second)

	for k, v := range counters {
		if handler.vs[k] != v {
			t.Errorf("key(%s) exp(%d) != got(%d)", k, v, handler.vs[k])
		}
	}

	worker.Stop()
}
