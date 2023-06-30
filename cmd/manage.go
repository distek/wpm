package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/distek/menu"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagAltScreen bool
	actions       = []string{"rename", "remove", "change path", "cancel"}

	flagPfxName string
	flagPfxPath string

	flagShowPath bool
)

type Prefix struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

func getPrefixes() []Prefix {
	var prefixes []Prefix

	err := viper.UnmarshalKey("prefixes", &prefixes)
	if err != nil {
		log.Fatal(err)
	}

	return prefixes
}

func getPrefix(name string) (Prefix, error) {
	for _, v := range getPrefixes() {
		if v.Name == name {
			return v, nil
		}
	}

	return Prefix{}, fmt.Errorf("could not find prefix with name: %s", name)
}

func addPrefix(name, path string) error {
	pfx := Prefix{
		Name: name,
		Path: path,
		UUID: uuid.New().String(),
	}

	prefixes := getPrefixes()
	viper.Set("prefixes", append(prefixes, pfx))

	return viper.WriteConfig()
}

func removePrefix(name string) error {
	var update []Prefix

	prefixes := getPrefixes()

	for _, v := range prefixes {
		if v.Name == name {
			continue
		}

		update = append(update, v)
	}

	viper.Set("prefixes", update)

	return viper.WriteConfig()
}

func renamePrefix(oldName, newName string) error {
	var update []Prefix

	prefixes := getPrefixes()

	for _, v := range prefixes {
		if v.Name == oldName {
			v.Name = newName
		}

		update = append(update, v)
	}

	viper.Set("prefixes", update)

	return viper.WriteConfig()
}

func changePathPrefix(newPath string, pfxName string) error {
	var update []Prefix

	prefixes := getPrefixes()

	for _, v := range prefixes {
		if v.Name == pfxName {
			v.Path = newPath
		}

		update = append(update, v)
	}

	viper.Set("prefixes", update)

	return viper.WriteConfig()
}

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "TUI-ish management of prefixes",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			prefixes := getPrefixes()

			var nameSlice []string

			for _, v := range prefixes {
				nameSlice = append(nameSlice, v.Name)
			}

			pfxNameModel := newMenu(nameSlice, "Select prefix", "")
			if pfxNameModel.Interrupt {
				return
			}

			pfxName := pfxNameModel.Selected

			actionModel := newMenu(actions, "Select action", "")
			if actionModel.Interrupt {
				return
			}

			switch actionModel.Selected {
			case "rename":
				input := menu.NewInput("Prefix name:", "", 1024, 50)
				ret, err := menu.Run(input, flagAltScreen)
				if err != nil {
					log.Fatal(err)
				}

				newName := ret.(menu.InputModel).Input.Value()

				err = renamePrefix(pfxName, newName)
				if err != nil {
					log.Fatal(err)
				}
			case "change path":
				pfx, _ := getPrefix(pfxName)

				input := menu.NewInput("Enter desired new path", fmt.Sprintf("Old path: %s", pfx.Path), 1024, 50)

				ret, err := menu.Run(input, flagAltScreen)
				if err != nil {
					log.Fatal(err)
				}

				newPath := ret.(menu.InputModel).Input.Value()

				err = changePathPrefix(newPath, pfxName)
				if err != nil {
					log.Fatal(err)
				}
			case "remove":
				yesNoModel := newMenu([]string{"no", "yes"}, fmt.Sprintf("Are you sure you want to delete: %s?", pfxName), "")

				yesNo := yesNoModel.Selected

				if yesNo == "no" {
					continue
				}

				err := removePrefix(pfxName)
				if err != nil {
					log.Fatal(err)
				}
			case "cancel":
				continue
			}
		}
	},
}

func newMenu(actions []string, title, message string) menu.SingleModel {
	actionMenu := menu.NewSingle(actions, "Select action", "")

	amPost, err := menu.Run(actionMenu, flagAltScreen)
	if err != nil {
		log.Fatal(err)
	}

	return amPost.(menu.SingleModel)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "add --name <name> --path <path>",
	Run: func(cmd *cobra.Command, args []string) {
		if flagPfxName == "" || flagPfxPath == "" {
			log.Println("Provide both --name and --path arguments")
			_ = cmd.Usage()
			os.Exit(1)
		}
		err := addPrefix(flagPfxName, flagPfxPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r"},
	Short:   "remove --name <name>",
	Run: func(cmd *cobra.Command, args []string) {
		if flagPfxName == "" {
			log.Println("Provide --name argument")
			_ = cmd.Usage()
			os.Exit(1)
		}
		err := removePrefix(flagPfxName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "list all prefixes (-p to show paths)",
	Run: func(cmd *cobra.Command, args []string) {
		prefixes := getPrefixes()
		for _, v := range prefixes {
			fmt.Print(v.Name)
			if flagShowPath {
				fmt.Printf(",%s", v.Path)
			}
			fmt.Print("\n")
		}

	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
	rootCmd.AddCommand(listCmd)
	manageCmd.AddCommand(addCmd)
	manageCmd.AddCommand(removeCmd)

	manageCmd.Flags().BoolVarP(&flagAltScreen, "alt-screen", "a", false, "Use alt screen for interactive menu")

	addCmd.Flags().StringVarP(&flagPfxName, "name", "n", "", "Name of prefix to add")
	addCmd.Flags().StringVarP(&flagPfxPath, "path", "p", "", "Path of prefix to add")

	listCmd.Flags().BoolVarP(&flagShowPath, "path", "p", false, "Show path when listing prefixes")
}
