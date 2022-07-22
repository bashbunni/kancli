package main

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
	status uint
)

const (
	todo status = iota
	inProgress
	done
)

type model struct {
	state    status
	tasks    []Task
	lists    []list.Model
	quitting bool
	err      error
}

type Task struct {
	status      status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

func (t *Task) Prev() {
	if t.status == todo {
		t.status = done
	} else {
		t.status--
	}
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

func (m *model) syncTasks() {
	todos := []list.Item{}
	inProgressItems := []list.Item{}
	doneItems := []list.Item{}
	for _, task := range m.tasks {
		switch task.status {
		case inProgress:
			inProgressItems = append(inProgressItems, task)
		case done:
			doneItems = append(doneItems, task)
		default:
			todos = append(todos, task)
		}
	}
	m.lists[todo].SetItems(todos)
	m.lists[todo].Title = "To Do"
	m.lists[inProgress].SetItems(inProgressItems)
	m.lists[inProgress].Title = "In Progress"
	m.lists[done].SetItems(doneItems)
	m.lists[done].Title = "Done"
}

func initialModel() model {
	m := model{state: todo}
	m.tasks = []Task{
		{status: todo, title: "buy milk", description: "strawberry milk"},
		{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
		{status: inProgress, title: "write code", description: "don't worry, it's go"},
		{status: done, title: "stay cool", description: "as a cucumber"},
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 40)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	m.syncTasks()
	log.Print(m.lists[todo].Items())
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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
		}
	case constants.ErrMsg:
		m.err = msg
	}

	currList, cmd := m.lists[m.state].Update(msg)
	m.lists[m.state] = currList
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	/*
		switch m.state {
		case todo:
			return m.lists[todo].View()
		case inProgress:
			return m.lists[inProgress].View()
		case done:
			return m.lists[done].View()
		}
	*/
	return lipgloss.JoinHorizontal(lipgloss.Left, m.lists[todo].View(), m.lists[inProgress].View(), m.lists[done].View())
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
