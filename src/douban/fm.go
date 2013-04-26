/* vim: se ts=2 sw=2 enc=utf-8: */
package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

func chkErr(err interface{}){
	if err !=nil {
		panic(err)
	}
}

func main(){
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://douban.fm", nil)
	chkErr(err)
	req.Header.Add("Cookie",`flag="ok"; ac="1366708199"; bid="aLWYjC+lTZ0"; __utma=58778424.206407853.1366708203.1366708203.1366708203.1; __utmb=58778424.3.9.1366708235994; __utmc=58778424; __utmz=58778424.1366708203.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); dbcl2="40361056:9am1EervRqE"; fmNlogin="y"`)
	req.Header.Add("Referer","http://douban.fm/")
	req.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31")
	resp, err := client.Do(req)
	chkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	chkErr(err)
	fmt.Println(string(body))
}
