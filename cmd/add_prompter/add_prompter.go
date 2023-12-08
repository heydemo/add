/*
Copyright © 2023 John De Mott
*/
package main

import (
	"os"
    "fmt"
	"github.com/spf13/cobra"
)

import add "heydemo/add/addmain"

var freshInstall bool
var configEnv *add.ConfigEnv
var promptFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "add_prompter",
	Short: "Collects user input based on promptables",
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("It works")

        var promptable add.Promptable = add.Promptable{
            Name: "test",
            Description: "This is a test",
            Type: "btdevice",
        }

        options := add.LoadPromptableOptions(promptable, promptFile)


        add.PrettyPrint(options)
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

    // Add a version flag
    rootCmd.Version = "0.0.1"

}

func main() {
    Execute()
}


