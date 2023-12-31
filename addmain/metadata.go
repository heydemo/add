package addmain

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Metadata struct {
	Name        string
	Description string
	Tags        []string
	Promptables []Promptable
}

func extractMetadata(filename string) Metadata {
	var metadata Metadata

	mdString, err := extractSection(filename)
	if err != nil {
		panic(err)
	}

	// TODO: parse metadata as yaml
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

	extractMarker := "#~~#"
	skipMarker := "#---#"

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

// Format metadata for insertion at top of script
func formatMetadata(metadata Metadata) string {
	var content string = "#!/bin/bash\n"
	content += "#-# description: " + metadata.Description + "\n"
	if len(metadata.Tags) > 0 {
		content += "#-# tags: " + strings.Join(metadata.Tags, ",") + "\n"
	}
	content += "\n"

	return content
}
