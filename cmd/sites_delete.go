package cmd

import (
	"fmt"

	"github.com/novshi-tech/atl-cli/internal/auth"
	"github.com/spf13/cobra"
)

var sitesDeleteCmd = &cobra.Command{
	Use:   "delete <alias>",
	Short: "Delete a configured site",
	Args:  cobra.ExactArgs(1),
	RunE:  runSitesDelete,
}

func init() {
	sitesCmd.AddCommand(sitesDeleteCmd)
}

func runSitesDelete(cmd *cobra.Command, args []string) error {
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

	defaultSite, _ := auth.GetDefaultSite(store)
	wasDefault := defaultSite == alias

	if err := auth.DeleteSite(store, alias); err != nil {
		return err
	}

	fmt.Printf("Site %q deleted.\n", alias)
	if wasDefault {
		fmt.Println("Warning: Deleted site was the default. Run 'atl sites use <alias>' to set a new default.")
	}
	return nil
}
