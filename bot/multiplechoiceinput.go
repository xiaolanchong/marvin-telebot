package bot

import (
	"strconv"
	"fmt"
)

type MultipleChoiceHandler struct{
	Sender				Sender
	Questions			[]MultipleChoiceQuestion
	CurrentQuestionNum	int
	CorrectAnswers		int
	TotalAnswers		int
	
	SystemHandler	DialogHandler
}

const skipQuestionKey 	= "skip"
const exitQuizKey 		= "exit"

func (handler *MultipleChoiceHandler) showQuestion() {
	if len(handler.Questions) > handler.CurrentQuestionNum {
		curQuestion := handler.Questions[handler.CurrentQuestionNum]
		
		layout := KeyboardLayout {
							make([]Key, len(curQuestion.Choice)),
							make([]Key, 2),
						}
		for i, choice := range(curQuestion.Choice) {
			layout[0][i] = Key{ Id: strconv.Itoa(i), Text: choice }
		}
		layout[1][0] = Key{ Id: skipQuestionKey, Text: "Пропустить" }
		layout[1][1] = Key{ Id: exitQuizKey, Text: "Выйти" }
		
		handler.Sender(OutMessage{Text:curQuestion.Text, Keyboard: layout})
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
}

func (handler *MultipleChoiceHandler) ProcessMessage(msg string) {
	handler.SystemHandler.ProcessMessage(msg)
}

func (handler *MultipleChoiceHandler) processAnswer(key string) {
	handler.TotalAnswers++
	curQuestion := handler.Questions[handler.CurrentQuestionNum]
	if curQuestion.Answer == key {
		handler.CorrectAnswers++
		handler.Sender(OutMessage{Text: "Верно!"})
	} else {
		outMsg := "Ошибка! Правильно " + curQuestion.Answer
		handler.Sender(OutMessage{Text: outMsg})
	}
}

func (handler *MultipleChoiceHandler) finishTest() {
	resultText := fmt.Sprintf("Тест окончен, кол-во верных ответов %d/%d",
	handler.CorrectAnswers, handler.TotalAnswers)
	handler.Sender(OutMessage{Text: resultText})
}

func (handler *MultipleChoiceHandler) ProcessKeyboard(key string, messageId int) {
	if len(handler.Questions) <= handler.CurrentQuestionNum {
		return;
	}
	
	if key == exitQuizKey {
		handler.finishTest()
		return
	}
	
	if key != skipQuestionKey {
	handler.processAnswer(key)
	}
	
	handler.CurrentQuestionNum++
	if len(handler.Questions) > handler.CurrentQuestionNum {
		handler.showQuestion()
	} else {
		handler.finishTest()
	}
}