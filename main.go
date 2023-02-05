package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	var links = fetchLinks("https://api.github.com/repos/inciner8r/sample_data/contents/")

	fmt.Println(links)

	for _, link := range links {
		wg.Add(1)

		go newFunction(link, wg)
	}
	wg.Wait()
}

// generated with help of https://mholt.github.io/json-to-go/
type links struct {
	Name string `json:"name"`
	// Path        string `json:"path"`
	// Sha         string `json:"sha"`
	// Size        int    `json:"size"`
	// URL         string `json:"url"`
	// HTMLURL     string `json:"html_url"`
	// GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	// Type        string `json:"type"`
	// Links       struct {
	// 	Self string `json:"self"`
	// 	Git  string `json:"git"`
	// 	HTML string `json:"html"`
	// } `json:"_links"`
}

func fetchLinks(link string) []links {

	var linkslist []links
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &linkslist); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	return linkslist
}

func newFunction(link links, wg sync.WaitGroup) {
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

	fmt.Println(n, "bytes downloaded.")
	wg.Done()

}
