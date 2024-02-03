package addmain

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Metadata struct {
	Name        string `yaml:"-"`
	Description string
	Tags        []string
	Promptables []Promptable
}

const (
	extractMarker = "#~~#"
	skipMarker    = "#---#"
)

func ReadMetadata(filename string) Metadata {
	var metadata Metadata

	mdString, err := extractSection(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(mdString), &metadata)
	if err != nil {
		panic(err)
	}

	return metadata
}

func PromptForMetadata() Metadata {
	var metadata Metadata
	metadata.Description = promptForString("Description: ")
	return metadata
}

func promptForString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	response, _ := reader.ReadString('\n')
	return strings.TrimSuffix(response, "\n")
}

// Extracts part of a file between two markers - defined by marker
// Removes the leading # from each line
func extractSection(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var extract bool = false
	var skip bool = false
	var result strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, extractMarker) {
			if !extract {
				extract = true
				continue
			} else {
				break
			}
		} else if strings.HasPrefix(line, skipMarker) {
			skip = true
		}

		if extract && !skip {
			result.WriteString(strings.Trim(line, "#") + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return result.String(), nil
}

func prependToLines(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// Format metadata for insertion at top of script
func FormatMetadata(metadata Metadata) string {

	content, err := yaml.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	return extractMarker + "\n" + prependToLines(string(content), "#") + "\n" +
		getPromptableUsageGuide(metadata.Promptables) +
		extractMarker + "\n"

}

func UpdateMetadata(metadata Metadata, script string) {
	b, err := os.ReadFile(script)
	if err != nil {
		panic(err)
	}

	var (
		start int
		end   int
		count int
	)

	content := string(b)

	count = strings.Count(content, extractMarker)

	if count == 0 {
		start = strings.Index(content, "\n") + 1
		end = start
	} else {
		start = strings.Index(content, extractMarker)
		end = strings.LastIndex(content, extractMarker)
	}

	newContent := content[:start] + FormatMetadata(metadata) + content[end+len(extractMarker):]

	os.WriteFile(script, []byte(newContent), 0644)

}

func getPromptableUsageGuide(promptables []Promptable) string {
	content := skipMarker + "\n"

	content += getPromptableBashVarMapping(promptables) + "\n\n"

	content += "# PROMPTABLE USAGE GUIDE\n"
	content += "#\n"
	content += "# The following variables can be used in your script below.\n"
	content += "#\n"

	content += getPromptableBashVarDescription(promptables)
	content += "#\n"
	content += "# For more infomation see https://github.com/heydemo/add\n"
	content += "#\n"
	content += skipMarker + "\n"

	return content

}

func getPromptableBashVarMapping(promptables []Promptable) string {
	var content string
	for i, p := range promptables {
		content += p.getVarName() + "=$" + fmt.Sprint(i+1) + ";"
	}

	return content + "\n"
}

func getPromptableBashVarDescription(promptables []Promptable) string {
	var content string
	for _, p := range promptables {
		content += "# $" + p.getVarName() + " — " + p.Description + "\n"
	}

	return content

}
