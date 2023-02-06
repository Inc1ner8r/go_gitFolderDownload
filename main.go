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

	outString := make(chan string)
	outInt := make(chan string)
	var links = fetchLinks("https://api.github.com/repos/inciner8r/sample_data/contents/")
	makedir()
	fmt.Println(links)

	for _, link := range links {

		go newFunction(link, outString, outInt)
	}
	for i := 0; i < len(links); i++ {
		msg := <-outString
		msgInt := <-outInt

		fmt.Println("filename - " + msg + "\nsize - " + string(msgInt))
	}
	close(outInt)
	close(outString)
	fmt.Println("here")
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

func newFunction(link links, out chan string, out1 chan string) {
	fmt.Println("exec start")
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

	out <- link.Name
	out1 <- fmt.Sprint(n)
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
