
package bot

import (
	//"strings"
	"testing"
	"time"
	"os"
	"path/filepath"
)

var testQuizFileName string
var testMultipleChoiceFileName string
var dataRootDir string

type MsgKeeper struct {
	OutMessages  []OutMessage
}

func (keeper *MsgKeeper) Send(outMsg OutMessage) {
	if keeper.OutMessages == nil {
		keeper.OutMessages = make([]OutMessage, 0, 20)
	}
	keeper.OutMessages = append(keeper.OutMessages, outMsg)
}

func init() {
	dataRootDir = os.Getenv("DATA_ROOT")
	testQuizFileName = filepath.Join(dataRootDir, "Intermediate Korean - a Grammar and Workbook", "exercises.yaml")
	testMultipleChoiceFileName = filepath.Join(dataRootDir, "topik2_2019_mock_test1", "test.yaml")
}

func TestDialog_GrammarBook(t *testing.T) {
	units, err := LoadUnits(testQuizFileName)
	if err != nil {
		t.Errorf("Failed to load unit data: %v", err)
		return
	}
	
	if len(units) != 1 {
		t.Errorf("Incorrect unit number: %v", len(units))
		return
	}
	
	exp_desc := 
`1.1 Поставьте глагол или прилагательное в скобках в просторечный стиль, затем переведите приложение. Пример:
학교에 (가다)/повелительное наклонение
= 학교에 가` // "Go to school."`
	if exp_desc != units[0].Exercise[0].Description {
		t.Errorf("Incorrect description, expected:\n%s\nactual:\n%s", exp_desc, units[0].Exercise[0].Description)
	}

	if 6 != len(units[0].Exercise[0].Question) {
		t.Errorf("Incorrect question: expected: %v, actual: %v", 6, len(units[0].Exercise[0].Question))
	}
	
	questions := units[0].Exercise[0].Question
	exp_quest := "공연을 일찍 (마치다)/предложное наклонение, 공연 выступление, 마치다 кончать."
	if exp_quest != questions[len(questions)-1].Text {
		t.Errorf("Incorrect question:\nexpected: %v, \nactual: %v", exp_quest, questions[len(questions)-1].Text)
	}
}

func TestDialog_Input_StartCommand(t *testing.T) {
	
	//msgs := make([]string, 0, 10)
	//sender := func (text string) { 
	///		   msgs = append(msgs, text)
	//		}
	sender := MsgKeeper{}
	hndl, _ := NewInputTestHandler(sender.Send, dataRootDir)
	dlg := NewDialog(sender.Send,
			time.Second * 5,
			"nemo", hndl)
	dlg.OnCommand("start", []string{})
	time.Sleep(time.Second * 1)
	if(len(sender.OutMessages) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(sender.OutMessages))
	}
	if(len(sender.OutMessages[0].Text) != 211) {
		t.Errorf("Incorrect message length: %v bytes, %+v", len(sender.OutMessages[0].Text), sender.OutMessages[0])
	}
}

func TestDialog_Input_TestCommand(t *testing.T) {
	sender := MsgKeeper{}
	hndl, _ := NewInputTestHandler(sender.Send, dataRootDir)
	dlg := NewDialog(sender.Send,
			time.Second * 5,
			"nemo", hndl)
	dlg.OnCommand("test", []string{ "1.1" })
	time.Sleep(time.Second * 1)
	if(len(sender.OutMessages) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(sender.OutMessages))
	}
	if(len(sender.OutMessages[0].Text) != 401) {
		t.Errorf("Incorrect message length: %v bytes, %+v", len(sender.OutMessages[0].Text), sender.OutMessages[0])
	}
}

func TestDialog_Input_Take1stTest(t *testing.T) {
	sender := MsgKeeper{}
	hndl, _ := NewInputTestHandler(sender.Send, dataRootDir)
	dlg := NewDialog(sender.Send,
			time.Second * 5,
			"nemo", hndl)
	dlg.OnCommand("test", []string{ "1.1" })
	time.Sleep(time.Second * 1)
	if(len(sender.OutMessages) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(sender.OutMessages))
	}
	if(len(sender.OutMessages[0].Text) != 401) {
		t.Errorf("Incorrect message length: %v bytes, %+v", len(sender.OutMessages[0].Text), sender.OutMessages[0])
	}
	sender.OutMessages = make([]OutMessage, 0, 20)
	for i := 0; i < 10; i++ {
		dlg.OnMessage("11")
		time.Sleep(time.Millisecond * 50)
	}
	if len(sender.OutMessages) != 10 {
		t.Errorf("Incorrect messages: %v bytes, \n\n%+v", len(sender.OutMessages), sender.OutMessages)
	}
}

//---------------------------

func Test_LoadMultipleChoiceTest(t *testing.T) {
	root, err := LoadMultipleChoiceTest(testMultipleChoiceFileName)
	if err != nil {
		t.Errorf("Failed to load multiple choice test: %+v", err)
		return
	}
	
	if len(root.Description) < 10 {
		t.Errorf("Incorrect description of multiple choice test: %+v, %+v", len(root.Description), root.Description)
		return
	}
	if len(root.Section) != 1 {
		t.Errorf("Incorrect section number: %+v, %+v", len(root.Section), root.Section)
		return
	}
	
	section := root.Section[0]
	if len(section.Title) <= 2 {
		t.Errorf("Incorrect section title: %+v, %+v", len(section.Title), section.Title)
		return
	}
	
	if len(section.Question) != 2 {
		t.Errorf("Incorrect question number: %+v, %+v", len(section.Question), section.Question)
		return
	}
}

func TestDialog_MultipleChoice_Take1stTest(t *testing.T) {
/*	msgs := make([]string, 0, 20)
	keyboards := make([]KeyboardLayout, 0, 20)
	sender := func (text string) { 
			   msgs = append(msgs, text)
			}
	kbdSender := func (text string, kbd KeyboardLayout) { 
			   msgs = append(msgs, text)
			}
	hndl, _ := NewInputTestHandler(sender, dataRootDir)
	dlg := NewDialog(sender,
			time.Second * 5,
			"nemo", hndl)
	dlg.OnCommand("test", []string{ "1.1" })
	time.Sleep(time.Second * 1)
	if(len(msgs) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(msgs))
	}
	if(len(msgs[0]) != 401) {
		t.Errorf("Incorrect message length: %v bytes, %+v", len(msgs[0]), msgs[0])
	}
	msgs = make([]string, 0, 20)
	for i := 0; i < 10; i++ {
		dlg.OnMessage("11")
		time.Sleep(time.Millisecond * 50)
	}
	if len(msgs) != 10 {
		t.Errorf("Incorrect messages: %v bytes, \n\n%+v", len(msgs), msgs)
	}*/
}