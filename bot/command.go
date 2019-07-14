package bot

import (
	"strings"
	"strconv"
	"fmt"
)

const (
	crossMarkEmoji 		= "\xE2\x9D\x8C"
	heavyCheckMarkEmoji = "\xE2\x9C\x94"
	
	cmdStart = "start"
	cmdHelp  = "help"
	cmdUnit  = "unit"
	cmdTest  = "test"
	cmdStop  = "stop"
)

type QuizState struct {
	Unit				uint
	Exercise			uint
	Question			uint

	CorrectAnswers		uint
	IncorrectAnswers	uint
}

var helpMessage = "Марвин - телеграм-бот для тестирования корейского языка.\n" +
				"Команды:\n" +
				"/" + cmdStart + "        Начать работу с ботом\n" +
				"/" + cmdHelp  + "        Показать подсказку\n" +
				"/" + cmdUnit  + " [N]    Показать теорию по разделу номер N (N=1..10)\n" +
				"/" + cmdTest  + " [N.M]  Начать тест N (N=1..10)\n" +
				"/" + cmdStop  + "        Прервать текущий тест\n" +
				""



func ProcessCommand(cmdText string, args []string, username string, unitData []Unit, quizState *QuizState) (string, *QuizState) {
	switch(cmdText) {
	case cmdStart:
		return fmt.Sprintf("언녕하세요, %s! Марвин - телеграм-бот для тестирования ваших знаний корейского языка. " +
						"Введите /help для вывода списка команд", username), nil
	case cmdHelp:
		return helpMessage, nil
	case cmdUnit:
		return "В разработке", nil
	case cmdTest:
		text, nextQuizState := ProcessStartTest(args, unitData)
		return text, nextQuizState
	case cmdStop:
		return "", nil
	default:
		return "Не знаю такой команды: " + cmdText, nil
	}
}

func ProcessMessage(msgText string, unitData []Unit, quizState *QuizState) (string, *QuizState) {
	if quizState == nil {
		return fmt.Sprintf("Давайте начнём тест командой /%s", cmdTest), nil
	}

	unitExercises := unitData[quizState.Unit].Exercise
	exercises := unitExercises[quizState.Exercise]
	question := exercises.Question[quizState.Question]
	
	response := ""
	if question.Answer != strings.TrimSpace(msgText) {
		quizState.IncorrectAnswers += 1
		response += crossMarkEmoji + " Ошибка!\nПравильно " + question.Answer
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

func ProcessStartTest(args []string, unitData []Unit) (string, *QuizState) {
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