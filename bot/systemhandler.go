package bot

import (
	"fmt"
)

var helpMessage = "Марвин - телеграм-бот для тестирования корейского языка.\n" +
				"Команды:\n" +
				"/" + cmdStart + "        Начать работу с ботом\n" +
				"/" + cmdHelp  + "        Показать подсказку\n" +
				"/" + cmdUnit  + " [N]    Показать теорию по разделу номер N (N=1..10)\n" +
				"/" + cmdTest  + " [N.M]  Начать тест N (N=1..10)\n" +
				"/" + cmdStop  + "        Прервать текущий тест\n" +
				""


type SystemDialogHandler struct {
	Sender	Sender
}

func (handler *SystemDialogHandler) ProcessCommand(cmdText string, args []string) {
	var outMsg string
	switch(cmdText) {
	case cmdStart:
		outMsg = "언녕하세요! Марвин - телеграм-бот для тестирования ваших знаний корейского языка. " +
						"Введите /help для вывода списка команд"
	case cmdHelp:
		outMsg = helpMessage
	default:
		outMsg = ( "Не знаю такой команды: " + cmdText )
	}
	handler.Sender(OutMessage{Text: outMsg})
}

func (handler *SystemDialogHandler) ProcessMessage(msg string) {
	outMsg := fmt.Sprintf("Давайте начнём тест командой /%s", cmdTest)
	handler.Sender(OutMessage{Text: outMsg})
}

func (handler *SystemDialogHandler) ProcessKeyboard(key string) {
}

