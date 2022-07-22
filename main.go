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
	state      status
	tasks      []Task
	lists      []list.Model
	todos      list.Model
	inProgress list.Model
	done       list.Model
	quitting   bool
	err        error
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

func (t Task) Next() {}

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
	m.todos.SetItems(todos)
	m.todos.Title = "To Do"
	m.inProgress.SetItems(inProgressItems)
	m.inProgress.Title = "In Progress"
	m.done.SetItems(doneItems)
	m.done.Title = "Done"
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
	m.todos = defaultList
	m.inProgress = defaultList
	m.done = defaultList
	m.syncTasks()
	log.Print(m.todos.Items())
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
		if msg.String() == "n" {
			m.state = inProgress
			return m, nil
		}
		if msg.String() == "t" {
			m.state = todo
			return m, nil
		}
		if msg.String() == "d" {
			m.state = done
			return m, nil
		}
	case constants.ErrMsg:
		m.err = msg
	}
	switch m.state {
	case inProgress:
		m.todos, cmd = m.todos.Update(msg)
	case done:
		m.todos, cmd = m.todos.Update(msg)
	default:
		m.todos, cmd = m.todos.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	/*
		switch m.state {
		case todo:
			return m.todos.View()
		case inProgress:
			return m.inProgress.View()
		case done:
			return m.done.View()
		}
	*/
	return lipgloss.JoinHorizontal(lipgloss.Left, m.todos.View(), m.inProgress.View(), m.done.View())
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
