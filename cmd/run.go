package cmd

import (
	"fmt"
	"os"

	"github.com/HenrikThoroe/ivy-adapter/internal/app/run"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type _runFlags struct {
	game    string
	color   string
	engine  string
	version string
}

var runFlags _runFlags

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the engine on the given game ID",
	Long: "Launches the engine given by name and version.\n" +
		"Will fail if the engine has not been previously installed.\n" +
		"The game is played by connecting against the game server which has to be available\n" +
		"on the configured address in the config file.\n" +
		"Use q or ctrl+c to exit at any time.",

	Run: func(cmd *cobra.Command, args []string) {
		version, vErr := mgmt.ParseVersion(runFlags.version, mgmt.DotVersionStyle)

		if vErr != nil {
			fmt.Println("Error parsing version: ", vErr)
			os.Exit(1)
		}

		eng, eErr := mgmt.BestMatch(runFlags.engine, *version)

		if eErr != nil {
			fmt.Println("Error parsing engine: ", eErr)
			os.Exit(1)
		}

		vm := run.BuildRunViewModel(runFlags.game, runFlags.color, eng)

		if _, err := tea.NewProgram(vm).Run(); err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runFlags.game, "game", "g", "", "Game ID")
	runCmd.Flags().StringVarP(&runFlags.color, "color", "c", "", "Color the selected engine should play as")
	runCmd.Flags().StringVarP(&runFlags.engine, "engine", "e", "", "Engine Name (must be installed)")
	runCmd.Flags().StringVarP(&runFlags.version, "version", "v", "", "Version of Engine (must be installed)")

	runCmd.MarkFlagRequired("game")
	runCmd.MarkFlagRequired("color")
	runCmd.MarkFlagRequired("engine")
	runCmd.MarkFlagRequired("version")

	conf.Load()
}
