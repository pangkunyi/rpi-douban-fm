/* vim: set ts=2 sw=2 enc=utf-8: */
package player 

import (
	"strings"
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

var (
	cmd *exec.Cmd
	inPipe io.WriteCloser
	outPipe *bufio.Reader
	ch = make(chan int)
	CurTrack Track
	WillLoadNextChannel = false
	running = make(chan bool, 1)
)

type PlayList interface {
	GetTracks() []Track
	ReLoad() error
	GetChannel() string
}

type Track interface {
	GetChannel() string
	GetTitle() string
	GetAlbumTitle() string
	GetAlbumCover() string
	GetArtistName() string
	GetLink() string
}

func PlayListAndWait(playList PlayList, fn func(Track)) error {
	if CurTrack != nil && playList.GetChannel() == CurTrack.GetChannel() {
		return nil
	}
	CurTrack = nil
	WillLoadNextChannel = true
	Stop()
	running <- true
	WillLoadNextChannel = false
	for {
		err := playList.ReLoad()
		if err != nil {
			return err
		}
		for _, track := range playList.GetTracks() {
			err := PlayTrackAndWait(track, fn)
			if err != nil {
				return err
			}
			if WillLoadNextChannel {
				WillLoadNextChannel = false
				goto outter
			}
		}
	}
	outter:
	<- running
	return nil
}

func PlayTrackAndWait(track Track, fn func(Track)) error {
	CurTrack = track
	fn(track)
	return PlayAudioAndWait(track.GetLink())
}

func PlayAudioAndWait(audio string) error {
	loadCmd :="l "+audio+"\n"
	_, err := inPipe.Write([]byte(loadCmd))
	if err !=nil {
		return err
	}

	result := <- ch
	for result == 1 { // wait for stop
		result = <- ch
	}
	return nil
}

func PauseOrResume() error {
	loadCmd :="p\n"
	_, err := inPipe.Write([]byte(loadCmd))
	if err !=nil {
		return err
	}

	return nil
}

func Stop() error {
	loadCmd :="s\n"
	_, err := inPipe.Write([]byte(loadCmd))
	if err !=nil {
		return err
	}

	return nil
}

func PlayNext() error{
	WillLoadNextChannel = false
	return Stop()
}

func process() {
	for{
		_line,_, err := outPipe.ReadLine()
		line := string(_line)
		if err !=nil {
			if err == io.EOF {
				ch <- 0
			}else{
				ch <- -1
			}
			continue
		}
		if strings.HasPrefix(line,"@P 0"){
			ch <- 0
		}else if strings.HasPrefix(line,"@P 1"){
			ch <- 1
		}
	}
}

func StartAndWait(){
	cmd = exec.Command("sudo", "mpg123", "-R")
	var err error
	inPipe,err =cmd.StdinPipe()
	if err !=nil {
		panic(err)
	}
	_outPipe,err :=cmd.StdoutPipe()
	if err !=nil {
		panic(err)
	}
	outPipe = bufio.NewReader(_outPipe)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("Waiting for command to finish...")
	go process()
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	fmt.Printf("mpg123 stoped.\n")
}
