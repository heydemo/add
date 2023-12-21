package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type searchList struct {
	executables []Executable
}

func (s *searchList) Init() tea.Cmd {
	return nil
}

func (s *searchList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s *searchList) View() string {
	output := ""
	for _, e := range s.executables {
		output += e.Name + "\n"
	}
	return output
}
