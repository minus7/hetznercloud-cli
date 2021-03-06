package cli

import "github.com/spf13/cobra"

func newNetworkCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "network",
		Short:                 "Manage networks",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runNetwork),
	}
	cmd.AddCommand(
		newNetworkListCommand(cli),
		newNetworkDescribeCommand(cli),
		newNetworkCreateCommand(cli),
		newNetworkUpdateCommand(cli),
		newNetworkDeleteCommand(cli),
		newNetworkChangeIPRangeCommand(cli),
		newNetworkAddRouteCommand(cli),
		newNetworkRemoveRouteCommand(cli),
		newNetworkAddSubnetCommand(cli),
		newNetworkRemoveSubnetCommand(cli),
		newNetworkAddLabelCommand(cli),
		newNetworkRemoveLabelCommand(cli),
		newNetworkEnableProtectionCommand(cli),
		newNetworkDisableProtectionCommand(cli),
	)
	return cmd
}

func runNetwork(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}
