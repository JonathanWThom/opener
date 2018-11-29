// Opener opens stores and then opens and closes programs for you.
// Jump start your day or shut down easily.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"sync"
)

var path = homeDir() + "/programs.json"
var wg = &sync.WaitGroup{}

type programs []string

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
	prog := programs{}
	err = decoder.Decode(&prog)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range prog {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			err := exec.Command("open", "-a", p).Run()
			if err != nil {
				fmt.Errorf("Could not open application: %s", p)
			}

		}(p)
	}
	wg.Wait()
}
