package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	add "heydemo/add/addmain"
)

type ExecFinder struct {
	configEnv   *add.ConfigEnv
	list        list.Model
	list_init   bool
	currentMsg  statusMsg
	editorStyle lipgloss.Style
	selectStyle lipgloss.Style
	mode        UpdateMode
	err         error
	viewport    viewport.Model
	ready       bool
	height      int
	msgs        []tea.Msg
	sizes       modelSizes
}

type UpdateMode int

type Executable struct {
	Name        string
	Collection  string
	Path        string
	Description string
}

func (e Executable) String() string {
	return e.Name
}

type loadListItemsMsg []list.Item

// Generic informational message
type statusMsg struct {
	success bool
	status  string
}

const (
	NormalMode = iota
	ConfirmMode
)

const (
	listHeight = 20
	listWidth  = 0
)

type listKeyMap struct {
	editItem key.Binding
	copyItem key.Binding
}

var keyBindings listKeyMap = listKeyMap{
	editItem: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	copyItem: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy to clipboard"),
	),
}

type ExecutableWrapper struct {
	executable Executable
}

func (e ExecutableWrapper) FilterValue() string {
	return e.executable.Name
}

func (e ExecutableWrapper) Title() string { return e.executable.Name }

var (
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ExecutableWrapper)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func newExecFinder(configEnv *add.ConfigEnv) *ExecFinder {
	return &ExecFinder{configEnv: configEnv}
}

func (e *ExecFinder) Init() tea.Cmd {

	return loadExecutablesFromDir(e.configEnv.Bin_dir)
}

type modelSizes struct {
	selectWidth  int
	editorWidth  int
	editorHeight int
	screenWidth  int
	screenHeight int
}

// Calculate sizes based on the terminal width and height
func calculateSizes(width, height int) modelSizes {
	height -= 10
	width -= 15

	selectWidth := int(math.Round(float64(width) * 0.3))
	editorWidth := width - selectWidth

	return modelSizes{
		selectWidth:  selectWidth,
		editorWidth:  editorWidth,
		editorHeight: height,
		screenWidth:  width,
		screenHeight: height,
	}

}

func (e *ExecFinder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	useHighPerformanceRenderer := false

	e.msgs = append(e.msgs, msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case e.mode == ConfirmMode:
			return e.updateConfirmDelete(msg)
		case key.Matches(msg, keyBindings.editItem) && !e.list.SettingFilter():
			// Potential panic problem
			selected := e.list.SelectedItem().(ExecutableWrapper)
			return e, openEditor(selected.executable.Path)
		case key.Matches(msg, keyBindings.copyItem) && !e.list.SettingFilter():
			selected := e.list.SelectedItem().(ExecutableWrapper)
			return e, copyToClipboard(selected.executable.Path)

		case msg.String() == "ctrl+d":
			e.mode = ConfirmMode
			return e, nil
		case msg.String() == "ctrl+j":
			e.viewport.LineDown(1)
			return e, nil
		case msg.String() == "ctrl+k":
			e.viewport.LineUp(1)
			return e, nil
		}
	case tea.WindowSizeMsg:
		if !e.ready {
			e.sizes = calculateSizes(msg.Width, msg.Height)
			e.viewport = viewport.New(e.sizes.editorWidth, e.sizes.editorHeight)
			e.viewport.YPosition = 0
			e.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			e.ready = true

			e.editorStyle = lipgloss.NewStyle().
				Height(e.sizes.editorHeight).
				Width(e.sizes.editorWidth).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFF")).
				Padding(0).
				Margin(0)

			e.viewport.Style = lipgloss.NewStyle().Margin(0).Padding(0, 1)

			e.selectStyle = lipgloss.NewStyle().
				Height(e.sizes.editorHeight).
				Width(e.sizes.selectWidth).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFF")).
				Padding(0).
				Margin(0).
				Foreground(lipgloss.Color("#FFF"))

		}
	case statusMsg:
		e.currentMsg = msg
		return e, nil
	case loadListItemsMsg:
		e.list_init = true
		e.list = list.New(msg, itemDelegate{}, listWidth, e.sizes.editorHeight)
		e.list.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{
				keyBindings.editItem,
				keyBindings.copyItem,
			}
		}
		e.list.SetShowHelp(false)
		e.list.SetFilteringEnabled(true)
		e.list.Title = "Scripts"
		e.list.Styles.Title = lipgloss.NewStyle().
			Background(lipgloss.Color("default")).
			Foreground(lipgloss.Color("#3AA"))
	}

	if e.list_init {
		newListModel, cmd := e.list.Update(msg)
		e.list = newListModel
		cmds = append(cmds, cmd)

		return e, tea.Batch(cmds...)
	}

	return e, nil

}

