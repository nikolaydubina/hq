_happy little queue_

97% covered, 80LOC, integration test, auto-cleaning, lightweight

![](doc/bobross.jpg)

When your Go code does not panic.  
When your infra does not fail.  
When data is small.  
When all you need is a happy little queue.

```go
// once you have a redis connection
// "github.com/go-redis/redis/v8"
rdb := redis.NewClient(...)             

...

// you can boot a lightweight worker
// interface { Handle(message []byte) error }
worker := hq.ReceiveWorker{
    Redis:      rdb,
    Queue:      "my-queue",
    PoolIdle:   time.Minute,            // recommended!
    PoolActive: time.Microsecond * 50,  // recommended! 
    Handler:    &handler,               
    NumWorkers: 10,                     // recommended!
}
go worker.Work()

...

// and you send something
sender := hq.Sender{
    Redis:  rdb,
    Queue:  "my-queue",
    MaxLen: 10,
    TTL:    time.Hour * 4,
}
sender.Send([]byte("my-bytes"))
```

P.S. "happy" because optimistic
