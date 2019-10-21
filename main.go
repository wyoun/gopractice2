package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

type RecordedTime struct {
	UTime int64  `json:"unixtime"`
	Ip    string `json:"client_ip"`
}

var count uint32
var lastTime int64
var ch chan string
var done chan bool

//var done chan bool = make(chan bool)

func main() {
	done = make(chan bool)
	ch = make(chan string)
	go getEulerTime()

	http.HandleFunc("/", RootHandler)
	http.ListenAndServe(":12345", nil)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {

	ipaddr := r.RemoteAddr
	Reqtime := time.Now().Unix()

	fmt.Fprintf(w, "<p>Client IP: %s</p>", ipaddr)
	fmt.Fprintf(w, "<p>Last Fetched Time: %d</p>", lastTime)
	fmt.Fprintf(w, "<p>Request Time: %d</p>", Reqtime)
	fmt.Fprintf(w, "<p>Total number of time requests to the api: %d</p>", count)

	str := fmt.Sprintf("%s-%d-%d\n", r.RemoteAddr, Reqtime, lastTime)

	go func() {
		f, err := os.OpenFile("logs", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.WriteString(str)
		if err != nil {
			panic(err)
		}
		done <- true
	}()

	<-done

}

func getEulerTime() {

	for {
		resp, err := http.Get("http://worldtimeapi.org/api/ip")
		if err != nil {
			fmt.Printf("failed getting time from worldtimeapi.org!")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			fmt.Printf("failed to read response body!")
			return
		}
		count = count + 1
		var Rectime RecordedTime
		err = json.Unmarshal(body, &Rectime)
		lastTime = Rectime.UTime

		time.Sleep(time.Duration(math.Round(math.E*1000000)) * time.Microsecond)
	}

}
