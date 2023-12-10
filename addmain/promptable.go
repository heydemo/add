// prompter reads promptable definitions and prompts the user
// for values to be used as arguments to the script being executed.
package addmain

import (
	"fmt"
	"github.com/erikgeiser/promptkit/selection"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// An argument to a script which is promptable
type Promptable struct {
	Name        string
	Description string
	Type        string
}

// Represents a selectable option
type PromptableOption struct {
	Type  string
	Label string
	Value string
	Props map[string]string
}

func (p PromptableOption) String() string {
	return p.Label
}

// Given a list of promptables and a list of arguments, prompt the user
// for the values of the promptables which are not already provided
func populatePromptables(promptables []Promptable, args []string, configEnv *ConfigEnv) (fargs, environ []string) {
	var final_args []string = make([]string, len(promptables))
	copy(final_args, args)

	var final_options []PromptableOption = make([]PromptableOption, len(promptables))

	env := os.Environ()

	for index, promptable := range promptables[len(args):] {
		option := promptForPromptable(promptable, configEnv)
		final_args[index] = option.Value
		final_options = append(final_options, option)
		env = append(env, getOptionEnvVars(option)...)
	}

	return final_args, env
}

func getOptionEnvVars(option PromptableOption) []string {
	var env []string

	for key, value := range option.Props {
		env = append(env, "p_"+option.Type+"_"+key+"="+value)
	}

	return env
}

func LoadPromptableOptions(promptable Promptable, filename string) []PromptableOption {
	contents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var values []PromptableOption

	err = yaml.Unmarshal(contents, &values)
	if err != nil {
		panic(err)
	}

	for index := range values {
		values[index].Type = promptable.Type
	}

	return values
}

func promptForPromptable(promptable Promptable, configEnv *ConfigEnv) PromptableOption {
	filename := filepath.Join(configEnv.Promptable_dir, promptable.Name) + ".yml"
	options := LoadPromptableOptions(promptable, filename)

	selectedOption := Prompt(promptable, options)
	return selectedOption

}

func Prompt(promptable Promptable, options []PromptableOption) PromptableOption {

	for _, option := range options {
		fmt.Println("option = " + option.String())
	}

	sp := selection.New("Select a value for the "+promptable.Name+" argument", options)
	sp.PageSize = 8

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the final choice
	PrettyPrint(choice)
	return choice

}
