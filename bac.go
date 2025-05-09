package main

import (
	"io"
	"net/http"
	"os"
)


var url string = "https://www.oxfordlearnersdictionaries.com/media/english/uk_pron/g/gra/grati/gratitude__gb_1.mp3"
func main(){

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	res, _ := client.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic("audio")
	}

	file, _ := os.Create("audio.mp3")
	io.Copy(file, res.Body)
	
}