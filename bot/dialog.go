package bot

import (
	"log"
	"time"
)

type command struct {
	Command		string
	Args		[]string
}

type action struct {
	Command		*command
	Message 	*InMessage
	Key			*string
	MessageId	int
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

func NewDialog(Sender Sender, timeout time.Duration, username string, dlgHandler DialogHandler) *Dialog {
	if dlgHandler == nil {
		return nil
	}
	
	this := &Dialog{
		Username: username,
		Sender: Sender,
		Timer: time.NewTimer(timeout),
		Timeout: timeout,
		ActionChannel: make(ActionChannel),
		DialogHandler: dlgHandler,
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
				this.Sender(OutMessage{Text: "Conversation ended"} )
				
			} else if(act.Message != nil) {
				log.Printf("[%s]Message id=%v, text=%s received", this.Username, act.Message.MessageId, act.Message.Text)
				this.stopTimer()
				this.DialogHandler.ProcessMessage(*act.Message)
				
			} else if(act.Key != nil) {
				this.DialogHandler.ProcessKeyboard(*act.Key, act.MessageId)
			} else if(act.Command != nil) {
				this.stopTimer()
				var textResponse string
				this.DialogHandler.ProcessCommand(act.Command.Command, act.Command.Args)
				log.Printf("[%s]Command %+v received", this.Username, *act.Command)
				if len(textResponse) != 0 {
					this.Sender(OutMessage{Text: textResponse})
				}
			}
		}
	}()

	return this
}

func (dialog *Dialog) OnCommand(cmd string, args []string) {
	dialog.ActionChannel <- action{ Command: &command {Command: cmd, Args: args } }
}

func (dialog *Dialog) OnMessage(msg InMessage) {
	dialog.ActionChannel <- action{ Message: &msg }
}

func (dialog *Dialog) OnKey(keyId string, messageId int) {
	dialog.ActionChannel <- action{ Key: &keyId, MessageId: messageId }
}

func (dialog *Dialog) stopTimer() {
	dialog.Timer.Stop()
}