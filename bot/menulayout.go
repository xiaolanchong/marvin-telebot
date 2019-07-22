package bot

import (
	"fmt"
)
/*
type MenuItem struct {
	Id		string
	Text	string
}


type OnMenuItem = func (id string)
*/

type MenuItemProc func()
type MenuItemId int

const RootMenuItemId    MenuItemId = 0
const InvalidMenuItemId MenuItemId = -1

type MenuIdAndText struct { 
	Id		MenuItemId
	Text	string
}

type InternalMenuItem struct {
	Id			MenuItemId
	Text		string
	Handler		MenuItemProc
	Children	[]InternalMenuItem
	Parent		*InternalMenuItem
}

type MenuLayout struct {
	Root		InternalMenuItem
	CurrentItem	*InternalMenuItem
	Counter		MenuItemId
}


// ---------

// Breadth first search
func findItemWithParent(parent *InternalMenuItem, itemToFind MenuItemId) *InternalMenuItem {
	if parent.Id == itemToFind {
		return parent
	}
	for i, _ := range(parent.Children) {
		if parent.Children[i].Id == itemToFind {
			return &parent.Children[i]
		}
	}
	for i, _ := range(parent.Children) {
		result := findItemWithParent(&parent.Children[i], itemToFind)
		if result != nil {
			return result
		}
	}
	return nil
}

//----------------------------------------

func NewMenuLayout() (*MenuLayout) {
	menu := &MenuLayout{
		Root: InternalMenuItem {
				Id: RootMenuItemId,
				Children: make([]InternalMenuItem, 0, 20),
			 },
		Counter: RootMenuItemId + 1,
	}
	menu.CurrentItem = &menu.Root
	return menu
}

func (menu *MenuLayout) AddItem(parentId MenuItemId, text string, handler MenuItemProc) (MenuItemId, error) {
	whereToAdd := findItemWithParent(&menu.Root, parentId)
	if whereToAdd == nil {
		return InvalidMenuItemId, fmt.Errorf("Item with id=%d not found ", parentId)
	}
	newItem := InternalMenuItem {
		Id:        menu.Counter,
		Children:  make([]InternalMenuItem, 0, 20),
		Text:      text,
		Handler:   handler,
		Parent:    whereToAdd,
	}
	menu.Counter++
	whereToAdd.Children = append(whereToAdd.Children, newItem)
	return newItem.Id, nil
}

func (menu *MenuLayout) GoTop() {
	menu.CurrentItem = &menu.Root
}

func (menu *MenuLayout) GoUp() {
	if menu.CurrentItem != nil && menu.CurrentItem.Parent != nil {
		menu.CurrentItem = menu.CurrentItem.Parent
	}
}

func (menu *MenuLayout) SelectItem(selectedId MenuItemId) error {
	if menu.CurrentItem == nil {
		return fmt.Errorf("No current item")
	}
	for i, child := range(menu.CurrentItem.Children) {
		if child.Id == selectedId {
			if len(child.Children) != 0 {
				menu.CurrentItem = &menu.CurrentItem.Children[i]
			}
			if child.Handler != nil {
				child.Handler()
				return nil
			} else if (len(child.Children) != 0) { // ok not to have handler for a parent
				return nil
			}
			return fmt.Errorf("No current item handler set for a menu item without children")
		}
	}
	return fmt.Errorf("Item with id=%d not found", selectedId)
}

func (menu *MenuLayout) GetCurrentLevel() []MenuIdAndText {
	if menu.CurrentItem == nil {
		return []MenuIdAndText{}
	}
	result := make([]MenuIdAndText, 0, len(menu.CurrentItem.Children))
	for _, child := range(menu.CurrentItem.Children) {
		result = append(result, MenuIdAndText{ Id: child.Id, Text: child.Text })
	}
	return result
}

func (menu *MenuLayout) Hide() {
	menu.CurrentItem = nil
}
