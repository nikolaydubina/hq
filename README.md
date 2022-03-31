## _happy little queue_

97% covered, 90LOC, 80_000RPS, integration test, auto-cleaning, lightweight

[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/hq.svg)](https://pkg.go.dev/github.com/nikolaydubina/hq)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/hq)](https://goreportcard.com/report/github.com/nikolaydubina/hq)

![](bobross.jpg)

When your Go code does not panic  
When your infra does not fail  
When your data is small  
When your data is temporary  
When all you need is a happy little queue  

```go
// once you have a redis connection
rdb := redis.NewClient(...)             // "github.com/go-redis/redis/v8"     

...

// you can boot a lightweight worker
worker := hq.ReceiveWorker{
    Redis:      rdb,
    Queue:      "my-queue",
    PoolIdle:   time.Minute,            // recommended!
    PoolActive: time.Millisecond * 50,  // recommended! 
    NumWorkers: 10,                     // recommended!
    Batch:      100,                    // recommended!
    Handler:    &handler,               // interface { Handle(message []byte) error }
}
go worker.Work()

...

// and send something
sender := hq.Sender{
    Redis:  rdb,
    Queue:  "my-queue",
    MaxLen: 10,
    TTL:    time.Hour * 4,
}
sender.Send([]byte("my-bytes"))

// in redis it is single list
// LLEN my-queue
```

It is as fast as Redis, so [should](https://www.digitalocean.com/community/tutorials/how-to-perform-redis-benchmark-tests) be around 80_000RPS.

P.S. "happy" because optimistic
