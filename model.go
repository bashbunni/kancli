package main

import (
	"log"

	"github.com/bashbunni/kancli/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	focus    status
	loaded   bool
	lists    []list.Model
	quitting bool
}

func (m *Model) Next() {
	if m.focus == done {
		m.focus = todo
	} else {
		m.focus++
	}
}

func (m *Model) Prev() {
	if m.focus == todo {
		m.focus = done
	} else {
		m.focus--
	}
}

func newModel() Model {
	m := Model{focus: todo, loaded: false}
	return m
}

func (m *Model) initLists(width, height int) {
	// init list model
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-divisor*2)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	// add list items
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	m.lists[todo].Title = "To Do"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's go"},
	})
	m.lists[inProgress].Title = "In Progress"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
	m.lists[done].Title = "Done"
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			models[input] = newForm(m.focus)
			return models[input].Update(nil)
			// Note: I don't need a list of models, I can just return a new
			// form model each time, but I'm keeping it in this case so you can
			// see what it looks like with a list of models }
		}
	case Task:
		task := msg
		log.Println(task.status)
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
	}
	currList, cmd := m.lists[m.focus].Update(msg)
	m.lists[m.focus] = currList
	return m, cmd
}

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focus].SelectedItem()
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focus].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

func (m Model) View() string {
	var cols []string
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgView := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		switch m.focus {
		case inProgress:
			cols = []string{
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				columnStyle.Render(doneView),
			}
		case done:
			cols = []string{
				columnStyle.Render(todoView),
				columnStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			}
		default:
			cols = []string{
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgView),
				columnStyle.Render(doneView),
			}
		}
		return lipgloss.JoinHorizontal(lipgloss.Left, cols...)
	} else {
		return "Loading..."
	}
}
