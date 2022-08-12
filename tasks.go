package main

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbletea-app-template/constants"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	focus    status
	loaded   bool
	lists    []list.Model
	quitting bool
	err      error
}

func (m *model) Next() {
	if m.focus == done {
		m.focus = todo
	} else {
		m.focus++
	}
}

func (m *model) Prev() {
	if m.focus == todo {
		m.focus = done
	} else {
		m.focus--
	}
}

func newModel() model {
	m := model{focus: todo, loaded: false}
	return m
}

func (m *model) initLists(width, height int) tea.Msg {
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
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// TODO: check if this is where they put custom messages in other examples
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
			case "n":
				// save state of current model before switching models
				models[tasks] = m
				return models[input].Update(nil)
			}
	case constants.ErrMsg:
		m.err = msg
	case Task:
		task := Task(msg)
		log.Print(m.lists)
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
	}
	currList, cmd := m.lists[m.focus].Update(msg)
	m.lists[m.focus] = currList
	return m, cmd
}

func (m *model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focus].SelectedItem()
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focus].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.loaded {
		switch m.focus {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left, 
				columnStyle.Render(m.lists[todo].View()), 
				focusedStyle.Render(m.lists[inProgress].View()), 
				columnStyle.Render(m.lists[done].View())) + "\n"
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left, 
				columnStyle.Render(m.lists[todo].View()), 
				columnStyle.Render(m.lists[inProgress].View()), 
				focusedStyle.Render(m.lists[done].View())) + "\n"
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left, 
				focusedStyle.Render(m.lists[todo].View()), 
				columnStyle.Render(m.lists[inProgress].View()), 
				columnStyle.Render(m.lists[done].View())) + "\n"
		}
	} else {
		return "Loading..."
	}
}


