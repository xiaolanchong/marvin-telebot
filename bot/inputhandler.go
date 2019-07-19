
package bot

import (
	"log"
	"path/filepath"
	"strings"
	"strconv"
	"fmt"
)

type QuizState struct {
	Unit				uint
	Exercise			uint
	Question			uint

	CorrectAnswers		uint
	IncorrectAnswers	uint
}

type InputTestHandler struct {
	SystemHandler	DialogHandler
	Sender			Sender
	QuizState		*QuizState
}

const quizFilePath = "Intermediate Korean - a Grammar and Workbook"
const quizFileName = "exercises.yaml"

var input_unitData []Unit

func NewInputTestHandler(sender Sender, dataRootDir string) (*InputTestHandler, error){
	systemHandler := &SystemDialogHandler{ Sender: sender }
	if input_unitData == nil {
		var err error
		fullName :=  filepath.Join(dataRootDir, quizFilePath, quizFileName)
		input_unitData, err = LoadUnits(fullName)
		if err != nil {
			log.Printf("Error on loading exercise data: %v", err)
			return nil, err
		}
	}

	return &InputTestHandler {
		SystemHandler: systemHandler,
		Sender:        sender,
	}, nil
}

func (handler *InputTestHandler) ProcessCommand(cmdText string, args []string) {
	var msg string
	var quizState *QuizState
	switch(cmdText) {
	case cmdTest:
		msg, quizState = input_ProcessStartTest(args, input_unitData)
	default:
		handler.SystemHandler.ProcessCommand(cmdText, args)
	}
	
	handler.QuizState = quizState
	if len(msg) != 0 {
		handler.Sender(OutMessage{Text: msg})
	}
}

func (handler *InputTestHandler) ProcessMessage(msg string) {
	if handler.QuizState == nil {
		msg := fmt.Sprintf("Давайте начнём тест командой /%s", cmdTest)
		handler.Sender(OutMessage{Text: msg})
		return
	}
	
	var msgOut string
	msgOut, handler.QuizState = input_ProcessMessage(msg, input_unitData, handler.QuizState)
	if len(msgOut) != 0 {
		handler.Sender(OutMessage{Text: msgOut})
	}
}

func (handler *InputTestHandler) ProcessKeyboard(key string, messageId int) {
	log.Printf("Handler of %s not implemented", key)
}

func input_ProcessStartTest(args []string, unitData []Unit) (string, *QuizState) {
	if len(args) != 1 {
		return "Неизвестные параметры, ожидается Раздел.Упражнение: 1.1, 2.3 и т.д.", nil
	}
	parts := strings.Split(args[0], ".")
	if len(parts) != 2 {
		return "Неизвестные параметры, ожидается Раздел.Упражнение: 1.1, 2.3 и т.д.", nil
	}
	var unitNumber, exerciseNumber uint64
	var err error
	unitNumber, err = strconv.ParseUint(parts[0], 10, 31)
	if err != nil {
		return fmt.Sprintf("Неправильный номер раздела: %s", parts[0]), nil
	}
	exerciseNumber, err = strconv.ParseUint(parts[1], 10, 31)
	if err != nil {
		return fmt.Sprintf("Неправильный номер упражнения: %s", parts[1]), nil
	}
	if int(unitNumber) > len(unitData) || unitNumber <= 0 {
		return fmt.Sprintf("Номер раздела должен быть 1-%d", len(unitData)), nil
	}
	unitExercises := unitData[unitNumber - 1].Exercise
	if int(exerciseNumber) > len(unitExercises) || exerciseNumber <= 0 {
		return fmt.Sprintf("Номер упражнения должен быть 1-%d", len(unitExercises)), nil
	}
	exercises := unitExercises[exerciseNumber - 1]
	questions := exercises.Question
	if len(questions) == 0 {
		return fmt.Sprintf("Нет вопросов в упражнении %d.%d", unitNumber, exerciseNumber), nil
	}
	text := exercises.Description + "\n" + "\n" + questions[0].Text
	return text, &QuizState{
					Unit: uint(unitNumber - 1),
					Exercise: uint(exerciseNumber - 1),
					Question: 0,
				}
}

func input_ProcessMessage(msgText string, unitData []Unit, quizState *QuizState) (string, *QuizState) {
	unitExercises := unitData[quizState.Unit].Exercise
	exercises := unitExercises[quizState.Exercise]
	question := exercises.Question[quizState.Question]
	
	response := ""
	if question.Answer != strings.TrimSpace(msgText) {
		answer := question.Answer
		if len(question.Show) != 0 {
			answer = strings.Replace(question.Show, answerPlaceholder, question.Answer, -1)
		}
		quizState.IncorrectAnswers += 1
		response += crossMarkEmoji + " Ошибка!\nПравильно " + answer
	} else {
		quizState.CorrectAnswers += 1
		response += heavyCheckMarkEmoji + " Верно!"
	}
	
	quizState.Question += 1
	if len(exercises.Question) <= int(quizState.Question) {
		response += fmt.Sprintf("\nТест закончен. Верных ответов: %d из %d", quizState.CorrectAnswers, len(exercises.Question))
		return response, nil
	} else {
		response += "\n\n" + exercises.Question[quizState.Question].Text
		return response, quizState
	}
}
