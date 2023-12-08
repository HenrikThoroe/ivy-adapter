package cmd

import (
	"fmt"
	"os"

	"github.com/HenrikThoroe/ivy-adapter/internal/app/run"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	"github.com/spf13/cobra"
)

type _runFlags struct {
	player  string
	exe     string
	engine  string
	version string
}

var runFlags _runFlags

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the engine on the given game ID",
	Long: "Launches the engine given by name and version or path.\n" +
		"A connection against the game host server will be created using the given player ID.\n" +
		"Once the game server sends a move request, the engine will be invoked to calculate a move.\n" +
		"The output from the engine and all prompts to the engine will be printed to the console.\n",

	Run: func(cmd *cobra.Command, args []string) {
		conf.Load()

		ifc, stdin, stdout, err := run.SetupEngineIfc(runFlags.exe, runFlags.engine, runFlags.version)

		if err != nil {
			fmt.Println("Error setting up engine interface: ", err)
			os.Exit(1)
		}

		pipe := func(c chan string, prefix string) {
			for s := range c {
				fmt.Println(prefix + s)
			}
		}

		go pipe(stdin, "> ")
		go pipe(stdout, "")

		run.Play(ifc, runFlags.player)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runFlags.player, "player", "p", "", "Player ID")
	runCmd.Flags().StringVarP(&runFlags.exe, "binary", "b", "", "Path to engine executable")
	runCmd.Flags().StringVarP(&runFlags.engine, "engine", "e", "", "Engine Name (must be installed)")
	runCmd.Flags().StringVarP(&runFlags.version, "version", "v", "", "Version of Engine (must be installed)")

	runCmd.MarkFlagRequired("player")

	conf.Load()
}
