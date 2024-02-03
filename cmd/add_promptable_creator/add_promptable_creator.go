package main

import (
	//tea "github.com/charmbracelet/bubbletea"
	// "golang.org/x/term"
	//"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"heydemo/add/addmain"
	add "heydemo/add/addmain"
)

const (
	ADD_NEW string = "ADD_NEW"
)

var freshInstall bool

func printTemplate() {
	content, err := add.Templates.ReadFile("templates/promptable-options.yml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", content)
}

func main() {
	var script string
	var number int

	add.Debug()

	freshInstall, configEnv := add.Bootstrap()
	if freshInstall {
		return
	}

	var rootCmd = &cobra.Command{
		Use:   "add_promptable_creator",
		Short: "Create promptables for a script",
		Run: func(cmd *cobra.Command, args []string) {
			count := 0

			metadata := add.ReadMetadata(script)
			promptableCount := max(number, len(metadata.Promptables))

			promptables := make([]add.Promptable, promptableCount)
			d := add.FormatMetadata(metadata)
			fmt.Printf("%s", d)
			for index, promptable := range metadata.Promptables {
				promptables[index] = promptable
			}

			for number > count {
				count++
				currentType := promptables[count-1].Type
				promptableType := selectPromptableType(count, currentType, configEnv)
				if promptableType == "" {
					break
				}
				promptables[count-1].Type = promptableType
				editPromptableForm(&promptables[count-1], count, configEnv)
			}

			for index, p := range promptables {

				if p.Type == "" {
					promptables = promptables[:index]
					break
				}

				filename := configEnv.Promptable_dir + "/" + p.Type + ".yml"
				fmt.Printf("Promptable $%d: %s\n", index+1, p.Type)
				fmt.Printf("You can edit promptable options here %s\n", filename)
			}
			metadata.Promptables = promptables
			add.UpdateMetadata(metadata, script)

		},
	}

	rootCmd.Flags().StringVarP(&script, "script", "s", "", "Name of the script for which to create promptables")
	rootCmd.MarkFlagRequired("script")

	rootCmd.Flags().IntVarP(&number, "number", "n", 1, "Number of promptables to create")
	rootCmd.MarkFlagRequired("number")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func editPromptableForm(promptable *add.Promptable, argNumber int, configEnv *add.ConfigEnv) {
	var (
		promptableName        string = promptable.Name
		promptableDescription string = promptable.Description
	)

	if promptableName == "" {
		promptableName = promptable.Type
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Promptable Name for $%d", argNumber)).
				Description(fmt.Sprintf("The name of the promptable, as used as the $%d argument to this script.", argNumber)).
				Value(&promptableName),
			huh.NewText().
				Title("Description").
				Description("The description of the promptable").
				Value(&promptableDescription),
		),
	)

	form.Run()

	promptable.Name = promptableName
	promptable.Description = promptableDescription

}

func selectPromptableType(count int, currentType string, configEnv *add.ConfigEnv) string {
	var (
		createNew bool = true
	)
	types := add.GetPromptableTypes(configEnv)
	var options []huh.Option[string]
	for _, t := range types {
		options = append(options, huh.Option[string]{Key: t, Value: t})
	}

	options = append(options, huh.Option[string]{Key: "Add New", Value: ADD_NEW})

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Create promptable for argument $"+fmt.Sprint(count)+"?").
				Value(&createNew)).WithHideFunc(func() bool {
			return currentType != ""
		}),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Promptable Type for argument $"+fmt.Sprint(count)).
				Description("The class of objects to which this argument belongs (i.e. server, directory, project, device). Types can be reused across scripts.").
				Options(options...).
				Value(&currentType),
		).WithHideFunc(func() bool {
			return !createNew
		}),
	)

	form.Run()

	if currentType == ADD_NEW {
		currentType = addPromptableType(configEnv)
	}

	return currentType

}

func addPromptableType(configEnv *add.ConfigEnv) string {
	var (
		promptableType string
		description    string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Add New Promptable Type").
				Description("The machine name of the promptable type (i.e. server, directory, project, device). Types can be reused across scripts.").
				Value(&promptableType).
				Validate(func(value string) error {
					switch {
					case value == "":
						return fmt.Errorf("Promptable type cannot be blank")
					case strings.Contains(value, " "):
						return fmt.Errorf("Promptable type cannot contain spaces")
					case strings.HasPrefix(value, "."):
						return fmt.Errorf("Promptable type cannot start with a period")
					default:
						return nil
					}
				}),

			huh.NewText().
				Title("Description").
				Description("A description of the promptable type, in general, not as used in this script specifically").
				Value(&description),
		),
	)

	form.Run()

	promptableType = strings.ToLower(promptableType)
	addmain.CreatePromptableType(promptableType, configEnv)

	return strings.ToLower(promptableType)
}
