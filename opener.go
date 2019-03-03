// Opener opens and closes applications for you.
// Jump start your day or shut down easily.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
var sys System

type applications map[string][]string

// System implements Run, which is a wrapper for exec.Command.Run
type System interface {
	Run(string, ...string) error
	Output(string, ...string) ([]byte, error)
}

// UserSystem is intantiated with the Run method
type UserSystem struct{}

func main() {
	sys = UserSystem{}
	Run()
}

// Run does all the work and is testable
func Run() {
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

// Run is wrapper around exec.Command.Run and can be mocked in tests
func (r UserSystem) Run(command string, args ...string) error {
	return exec.Command(command, args...).Run()
}

func (r UserSystem) Output(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}

func homeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func closeApp(app string) error {
	cmd := fmt.Sprintf("ps cax | grep '%s'", app)
	out, _ := sys.Output("bash", "-c", cmd)
	if string(out) == "" {
		return nil
	}

	return sys.Run(pkill, "-SIGINT", "-a", app)
}

func openApp(app string, apps applications) error {
	args := []string{"-a", app}

	for _, site := range apps[app] {
		args = append(args, site)
	}

	return sys.Run(open, args...)
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
	n := len(opts)
	s := string(opts[:n])
	fmt.Printf("Your current current configuration:\n\n")
	fmt.Printf("%s\n\n", s)
	fmt.Println("* Add an app to your default group by typing 'a AppName'")
	fmt.Println("* Remove an app from your default group by typing 'd AppName'")
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

		// Parse options
		arg := input[:1]
		if arg == "a" || arg == "d" {
			err := modifyApps(input, &apps, arg)
			if err != nil {
				fmt.Println(err)
			}
		} else if input == "quit" {
			fmt.Println("Goodbye")
			return
		} else {
			fmt.Printf("Could not recognize command: %s\n", input)
		}
	}
}

func modifyApps(input string, apps *applications, arg string) error {
	if len(input) > 1 {
		search := input[2:]
		switch arg {
		case "a":
			err := addApp(search, apps)
			if err != nil {
				return err
			}
		case "d":
			err := removeApp(search, apps)
			if err != nil {
				return err
			}
		}
	} else {
		msg := fmt.Sprintf("Must pass app name as argument to '%s'", arg)
		return errors.New(msg)
	}

	return nil
}

// addApp adds an app to the default group
func addApp(search string, apps *applications) error {
	_, found := stringInSlice(search, (*apps)["default"])
	if found == true {
		msg := fmt.Sprintf("%s already exists in default opener group", search)
		return errors.New(msg)
	}

	(*apps)["default"] = append((*apps)["default"], search)
	msg := fmt.Sprintf("Successfully added %s to default opener group", search)
	err := rewriteApps(apps, msg)
	if err != nil {
		return err
	}

	return nil
}

// removeApp looks for the application in the default group, and removes it if found
func removeApp(search string, apps *applications) error {
	if i, found := stringInSlice(search, (*apps)["default"]); found == true {
		(*apps)["default"] = append((*apps)["default"][:i], (*apps)["default"][i+1:]...)
		msg := fmt.Sprintf("Successfully removed %s from default opener group", search)
		err := rewriteApps(apps, msg)
		if err != nil {
			return err
		}
	} else {
		msg := fmt.Sprintf("Could not find %s in configuration\n", search)
		return errors.New(msg)
	}

	return nil
}

func rewriteApps(apps *applications, msg string) error {
	json, err := json.Marshal(apps)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, json, 0755)
	if err != nil {
		return err
	}

	fmt.Println(msg)
	file := openFile()
	opts := readFile(file)
	printConfig(opts)

	return nil
}

func stringInSlice(a string, list []string) (int, bool) {
	for i, b := range list {
		if b == a {
			return i, true
		}
	}
	return 0, false
}
