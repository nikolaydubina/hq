## _happy little queue_

97% covered, 80LOC, integration test, auto-cleaning, lightweight

[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/hq.svg)](https://pkg.go.dev/github.com/nikolaydubina/hq)

![](doc/bobross.jpg)

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
    PoolActive: time.Microsecond * 50,  // recommended! 
    Handler:    &handler,               // interface { Handle(message []byte) error }
    NumWorkers: 10,                     // recommended!
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

P.S. "happy" because optimistic
