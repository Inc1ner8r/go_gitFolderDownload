package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var links = fetchLinks("https://api.github.com/repos/jonascarpay/Wallpapers/contents/papes")
	downloadFiles(links)
}

// generated with help of https://mholt.github.io/json-to-go/
type links struct {
	// Name        string `json:"name"`
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
	// Loop through the data node for the downloadURL
	// for _, rec := range links {
	// 	//fmt.Println(rec)
	// 	dloadLinks = append(dloadLinks, rec.DownloadURL)
	// }
	//fmt.Println(dloadLinks)
	return linkslist
}

func downloadFiles(links []links) {
	for _, link := range links {
		fmt.Println(link)
	}
}
