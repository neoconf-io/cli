package plugins

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/abenz1267/neoconf/structure"
)

type updated struct {
	sync.RWMutex
	list []plugin
}

func (u *updated) append(p plugin) {
	u.Lock()
	defer u.Unlock()

	u.list = append(u.list, p)
}

func Update() {
	var wg sync.WaitGroup

	items := &updated{}

	i := getJSON()

	p := getPlugins(i)

	if l := len(p); l > 0 {
		wg.Add(l)

		for _, v := range p {
			if !structure.Exists(structure.GetPluginDir(string(v.dir), v.opt)) {
				wg.Done()

				continue
			}

			go update(v, items, &wg)
		}
	}

	wg.Wait()

	n := len(items.list)
	if n > 0 && confirmation(fmt.Sprintf("%d packages have been updated. Show info?", n)) {
		for _, plugin := range items.list {
			showUpdateInfo(plugin)
		}
	}
}

func showUpdateInfo(p plugin) {
	cmd := exec.Command("git", "log", "--pretty=format:- %s", "@{1}..")
	cmd.Dir = structure.GetPluginDir(string(p.dir), p.opt)

	o, err := cmd.Output()
	if err == nil {
		fmt.Printf("%s:\n", strings.Replace(filepath.Base(cmd.Dir), "_", "/", 1))
		fmt.Println(string(o))
		fmt.Println()
	}
}

func update(p plugin, items *updated, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("git", "pull")
	cmd.Dir = structure.GetPluginDir(string(p.dir), p.opt)

	b := filepath.Base(cmd.Dir)

	o, err := cmd.Output()
	if err != nil {
		fmt.Printf("Updating '%s': %s", b, err)

		return
	}

	if res := string(o); strings.Contains(res, "Already up to date") {
		fmt.Printf("Updating '%s': %s", strings.Replace(b, "_", "/", 1), res)

		return
	}

	processInstallCmds(p)

	items.append(p)
}
