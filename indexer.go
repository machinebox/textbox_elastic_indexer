package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/machinebox/sdk-go/textbox"
)

var (
	client *http.Client
	tbox   *textbox.Client
)

func main() {
	var dataset, es, textboxAddr string
	flag.StringVar(&dataset, "dataset", "./bbcsport", "directory where the root dir of the dataset is")
	flag.StringVar(&es, "es", "http://localhost:9200", "Elastic Search address")
	flag.StringVar(&textboxAddr, "textbox", "http://localhost:8080", "Textbox address")
	flag.Parse()

	client = &http.Client{
		Timeout: 5 * time.Second,
	}

	tbox = textbox.New(textboxAddr)
	index := "news_textbox/articles"

	log.Println("[INFO]: Using ES on ", es)
	log.Println("[INFO]: Using Textbox on ", textboxAddr)
	log.Println("[INFO]: Start indexing articles from ", dataset, "to the index/type", index)

	// Walks the ./dataset inserting any .txt into Elastic Search
	// the index and type will be "/news_raw/articles"
	filepath.Walk(dataset, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !strings.HasSuffix(info.Name(), ".txt") {
			return nil
		}
		err = insertWithTextboxES(es, index, path)
		if err != nil {
			log.Println("[ERROR]: Inserting article", path, err)
		}
		return nil
	})

	log.Println("[INFO]: Finished")
}

// inserts an article pre-procesing it with textbox on Elastic Search with this structure
// {
//	 id: "xxxxxx",
//   title: "title of the article"
//   content: "content of the article",
//   keywords: "<most relevant keywords>",
//   people: "<people named in the content>"
//   places: "<places named in the content>"
// }
func insertWithTextboxES(es, index, path string) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	text := string(b)

	resp, err := tbox.Check(strings.NewReader(text))
	if err != nil {
		return errors.New("textbox: " + err.Error())
	}

	keywords := []string{}
	people := []string{}
	places := []string{}
	for _, k := range resp.Keywords {
		keywords = append(keywords, k.Keyword)
	}
	for _, s := range resp.Sentences {
		for _, ent := range s.Entities {
			if ent.Type == "person" {
				people = append(people, ent.Text)
			}
			if ent.Type == "place" {
				places = append(places, ent.Text)
			}
		}
	}
	// split to get the title and the content
	split := strings.SplitN(text, "\n", 2)
	body := map[string]interface{}{
		"title":    split[0],
		"content":  split[1],
		"keywords": keywords,
		"people":   people,
		"places":   places,
	}
	return postES(es, index, path, body)
}

// Post to Elastic Search, and returns an error in case is not success
// (there are good elastic search Go libs to do that, but is easy enoughs)
func postES(es string, index string, path string, body map[string]interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	u, err := url.Parse(es)
	if err != nil {
		return err
	}
	u.Path = index
	reader := bytes.NewReader(b)
	r, err := http.NewRequest(http.MethodPost, u.String(), reader)
	if err != nil {
		return err
	}

	resp, err := client.Do(r)
	if err != nil {
		return errors.New("ES error:" + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("[ERROR]: reading ES response", path, err)
		} else {
			log.Println("[ERROR]: ES error ", string(respBody))
		}
		return errors.New("Error creating article on Elastic Search " + resp.Status)
	}
	log.Println("[INFO]: Created article from ", path)
	return nil
}
