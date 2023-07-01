package cmd

import (
	"log"
	"os"
	"os/exec"
	"regexp"

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
	_ = execCmd.RegisterFlagCompletionFunc("prefix", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		var ret []string

		compRx, err := regexp.Compile("^" + toComplete)
		if err != nil {
			return nil, 0
		}

		pfx := getPrefixes()
		for _, v := range pfx {
			if compRx.MatchString(v.Name) {
				ret = append(ret, v.Name)
			}
		}

		return ret, 0
	})

}
