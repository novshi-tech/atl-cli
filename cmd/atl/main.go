package main

import (
	atlcli "github.com/novshi-tech/atl-cli"
	"github.com/novshi-tech/atl-cli/cmd"
)

func main() {
	cmd.SetSkillsFS(atlcli.SkillsFS)
	cmd.Execute()
}
