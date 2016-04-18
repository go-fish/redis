package main

import (
	"fmt"
	"testing"

	"github.com/HearingFish/redis"
)

var FishClient *redis.Conn

func TestFishRedisConnection(t *testing.T) {
	var err error
	if FishClient, err = redis.NewConnection(fmt.Sprintf("%s:%d", host, port)); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	s, err := FishClient.Ping()
	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	if !s {
		t.Fatalf("Failed")
	}
}

func BenchmarkFishRedisPing(b *testing.B) {
	var err error
	FishClient.Del("hello")
	for i := 0; i < b.N; i++ {
		_, err = FishClient.Ping()
		if err != nil {
			b.Fatalf(err.Error())
			break
		}
	}
}

func BenchmarkFishRedisSet(b *testing.B) {
	var err error
	FishClient.Del("hello")
	for i := 0; i < b.N; i++ {
		err = FishClient.Set("hello", 1)
		if err != nil {
			b.Fatalf(err.Error())
			break
		}
	}
}

func BenchmarkFishRedisGet(b *testing.B) {
	var err error
	FishClient.Del("hello")
	for i := 0; i < b.N; i++ {
		_, err = FishClient.Get("hello")
		if err != nil {
			b.Fatalf(err.Error())
			break
		}
	}
}

func BenchmarkFishRedisIncr(b *testing.B) {
	var err error
	FishClient.Del("hello")
	for i := 0; i < b.N; i++ {
		_, err = FishClient.Incr("hello")
		if err != nil {
			b.Fatalf(err.Error())
			break
		}
	}
}

func BenchmarkFishRedisLPush(b *testing.B) {
	var err error
	FishClient.Del("hello")
	for i := 0; i < b.N; i++ {
		_, err = FishClient.Lpush("hello", i)
		if err != nil {
			b.Fatalf(err.Error())
			break
		}
	}
}

func BenchmarkFishRedisLRange10(b *testing.B) {
	var err error
	FishClient.Del("hello")
	for i := 0; i < b.N; i++ {
		_, err = FishClient.Lrange("hello", 0, 10)
		if err != nil {
			b.Fatalf(err.Error())
			break
		}
	}
}

func BenchmarkFishRedisLRange100(b *testing.B) {
	var err error
	FishClient.Del("hello")
	for i := 0; i < b.N; i++ {
		_, err = FishClient.Lrange("hello", 0, 100)
		if err != nil {
			b.Fatalf(err.Error())
			break
		}
	}
}
