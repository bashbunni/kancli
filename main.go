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
	todos      list.Model
	inProgress list.Model
	done       list.Model
	quitting   bool
	err        error
}

type Task struct {
	status      status
	description string
}

func (t Task) FilterValue() string {
	return t.description
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
	m.todos.SetItems(todos)
	m.inProgress.SetItems(inProgressItems)
	m.done.SetItems(doneItems)
}

// TODO: organize by Task status
func initialModel() model {
	m := model{state: todo}
	m.tasks = []Task{
		{status: todo, description: "buy milk"},
		{status: todo, description: "eat sushi"},
		{status: todo, description: "fold laundry"},
		{status: inProgress, description: "write code"},
		{status: done, description: "stay cool"},
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 10, 10)
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
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, constants.QuitKeys) {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	case constants.ErrMsg:
		m.err = msg
		return m, nil
	default:
		var cmd tea.Cmd
		// do nothing
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.done.View(), m.inProgress.View(), m.done.View())
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
