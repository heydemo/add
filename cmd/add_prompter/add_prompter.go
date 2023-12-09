/*
Copyright © 2023 John De Mott
*/
package main

import (
	"fmt"
	"os"
	"strings"
    "path/filepath"
	"github.com/spf13/cobra"
    "bufio"

	add "heydemo/add/addmain"
)

var freshInstall bool
var configEnv *add.ConfigEnv
var promptFile string
var promptName string
var promptType string
var promptDescription string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "add_prompter",
	Short: "Collects user input based on promptables",
	Run: func(cmd *cobra.Command, args []string) {

        if promptFile == "" {
            fmt.Println("No prompt file specified")
            os.Exit(1)
        }

        if promptName == "" {
            fmt.Println("No prompt name specified")
            os.Exit(1)
        }

        promptType := strings.TrimSuffix(filepath.Base(promptFile), filepath.Ext(promptFile))

		var promptable add.Promptable = add.Promptable{
			Name:        promptName,
			Description: promptDescription,
			Type:        promptType,
		}

		options := add.LoadPromptableOptions(promptable, promptFile)

		add.PrettyPrint(options)

        Prompt(promptable, options)
		fmt.Println("All done")

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	freshInstall, configEnv = add.Bootstrap()

	if freshInstall {
		return
	}

	rootCmd.PersistentFlags().BoolP("force", "f", false, "Force delete")
	rootCmd.PersistentFlags().StringVarP(&promptFile, "prompt", "p", "", "Prompt file")
	rootCmd.PersistentFlags().StringVarP(&promptName, "name", "n", "", "The name of the promptable to load")
	rootCmd.PersistentFlags().StringVarP(&promptDescription, "description", "d", "", "Prompt description")

	// Add a version flag
	rootCmd.Version = "0.0.1"

}

func main() {
    fmt.Println("Waiting for debugger to attach. Press ENTER to continue...")
    fmt.Printf("Current Process ID: %d\n", os.Getpid())
    bufio.NewReader(os.Stdin).ReadBytes('\n')
	Execute()
}
