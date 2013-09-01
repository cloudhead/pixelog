package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

func main() {
	l := flag.String("l", ":8080", "listen addr")
	file := flag.String("file", "t.gif", "tracking pixel filename")

	flag.Parse()

	pixel, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/"+*file, handlePixel(pixel))
	http.ListenAndServe(*l, nil)
}

func handlePixel(pixel []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "-"
		}
		agent := r.Header.Get("User-Agent")
		if agent == "" {
			agent = "-"
		}
		remote, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remote = "-"
		}
		log.Println(remote, r.URL.String(), strconv.QuoteToASCII(agent), referer)

		w.Header().Set("Cache-Control", "private, max-age=0")
		w.WriteHeader(http.StatusOK)
		w.Write(pixel)
	}
}
