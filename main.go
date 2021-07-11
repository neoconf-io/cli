package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/abenz1267/neoconf/plugins"
	"github.com/abenz1267/neoconf/structure"
)

//nolint
//go:embed files/**
var f embed.FS

const minLength = 2

type cmd struct {
	long        string
	short       string
	example     string
	description string
	do          func()
}

func (c cmd) check(arg string) bool {
	if c.long == arg || c.short == arg {
		c.do()

		return true
	}

	return false
}

func getCmds() []cmd {
	init := cmd{
		long:        "init",
		short:       "init",
		example:     "",
		description: "Creates file structure and installs plugins from 'nvim/plugins.json'.",
		do: func() {
			structure.CheckFolders()
			structure.CheckFiles()
			plugins.Install([]string{})
		},
	}

	install := cmd{
		long:        "install",
		short:       "i",
		example:     "neoconf install neovim/nvim-lspconfig nvim-telescope/telescope.nvim",
		description: "Installs all plugins provided.",
		do: func() {
			plugins.Install(os.Args[2:])
		},
	}

	update := cmd{
		long:        "update",
		short:       "u",
		example:     "",
		description: "Updates all installed plugins.",
		do: func() {
			plugins.Update()
		},
	}

	remove := cmd{
		long:        "remove",
		short:       "r",
		example:     "neoconf remove <enter number(s) from list>",
		description: "Shows a list of all installed plugins. Prompts for number of plugin to delete. Space-separated for multiple deletions.",
		do: func() {
			plugins.RemoveN()
		},
	}

	clean := cmd{
		long:        "clean",
		short:       "c",
		example:     "",
		description: "Removes config files for missing plugins.",
		do: func() {
			plugins.Clean()
		},
	}

	list := cmd{
		long:        "list",
		short:       "l",
		example:     "",
		description: "Lists all installed plugins",
		do: func() {
			plugins.List()
		},
	}

	help := cmd{
		long:        "help",
		short:       "h",
		example:     "",
		description: "",
		do: func() {
			for _, v := range []cmd{list, clean, remove, update, install, init} {
				fmt.Printf("    --%s, %s:\n", v.long, v.short)
				fmt.Printf("      %s \n\n", v.description)
				if v.example != "" {
					fmt.Printf("      Example: %s \n\n", v.example)
				}
			}
		},
	}

	return []cmd{list, clean, remove, update, install, init, help}
}

func main() {
	checkGit()
	checkYarn()

	structure.SetFilesystem(f)
	structure.SetFolders("", "")
	structure.SetFiles()

	if len(os.Args) < minLength {
		// default action
		return
	}

	cmds := getCmds()

	job := false
	for _, v := range cmds {
		job = v.check(os.Args[1])

		if job {
			break
		}
	}

	if !job {
		fmt.Println("unknown command")
	}
}

const status = 1

func checkGit() {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Missing 'git'. Needed to clone and update plugins.")
		os.Exit(status)
	}
}

func checkYarn() {
	if _, err := exec.LookPath("yarn"); err != nil {
		fmt.Println("Missing 'yarn'. Needed for some post-install commands.")
		os.Exit(status)
	}
}
