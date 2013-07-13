package xiami

import (
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"strconv"
	"math"
	"net/url"
	"strings"
	"player"
	"math/rand"
	"fmt"
)

type PlayList struct{
	XMLName xml.Name `xml:"playList"`
	TrackList TrackList `xml:"trackList"`
	Channel string
}

type TrackList struct{
	Tracks []Track `xml:"track"`
}

type Track struct{
	SongName string `xml:"song_name"`
	AlbumName string `xml:"album_name"`
	AlbumCover string `xml:"album_cover"`
	ArtistName string `xml:"artist_name"`
	Location string `xml:"location"`
	Channel string
}

func (this *Track) GetChannel() string{
	return this.Channel
}

func (this *Track) GetTitle() string{
	return this.SongName
}

func (this *Track) GetAlbumTitle() string{
	return this.AlbumName
}

func (this *Track) GetAlbumCover() string{
	return this.AlbumCover
}

func (this *Track) GetArtistName() string{
	return this.ArtistName
}

func (this *Track) GetLink() string{
	return this.Location
}

func (this *PlayList) GetChannel() string {
	return this.Channel
}

func (this *PlayList) GetTracks() []player.Track{
	tracks := []player.Track{}
	for _, track := range this.TrackList.Tracks {
		_track := track
		_track.Channel = this.Channel
		tracks = append(tracks, &_track)
	}
	return tracks
}

func (this *PlayList) ReLoad() error {
	url := fmt.Sprintf(this.Channel,rand.Int63())
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(body, this)
	if err != nil {
		return err
	}
	this.Decode()
	return nil
}

func (this *PlayList) Decode() {
	for idx, track:= range this.TrackList.Tracks {
		this.TrackList.Tracks[idx].Location = decode(track.Location)
	}
}

func decode(loc string) string{
	var _local10 int
	size, _ := strconv.Atoi(loc[:1])
	mainBody := loc[1:]
	factor := int(len(mainBody) / size)
	modFactor := len(mainBody) % size
	arr := make([]string, int(math.Max(float64(modFactor),float64(size))))
	idx := 0
	for idx < modFactor{
		sIdx := (factor+1) * idx
		arr[idx] = mainBody[sIdx:sIdx+factor+1]
		idx++
	}
	idx = modFactor
	for idx < size {
		sIdx := ((factor * (idx - modFactor)) + ((factor + 1) * modFactor))
		arr[idx] = mainBody[sIdx:sIdx+factor];
		idx++
	}

	var _local8 string
	idx = 0
	for idx < len(arr[0]) {
		_local10 = 0
		for _local10 < len(arr) {
			if idx < len(arr[_local10]){
				_local8 = (_local8 + arr[_local10][idx:idx+1])
			}
			_local10++
		}
		idx++
	}

	_local8,_ = url.QueryUnescape(_local8)
	var _local9 string
	_local9 = strings.Replace(_local8,"^", "0", -1)
	_local9 = strings.Replace(_local9,"+", " ", -1)
	return _local9
}
