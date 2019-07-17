package bot

import (

)

const (
	crossMarkEmoji 		= "\xE2\x9D\x8C"
	heavyCheckMarkEmoji = "\xE2\x9C\x94"
	
	cmdStart = "start"
	cmdHelp  = "help"
	cmdUnit  = "unit"
	cmdTest  = "test"
	cmdStop  = "stop"
	
	answerPlaceholder = "{{answer}}"
)

type DialogHandler interface {
	ProcessCommand(cmdText string, args []string)
	ProcessMessage(msg string)
	ProcessKeyboard(key string)
}

type KeyboardLayout = [][]string

type OutMessage struct {
	Text		string
	Keyboard	KeyboardLayout
}

type Sender func(msg OutMessage)


