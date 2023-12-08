package addmain

import (
    "fmt"
	"os"
	"path/filepath"
    "os/exec"
    "github.com/spf13/viper"
)

type ConfigEnv struct {
	Config_dir     string
	User_bin_dir   string
	Public_bin_dir string
    Core_bin_dir   string
	State_dir      string
    Promptable_dir string
}

func GetConfigEnv() *ConfigEnv {
    initConfig()
    var config_env ConfigEnv

	config_env.State_dir = viper.GetString("state_dir")
	config_env.User_bin_dir = viper.GetString("user_bin_dir")
	config_env.Public_bin_dir = viper.GetString("public_bin_dir")
	config_env.Core_bin_dir = viper.GetString("core_bin_dir")
    config_env.Config_dir = viper.GetString("config_dir")
    config_env.Promptable_dir = viper.GetString("promptable_dir")

	ensureDirExists(config_env.Config_dir)
	ensureDirExists(config_env.State_dir)
	ensureDirExists(config_env.User_bin_dir)
	ensureDirExists(config_env.Public_bin_dir)
	ensureDirExists(config_env.Core_bin_dir)

	return &config_env
}

func GetEditor() string {
    return getFirstDefinedEnvVars(
        []string{"VISUAL", "EDITOR"},
        getFirstInstalledExecutable([]string{"nvim", "vim", "vi", "emacs", "nano" }))
}

func initConfig() {
    viper.SetConfigName("config")
    //viper.AddConfigPath(getConfigDir())
    viper.AddConfigPath("$HOME/.config/add")

    // Find config path
    if v := os.Getenv("ADD_CONFIG_DIR"); v != "" {
        viper.AddConfigPath(v)
    } else if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
        viper.AddConfigPath(filepath.Join(v, "add"))
    } else {
        viper.AddConfigPath("$HOME/.config/add")
    }

    viper.SetDefault("editor", GetEditor())
    viper.SetDefault("state_dir", getStateDirDefault())
    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    configDir := filepath.Dir(viper.ConfigFileUsed())
    viper.Set("config_dir", configDir)
    viper.SetDefault("user_bin_dir", filepath.Join(configDir, "bin", "user"))
    viper.SetDefault("public_bin_dir", filepath.Join(configDir, "bin", "user"))
    viper.SetDefault("core_bin_dir", filepath.Join(configDir, "bin", "user"))
    viper.SetDefault("promptable_dir", filepath.Join(configDir, "promptables"))

    if err != nil {
        panic(err)
    }

}

func getFirstInstalledExecutable(exes []string) string {
  // Check if nvim is installed
  for _, exe := range exes {
    if _, err := exec.LookPath(exe); err == nil {
      return exe
    }
  }
  return ""

}

func ensureDirExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func getUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func getStateDirDefault() string {
    if v:= os.Getenv("ADD_STATE_DIR"); v != "" {
        return v
    } else if v := os.Getenv("XDG_STATE_HOME"); v != "" {
        return filepath.Join(v, "add")
    } else {
        return filepath.Join(getUserHomeDir(), ".local", "state", "add")
    }
}

func getFirstDefinedEnvVars(env_vars []string, defaultVal string) string {
	for _, env_var := range env_vars {
		if v := os.Getenv(env_var); v != "" {
			return v
		}
	}
	return defaultVal
}

