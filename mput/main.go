// mput - create an object
//
// http://apidocs.joyent.com/manta/mput.html
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/davecheney/manta"
)

func init() {
	manta.Flags.Parse(os.Args[1:])
}

func main() {
	if len(manta.Flags.Args()) < 1 {
		log.Fatal("remote path must be supplied")
	}
	client := manta.DefaultClient()
	req, err := client.NewRequest("PUT", manta.Flags.Args()[0], os.Stdin)
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
	if resp.StatusCode != 204 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalf("%s", body)
	}
}
