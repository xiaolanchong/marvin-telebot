package bot

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
)

type NavigationHandler struct{
	Sender			Sender
	Menu			*MenuLayout
}

func (handler *NavigationHandler) startTest(questions []MultipleChoiceQuestion, title string) {
	handler.Sender(OutMessage{Text: "Selected " + title})
}

func NewNavigationHandler(sender Sender,
		dataRootDir string) (*NavigationHandler, error) {
	fileName := filepath.Join(dataRootDir, "topik2_2019_mock_test1", "test.yaml")
	dataRoot, err := LoadMultipleChoiceTest(fileName)
	if err != nil {
		log.Printf("Failed to load test data: %+v", err)
		return nil, err
	}
	
	navHandler := &NavigationHandler{
		Sender: sender,
		Menu: NewMenuLayout(),
	}
	
	firstLevel, errL1 := navHandler.Menu.AddItem(RootMenuItemId, dataRoot.Description, nil)
	if errL1 != nil {
		log.Printf("Failed to add item: %+v", errL1)
		return nil, errL1
	}
	
	for _, section := range(dataRoot.Section) {
		navHandler.Menu.AddItem(firstLevel, section.Title, func() { 
			navHandler.startTest(section.Question, section.Title) 
		})
	}
	navHandler.Menu.AddItem(firstLevel, "< Back", func() { 
		navHandler.Menu.GoUp()
		navHandler.Menu.GoUp()
	//	log.Printf("Go up: %+v", "zz")
	})

	return navHandler, nil
}

func menuToKeyboard(menuLayout []MenuIdAndText) KeyboardLayout {
	result := make(KeyboardLayout, 0, len(menuLayout))
	for _, item := range(menuLayout) {
		result = append(result, []Key { 
					Key{ Id: strconv.Itoa(int(item.Id)),
					Text: item.Text,
					},
				})
	}
	return result
}

func (handler *NavigationHandler) ProcessCommand(cmdText string, args []string) {
	var outMsg string
	var keyboard KeyboardLayout
	switch(cmdText) {
	case cmdStart:
		outMsg = "언녕하세요! Марвин - телеграм-бот для тестирования TOPIK. " +
						"Выберите тест или введите /help для вывода списка команд"
		handler.Menu.GoTop()
		keyboard = menuToKeyboard(handler.Menu.GetCurrentLevel())
	case cmdHelp:
		outMsg = helpMessage
	case cmdStop:
		outMsg = "Stopped"
		handler.Menu.GoTop()
	default:
		outMsg = "Не знаю такой команды: " + cmdText
	}
	handler.Sender(OutMessage{Text: outMsg, Keyboard: keyboard})
}

func (handler *NavigationHandler) ProcessMessage(msg string) {
	outMsg := fmt.Sprintf("Давайте начнём тест командой /%s", cmdTest)
	handler.Sender(OutMessage{Text: outMsg})
}

func (handler *NavigationHandler) ProcessKeyboard(keyId string, messageId int) {
	intId, err := strconv.Atoi(keyId)
	if err != nil {
		log.Printf("Error converting key id to int: %v", err)
		return
	}
	handler.Menu.SelectItem(MenuItemId(intId))
	keyboard := menuToKeyboard(handler.Menu.GetCurrentLevel())
	outMsg := OutMessage{
				Keyboard:         keyboard,
				ReplyToMessageId: messageId,
				IsKeyboardMsg:    true,
			}
	handler.Sender(outMsg)
}