package addmain

import (
    "fmt"
    "os"
    "strings"
    "bufio"
    "gopkg.in/yaml.v2"
)

type Metadata struct {
    Name string
    Description string
    Tags []string
    Promptables []Promptable
}

func extractMetadata(filename string) Metadata {
    var metadata Metadata

    mdString, err := extractSection(filename, "#~~#")
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
func extractSection(filename, marker string) (string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var extract bool = false
    var result strings.Builder

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println("line: " + line)

        if strings.HasPrefix(line, marker) {
            if !extract {
                extract = true
                continue
            } else {
                break
            }
        }

        if extract {
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

    return content
}

