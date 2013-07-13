/* vim: se ts=2 sw=2 enc=utf-8: */
package douban

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"encoding/xml"
	"player"
	"math/rand"
)

type PlayList struct {
	R int
	Song []Song
	Channel string
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
	Channel string
}

type AlbumEntry struct {
    XMLName     xml.Name `xml:"entry"`
    Summary string   `xml:"summary"`
}

func (this *Song) GetChannel() string{
	return this.Channel
}

func (this *Song) GetTitle() string{
	return this.Title
}

func (this *Song) GetAlbumTitle() string{
	return this.Album
}

func (this *Song) GetAlbumCover() string{
	return this.Picture
}

func (this *Song) GetArtistName() string{
	return this.Artist
}

func (this *Song) GetLink() string{
	return this.Url
}

func (this *PlayList) GetChannel() string {
	return this.Channel
}

func (this *PlayList) GetTracks() []player.Track{
	tracks := []player.Track{}
	for _, track := range this.Song {
		_track := track
		_track.Channel = this.Channel
		tracks = append(tracks, &_track)
	}
	return tracks
}

func (this *AlbumEntry) LoadAlbumInfo(albumId string) error{
	url := fmt.Sprintf("http://api.douban.com/music/subject/%s", albumId)
	resp,err := http.Get(url)
	if err!= nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!= nil {
		return err
	}
	err = xml.Unmarshal(body, this)
	if err!= nil {
		return err
	}
	return nil
}

func (this *PlayList) ReLoad() error{
	url := fmt.Sprintf(this.Channel,rand.Int63())
	resp,err := http.Get(url)
	if err!= nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!= nil {
		return err
	}
	return json.Unmarshal(body, this)
}

