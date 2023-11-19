package addmain

import (
    "testing"
    "fmt"
)

func TestBootstrap(t *testing.T) {
    _, configEnv := Bootstrap()
    if configEnv == nil {
        t.Errorf("ConfigEnv is nil")
    }

    aliases := getAliasIncludeContent(configEnv)
    fmt.Println(aliases)
    if aliases == "" {
        t.Errorf("Alias content is empty")
    }

}



