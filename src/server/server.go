/* vim: set ts=2 sw=2 enc=utf-8: */
package main

import (
	"code.google.com/p/go.net/websocket"
	"html/template"
	"strings"
	"net/http"
	"log"
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/mux"
	"bytes"
	"io"
	"bufio"
	"io/ioutil"
	"fmt"
	"os/exec"
	"sync/atomic"
	"time"
)

type PlayList struct {
	R int
	Song []Song
}

type Song struct {
	Album string
	Picture string
	Ssid string
	Artist string
	Url string
	Company string
	Title string
	Rating_avg float64
	Length	int
	Subtype string
	Public_time string
	Sid string
	Aid string
	Kbps string
	Albumtitle string
	Like int
	AlbumInfo AlbumEntry
}

type AlbumEntry struct {
    XMLName     xml.Name `xml:"entry"`
    Summary string   `xml:"summary"`
}

func backHandler(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("res/tpls/index.gtpl")
	t.Execute(w, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	if curChannel == "unknown" {
		t,_ := template.ParseFiles("res/tpls/index.gtpl")
		t.Execute(w, nil)
	}else {
		t,_ := template.ParseFiles("res/tpls/channel.gtpl")
		loadAlbumInfo()
		t.Execute(w, curSong)
	}
}
func loadAlbumInfo(){
	if(albumInfoLoaded){
		return
	}
	url := "http://api.douban.com/music/subject/"+curSong.Aid;
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp,_ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(body, &curSong.AlbumInfo)
	fmt.Printf("album: %v\n", curSong.AlbumInfo)
	albumInfoLoaded=true
}

var songChangeChan = make(chan string)
func songHandler(ws *websocket.Conn) {
	for {
		oldSongVersion := songVersion
		fmt.Println("song handler, ->",oldSongVersion,":",songVersion)
		for oldSongVersion == songVersion {
			time.Sleep(1 * time.Second)
		}
		fmt.Println("song version changed, ->",oldSongVersion,":",songVersion)
		t,_ := template.ParseFiles("res/tpls/song.gtpl")
		loadAlbumInfo()
		var msg bytes.Buffer
		t.Execute(&msg, curSong)
		if err := websocket.Message.Send(ws, msg.String()); err != nil{
			fmt.Printf("error send websocket msg: %v\n", msg)
			break
		}
	}
}

var curChannel = "unknown"
var channelVersion int32
var songVersion int32
var playlist PlayList
var curMusicIdx int
var curSong Song
var done chan bool
var albumInfoLoaded = false

func logStat(){
	fmt.Printf("stat:\n curChannel: %v\n channelVersion: %v\n songVersion: %v\n curMusicIdx: %v\n curSong: %v\n playlist: %v\n",
		curChannel, channelVersion, songVersion, curMusicIdx, curSong, playlist)
}

func togglePauseHandler(w http.ResponseWriter, r *http.Request){
	loadCmd :="P\n"
	fmt.Printf("pause music.\n")
	logStat()
	_, err := inPipe.Write([]byte(loadCmd))
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
	cv :=atomic.AddInt32(&channelVersion, 1)
	fmt.Printf("channel will be changed: %v -> %v , new channel version: %v\n", oldChannel, curChannel, cv)
	go loopPlay(curChannel, cv)
	t,_ := template.ParseFiles("res/tpls/channel.gtpl")
	done = make(chan bool, 1)
	<- done
	done <- true
	loadAlbumInfo()
	t.Execute(w, curSong)
}

func loadPlayList(channel string){
	url := "http://douban.fm/j/mine/playlist?type=n&sid=&pt=0.0&channel="+channel+"&from=mainsite&r=daab079b3c"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp,_ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &playlist)
	fmt.Printf("playlist : %v\n", playlist)
}

func loopPlay(channel string, cv int32){
	for {
		loadPlayList(channel)
		curMusicIdx =0
		size :=len(playlist.Song)
		if size > 0{
			for {
				if channelVersion != cv {
					fmt.Printf("channel changed: %v -> %v\n", channel, curChannel)
		//			goto outter
				}
				if curMusicIdx >= size {
					break
				}
				err :=play()
				<- done
				if err !=nil {
					fmt.Printf("play: %v\n", err)
					goto outter
				}
			}
		}
	}
	outter:
}

var cmd *exec.Cmd
var inPipe io.WriteCloser
var outPipe *bufio.Reader

func play() error{
	atomic.AddInt32(&songVersion, 1)
	albumInfoLoaded=false
	curSong = playlist.Song[curMusicIdx]
	done <- true
	curMusicIdx ++
	mp3 := curSong.Url
	fmt.Printf("ready to play: %v\n", mp3)
	loadCmd :="load "+mp3+"\n"
	fmt.Printf("mpg123 %v\n", loadCmd)
	_, err := inPipe.Write([]byte(loadCmd))
	if err !=nil {
		return err
	}

	fmt.Printf("ready output....\n")
	for{
		_line,_, err := outPipe.ReadLine()
		line := string(_line)
		if err !=nil {
			if err == io.EOF {
				fmt.Println("eof line: ", line)
				return nil
			}
			fmt.Println("err line: ", line)
			return err
		}
		if strings.HasPrefix(line,"@P 0"){
			fmt.Println("exit line: ", line)
			return nil
		}
	}
	return nil
}

func startMpg123(){
	cmd = exec.Command("sudo", "mpg123", "-R")
	var err error
	inPipe,err =cmd.StdinPipe()
	if err !=nil {
		log.Fatal("mpg123 stdin: ", err)
	}

	_outPipe,err :=cmd.StdoutPipe()
	if err !=nil {
		log.Fatal("mpg123 stdout: ", err)
	}
	outPipe = bufio.NewReader(_outPipe)
	err = cmd.Start()
	if err != nil {
		log.Fatal("start mpg123: ", err)
	}
	fmt.Println("Waiting for command to finish...")
	err = cmd.Wait()
	fmt.Printf("Command finished with error: %v\n", err)
}

func main(){
	go startMpg123()
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/index.html", indexHandler)
	r.HandleFunc("/back.html", backHandler)
	r.HandleFunc("/next.html", nextHandler)
	r.Handle("/song.html", websocket.Handler(songHandler))
	r.HandleFunc("/togglePause.html", togglePauseHandler)
	r.HandleFunc("/music/{channel}.html", channelHandler)

	http.Handle("/static/", http.FileServer(http.Dir("res/assets")))
	http.Handle("/", r)
	err := http.ListenAndServe(":9001", nil)
	if err != nil{
		log.Fatal("ListenAndServe: ", err)
	}
}
