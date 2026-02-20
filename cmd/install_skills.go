package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var skillsFS fs.FS

// SetSkillsFS sets the embedded filesystem containing skill definitions.
func SetSkillsFS(fsys fs.FS) {
	skillsFS = fsys
}

var installSkillsCmd = &cobra.Command{
	Use:   "install-skills [skill...]",
	Short: "Install embedded skills to the local skills directory",
	Long: `Install skill definitions bundled in the atl binary.

Without arguments, all available skills are installed.
Specify skill names to install only selected skills.`,
	RunE: runInstallSkills,
}

func init() {
	installSkillsCmd.Flags().String("path", "", "Target directory for skill installation")
	installSkillsCmd.Flags().Bool("dry-run", false, "Preview without writing files")
	installSkillsCmd.Flags().Bool("list", false, "List available skills")
	rootCmd.AddCommand(installSkillsCmd)
}

func runInstallSkills(cmd *cobra.Command, args []string) error {
	list, _ := cmd.Flags().GetBool("list")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	pathFlag, _ := cmd.Flags().GetString("path")

	available, err := listSkills()
	if err != nil {
		return fmt.Errorf("reading embedded skills: %w", err)
	}

	if list {
		for _, name := range available {
			fmt.Println(name)
		}
		return nil
	}

	destDir, err := resolveDestDir(pathFlag)
	if err != nil {
		return err
	}

	targets := available
	if len(args) > 0 {
		avail := make(map[string]bool, len(available))
		for _, s := range available {
			avail[s] = true
		}
		for _, a := range args {
			if !avail[a] {
				return fmt.Errorf("unknown skill %q; run 'atl install-skills --list' to see available skills", a)
			}
		}
		targets = args
	}

	for _, skill := range targets {
		if err := installSkill(skill, destDir, dryRun); err != nil {
			return fmt.Errorf("installing skill %q: %w", skill, err)
		}
	}

	if dryRun {
		fmt.Println("\n(dry-run: no files were written)")
	} else {
		fmt.Printf("\nInstalled %d skill(s) to %s\n", len(targets), destDir)
	}
	return nil
}

func listSkills() ([]string, error) {
	entries, err := fs.ReadDir(skillsFS, "skills")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func resolveDestDir(pathFlag string) (string, error) {
	if pathFlag != "" {
		return pathFlag, nil
	}
	return defaultSkillsPath()
}

func defaultSkillsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	for _, dir := range []string{".claude", ".codex", ".cursor"} {
		p := filepath.Join(home, dir)
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return filepath.Join(p, "skills"), nil
		}
	}
	return filepath.Join(home, ".claude", "skills"), nil
}

func installSkill(skill, destDir string, dryRun bool) error {
	srcRoot := "skills/" + skill
	return fs.WalkDir(skillsFS, srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// path is "skills/<skill>/..." — strip the leading "skills/" to get the relative path
		// so files end up at destDir/<skill>/...
		rel := strings.TrimPrefix(path, "skills/")
		dest := filepath.Join(destDir, filepath.FromSlash(rel))

		if d.IsDir() {
			if dryRun {
				fmt.Printf("  mkdir %s\n", dest)
				return nil
			}
			return os.MkdirAll(dest, 0755)
		}

		data, err := fs.ReadFile(skillsFS, path)
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Printf("  write %s (%d bytes)\n", dest, len(data))
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}

		perm := os.FileMode(0644)
		if strings.HasSuffix(path, ".sh") {
			perm = 0755
		}

		fmt.Printf("  %s\n", dest)
		return os.WriteFile(dest, data, perm)
	})
}
