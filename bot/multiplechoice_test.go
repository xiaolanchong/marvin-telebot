package bot

import (
	"testing"
	"strings"
)

func TestMultipleChoice_Start(t *testing.T) {
	questions := []MultipleChoiceQuestion {
		MultipleChoiceQuestion{Text: "q1", Choice: []string{"1", "2"}, Answer: "2", Comment: "c1" },
		MultipleChoiceQuestion{Text: "q2", Choice: []string{"A", "B"}, Answer: "A", Comment: "c2" },
	}
	sender := MsgKeeper{}
	handler, _ := NewMultipleChoiceHandler(sender.Send, questions)
	
	if(len(sender.OutMessages) != 1) {
		t.Errorf("1 message expected: %v message(s)", len(sender.OutMessages))
	}
	if(sender.OutMessages[0].Text != "q1") {
		t.Errorf("Incorrect text: %v", sender.OutMessages[0].Text)
	}
	
	handler.ProcessKeyboard("1", 1)
}

func TestMultipleChoice_DoQuiz(t *testing.T) {
	questions := []MultipleChoiceQuestion {
		MultipleChoiceQuestion{Text: "q1", Choice: []string{"1", "2"}, Answer: "2", Comment: "c1" },
		MultipleChoiceQuestion{Text: "q2", Choice: []string{"A", "B"}, Answer: "A", Comment: "c2" },
	}
	sender := MsgKeeper{}
	handler, _ := NewMultipleChoiceHandler(sender.Send, questions)
	sender.OutMessages = make([]OutMessage, 0, 20)

	handler.ProcessKeyboard("0", 1)
	if(len(sender.OutMessages) != 3) { // remove keyboard, result, next question
		t.Errorf("3 messages expected, actually %v message(s), %+v", len(sender.OutMessages), sender.OutMessages)
		return
	}
	if strings.Index(sender.OutMessages[1].Text, "ошибка") == -1 {
		t.Errorf("Incorrect text: %+v", sender.OutMessages[0].Text)
	}
	if strings.Index(sender.OutMessages[1].Text, "c1") == -1 {
		t.Errorf("Incorrect text: %+v", sender.OutMessages[0].Text)
	}
	if sender.OutMessages[1].IsKeyboardMsg {
		t.Errorf("Incorrect msg type: %+v", sender.OutMessages[0])
	}
	if sender.OutMessages[2].Text != "q2" {
		t.Errorf("Incorrect text: %+v", sender.OutMessages[1].Text)
	}
	if sender.OutMessages[2].IsKeyboardMsg {
		t.Errorf("Incorrect msg type: %+v", sender.OutMessages[0])
	}

	sender.OutMessages = make([]OutMessage, 0, 20)
	handler.ProcessKeyboard("0", 2)
	if len(sender.OutMessages) != 3 {
		t.Errorf("3 message expected, actually %v message(s)", len(sender.OutMessages))
		return
	}
	if strings.Index(sender.OutMessages[1].Text, "верно!") == -1 {
		t.Errorf("Incorrect message: %v", sender.OutMessages[1].Text)
	}
	expSubText := "Тест окончен, кол-во верных ответов 1/2"
	if sender.OutMessages[2].Text[0:len(expSubText)] != expSubText {
		t.Errorf("Incorrect message: %v", sender.OutMessages[2].Text[0:len(expSubText)])
	}
}