package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	outName := make(chan string)
	outSize := make(chan string)
	var links = fetchLinks("https://api.github.com/repos/inciner8r/sample_data/contents/")
	makedir()

	for _, link := range links {

		go download(link, outName, outSize)
	}
	for i := 0; i < len(links); i++ {
		name := <-outName
		size := <-outSize

		fmt.Println("filename - " + name + "\nsize - " + size)
	}
	close(outName)
	close(outSize)
}

type links struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

func fetchLinks(link string) []links {

	var linkslist []links
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if err := json.Unmarshal(body, &linkslist); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	return linkslist
}

func download(link links, outName chan string, outSize chan string) {
	output, err := os.Create("./downloads/" + link.Name)
	if err != nil {
		fmt.Println("Error while creating", "ok.jpg", "-", err)
		log.Fatal()
	}
	defer output.Close()
	response, err := http.Get(link.DownloadURL)
	if err != nil {
		fmt.Println("Error while downloading", link.DownloadURL, "-", err)
		log.Fatal()

	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", link.DownloadURL, "-", err)
		log.Fatal()
	}

	outName <- link.Name
	outSize <- fmt.Sprint(n)
}

func makedir() {
	path := "./downloads"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
