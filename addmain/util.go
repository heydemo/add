package addmain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

// Return the text of a file, excluding any lines that start with a #
func ReadBashScriptWithoutComments(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var result strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip lines starting with '#'
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			result.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return result.String(), nil

}
