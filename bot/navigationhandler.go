package bot

import (
	"fmt"
	"log"
	"path/filepath"
)

type TerminalNode struct{
	Name			string
	DialogHandler	DialogHandler
}

type NonterminalNode struct{
	Name		string
	Children	[]NonterminalNode
	Leaf		[]TerminalNode
}

type NavigationHandler struct{
	Sender			Sender
	CurrentNode		*NonterminalNode	
}

var rootNavigation *NonterminalNode

func NewNavigationHandler(sender Sender,
		dataRootDir string) (*NavigationHandler, error) {
	var currentNode *NonterminalNode
	if rootNavigation == nil {
		fileName := filepath.Join(dataRootDir, "topik2_2019_mock_test1", "test.yaml")
		dataRoot, err := LoadMultipleChoiceTest(fileName)
		if err != nil {
			log.Printf("Failed to load test data: %+v", err)
			return nil, err
		}
		
		multipleChoiceHandler, mchErr := NewMultipleChoiceHandler(sender, dataRoot.Section[0].Question)
		if mchErr != nil {
			return nil, mchErr
		}
		currentNode = &NonterminalNode{
			Name:     "TOPIK 2 Mock Test",
			Leaf: []TerminalNode{ 
							TerminalNode{
								Name: dataRoot.Section[0].Title,
								DialogHandler: multipleChoiceHandler,
							},
						},
			}
	}
	
	
	return &NavigationHandler{
		Sender: sender,
		CurrentNode: currentNode,
	}, nil
}

func (handler *NavigationHandler) ProcessCommand(cmdText string, args []string) {
	switch(cmdText) {
	case cmdStart:
	//	handler.Sender( "언녕하세요! Марвин - телеграм-бот для тестирования ваших знаний корейского языка. " +
	//					"Введите /help для вывода списка команд")
	case cmdHelp:
	//	handler.Sender( helpMessage)
	default:
	//	handler.Sender( "Не знаю такой команды: " + cmdText )
	}
}

func (handler *NavigationHandler) ProcessMessage(msg string) {
	outMsg := fmt.Sprintf("Давайте начнём тест командой /%s", cmdTest)
	handler.Sender(OutMessage{Text: outMsg})
}

func (handler *NavigationHandler) ProcessKeyboard(key string) {
}