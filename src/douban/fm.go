/* vim: se ts=2 sw=2 enc=utf-8: */
package douban

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"encoding/xml"
)

var(
	AlbumInfoLoaded = false
	CurSong Song
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

func LoadAlbumInfo() error{
	if(AlbumInfoLoaded){
		return nil
	}
	url := fmt.Sprintf("http://api.douban.com/music/subject/%s",CurSong.Aid)
	resp,err := http.Get(url)
	if err!= nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!= nil {
		return err
	}
	err = xml.Unmarshal(body, &CurSong.AlbumInfo)
	if err!= nil {
		return err
	}
	AlbumInfoLoaded=true
	return nil
}

func (this *PlayList) LoadPlayList(channel string) error{
	url := fmt.Sprintf("http://douban.fm/j/mine/playlist?type=n&sid=&pt=0.0&channel=%s&from=mainsite&r=daab079b3c", channel)
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

