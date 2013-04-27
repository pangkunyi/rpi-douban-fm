/* vim: set ts=2 sw=2 enc=utf-8: */
package main

import (
	"html/template"
	"strings"
	"net/http"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"bufio"
	"io/ioutil"
	"fmt"
	"os/exec"
//	"time"
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
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("res/tpls/index.gtpl")
	t.Execute(w, nil)
}

func channelHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	go loopPlay(vars["channel"])
	t,_ := template.ParseFiles("res/tpls/channel.gtpl")
	t.Execute(w, nil)
}

func loadPlayList(channel string) PlayList{
	url := "http://douban.fm/j/mine/playlist?type=n&sid=&pt=0.0&channel="+channel+"&from=mainsite&r=daab079b3c"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp,_ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var playlist PlayList
	json.Unmarshal(body, &playlist)
	fmt.Printf("playlist : %v\n", playlist)
	return playlist
}

func loopPlay(channel string){
	for {
		playlist :=loadPlayList(channel)
		cur :=0
		size :=len(playlist.Song)
		if size > 0{
			for {
				if cur >= size {
					break
				}
				err :=play(cur, &playlist)
				cur ++
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

func play(cur int , playlist *PlayList) error{
	mp3 := playlist.Song[cur].Url
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
//		fmt.Println(line)
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
	r.HandleFunc("/music/{channel}.html", channelHandler)

	http.Handle("/static/", http.FileServer(http.Dir("res/assets")))
	http.Handle("/", r)
	err := http.ListenAndServe(":9001", nil)
	if err != nil{
		log.Fatal("ListenAndServe: ", err)
	}
}
