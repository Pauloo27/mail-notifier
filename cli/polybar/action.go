package polybar

import "fmt"

type MouseIndex uint

// From https://github.com/polybar/polybar/wiki/Formatting#action-a
const (
	LeftClick = MouseIndex(iota + 1)
	MiddleClick
	RightClick
	ScrollUp
	ScrollDown
	// Double click is kinda "meh", so avoid it
	DoubleLeftClick
	DoubleMiddleClick
	DoubleRightClick
)

type ActionButton struct {
	Index                            MouseIndex
	Display, Command, UnderlineColor string
}

func (a ActionButton) String() string {
	button := fmt.Sprintf("%%{A%d:%s:}%s%%{A}", a.Index, a.Command, a.Display)
	if a.UnderlineColor == "" {
		return button
	}
	return fmt.Sprintf("%%{u%s}%%{+u}%s%%{-u}", a.UnderlineColor, button)
}

func ActionOver(a ActionButton, index MouseIndex, command string) ActionButton {
	return ActionButton{
		Index: index, Display: a.String(), Command: command,
	}
}
