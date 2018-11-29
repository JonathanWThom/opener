// Opener opens stores and then opens and closes programs for you.
// Jump start your day or shut down easily.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/user"
	"sync"
)

var path = homeDir() + "/applications.json"
var wg = &sync.WaitGroup{}
var close = flag.Bool("c", false, "add c flag to close files")

type applications []string

func homeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func main() {
	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	apps := applications{}
	err = decoder.Decode(&apps)
	if err != nil {
		log.Fatal(err)
	}

	for _, app := range apps {
		wg.Add(1)
		go func(app string) {
			defer wg.Done()
			var cmd string
			if *close {
				cmd = "pkill"
			} else {
				cmd = "open"
			}

			err := exec.Command(cmd, "-a", app).Run()
			if err != nil {
				log.Fatal(err)
			}

		}(app)
	}
	wg.Wait()
}
