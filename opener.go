// Opener opens stores and then opens and closes programs for you.
// Jump start your day or shut down easily.
package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/user"
)

var path = homeDir() + "/programs.json"

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

	// could these be done concurrently?
	for _, p := range prog {
		err := exec.Command("open", "-a", p).Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
