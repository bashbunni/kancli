package main

/*
Qs
- don't need to store list of Items and data structure - should be able to wrap each of them when necessary
- cast it as type, modify it, sub it back in as interface
*/

/* functionality
- add tasks to current list
- edit selected task
- move selected task to next board
*/

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbletea-app-template/constants"
	"github.com/charmbracelet/lipgloss"
)

type (
	status       uint
	MovedTaskMsg bool
)

const (
	todo status = iota
	inProgress
	done
)

const divisor = 4

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)

	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
)

type model struct {
	state    status
	loaded   bool
	lists    []list.Model
	quitting bool
	err      error
}

func (m *model) Next() {
	if m.state == done {
		m.state = todo
	} else {
		m.state++
	}
}

func (m *model) Prev() {
	if m.state == todo {
		m.state = done
	} else {
		m.state--
	}
}

func initialModel() model {
	m := model{state: todo, loaded: false}
	return m
}

func (m *model) initLists(width, height int) {
	// init list model
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/divisor)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	// add list items
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's go"},
	})
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		if key.Matches(msg, constants.QuitKeys) {
			m.quitting = true
			return m, tea.Quit
		}
		switch msg.String() {
		case "right":
			m.Next()
		case "left":
			m.Prev()
		case "enter":
			return m, m.MoveToNext
		}
	case constants.ErrMsg:
		m.err = msg
	}
	currList, cmd := m.lists[m.state].Update(msg)
	m.lists[m.state] = currList
	return m, cmd
}

func (m *model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.state].SelectedItem()
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.state].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.loaded {
		switch m.state {
		case inProgress:
			return lipgloss.JoinHorizontal(lipgloss.Left, columnStyle.Render(m.lists[todo].View()), focusedStyle.Render(m.lists[inProgress].View()), columnStyle.Render(m.lists[done].View())) + "\n"
		case done:
			return lipgloss.JoinHorizontal(lipgloss.Left, columnStyle.Render(m.lists[todo].View()), columnStyle.Render(m.lists[inProgress].View()), focusedStyle.Render(m.lists[done].View())) + "\n"
		default:
			return lipgloss.JoinHorizontal(lipgloss.Left, focusedStyle.Render(m.lists[todo].View()), columnStyle.Render(m.lists[inProgress].View()), columnStyle.Render(m.lists[done].View())) + "\n"
		}
	} else {
		return "Loading..."
	}
}

func main() {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
