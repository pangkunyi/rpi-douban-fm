/* vim: set ts=2 sw=2 enc=utf-8: */
package main

import (
	"code.google.com/p/go.net/websocket"
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"bytes"
	"douban"
	"xiami"
	"player"
	"time"
)

const (
	xiami_url string = `http://www.xiami.com/radio/xml/type/%v/id/%v?v=%v`
	douban_url string = `http://douban.fm/j/mine/playlist?type=n&sid=&pt=0.0&channel=%v&from=mainsite&r=%v`
)

var (
	trackChan = make(chan bool)
)

type Track struct {
	Title string
	AlbumTitle string
	AlbumCover string
	ArtistName string
	Link string
}

func (this *Track) load(track player.Track){
	if track == nil {
		return
	}
	this.Title = track.GetTitle()
	this.AlbumTitle = track.GetAlbumTitle()
	this.AlbumCover = track.GetAlbumCover()
	this.ArtistName = track.GetArtistName()
	this.Link = track.GetLink()
}

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
		time.Sleep(1 * time.Second)
		if player.CurTrack == nil {
			continue
		}
		t,_ := template.ParseFiles("res/tpls/song.gtpl")
		var msg bytes.Buffer
		var track Track
		track.load(player.CurTrack)
		t.Execute(&msg, track)
		if err := websocket.Message.Send(ws, msg.String()); err != nil{
			fmt.Printf("error send websocket msg: %v\n", msg)
			break
		}
	}
}

func jsonHandle(err error, w http.ResponseWriter){
	if err !=nil {
		w.Write([]byte(`{"success":false}`))
	}
	w.Write([]byte(`{"success":true}`))
}

func togglePauseHandler(w http.ResponseWriter, r *http.Request){
	err := player.PauseOrResume()
	jsonHandle(err, w)
}

func nextHandler(w http.ResponseWriter, r *http.Request){
	err := player.PlayNext()
	jsonHandle(err, w)
}
func doubanChannelHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	channel := fmt.Sprintf(douban_url,vars["channel"],"%v")
	playList := &douban.PlayList{Channel:channel}
	viewChannel(playList, w, r)
}

func xiamiChannelHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	channel := fmt.Sprintf(xiami_url,vars["type"],vars["id"],"%v")
	playList := &xiami.PlayList{Channel:channel}
	viewChannel(playList, w, r)
}

func viewChannel(playList player.PlayList, w http.ResponseWriter, r *http.Request){
	if player.CurTrack == nil || playList.GetChannel() != player.CurTrack.GetChannel() {
		player.CurTrack = nil
		go player.PlayListAndWait(playList, trackInspector)
		for player.CurTrack == nil {
			time.Sleep(1 * time.Second)
		}
	}
	t,_ := template.ParseFiles("res/tpls/channel.gtpl")
	var track Track
	track.load(player.CurTrack)
	t.Execute(w, track)
}

func trackInspector(track player.Track){
	fmt.Println(track)
}
