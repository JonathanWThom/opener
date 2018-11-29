// Opener opens stores and then opens and closes programs for you.
// Jump start your day or shut down easily.
package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/user"
	"sync"
)

var path = homeDir() + "/applications.json"
var wg = &sync.WaitGroup{}

type applications []string

func homeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func main() {
	// add flags for open vs close
	// add more descriptive errors
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
			err := exec.Command("open", "-a", app).Run()
			// TODO: Handle the application not being found better
			if err != nil {
				log.Fatal(err)
			}

		}(app)
	}
	wg.Wait()
}
