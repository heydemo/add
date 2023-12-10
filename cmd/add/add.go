/*
Copyright © 2023 John De Mott
*/
package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

import add "heydemo/add/addmain"

var freshInstall bool
var configEnv *add.ConfigEnv

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "add",
	Short: "Command line snippet and script manager",
	Long:  `Add makes it easy to manage your command line snippets and scripts`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Add is a command line snippet and script manager")

		add.Add(args[0], configEnv)
	},
	Args: cobra.ExactArgs(1),
}

var xCmd = &cobra.Command{
	Use:   "x",
	Short: "Execute command",
	Long:  "Executes an add script",
	Run: func(cmd *cobra.Command, args []string) {
		add.Execute(args[0], args[1:], configEnv)
	},
}

var gaCmd = &cobra.Command{
	Use:   "ga",
	Short: "Force alias generation",
	Long:  "Generates aliases for all scripts",
	Run: func(cmd *cobra.Command, args []string) {
		add.GenerateAliases(configEnv)
	},
}

var lCmd = &cobra.Command{
	Use:   "l",
	Short: "List scripts",
	Long:  "List all scripts",
	Run: func(cmd *cobra.Command, args []string) {
		add.Subproc("ls", configEnv.User_bin_dir)
	},
}

var cCmd = &cobra.Command{
	Use:   "c",
	Short: "Cat script",
	Long:  "Output the full contents of a script to the console",
	Run: func(cmd *cobra.Command, args []string) {
		path := add.FindExecutable(args[0], configEnv)
		if path == "" {
			fmt.Println("Script not found")
			os.Exit(1)
		}
		add.Subproc("cat", path)
	},
}

var dCmd = &cobra.Command{
	Use:   "d",
	Short: "Delete script",
	Long:  "Delete the named script",
	Run: func(cmd *cobra.Command, args []string) {
		path := add.FindExecutable(args[0], configEnv)
		if path == "" {
			fmt.Println("Script not found")
			os.Exit(1)
		}

		if cmd.Flag("force").Value.String() == "false" {
			if !add.ConfirmPrompt("Are you sure you want to delete " + args[0] + "?") {
				os.Exit(0)
			}
		}

		add.Subproc("rm", "-f", path)
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(xCmd)
	rootCmd.AddCommand(lCmd)
	rootCmd.AddCommand(cCmd)
	rootCmd.AddCommand(dCmd)
	rootCmd.AddCommand(gaCmd)

	rootCmd.PersistentFlags().BoolP("force", "f", false, "Force delete")

	// Add a version flag
	rootCmd.Version = "0.0.1"
}

func main() {
	// if --debug is in the args, wait for debugger to attach
	fmt.Println(os.Args)
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("Waiting for debugger to attach. Press ENTER to continue...")
		fmt.Printf("Current Process ID: %d\n", os.Getpid())
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	Execute()
}
