package cmd

import (
	"fmt"
	"os"

	"github.com/HenrikThoroe/ivy-adapter/internal/app/test"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type _testFlags struct {
	config string
}

var testFlags _testFlags

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run headless tests and sync with server",
	Long: "The test command connects to the test server based on the configuration file.\n" +
		"The programm will report the system stats to the test server and wait until the server requests a test to run.\n" +
		"A test will download the requested engines and play a batch of games.\n" +
		"The number of games played depends on the number of cores and memory available on the system.\n" +
		"Use q or ctrl+c to exit at any time.",
	Run: func(cmd *cobra.Command, args []string) {
		conf.Load(testFlags.config)
		model := test.BuildTestViewModel()

		if _, err := tea.NewProgram(model).Run(); err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&testFlags.config, "config", "c", "", "The path to the configuration file")
}
