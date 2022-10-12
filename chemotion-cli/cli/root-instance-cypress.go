package cli

import (
	"github.com/spf13/cobra"
)

var cypressCmdTable = make(cmdTable)

var cypressCmd = &cobra.Command{
	Use:     "cypress",
	Aliases: []string{"cy", "cypress"},
	Args:    cobra.NoArgs,
	Short:   "Perform cypress tests to validate UI components (i,e. React and native Html components) and also test their functionalities" + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive(true)
		var acceptedOpts []string
		if elementInSlice(instanceStatus(currentInstance), &[]string{"Exited", "Created"}) == -1 { // checks if the instance is running
			acceptedOpts = []string{"setup", "pull changes", "start"}
			cypressCmdTable["setup"] = setupCypressInstanceRootCmd.Run
			cypressCmdTable["pull changes"] = pullChangesCypressInstanceRootCmd.Run
			cypressCmdTable["start"] = startCypressInstanceRootCmd.Run
		} else {
			acceptedOpts = []string{"logs"}
			cypressCmdTable["logs"] = logInstanceRootCmd.Run
		}
		cypressCmdTable[selectOpt(acceptedOpts, "")](cmd, args)
	},
}

func init() {
	instanceRootCmd.AddCommand(cypressCmd)
}
