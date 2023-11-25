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
        os.WriteFile(path, []byte("#!/bin/bash\n"), 0644)
    }

    Subproc(editor, path)
    // make path executable
    os.Chmod(path, 0755)
    fmt.Println("Wrote file to", path)

    GenerateAliases(configEnv)

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
