package bot

import (
)

type MultipleChoiceHandler struct{
	Sender			Sender
	Questions		[]MultipleChoiceQuestion
	CurrentQuestionNum	int
	CorrectAnswers		int
	TotalAnswers		int
	
	SystemHandler	DialogHandler
}

func (handler *MultipleChoiceHandler) showQuestion() {
	if len(handler.Questions) > handler.CurrentQuestionNum {
		curQuestion := handler.Questions[handler.CurrentQuestionNum]
		handler.Sender(OutMessage{Text:curQuestion.Answer})  // <-- kbd
	}
	
}

func NewMultipleChoiceHandler(sender Sender,
			questions []MultipleChoiceQuestion) (*MultipleChoiceHandler, error){
	systemHandler := &SystemDialogHandler{ Sender: sender }
	newHandler := &MultipleChoiceHandler{
			Sender: sender,
			Questions: questions,
			CurrentQuestionNum: 0,
			CorrectAnswers: 0,
			TotalAnswers: 0,
			SystemHandler: systemHandler,
		}
	

	if len(questions) != 0 {
		newHandler.showQuestion()
	}
	return newHandler, nil
}

func (handler *MultipleChoiceHandler) ProcessCommand(cmdText string, args []string) {
	handler.SystemHandler.ProcessCommand(cmdText, args)
/*	switch(cmdText) {
	case cmdStart:
		//
	default:
		handler.Sender( "Не знаю такой команды: " + cmdText )
	}*/
}

func (handler *MultipleChoiceHandler) ProcessMessage(msg string) {
	handler.SystemHandler.ProcessMessage(msg)
}

func (handler *MultipleChoiceHandler) ProcessKeyboard(key string) {
	if len(handler.Questions) > handler.CurrentQuestionNum {
		handler.TotalAnswers++
		curQuestion := handler.Questions[handler.CurrentQuestionNum]
		if curQuestion.Answer == key {
			handler.CorrectAnswers++
			handler.Sender(OutMessage{Text: "Right!"})
		} else {
			outMsg := "Wrong!\n" + curQuestion.Answer
			handler.Sender(OutMessage{Text: outMsg})
		}
		
		handler.CurrentQuestionNum++
		if len(handler.Questions) > handler.CurrentQuestionNum {
			handler.showQuestion()
		}
	}
}