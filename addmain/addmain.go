package addmain

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Add a new script to our local script library
func Add(filename string, configEnv *ConfigEnv) {
	path := filepath.Join(configEnv.User_bin_dir, filename)
	editor := GetEditor()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		metadata := PromptForMetadata()
		content := "#!/bin/bash\n"
		content += FormatMetadata(metadata)
		os.WriteFile(path, []byte(content), 0644)
	}

	Subproc(editor, path)
	os.Chmod(path, 0755)

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

func Execute(script string, args []string, configEnv *ConfigEnv) {
	path := FindExecutable(script, configEnv)
	if path == "" {
		fmt.Println("Script not found")
		os.Exit(1)
	}

	metadata := ReadMetadata(path)
	final_args, env := populatePromptables(metadata.Promptables, args, configEnv)

	SubprocWithEnv(path, env, final_args...)
}
