package plugins

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abenz1267/neoconf/structure"
)

func Opt(opt bool) {
	p := getPlugins(getJSON())

	filtered := []plugin{}

	for _, v := range p {
		if opt != v.opt {
			filtered = append(filtered, v)
		}
	}

	ListSpecial(filtered)

	for _, v := range getSelections() {
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Printf("Couldn't process '%s'\n", v)

			continue
		}

		if i > len(p) {
			continue
		}

		item := filtered[i-1]

		for i, v := range p {
			if v.repo == item.repo {
				p[i].opt = !p[i].opt
			}
		}

		err = os.Rename(structure.GetPluginDir(string(item.dir), item.opt), structure.GetPluginDir(string(item.dir), !item.opt))
		if err != nil {
			panic(err)
		}
	}

	writeList(p)
}
