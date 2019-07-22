package bot

import (
	"strconv"
	"log"
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
	if len(handler.Questions) <= handler.CurrentQuestionNum {
		return
	}

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
	
	if len(curQuestion.Audio) != 0 {
		handler.Sender(OutMessage{Text:curQuestion.Text})
		handler.Sender(OutMessage{Audio: curQuestion.Audio, Keyboard: layout})
	} else {
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

func (handler *MultipleChoiceHandler) ProcessMessage(msg InMessage) {
	handler.SystemHandler.ProcessMessage(msg)
}

func (handler *MultipleChoiceHandler) processAnswer(keyId string) {
	choiceIndex, errConv := strconv.Atoi(keyId)
	if errConv != nil {
		log.Printf("Error converting key id %s: %v, key sent a wrong id?", keyId, errConv)
		return
	}
	
	curQuestion := handler.Questions[handler.CurrentQuestionNum]
	if choiceIndex < 0 || choiceIndex >= len(curQuestion.Choice) {
		log.Printf("Choice index %d is out of range [%d, %d], key sent a wrong id?",
				choiceIndex, 0, len(curQuestion.Choice))
		return
	}
	userAnswer := curQuestion.Choice[choiceIndex]
	handler.TotalAnswers++

	comment := curQuestion.Comment
	var outTxt string
	if len(curQuestion.Choice) > 1 {		// don't count a question if no answer alternatives
		if curQuestion.Answer == userAnswer {
			handler.CorrectAnswers++
			outTxt = fmt.Sprintf("%s %s - верно!", OkEmoji, userAnswer)
		} else {
			outTxt = fmt.Sprintf("%s %s - ошибка! Правильно %s", ErrEmoji, userAnswer, curQuestion.Answer)
		}
	}
	if len(comment) != 0 {
		if len(outTxt) != 0 {
			outTxt += "\n"
		}
		outTxt += comment
	}
	handler.Sender(OutMessage{Text: outTxt}) 
}

func (handler *MultipleChoiceHandler) finishTest() {
	resultText := fmt.Sprintf("Тест окончен, кол-во верных ответов %d/%d.\nВведите /%s для выбора теста", 
						handler.CorrectAnswers, handler.TotalAnswers, cmdStart)
	handler.Sender(OutMessage{Text: resultText})
}

func (handler *MultipleChoiceHandler) ProcessKeyboard(keyId string, messageId MessageId) {
	HideMessageKeyboard(messageId, handler.Sender)
	
	if len(handler.Questions) <= handler.CurrentQuestionNum {
		return
	}
	
	if keyId == exitQuizKey {
		handler.finishTest()
		return
	}
	
	if keyId != skipQuestionKey {
		handler.processAnswer(keyId)
	}
	
	handler.CurrentQuestionNum++
	if len(handler.Questions) > handler.CurrentQuestionNum {
		handler.showQuestion()
	} else {
		handler.finishTest()
	}
}