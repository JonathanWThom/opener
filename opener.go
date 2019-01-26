// Opener opens and closes applications for you.
// Jump start your day or shut down easily.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

			if *close {
				err = closeApp(app)
			} else {
				err = openApp(app, apps)
			}

			if err != nil {
				fmt.Println(err)
			}
		}(app)
	}

	wg.Wait()
}

func closeApp(app string) error {
	cmd := fmt.Sprintf("ps cax | grep '%s'", app)
	out, _ := exec.Command("bash", "-c", cmd).Output()
	if string(out) == "" {
		return nil
	}

	return exec.Command(pkill, "-SIGINT", "-a", app).Run()
}

func openApp(app string, apps applications) error {
	cmd := fmt.Sprintf("ps cax | grep '%s'", app)
	out, _ := exec.Command("bash", "-c", cmd).Output()
	if string(out) != "" {
		return nil
	}

	args := []string{"-a", app}

	for _, site := range apps[app] {
		args = append(args, site)
	}

	return exec.Command(open, args...).Run()
}
