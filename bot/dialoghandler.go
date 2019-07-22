package bot

import (

)

const (
	crossMarkEmoji 		= "\xE2\x9D\x8C"		// x
	whiteHeavyCheckMarkEmoji = "\xE2\x9C\x85"   // v
	
	OkEmoji = whiteHeavyCheckMarkEmoji
	ErrEmoji = crossMarkEmoji
	
	cmdStart = "start"
	cmdHelp  = "help"
	cmdFeedback  = "feedback"
	
	cmdUnit  = "unit"
	cmdTest  = "test"
	cmdStop  = "stop"
	
	answerPlaceholder = "{{answer}}"
)

type Key struct {
	Id		string
	Text	string
}

type KeyboardLayout = [][]Key

type MessageId = int

// Incoming message (tg -> bot)
type InMessage struct {
	MessageId			int
	Text 				string
}

// Outgoing message (bot -> tg)
type OutMessage struct {
	Text				string			// mutually excluded with Audio field
	Keyboard			KeyboardLayout	// include incline 
	IsKeyboardMsg		bool			// change only the msg keyboard, ReplyToMessageId must be set
	ReplyToMessageId	MessageId
	Audio				string 			// sent audio instead of text
}

type DialogHandler interface {
	ProcessCommand(cmdText string, args []string)
	ProcessMessage(msg InMessage)
	ProcessKeyboard(key string, fromMessageId MessageId)
}

type Sender func(msg OutMessage)

func HideMessageKeyboard(messageId MessageId, sender Sender) {
	hideKeyboardMsg := OutMessage{
				Keyboard:         KeyboardLayout{},
				ReplyToMessageId: messageId,
				IsKeyboardMsg:    true,
			}
	sender(hideKeyboardMsg)
}
