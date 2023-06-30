package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "wpm",
	Short: "wine prefix management cli",
	// Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/wpm.json)")
}

func initConfig() {
	configHome, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	configDir := configHome + "/wpm/"

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		viper.AddConfigPath(configDir)
		viper.SetConfigType("json")
		viper.SetConfigName("wpm")

		err = os.MkdirAll(configDir, 0750)
		if err != nil {
			log.Fatal(err)
		}

		// Create new config
		if _, err := os.Stat(configDir + "wpm.json"); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				viper.Set("prefixes", make([]Prefix, 0))

				err := viper.SafeWriteConfig()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err.(viper.ConfigFileNotFoundError), viper.ConfigFileNotFoundError{}) {
			_, err = os.Create(configDir + "wpm.json")

			if err != nil {
				log.Println(err)
			}
		}
	}
}
