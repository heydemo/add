package addmain

import (
	"os"
	"path/filepath"
    "os/exec"
)

type ConfigEnv struct {
	Config_dir     string
	User_bin_dir   string
	Public_bin_dir string
    Core_bin_dir   string
	State_dir      string
}

func GetConfigEnv() *ConfigEnv {
	var config_env ConfigEnv
	config_env.Config_dir = getConfigDir()
	config_env.State_dir = getStateDir()
	config_env.User_bin_dir = filepath.Join(config_env.Config_dir, "bin", "user")
	config_env.Public_bin_dir = filepath.Join(config_env.Config_dir, "bin", "public")
	config_env.Core_bin_dir = filepath.Join(config_env.Config_dir, "bin", "core")

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

func getConfigDir() string {
	return getFirstDefinedEnvVars(
		[]string{"ADD_CONFIG_DIR", "XDG_CONFIG_HOME"},
		filepath.Join(getUserHomeDir(), ".config", "add"))
}

func getStateDir() string {
	return getFirstDefinedEnvVars(
		[]string{"ADD_STATE_DIR", "XDG_STATE_HOME"},
		filepath.Join(getUserHomeDir(), ".local", "state", "add"))
}

func getFirstDefinedEnvVars(env_vars []string, defaultVal string) string {
	for _, env_var := range env_vars {
		if v := os.Getenv(env_var); v != "" {
			return v
		}
	}
	return defaultVal
}
