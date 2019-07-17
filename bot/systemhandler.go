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
	switch(cmdText) {
	case cmdStart:
		handler.Sender( "언녕하세요! Марвин - телеграм-бот для тестирования ваших знаний корейского языка. " +
						"Введите /help для вывода списка команд")
	case cmdHelp:
		handler.Sender( helpMessage)
	default:
		handler.Sender( "Не знаю такой команды: " + cmdText )
	}
}

func (handler *SystemDialogHandler) ProcessMessage(msg string) {
	outMsg := fmt.Sprintf("Давайте начнём тест командой /%s", cmdTest)
	handler.Sender(outMsg)
}

func (handler *SystemDialogHandler) ProcessKeyboard(key string) {
}

