package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	flagPrefix string
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute command with prefix (--prefix <prefix>)",
	Long: `

	`,
	Run: func(cmd *cobra.Command, args []string) {
		if flagPrefix == "" {
			log.Println("Prefix not provided (this failure is intentional)")
			_ = cmd.Usage()
			os.Exit(1)
		}

		pfx, err := getPrefix(flagPrefix)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Setenv("WINEPREFIX", pfx.Path)

		if err != nil {
			log.Println("Could not set WINEPREFIX environmental variable.")
			_ = cmd.Usage()
			os.Exit(1)
		}

		if len(args) == 0 {
			log.Println("Please provide a command to run.")
			_ = cmd.Usage()
			os.Exit(1)
		}

		c := exec.Command("sh", append([]string{"-c"}, args...)...)

		runCmd(c)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.PersistentFlags().StringVarP(&flagPrefix, "prefix", "p", "", "Prefix to use")
}
