package main

import (
    "fmt"
    "time"
	"log"
	"strings"
    "encoding/json"
    "os"
    "os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    content string
    left int
    top int
}

var memodel *model


func New(content string) *model {
    return &model{content: content, left: 0, top: 0}
}

func (m model) Init() tea.Cmd {
    return tickEvery()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    log.Println("msg", msg)
    // *memodel = m
    switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
              case "ctrl+c":
                  return m, tea.Quit
              case "up":
                  if m.top > 0 { m.top-- }
                  return m, nil
              case "down":
                  m.top++
                  return m, nil
              case "v":
                  //syscall.Exec("/usr/bin/vim", []string{"/usr/bin/vim", "debug.log"}, os.Environ())
                  interactiveExec("nvim", "debug.log")
                  return m, tea.Quit
              case "left":
                  if m.left > 0 { m.left-- }
              case "right":
                  m.left++
              default:
                  m.content = msg.String()
                  return m, nil
            }
        case TickMsg:
            m.top++
            return m, tickEvery()
    }
    return m, nil
}

func (m model) View() string {
    return strings.Repeat("\n", m.top) + strings.Repeat(" ", m.left) + m.content
}

type TickMsg time.Time

// Send a message every second.
func tickEvery() tea.Cmd {
    return tea.Every(time.Second / 3, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

func interactiveExec(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return nil
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
		} else {
			fmt.Println("Command finished successfully")
		}
	}()

	return cmd
}




func RunPrompterUI() {
    f, err := tea.LogToFile("debug.log", "debug")
    if err != nil {
        log.Fatal("err %w", err)
    }

    defer f.Close()

    memodel := &model{content: "hello world", left: 0, top: 0}
    // memodel = New("hello world")

    p := tea.NewProgram(memodel, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }

    println("content: " + memodel.content)
    fmt.Println("top: ", memodel.top)

    println("now our story comes to an end, my friend")
    jsonBytes, _ := json.MarshalIndent(*p, "", "  ")
    fmt.Println(string(jsonBytes))
    println(p)


}


