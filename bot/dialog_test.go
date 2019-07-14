
package bot

import (
	//"strings"
	"testing"
	"time"
)

func TestDialog_GrammarBook(t *testing.T) {
	units, err := LoadUnits("data/Intermediate Korean - a Grammar and Workbook.yaml")
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
학교에 (가다) / повелительное наклонение
= 학교에 가` // "Go to school."`
	if exp_desc != units[0].Exercise[0].Description {
		t.Errorf("Incorrect description, expected:\n%s\nactual:\n%s", exp_desc, units[0].Exercise[0].Description)
	}

	if 6 != len(units[0].Exercise[0].Question) {
		t.Errorf("Incorrect question: expected: %v, actual: %v", 6, len(units[0].Exercise[0].Question))
	}
	
	questions := units[0].Exercise[0].Question
	exp_quest := "공연을 일찍 (마치다)/propositive"
	if exp_quest != questions[len(questions)-1].Text {
		t.Errorf("Incorrect question: expected: %v, actual: %v", exp_quest, questions[len(questions)-1].Text)
	}
}

func TestDialog_StartCommand(t *testing.T) {
	msgs := make([]string, 0, 10)
	dlg := NewDialog(func (text string) { 
			   msgs = append(msgs, text)
			},
			time.Second * 5,
			"nemo")
	dlg.OnCommand("start", []string{})
	time.Sleep(time.Second * 1)
	if(len(msgs) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(msgs))
	}
	if(len(msgs[0]) != 217) {
		t.Errorf("Incorrect message length: %v bytes", len(msgs[0]))
	}
}

func TestDialog_TestCommand(t *testing.T) {
	msgs := make([]string, 0, 10)
	dlg := NewDialog(func (text string) { 
			   msgs = append(msgs, text)
			},
			time.Second * 5,
			"nemo")
	dlg.OnCommand("test", []string{ "1.1" })
	time.Sleep(time.Second * 1)
	if(len(msgs) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(msgs))
	}
	if(len(msgs[0]) != 353) {
		t.Errorf("Incorrect message length: %v bytes, %+v", len(msgs[0]), msgs[0])
	}
}

func TestDialog_Take1stTest(t *testing.T) {
	msgs := make([]string, 0, 20)
	dlg := NewDialog(func (text string) { 
			   msgs = append(msgs, text)
			},
			time.Second * 5,
			"nemo")
	dlg.OnCommand("test", []string{ "1.1" })
	time.Sleep(time.Second * 1)
	if(len(msgs) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(msgs))
	}
	if(len(msgs[0]) != 353) {
		t.Errorf("Incorrect message length: %v bytes, %+v", len(msgs[0]), msgs[0])
	}
	msgs = make([]string, 0, 20)
	for i := 0; i < 10; i++ {
		dlg.OnMessage("11")
		time.Sleep(time.Millisecond * 50)
	}
	if len(msgs) != 10 {
		t.Errorf("Incorrect messages: %v bytes, \n\n%+v", len(msgs), msgs)
	}
}