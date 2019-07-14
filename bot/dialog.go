package bot

import (
	"log"
	"time"
)
/*
type MessageType = int8

const (
	Ordinal		MessageType = 0
	DialogEnd	MessageType = 1
)
*/

type command struct {
	Command		string
	Args		[]string
}

type action struct {
	Command		*command
	Message 	*string
	Timeout 	bool
}

type ActionChannel chan action

type Sender func(text string)

type Dialog struct {
	Username		string
	Sender			Sender
	ActionChannel	ActionChannel
	Timeout			time.Duration
	Timer			*time.Timer
	QuizState		*QuizState
}

var unitData []Unit
const quizFilePath = "data/Intermediate Korean - a Grammar and Workbook.yaml"
const altQuizFilePath = "bot/" + quizFilePath

func NewDialog(Sender Sender, timeout time.Duration, username string) *Dialog {
	if unitData == nil {
		var err error
		unitData, err = LoadUnits(quizFilePath)
		if err != nil {
			log.Printf("Error on loading exercise data: %v", err)
			unitData, err = LoadUnits(altQuizFilePath)
			if err != nil {
				log.Printf("Error on loading exercise data: %v", err)
			}
		}
	}

	this := &Dialog{
		Username: username,
		Sender: Sender,
		Timer: time.NewTimer(timeout),
		Timeout: timeout,
		ActionChannel: make(ActionChannel),
	}
	this.Timer.Stop()

	go func() {
		for range this.Timer.C {
			this.ActionChannel <- action { Timeout: true }
		}
	}()
	
	go func() {
		for act := range this.ActionChannel {
			if(act.Timeout) {
				log.Printf("[%s]Timeout received", this.Username)
				this.Sender("Conversation ended")
				
			} else if(act.Message != nil) {
				log.Printf("[%s]Message %s received", this.Username, *act.Message)
				this.stopTimer()
				var textResponse string
				textResponse, this.QuizState = ProcessMessage(*act.Message, unitData, this.QuizState)
				this.Sender(textResponse)
				
			} else if(act.Command != nil) {
				this.stopTimer()
				var textResponse string
				textResponse, this.QuizState = ProcessCommand(act.Command.Command, act.Command.Args, this.Username, unitData, this.QuizState)
				log.Printf("[%s]Command %+v received", this.Username, *act.Command)
				if len(textResponse) != 0 {
					this.Sender(textResponse)
				}
			}
		}
	}()

	return this
}

func (dialog *Dialog) OnCommand(cmd string, args []string) {
	dialog.ActionChannel <- action{ Command: &command {Command: cmd, Args: args } }
}

func (dialog *Dialog) OnMessage(text string) {
	dialog.ActionChannel <- action{ Message: &text }
}

func (dialog *Dialog) stopTimer() {
	dialog.Timer.Stop()
}