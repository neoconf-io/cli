package plugins

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

type plugin struct {
	repo   repo
	dir    dir
	cfg    cfg
	branch string
	opt    bool
}

type (
	cfg  string
	dir  string
	repo string
)

func getPlugins(i []string) []plugin {
	l := []plugin{}

	for _, s := range i {
		l = append(l, parseRepo(s))
	}

	return l
}

func getJSON() []string {
	f, err := ioutil.ReadFile(structure.Files.Plugins.O)
	if err != nil {
		panic(err)
	}

	p := []string{}

	err = json.Unmarshal(f, &p)
	if err != nil {
		panic(err)
	}

	sort.Strings(p)

	return p
}

func (i dir) cfg() cfg {
	return cfg(strings.ReplaceAll(string(i), ".", "-"))
}

func (i repo) dir() dir {
	return dir(strings.Split(string(i), "/")[1])
}

func parseRepo(i string) plugin {
	r, b, o := parsePluginString(i)
	p := plugin{}
	p.repo = r
	p.dir = p.repo.dir()
	p.cfg = p.dir.cfg()
	p.branch = b
	p.opt = o

	return p
}

const minSplit = 1

func parsePluginString(i string) (r repo, b string, o bool) {
	withopt := strings.Split(i, ":")

	if len(withopt) > minSplit {
		o = true
	}

	s := strings.Split(withopt[0], "@")

	if len(s) > minSplit {
		b = s[1]
	}

	r = repo(s[0])

	return r, b, o
}