type editorFinishedMsg struct{ err error }

func openEditor(filename string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, filename) //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func copyToClipboard(filename string) tea.Cmd {
	return func() tea.Msg {
		contents, err := add.ReadBashScriptWithoutComments(filename)
		if err != nil {
			return statusMsg{false, err.Error()}
		}
		contents = strings.Trim(contents, "\n ")
		err = clipboard.WriteAll(string(contents))
		if err != nil {
			return statusMsg{false, err.Error()}
		}

		return statusMsg{true, "Copied to clipboard"}
	}
}

func (e *ExecFinder) updateConfirmDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			item := e.list.SelectedItem()
			wrapper := item.(ExecutableWrapper)
			e.list.RemoveItem(e.list.Index())

			err := os.Remove(wrapper.executable.Path)
			if err != nil {
				e.err = err
			}

			e.mode = NormalMode
			return e, nil
		default:
			e.mode = NormalMode
			e.err = nil
			return e, nil
		}
	}

	return e, nil

}

func removeNthElement(slice []int, n int) []int {
	return append(slice[:n], slice[n+1:]...)
}

func (e *ExecFinder) helpView() string {
	return e.list.Styles.HelpStyle.Render(e.list.Help.View(e.list))
}

func (e *ExecFinder) statusView() string {
	var color string

	if e.currentMsg.status == "" {
		return ""
	}

	if e.currentMsg.success {
		color = "#0F0"
	} else {
		color = "#F00"
	}

	return lipgloss.NewStyle().
		Width(e.sizes.screenWidth).
		Border(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color(color)).
		Render(e.currentMsg.status)

}

func (e *ExecFinder) View() string {
	list := e.list.View()

	editorContent := e.getEditorView()
	wrapped := lipgloss.NewStyle().Width(e.sizes.editorWidth - 5).Height(e.sizes.editorHeight - 20).Render(editorContent)
	e.viewport.SetContent(wrapped)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		e.statusView(),
		lipgloss.JoinHorizontal(lipgloss.Top,
			e.selectStyle.Render(list),
			e.editorStyle.Render(e.viewport.View()),
		),
		e.getFooter(),
		e.helpView(),
	)

	return content
}

func (e *ExecFinder) getEditorView() string {
	currentItem, ok := e.list.SelectedItem().(ExecutableWrapper)

	if !ok {
		return "No file to load rn"
	}

	return e.getContents(currentItem.executable)

}

func (e *ExecFinder) getFooter() string {
	if e.mode == ConfirmMode {
		return "Are you sure you want to delete this script? (y/n)"
	}
	return ""
}

func (e *ExecFinder) getContents(exec Executable) string {
	filename := getFilePath(exec, e.configEnv)
	//contents, err := os.ReadFile(filename)
	args := strings.Split("-O truecolor -S bash -s baycomb "+filename, " ")
	contents, err := add.SubprocAndOutput("highlight", args...)
	if err != nil {
		return err.Error()
	}
	return contents

}

func getFilePath(e Executable, configEnv *add.ConfigEnv) string {
	return configEnv.Bin_dir + "/" + e.Collection + "/" + e.Name
}

func main() {

	if os.Getenv("ADD_DEBUG") == "true" {
		fmt.Println("Waiting for debugger to attach. Press ENTER to continue...")
		fmt.Printf("Current Process ID: %d\n", os.Getpid())
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	freshInstall, configEnv := add.Bootstrap()
	if freshInstall {
		return
	}

	model := newExecFinder(configEnv)

	p := tea.NewProgram(model)

	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)

		os.Exit(1)
	}
}

func loadExecutablesFromDir(binDir string) tea.Cmd {
	return func() tea.Msg {
		var executables []list.Item
		directories, error := os.ReadDir(binDir)
		if error != nil {
			panic(error)
		}
		for _, directory := range directories {
			if directory.IsDir() {
				files, error := os.ReadDir(binDir + "/" + directory.Name())
				if error != nil {
					panic(error)
				}
				for _, file := range files {
					if !file.IsDir() {
						executables = append(executables, ExecutableWrapper{
							executable: Executable{
								Name:       file.Name(),
								Collection: directory.Name(),
								Path:       binDir + "/" + directory.Name() + "/" + file.Name(),
							},
						})
					}
				}
			}
		}
		return loadListItemsMsg(executables)
	}
}
