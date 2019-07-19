package bot

import (
	"testing"
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

	handler.ProcessKeyboard("1", 1)
	if(len(sender.OutMessages) != 2) {
		t.Errorf("1 message expected, actually %v message(s)", len(sender.OutMessages))
	}
	if(sender.OutMessages[0].Text != "Ошибка! Правильно 2") {
		t.Errorf("Incorrect text: %+v", sender.OutMessages[0].Text)
	}
	if sender.OutMessages[0].IsKeyboardMsg {
		t.Errorf("Incorrect msg type: %+v", sender.OutMessages[0])
	}
	if sender.OutMessages[1].Text != "q2" {
		t.Errorf("Incorrect text: %+v", sender.OutMessages[1].Text)
	}
	if sender.OutMessages[1].IsKeyboardMsg {
		t.Errorf("Incorrect msg type: %+v", sender.OutMessages[0])
	}

	sender.OutMessages = make([]OutMessage, 0, 20)
	handler.ProcessKeyboard("A", 2)
	if(len(sender.OutMessages) != 2) {
		t.Errorf("1 message expected, actually %v message(s)", len(sender.OutMessages))
	}
	if( sender.OutMessages[0].Text != "Верно!") {
		t.Errorf("Incorrect message: %v", sender.OutMessages[0].Text)
	}
	if( sender.OutMessages[1].Text != "Тест окончен, кол-во верных ответов 1/2") {
		t.Errorf("Incorrect message: %v", sender.OutMessages[1].Text)
	}
}