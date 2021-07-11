package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/abenz1267/neoconf/structure"
)

// Clean removes orphaned plugin configs and folders.
func Clean() {
	f, err := ioutil.ReadDir(structure.Dir.PluginCfg)
	if err != nil {
		panic(err)
	}

	d := []string{}

	p := getPlugins(getJSON())

	for _, c := range f {
		exists := false
		n := c.Name()

		if n == "init.lua" {
			continue
		}

		for _, v := range p {
			if n == string(v.cfg)+".lua" {
				exists = true

				break
			}
		}

		if !exists {
			d = append(d, n)
		}
	}

	s, err := ioutil.ReadDir(structure.Dir.PStart)
	if err != nil {
		panic(err)
	}

	pf := []string{}

	for _, c := range s {
		found := false

		for _, v := range p {
			if string(v.repo.dir()) == c.Name() {
				found = true

				break
			}
		}

		if !found {
			pf = append(pf, c.Name())
		}
	}

	for _, v := range d {
		err := os.Remove(filepath.Join(structure.Dir.PluginCfg, v))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Removed '%s'\n", v)
	}

	for _, v := range pf {
		err := os.RemoveAll(filepath.Join(structure.Dir.PStart, v))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Removed '%s'\n", v)
	}

	updateCfgInit()
}
