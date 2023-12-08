package addmain

import (
	"fmt"
	"os"
	"path/filepath"
    "strings"
)

func Add(filename string, configEnv *ConfigEnv) {
    path := filepath.Join(configEnv.User_bin_dir, filename)
    editor := GetEditor()

    // If file does not exist, create it
    if _, err := os.Stat(path); os.IsNotExist(err) {
        metadata := PromptForMetadata()
        os.WriteFile(path, []byte(formatMetadata(metadata)), 0644)
        PrettyPrint(metadata)
    }

    Subproc(editor, path)
    // make path executable
    os.Chmod(path, 0755)
    //fmt.Println("Wrote file to", path)

    //GenerateAliases(configEnv)
    PrettyPrint(editor)

}

func FindExecutable(name string, configEnv *ConfigEnv) string {
    dirs := []string{configEnv.User_bin_dir, configEnv.Core_bin_dir, configEnv.Public_bin_dir}

    for _, dir := range dirs {
        path := filepath.Join(dir, name)
        if _, err := os.Stat(path); err == nil {
            return path
        }
    }
    return ""
}

func ConfirmPrompt(prompt string) bool {
    fmt.Println(prompt + " [y/n]")
    var response string
    fmt.Scanln(&response)
    return strings.ToLower(response) == "y"
}



func Execute(script string, args []string, configEnv *ConfigEnv) {
    fmt.Println("Executing", script)
    path := FindExecutable(script, configEnv)
    if path == "" {
        fmt.Println("Script not found")
        os.Exit(1)
    }

    metadata := extractMetadata(path)
    final_args := populatePromptables(metadata.Promptables, args, configEnv)
    PrettyPrint(metadata)
    PrettyPrint(final_args)

    Subproc(path, args...)
}
