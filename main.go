package main

import (
	"container/list"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"tea_demo/api"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table     table.Model
	fileInfos []os.FileInfo
	path      string
	stack     *list.List
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[0]),
			)
		case "c":
			s := m.fileInfos[m.table.Cursor()]
			if s != nil && s.IsDir() {
				m.stack.PushFront(m.path)
				m.path += "\\" + s.Name()
				m.fileInfos = api.TreeFiles(m.path)
				m.table.SetCursor(0)
				m.table.SetRows(getRows(m.fileInfos))
			} else {
				return m, tea.Batch(
					tea.Printf("this is not a dir"),
				)
			}
		case "b":
			front := m.stack.Front()
			if front != nil {
				m.stack.Remove(front)
				s := front.Value.(string)
				m.path = s
				m.fileInfos = api.TreeFiles(m.path)
				m.table.SetCursor(0)
				m.table.SetRows(getRows(m.fileInfos))

			}
		case "d":
			s := m.fileInfos[m.table.Cursor()]
			if s != nil && s.IsDir() {
				m.path += "\\" + s.Name()
				api.Drop(m.path)
				m.fileInfos = api.TreeFiles(m.path)
				m.table.SetCursor(0)
				m.table.SetRows(getRows(m.fileInfos))
			}
			return m, tea.Batch(
				tea.Printf("dropped!"),
			)
		}

	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Size", Width: 20},
		{Title: "IsDir", Width: 10},
		{Title: "ModeTime", Width: 20},
	}

	var path = "C:\\Users\\Administrator\\Desktop\\zookeeper\\test"
	if len(os.Args) >= 2 {
		param := os.Args[1]
		fmt.Println(param)
		if param != "" {
			path = param
		}
	}
	files := api.TreeFiles(path)
	rows := getRows(files)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t, files, path, list.New()}
	m.stack.PushBack(path)
	m.fileInfos = api.TreeFiles(path)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getRows(files []fs.FileInfo) []table.Row {
	var rows []table.Row
	for _, file := range files {
		row := table.Row{}
		rows = append(rows, append(row, file.Name(), strconv.FormatInt(file.Size()/1000, 10), strconv.FormatBool(file.IsDir()), file.ModTime().Format("2006-01-02 15:04:05")))

	}
	return rows
}
