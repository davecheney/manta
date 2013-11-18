package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/davecheney/manta"
)

func main() {
	if len(flag.Args()) < 1 {
		log.Fatal("remote path must be supplied")
	}
	client := manta.Client{
		User:  MANTA_USER,
		KeyId: MANTA_KEY_ID,
		Key:   MANTA_KEY,
		Url:   MANTA_URL,
	}
	url := fmt.Sprintf("%s%s", client.Url, flag.Args()[0])
	req, err := http.NewRequest("GET", url, nil)
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
