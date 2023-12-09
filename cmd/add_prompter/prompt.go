package main

import (
    "github.com/erikgeiser/promptkit/selection"
    "os"
    "fmt"
	add "heydemo/add/addmain"
)

//type Choice selection.Choice[add.PromptableOption]


func Prompt(promptable add.Promptable, options []add.PromptableOption) {

    for _, option := range options {
        fmt.Println("option = " + option.String())
    }

    sp := selection.New("Select a value for da " + promptable.Name + " argument", options)
	sp.PageSize = 8

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		os.Exit(1)
	}

	// do something with the final choice
	add.PrettyPrint(choice)


}
