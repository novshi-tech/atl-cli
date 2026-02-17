package main

import "github.com/novshi-tech/atl-cli/cmd"

func main() {
	cmd.SetSkillsFS(SkillsFS)
	cmd.Execute()
}
