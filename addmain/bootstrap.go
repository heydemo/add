package addmain

import (
	"fmt"
	"os"
    "os/user"
	"path/filepath"
	"strings"
)

func getIncludeShContent(configEnv *ConfigEnv) string {
    var content string = ""

    //	content += fmt.Sprintf("export PATH=$PATH:%s:%s\n",
    //                            configEnv.User_bin_dir,
    //                            configEnv.public_bin_dir)

	content += "export ADD_INSTALLED=1\n"

    content += "source " + filepath.Join(configEnv.State_dir, "aliases.sh") + "\n"

    return content

}

func getAliasIncludeContent(configEnv *ConfigEnv) string {
    var content string

    files := getExecutables(configEnv)

    for _, file := range files {
        content += fmt.Sprintf("alias %s='add x %s'\n", file, file)
    }

    return content

}

func getExecutables(configEnv *ConfigEnv) []string {
    dirs := []string{configEnv.User_bin_dir, configEnv.Core_bin_dir, configEnv.Public_bin_dir}
    var executables []string
    for _, dir := range dirs {
        files, err := os.ReadDir(dir)
        if err != nil {
            panic(err)
        }
        for _, file := range files {
            if !file.IsDir() {
                executables = append(executables, file.Name())
            }
        }
    }
    return executables
}

func writeContent(filePath string, content string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	file.WriteString(content)
}


func determineProfileFile() string {
    usr, err := user.Current()
    if err != nil {
        panic(err)
    }
	candidates := []string{".bashrc", ".bash_profile", ".profile", ".zshrc"}
	for _, candidate := range candidates {
		expandedPath := usr.HomeDir + "/" + candidate
		if _, err := os.Stat(expandedPath); err == nil {
			return expandedPath
		}
	}
    panic("Could not find a profile file to install to.")
}

func confirmInstallPrompt(profileFile string) bool {
	fmt.Printf("Welcome to ADD!\n\nTo work properly, you will need to source ADDs include.sh file\n\nShall we add to your profile file (%s)? [y/n] ", profileFile)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

func GenerateAliases(configEnv *ConfigEnv) {
    filePath := filepath.Join(configEnv.State_dir, "aliases.sh")
    writeContent(filePath, getAliasIncludeContent(configEnv))
}

func Bootstrap() (bool, *ConfigEnv) {
	configEnv := GetConfigEnv()
    freshInstall := false

	isInstalled := os.Getenv("ADD_INSTALLED") == "1"
	includeShPath := filepath.Join(configEnv.State_dir, "include.sh")
	aliasesPath := filepath.Join(configEnv.State_dir, "aliases.sh")

	if _, err := os.Stat(includeShPath); os.IsNotExist(err) {
	    filePath := filepath.Join(configEnv.State_dir, "include.sh")
		writeContent(filePath, getIncludeShContent(configEnv))
	}

    if _, err := os.Stat(aliasesPath); os.IsNotExist(err) {
        GenerateAliases(configEnv)
    }

	if !isInstalled {
        freshInstall = true
		profileFile := determineProfileFile()
		if profileFile == "" {
			fmt.Println("Could not find a profile file to install to")
			os.Exit(1)
		}
		if confirmInstallPrompt(profileFile) {
			file, err := os.OpenFile(profileFile, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error:", err)
				return freshInstall, configEnv
			}
			defer file.Close()
			file.WriteString(fmt.Sprintf("\nsource %s/include.sh #ADD INSTALL\n", configEnv.State_dir))
			fmt.Printf("Added to %s\n", profileFile)
			fmt.Printf("Please restart your shell or run `source %s` to use ADD\n", profileFile)
		} else {
			fmt.Println("Aborting install")
		}
	}

	return freshInstall, configEnv
}
