/* vim: set ts=2 sw=2 enc=utf-8: */
package main

import (
	"code.google.com/p/go.net/websocket"
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"bytes"
	"fmt"
	"sync/atomic"
	"time"
	"douban"
	"player"
)

const (
	xiami_url string = `http://www.xiami.com/radio/xml/type/%v/id/%v?v=%v`
)

func backHandler(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("res/tpls/index.gtpl")
	t.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("res/tpls/index.gtpl")
	t.Execute(w, nil)
}

//func songHandler(ws *websocket.Conn) {
//	for {
//		oldSongVersion := songVersion
//		fmt.Println("song handler, ->",oldSongVersion,":",songVersion)
//		for oldSongVersion == songVersion {
//			time.Sleep(1 * time.Second)
//		}
//		fmt.Println("song version changed, ->",oldSongVersion,":",songVersion)
//		t,_ := template.ParseFiles("res/tpls/song.gtpl")
//		douban.LoadAlbumInfo()
//		var msg bytes.Buffer
//		t.Execute(&msg, douban.CurSong)
//		if err := websocket.Message.Send(ws, msg.String()); err != nil{
//			fmt.Printf("error send websocket msg: %v\n", msg)
//			break
//		}
//	}
//}

func togglePauseHandler(w http.ResponseWriter, r *http.Request){
	err := player.PauseOrResume()
	if err !=nil {
		w.Write([]byte(`{"success":false}`))
	}
	w.Write([]byte(`{"success":true}`))
}

//func nextHandler(w http.ResponseWriter, r *http.Request){
//	loadCmd :="S\n"
//	fmt.Printf("stop music.\n")
//	logStat()
//	_, err := inPipe.Write([]byte(loadCmd))
//	if err !=nil {
//		w.Write([]byte(`{"success":false}`))
//	}
//	w.Write([]byte(`{"success":true}`))
//}

func xiamiChannelHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	channel := fmt.Sprintf(xiami_url,vars["type"],vars["id"],"%v")
	
	t,_ := template.ParseFiles("res/tpls/channel.gtpl")
	t.Execute(w, douban.CurSong)
}
