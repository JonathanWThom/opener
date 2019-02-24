// Opener opens and closes applications for you.
// Jump start your day or shut down easily.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
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
var set = flag.Bool("s", false, "set applications.json file from command line")

type applications map[string][]string

func homeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func main() {
	flag.Parse()

	// Open or create file at ~/applications.json
	file := openFile()

	// If file is empty write initial JSON to it
	fi, _ := file.Stat()
	if fi.Size() == 0 {
		_, err := file.WriteString("{}")
		if err != nil {
			log.Fatal(err)
		}

		// Reopen it to get updates
		file = openFile()
	}

	// Read file and parse into applications struct
	opts := readFile(file)
	apps := make(applications)
	err := json.Unmarshal(opts, &apps)
	if err != nil {
		log.Fatal(err)
	}

	// Handle reading and setting configuration if -s flag passed
	if *set {
		setConfig(opts, apps)
		return
	}

	// Loop through applications and either open or close them
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
	args := []string{"-a", app}

	for _, site := range apps[app] {
		args = append(args, site)
	}

	return exec.Command(open, args...).Run()
}

func openFile() *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func readFile(file *os.File) []byte {
	opts, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return opts
}

func printConfig(opts []byte) {
	// Idea: Change this to just read applications, and format a bit nicer
	n := len(opts)
	s := string(opts[:n])
	fmt.Println("Your current current configuration:\n")
	fmt.Printf("%s\n\n", s)
	fmt.Println("* Add to default list by typing 'a AppName'")
	fmt.Println("* Quit interactive session with 'quit'")
}

func setConfig(opts []byte, apps applications) {
	printConfig(opts)

	// Start repl
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		input = strings.TrimSuffix(input, "\n")

		if input[:1] == "a" {
			// Collect app name and add to default list
			if len(input) > 1 {
				apps["default"] = append(apps["default"], input[2:])
				json, _ := json.Marshal(apps)
				err := ioutil.WriteFile(path, json, 0755)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Succesfully added %v to default opener list\n", input[2:])
					file := openFile()
					opts := readFile(file)
					printConfig(opts)
				}
			} else {
				fmt.Println("Must pass app name as argument to 'a'")
			}
		}

		// End repl session
		if input == "quit" {
			fmt.Println("Goodbye")
			return
		}
	}
}
