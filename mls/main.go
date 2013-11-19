// mls - list directory contents.
// http://apidocs.joyent.com/manta/mls.html
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecheney/manta"
)

var (
	flags = manta.Flags()
	l     bool
)

func init() {
	flags.BoolVar(&l, "l", false, "Use a long listing format.")
	flags.Parse(os.Args[1:])
}

func main() {
	if len(flags.Args()) < 1 {
		log.Fatal("remote path must be supplied")
	}
	client := manta.DefaultClient()
	resp, err := client.Get(flags.Arg(0))
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
