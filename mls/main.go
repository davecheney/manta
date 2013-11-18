package main

import (
	"io"
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
	io.Copy(os.Stdout, resp.Body)
}
