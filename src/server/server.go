/* vim: set ts=2 sw=2 enc=utf-8: */
package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"player"
)

func main(){
	go player.StartAndWait()
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/index.html", indexHandler)
	r.HandleFunc("/back.html", backHandler)
	r.HandleFunc("/next.html", nextHandler)
	r.Handle("/song.html", websocket.Handler(songHandler))
	r.HandleFunc("/togglePause.html", togglePauseHandler)
	r.HandleFunc("/douban/{channel}.html", doubanChannelHandler)
	r.HandleFunc("/xiami/{channel}.html", xiamiChannelHandler)

	http.Handle("/static/", http.FileServer(http.Dir("res/assets")))
	http.Handle("/", r)
	err := http.ListenAndServe(":9001", nil)
	if err != nil{
		panic(err)
	}
}
