// Opener opens and closes applications for you.
// Jump start your day or shut down easily.
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os/exec"
	"os/user"
	"sync"
)

const (
	pkill = "pkill"
	open  = "open"
)

var path = homeDir() + "/applications.json"
var wg = &sync.WaitGroup{}
var close = flag.Bool("c", false, "add c flag to close files")
var group = flag.String("g", "default", "specify which group of files to open or close")

type applications map[string][]string

func homeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func main() {
	flag.Parse()

	opts, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	apps := make(applications)
	err = json.Unmarshal(opts, &apps)
	if err != nil {
		log.Fatal(err)
	}

	for _, app := range apps[*group] {
		wg.Add(1)
		go func(app string) {
			defer wg.Done()
			var cmd string
			if *close {
				cmd = pkill
			} else {
				cmd = open
			}

			args := []string{"-a", app}

			if cmd == open {
				for _, site := range apps[app] {
					args = append(args, site)
				}
			}
			err = exec.Command(cmd, args...).Run()

			if err != nil {
				log.Fatal(err)
			}
		}(app)
	}

	wg.Wait()
}
