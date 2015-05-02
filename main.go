package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	pool        *redis.Pool
	redisServer = flag.String("redisServer", ":6379", "")
	//redisPassword = flag.String("redisPassword", "", "")
)

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			/*if _, err := c.Do("AUTH", password); err != nil {
			    c.Close()
			    return nil, err
			}*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

//the appError object holds the error data
type appError struct {
	Error   error
	Message string
	Code    int
}

//the shorter version of the error is passed back to the user
type shortError struct {
	Message string
	Code    int
}

//holds the word and its definition
type dictEntry struct {
	Word       string
	Definition string
}

/*
define new handler function that returns an error if something goes wrong --
the http package doesn't understand functions that return error
*/
type customHandler func(http.ResponseWriter, *http.Request) *appError

/*
implement the http.Handler interface's ServeHTTP method on appHandler, so
appHandler can be passed to http.Handle
*/
func (fn customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		//todo log full error to file
		errorMsg, _ := json.Marshal(shortError{e.Message, e.Code})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Code)
		w.Write(errorMsg)
	}
}

/*
the autocomplete function
*/
func autocomp(w http.ResponseWriter, r *http.Request) *appError {
	conn := pool.Get()
	defer conn.Close()

	query := r.URL.Path[len("/autocomp/"):]

	redisReply, err := conn.Do("ZRANGEBYLEX", "myzset", "["+query, "["+query+"\xff")
	if err != nil {
		return &appError{err, "Database connection failed", 500}
	}

	wordList, err := redis.Strings(redisReply, err)
	if err != nil {
		return &appError{err, "Database reply conversion failed", 500}
	}

	jsonReply, err := json.Marshal(wordList)
	if err != nil {
		return &appError{err, "JSON encoding failed", 500}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonReply)
	return nil
}

func getDefinition(w http.ResponseWriter, r *http.Request) *appError {
	conn := pool.Get()
	defer conn.Close()

	query := r.URL.Path[len("/get/"):]

	redisReply, err := conn.Do("GET", query)
	if err != nil {
		return &appError{err, "Database connection failed", 500}
	}
	if redisReply == nil {
		return &appError{err, "Database entry not found", 404}
	}

	replyString, err := redis.String(redisReply, err)
	if err != nil {
		return &appError{err, "Database reply conversion failed", 500}
	}

	jsonReply, err := json.Marshal(dictEntry{Word: query, Definition: replyString})
	if err != nil {
		return &appError{err, "JSON encoding failed", 500}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonReply)
	return nil

}

func main() {
	//check if the static path is provided in the command line arguments
	args := os.Args
	if len(args) < 2 {
		err := fmt.Errorf("Error: please specify the static path as argument")
		log.Fatal(err)
	}

	//create the redis connection pool
	flag.Parse()
	pool = newPool(*redisServer)

	//start the web server
	staticPath := args[1]
	fs := http.FileServer(http.Dir(staticPath))

	http.Handle("/", fs)
	http.Handle("/autocomp/", customHandler(autocomp))
	http.Handle("/get/", customHandler(getDefinition))
	http.ListenAndServe(":8080", nil)
}
