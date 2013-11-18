package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/davecheney/manta"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("remote path must be supplied")
	}
	client := manta.DefaultClient()
	req, err := client.NewRequest("GET", os.Args[1], nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.SignRequest(req); err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalf("%s", body)
	}
	s := struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Mtime string `json:"mtime"`
	}{}
	d := json.NewDecoder(resp.Body)
	for {
		if err := d.Decode(&s); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if s.Type == "directory" {
			fmt.Println(s.Name + "/")
		} else {
			fmt.Println(s.Name)
		}
	}
}
