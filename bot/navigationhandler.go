package bot

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
)

const (
	MockTest1Dir 			= "topik2_2019_mock_test1"
	ListeningTestFileName	= "listening.yaml"
	ReadingTestFileName		= "reading.yaml"
	WritingTestFileName		= "writing.yaml"
)

type NavigationHandler struct{
	Sender			Sender
	Menu			*MenuLayout
	ChildHandler	DialogHandler
}

func (handler *NavigationHandler) startTest(questions []MultipleChoiceQuestion, title string) {
	handler.ChildHandler, _ = NewMultipleChoiceHandler(handler.Sender, questions)
}

// Allows editing and reloading test files in runtime w/o restarting
func (handler *NavigationHandler) reloadAndStartTest(fullFileName string) error {
	dataRoot, err := LoadMultipleChoiceTest(fullFileName)
	if err != nil {
		log.Printf("Failed to load test data: %+v", err)
		return err
	}
	handler.startTest(dataRoot.Question, dataRoot.Title)
	return nil
}

func (handler *NavigationHandler) loadTestFile(fullFileName string, parentId MenuItemId) error {
	dataRoot, err := LoadMultipleChoiceTest(fullFileName)
	if err != nil {
		log.Printf("Failed to load test data: %+v", err)
		return err
	}

	handler.Menu.AddItem(parentId, dataRoot.Title, func() {
		err := handler.reloadAndStartTest(fullFileName)
		if err == nil {
			handler.Menu.Hide()
		} else {
			log.Printf("Failed to load test data from file %s : %+v", fullFileName, err)
		}
	})
	return nil
}

func (handler *NavigationHandler) loadDirFiles(dataRootDir string,
						menuName string, dateTestDir string, fileNames []string) error {
	top1Id, errL1 := handler.Menu.AddItem(RootMenuItemId, menuName, nil)
	if errL1 != nil {
		log.Printf("Failed to add item: %+v", errL1)
		return errL1
	}
	for _, fileName := range fileNames{
		fullFileName := filepath.Join(dataRootDir, dateTestDir, fileName)
		err := handler.loadTestFile(fullFileName, top1Id)
		if err != nil {
			return err
		}
	}

	// must be the last menu item
	handler.Menu.AddItem(top1Id, "< Назад", func() { 
		handler.Menu.GoUp()
	})

	return nil
}

func NewNavigationHandler(sender Sender, dataRootDir string) (*NavigationHandler, error) {
	handler := &NavigationHandler{
		Sender: sender,
		Menu: NewMenuLayout(),
	}
	
	errT := handler.loadDirFiles(dataRootDir,
				"TOPIK 2 пробный тест 1 (2019)", MockTest1Dir,
				[]string { ListeningTestFileName, WritingTestFileName, ReadingTestFileName,})
	if errT != nil {
		return nil, errT
	}
	
	errT60 := handler.loadDirFiles(dataRootDir, 
				"TOPIK 1 N60 (2018)", "topik1_2018_60",
				[]string { "listening.yaml", "reading.yaml",})
	if errT60 != nil {
		return nil, errT60
	}

	errD := handler.loadDirFiles(dataRootDir, 
				"日本語パワードリル [N1 文法]", "NihongoPowerDrillN1",
				[]string { "01.yaml", })
	if errD != nil {
		return nil, errD
	}
	
	// must be the last menu item
	handler.Menu.AddItem(RootMenuItemId, "Выход", func() { 
		handler.Menu.Hide()
	})

	return handler, nil
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
	handler.ChildHandler = nil

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

func (handler *NavigationHandler) ProcessMessage(msg InMessage) {
	outMsg := fmt.Sprintf("Давайте начнём тест командой /%s", cmdTest)
	handler.Sender(OutMessage{Text: outMsg})
}

func (handler *NavigationHandler) processMyKey(keyId string, messageId MessageId) {
	// TODO hide keyboards from unknown message
	intId, err := strconv.Atoi(keyId)
	if err != nil {
		log.Printf("Error converting key id to int: %v, key=%v, messageId=%v", err, keyId, messageId)
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

func (handler *NavigationHandler) ProcessKeyboard(keyId string, messageId MessageId) {
	if handler.ChildHandler != nil {
		handler.ChildHandler.ProcessKeyboard(keyId, messageId)
	} else {
		handler.processMyKey(keyId, messageId)
	}
}