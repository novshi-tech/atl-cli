package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bbMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show the currently authenticated Bitbucket user",
	RunE:  runBBMe,
}

func init() {
	bitbucketCmd.AddCommand(bbMeCmd)
}

func runBBMe(cmd *cobra.Command, args []string) error {
	client, err := newBitbucketClient(cmd)
	if err != nil {
		return err
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		return err
	}

	if jsonMode(cmd) {
		return printJSON(JSONBBUserItem{
			AccountID:   user.AccountID,
			DisplayName: user.DisplayName,
			Nickname:    user.Nickname,
			UUID:        user.UUID,
			CreatedOn:   user.CreatedOn,
		})
	}

	fmt.Printf("Account ID:    %s\n", user.AccountID)
	fmt.Printf("Display Name:  %s\n", user.DisplayName)
	fmt.Printf("Nickname:      %s\n", user.Nickname)
	fmt.Printf("UUID:          %s\n", user.UUID)
	fmt.Printf("Created On:    %s\n", user.CreatedOn)
	return nil
}
