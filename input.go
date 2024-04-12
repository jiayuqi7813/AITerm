package main

import "slices"

type InputCommand struct {
	Command     []byte
	Cursor      int
	controlMode bool
}

func (i *InputCommand) Send(b byte) error {
	if i.controlMode {
		switch b {
		case '[':
			return nil
		case 'D': // 左箭头
			i.Move(-1)
		case 'C':
			i.Move(1) // 右箭头
		}
	}

	switch b {
	case '\n':
	// TODO
	case '\x1B':
		i.controlMode = true
	default:
		slices.Insert(i.Command, i.Cursor, b)
		i.Cursor++
	}
	return nil
}

func (i *InputCommand) Move(l int) {
	i.Cursor += l
	if i.Cursor > len(i.Command) {
		i.Cursor = len(i.Command)
	}
	if i.Cursor < 0 {
		i.Cursor = 0
	}
}

func (i *InputCommand) String() string {
	return string(i.Command)
}
