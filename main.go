package main

import (
	"container/list"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
	"tea_demo/api"
)

type model struct {
	choices  []os.FileInfo
	cursor   int
	selected map[int]struct{}
	path     string
	stack    *list.List
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "c":
			s := m.choices[m.cursor]
			if s != nil && s.IsDir() {
				m.stack.PushBack(m.path)
				m.path += "\\" + s.Name()
				m.choices = api.TreeFiles(m.path)
				m.cursor = 0
			}

		case "b":
			front := m.stack.Front()
			if front != nil {
				m.stack.Remove(front)
				s := front.Value.(string)
				m.path = s
				m.choices = api.TreeFiles(m.path)
				m.cursor = 0
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Files:" + m.path + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += getString(fmt.Sprintf("%s [%s] %s-----%v \n", cursor, checked, choice.Name(), choice.IsDir()))
	}

	s += "\nPress q to quit.\n"
	return s
}

var initModel = model{
	choices:  []os.FileInfo{},
	selected: make(map[int]struct{}),
	stack:    list.New(),
}

func main() {
	var path = "M:\\OneDrive - Liu Mi Tech\\go_project\\tea_demo"
	m := initModel
	m.path = path
	m.stack.PushBack(path)
	m.choices = api.TreeFiles(path)
	cmd := tea.NewProgram(m)
	if _, err := cmd.Run(); err != nil {
		fmt.Println("start failed:", err)
		os.Exit(1)
	}
}

func getString(str string) string {
	l := len(str)
	if l < 100 {
		builder := strings.Builder{}
		for i := 0; i < 100-l; i++ {
			builder.WriteString(" ")
		}
		return strings.ReplaceAll(str, "-----", builder.String())

	} else {
		return str
	}
}
