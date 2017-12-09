package main

import (
	"fmt"
	"net/http"
)

func main() {
	
	http.Handle("/", indexHandler{})
	http.Handle("/redis/", redisHandler{})
	http.Handle("/s3/", s3Handler{})
	http.Handle("/db/", dbHandler{})

	http.ListenAndServe(":8080", nil)
	

}

type indexHandler struct{}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func hello() string{
	return "Hello"
}