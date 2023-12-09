// prompter reads promptable definitions and prompts the user
// for values to be used as arguments to the script being executed.
package addmain

import (
    "os"
    "path/filepath"
    "gopkg.in/yaml.v2"
    "fmt"
)


// An argument to a script which is promptable
type Promptable struct {
    Name string
    Description string
    Type string
}

// Represents a selectable option
type PromptableOption struct {
    Type string
    Label string
    Value string
    Props map[string]string
}

func (p PromptableOption) String() string {
    return p.Label
}

// #
// # Promptable creation
// #

func populatePromptables(promptables []Promptable, args []string, configEnv *ConfigEnv) []string {
    var final_args []string = make([]string, len(promptables))
    copy(final_args, args)

    for _, promptable := range promptables[len(args):] {
        final_args = append(final_args, promptForPromptable(promptable, configEnv))
    }
    return final_args
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

    return values
}

func promptForPromptable(promptable Promptable, configEnv *ConfigEnv) string {
    filename := filepath.Join(configEnv.Promptable_dir, promptable.Name) + ".yml"
    values := LoadPromptableOptions(promptable, filename)

    for _, option := range values {
        fmt.Println("Enter value for " + promptable.Name)
        fmt.Println("Description: " + promptable.Description)
        fmt.Print(option.Label + ": ")
    }
    return ""

}
