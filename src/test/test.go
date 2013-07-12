package main

import(
	"fmt"
	"xiami"
	"player"
)

const (
	url string = `http://www.xiami.com/radio/xml/type/17/id/3?v=%v`
)

func main(){
	go player.StartAndWait()
	playList := &xiami.PlayList{Channel:url}
	err := playList.ReLoad()
	if err != nil {
		panic(err)
	}
	player.PlayListAndWait(playList, test)
}

func test(track player.Track){
	fmt.Println("track: ", track)
}
