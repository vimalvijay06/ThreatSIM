package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available attack plugins",
		Run: func(cmd *cobra.Command, args []string) {
			pluginList := registry.List()

			header := color.New(color.FgCyan, color.Bold)
			header.Printf("\n  Available Attack Plugins (%d)\n", len(pluginList))
			fmt.Println("  " + "─────────────────────────────────────────────")

			for _, p := range pluginList {
				name := color.New(color.FgWhite, color.Bold).Sprintf("%-30s", p.Name())
				id := color.New(color.FgYellow).Sprintf("[%s]", p.ID())

				fmt.Printf("  %s %s\n", name, id)

				desc := color.New(color.Faint)
				desc.Printf("    %s\n\n", p.Description())
			}

			fmt.Println("  " + "─────────────────────────────────────────────")
			color.New(color.Faint).Println("  Run: threatsim simulate <plugin-id>")
			fmt.Println()
		},
	}
}
