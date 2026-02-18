package cmd

import (
	"fmt"

	"github.com/novshi-tech/atl-cli/internal/auth"
	"github.com/spf13/cobra"
)

var sitesUseCmd = &cobra.Command{
	Use:   "use <alias>",
	Short: "Set a site as the default",
	Args:  cobra.ExactArgs(1),
	RunE:  runSitesUse,
}

func init() {
	sitesCmd.AddCommand(sitesUseCmd)
}

func runSitesUse(cmd *cobra.Command, args []string) error {
	alias := args[0]

	store, err := auth.NewStore()
	if err != nil {
		return err
	}

	sites, err := auth.ListSites(store)
	if err != nil {
		return err
	}

	found := false
	for _, s := range sites {
		if s == alias {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("site %q not found", alias)
	}

	if err := auth.SetDefaultSite(store, alias); err != nil {
		return err
	}

	fmt.Printf("Site %q set as default.\n", alias)
	return nil
}
