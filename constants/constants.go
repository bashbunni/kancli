package constants

import "github.com/charmbracelet/bubbles/key"

type (
	ErrMsg error
)

var QuitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)
