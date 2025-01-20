package main

import (
	"log"
	"time"

	limiter "github.com/danielnegreiros/simple-sys-design-go/pkg/rate_limitier/bucket/cmd/limitier"
)


func main()  {
	refillRateMinute := 60
	size := 15
	key := "app"

	bl, err := limiter.NewBucketLimiter(refillRateMinute, size)
	if err != nil {
		log.Fatal(err)
	}

	for range 500 {
		if ok := bl.RequestToken(key); ok {
			log.Printf("request being served, %d tokens remaining", bl.GetCurrentUnits(key))
		}else{
			log.Printf("request denied, %d tokens remaining", bl.GetCurrentUnits(key))
		}
		time.Sleep(500 * time.Millisecond)
	}
}


