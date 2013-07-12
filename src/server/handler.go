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

func backHandler(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("res/tpls/index.gtpl")
	t.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("res/tpls/index.gtpl")
	t.Execute(w, nil)
}

func songHandler(ws *websocket.Conn) {
	for {
		oldSongVersion := songVersion
		fmt.Println("song handler, ->",oldSongVersion,":",songVersion)
		for oldSongVersion == songVersion {
			time.Sleep(1 * time.Second)
		}
		fmt.Println("song version changed, ->",oldSongVersion,":",songVersion)
		t,_ := template.ParseFiles("res/tpls/song.gtpl")
		douban.LoadAlbumInfo()
		var msg bytes.Buffer
		t.Execute(&msg, douban.CurSong)
		if err := websocket.Message.Send(ws, msg.String()); err != nil{
			fmt.Printf("error send websocket msg: %v\n", msg)
			break
		}
	}
}


func logStat(){
	fmt.Printf("stat:\n curChannel: %v\n channelVersion: %v\n songVersion: %v\n curMusicIdx: %v\n curSong: %v\n playlist: %v\n",
		curChannel, channelVersion, songVersion, curMusicIdx, douban.CurSong, playlist)
}

func togglePauseHandler(w http.ResponseWriter, r *http.Request){
	err := player.PauseOrResume()
	if err !=nil {
		w.Write([]byte(`{"success":false}`))
	}
	w.Write([]byte(`{"success":true}`))
}

func nextHandler(w http.ResponseWriter, r *http.Request){
	loadCmd :="S\n"
	fmt.Printf("stop music.\n")
	logStat()
	_, err := inPipe.Write([]byte(loadCmd))
	if err !=nil {
		w.Write([]byte(`{"success":false}`))
	}
	w.Write([]byte(`{"success":true}`))
}

func channelHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	oldChannel := curChannel
	curChannel = vars["channel"]
	if oldChannel != curChannel {
		cv :=atomic.AddInt32(&channelVersion, 1)
		fmt.Printf("channel will be changed: %v -> %v , new channel version: %v\n", oldChannel, curChannel, cv)
		go loopPlay(curChannel, cv)
		done = make(chan bool, 1)
		<- done
		done <- true
	}
	t,_ := template.ParseFiles("res/tpls/channel.gtpl")
	douban.LoadAlbumInfo()
	t.Execute(w, douban.CurSong)
}
