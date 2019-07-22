package bot

import (
	"testing"
	"os"
	"path/filepath"
)

var testMultipleChoiceFileName string

func init() {
	dataRootDir := os.Getenv("DATA_ROOT")
	testMultipleChoiceFileName = filepath.Join(dataRootDir, MockTest1Dir, ListeningTestFileName)
}

func Test_LoadMultipleChoiceTest(t *testing.T) {
	root, err := LoadMultipleChoiceTest(testMultipleChoiceFileName)
	if err != nil {
		t.Errorf("Failed to load multiple choice test: %+v", err)
		return
	}
	
	if len(root.Title) <= 2 {
		t.Errorf("Incorrect section title: %+v, %+v", len(root.Title), root.Title)
		return
	}
	
	if len(root.Question) < 2 {
		t.Errorf("Incorrect question number: %+v, %+v", len(root.Question), root.Question)
		return
	}
}

func TestMenu_Empty(t *testing.T) {
	menu := NewMenuLayout()
	
	var layout []MenuIdAndText
	layout = menu.GetCurrentLevel()
	if len(layout) != 0 {
		t.Errorf("Incorrect menu item number: %+v", len(layout))
	}
	
	menu.GoTop()
	layout = menu.GetCurrentLevel()
	if len(layout) != 0 {
		t.Errorf("Incorrect menu item number: %+v", len(layout))
	}
	
	err := menu.SelectItem(-333)
	if err == nil {
		t.Errorf("Incorrect menu id: %+v", "")
	}
	
	err = menu.SelectItem(RootMenuItemId)
	if err == nil {
		t.Errorf("Incorrect menu id: %+v", "")
	}
}

func TestMenu_OneLevel(t *testing.T) {
	menu := NewMenuLayout()
	selected1st := false
	
	var layout []MenuIdAndText
	id1, err1 := menu.AddItem(RootMenuItemId, "Item1", func() { selected1st = true })
	_, err2 := menu.AddItem(RootMenuItemId, "Item2", nil)
	
	if err1 != nil {
		t.Errorf("Error on adding item: %+v", err1)
	}
	if err2 != nil {
		t.Errorf("Error on adding item: %+v", err2)
	}
	
	layout = menu.GetCurrentLevel()
	if len(layout) != 2 {
		t.Errorf("Incorrect menu item number: %+v", len(layout))
	}
	if layout[0].Text != "Item1" {
		t.Errorf("Incorrect menu item number: %+v", layout[0].Text)
	}
	if layout[1].Text != "Item2" {
		t.Errorf("Incorrect menu item number: %+v", layout[0].Text)
	}
	
	selErr := menu.SelectItem(id1)
	if err1 != nil {
		t.Errorf("Error on selecting item: %+v", selErr)
	}
	if !selected1st {
		t.Errorf("item not selected: %+v", id1)
	}
}

func TestMenu_TwoLevel(t *testing.T) {
	menu := NewMenuLayout()
	selected1st := false
	
	var layout []MenuIdAndText
	id1, _ := menu.AddItem(RootMenuItemId, "Item1", nil)
	menu.AddItem(RootMenuItemId, "Item2", nil)

	id1_1, err3 := menu.AddItem(id1, "Item2_1", func() { selected1st = true })
	if err3 != nil {
		t.Errorf("Error on adding item: %+v", err3)
	}
	
	selErr := menu.SelectItem(id1)
	if selErr != nil {
		t.Errorf("Error on adding item: %+v", selErr)
	}

	layout = menu.GetCurrentLevel()
	if len(layout) != 1 {
		t.Errorf("Incorrect menu item number: %+v, %+v", len(layout), menu)
		return
	}
	if layout[0].Text != "Item2_1" {
		t.Errorf("Incorrect menu item number: %+v", layout[0].Text)
	}
	
	selErr = menu.SelectItem(id1_1)
	if selErr != nil {
		t.Errorf("Error on adding item: %+v", selErr)
	}
	if !selected1st {
		t.Errorf("item not selected: %+v", id1)
	}
	layout = menu.GetCurrentLevel()
	if len(layout) != 1 {
		t.Errorf("Incorrect menu item number: %+v, \n%+v, \n%+v", len(layout), layout, menu)
	}
	
	menu.GoUp() // top
	layout = menu.GetCurrentLevel()
	if len(layout) != 2 {
		t.Errorf("Incorrect menu item number: %+v, %+v", len(layout), menu)
		return
	}
	if layout[0].Id != id1 {
		t.Errorf("Incorrect menu item number: %+v, %+v", layout[0].Id, menu)
		return
	}
}