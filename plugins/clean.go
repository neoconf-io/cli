package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/abenz1267/neoconf/structure"
)

// Clean removes plugin configs for missing plugins.
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

	for _, v := range d {
		err := os.Remove(filepath.Join(structure.Dir.PluginCfg, v))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Removed '%s'\n", v)
	}

	updateCfgInit()
}
