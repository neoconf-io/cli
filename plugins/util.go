package plugins

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

func processInstallCmds(p plugin) {
	d := structure.GetPluginDir(string(p.dir), p.opt)
	if hasReadme(d) {
		b, err := ioutil.ReadFile(filepath.Join(d, "README.md"))
		if err != nil {
			panic(err)
		}

		runPostInstallCmd(findCmd(p, b), p.repo)
	}
}

func runPostInstallCmd(cmd *exec.Cmd, r repo) {
	if cmd == nil {
		return
	}

	o, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Running post-install command for '%s':\n%s", r, string(o))
}

func findCmd(p plugin, b []byte) *exec.Cmd {
	re := regexp.MustCompile(`(cd.*&&.)?yarn.install`)

	res := re.Find(b)
	if len(res) == 0 {
		return nil
	}

	args := strings.Split(strings.TrimSpace(string(res)), " ")

	var dir string
	if args[0] == "cd" {
		dir = args[1]
	}

	cmd := exec.Command("yarn", "install")
	cmd.Dir = filepath.Join(structure.GetPluginDir(string(p.dir), p.opt), dir)

	return cmd
}

func hasReadme(d string) bool {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		panic(err)
	}

	for _, v := range files {
		if v.Name() == "README.md" {
			return true
		}
	}

	return false
}

func writeList(p []plugin) {
	r := []string{}

	for _, v := range p {
		if v.repo == "" {
			continue
		}

		pluginString := string(v.repo)
		if v.branch != "" {
			pluginString = strings.Join([]string{string(v.repo), v.branch}, "@")
		}

		if v.opt {
			pluginString = pluginString + ":opt"
		}

		r = append(r, pluginString)
	}

	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(structure.Files.Plugins.O, b, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func updateCfgInit() {
	d, err := ioutil.ReadDir(structure.Dir.PluginCfg)
	if err != nil {
		panic(err)
	}

	f := []string{}

	for _, v := range d {
		if v.Name() == "init.lua" {
			continue
		}

		f = append(f, strings.TrimSuffix(v.Name(), ".lua"))
	}

	structure.WriteTmpl(structure.Files.PluginsInit, f)
}

func confirmation(msg string) bool {
	var response string

	fmt.Printf("%s (y/n) ", msg)

	if _, err := fmt.Scanln(&response); err != nil {
		panic(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Wrong input.")

		return confirmation(msg)
	}
}

const (
	minL   = 1
	offset = 1
)

func List() {
	p := getPlugins(getJSON())

	if len(p) < minL {
		fmt.Println("No plugins installed")

		return
	}

	txt := "%d: %s\n"
	txt2 := "%d: %s -- optional\n"

	for k, v := range p {
		if v.opt {
			fmt.Printf(txt2, k+offset, strings.Split(string(v.repo), "/")[1])

			continue
		}

		fmt.Printf(txt, k+offset, strings.Split(string(v.repo), "/")[1])
	}
}

func getSelections() []string {
	fmt.Print("Enter a number: ")

	reader := bufio.NewReader(os.Stdin)

	s, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)

	return strings.Split(s, " ")
}
