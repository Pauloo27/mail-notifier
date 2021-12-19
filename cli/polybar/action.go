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
	Index            MouseIndex
	Display, Command string
}

func (a ActionButton) String() string {
	return fmt.Sprintf("%%{A%d:%s:}%s%%{A}", a.Index, a.Command, a.Display)
}

func ActionOver(a ActionButton, index MouseIndex, command string) ActionButton {
	return ActionButton{index, a.String(), command}
}
