package execution

import "github.com/spf13/cobra"

// fakeCobraCommand returns a minimal *cobra.Command with given use string.
// We only need the Use field and flag handling for Parse tests.
func fakeCobraCommand(use string) *cobra.Command {
	cmd := &cobra.Command{Use: use}
	return cmd
}
