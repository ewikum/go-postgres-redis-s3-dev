package main

import (
	"fmt"
	"os"
	 "strconv"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

type redisHandler struct{}

func (h redisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v:=getViews()
	v++
	updateViews(v)
	fmt.Fprintf(w, "Viewed count from redis : "+strconv.Itoa(v))
}

func redisConnect() redis.Conn {
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_PORT_6379_TCP_ADDR"),
		os.Getenv("REDIS_PORT_6379_TCP_PORT"),
	))
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func updateViews(cnt int) {

	c := redisConnect()
	defer c.Close()

	// Save cnt to Redis
	reply, err := c.Do("SET", "viewedcount",cnt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GET ", reply)
}

func getViews() int {

	c := redisConnect()
	defer c.Close()
	// Get views from Redis
	reply, err := c.Do("GET", "viewedcount")
	if err != nil {
		log.Fatal(err)
	}
	if reply != nil {
		s := string(reply.([]byte))
		log.Println("GET ",s )
		i,_:= strconv.Atoi(s)
		return i
	}

	return 0
}