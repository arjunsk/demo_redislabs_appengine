package main

import (
	"fmt"
	"net/http"
	"github.com/garyburd/redigo/redis"
	"google.golang.org/appengine"
)

var redisPool *redis.Pool

func main() {
	redisAddr := "redis-XXXXX.c9.us-east-1-2.ec2.cloud.redislabs.com:XXXXX"

	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisAddr)
			return conn, err
		},
	}

	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	redisConn := redisPool.Get()
	defer redisConn.Close()

	count, err := redisConn.Do("INCR", "count")
	if err != nil {
		msg := fmt.Sprintf("Could not increment count: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Count: %d", count)
}
