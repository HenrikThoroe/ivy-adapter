package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/HenrikThoroe/ivy-adapter/internal/app/instl"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	tea "github.com/charmbracelet/bubbletea"
)

type _installFlags struct {
	engine  string
	version string
	config  string
}

var installFlags _installFlags

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install an engine",
	Long:  `Installs an engine which is managed by the 'Engine Version Control' system. Use q or ctrl+c to exit at any time.`,

	Run: func(cmd *cobra.Command, args []string) {
		conf.Load(installFlags.config)

		model, engine := instl.BuildInstallationViewModel(installFlags.engine, installFlags.version)

		if _, err := tea.NewProgram(model).Run(); err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}

		fmt.Println("Path to engine binary: ", engine.Path())
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&installFlags.engine, "engine", "e", "", "The engine to install")
	installCmd.Flags().StringVarP(&installFlags.version, "version", "v", "", "The version of the engine to install")
	installCmd.Flags().StringVarP(&installFlags.config, "config", "c", "", "The path to the configuration file")
}
