package main

/*
Qs
- don't need to store list of Items and data structure - should be able to wrap each of them when necessary
- cast it as type, modify it, sub it back in as interface
*/

/* functionality
- add tasks to current list
- edit selected task
*/

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	status uint
	page   int
)

const (
	todo status = iota
	inProgress
	done
)

const (
	tasks page = iota
	input
)

const (
	divisor = 4
)

var (
	models      []tea.Model
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)

	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

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
	// initialize models
	models = []tea.Model{newModel(), newForm(todo)}

	p := tea.NewProgram(models[tasks])
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
