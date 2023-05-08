package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/ui/instl"
	tea "github.com/charmbracelet/bubbletea"
)

type _installFlags struct {
	engine  string
	version string
}

var installFlags _installFlags

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install an engine",
	Long:  `Installs an engine which is managed by the 'Engine Version Control' system. Use q or ctrl+c to exit at any time.`,

	Run: func(cmd *cobra.Command, args []string) {
		model := instl.BuildInstallationViewModel(installFlags.engine, installFlags.version)

		if _, err := tea.NewProgram(model).Run(); err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&installFlags.engine, "engine", "e", "", "The engine to install")
	installCmd.Flags().StringVarP(&installFlags.version, "version", "v", "", "The version of the engine to install")
	conf.Load()
}
