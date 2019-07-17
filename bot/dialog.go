package bot

import (
	"log"
	"time"
//	"path/filepath"
)

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

type Dialog struct {
	Username		string
	Sender			Sender
	ActionChannel	ActionChannel
	Timeout			time.Duration
	Timer			*time.Timer
	DialogHandler	DialogHandler
}


func NewDialog(Sender Sender, timeout time.Duration, username string, dataRootDir string) *Dialog {
	this := &Dialog{
		Username: username,
		Sender: Sender,
		Timer: time.NewTimer(timeout),
		Timeout: timeout,
		ActionChannel: make(ActionChannel),
		//DialogHandler: dialogHandler,
	}
	this.Timer.Stop()
	systemHandler := &SystemDialogHandler{Sender: Sender}
	var err error
	this.DialogHandler, err = NewInputTestHandler(systemHandler, Sender, dataRootDir)
	if err != nil {
		return nil
	}

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
				//var textResponse string
				this.DialogHandler.ProcessMessage(*act.Message)
				//this.Sender(textResponse)
				
			} else if(act.Command != nil) {
				this.stopTimer()
				var textResponse string
				//textResponse, this.QuizState = ProcessCommand(act.Command.Command, act.Command.Args, this.Username, unitData, this.QuizState)
				this.DialogHandler.ProcessCommand(act.Command.Command, act.Command.Args)
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