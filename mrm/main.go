// mput - create an object
//
// http://apidocs.joyent.com/manta/mput.html
package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/davecheney/manta"
)

var flags = manta.Flags()

func init() {
	flags.Parse(os.Args[1:])
}

func main() {
	if len(flags.Args()) < 1 {
		log.Fatal("remote path must be supplied")
	}
	client, err := manta.DefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Delete(flags.Arg(0))
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
