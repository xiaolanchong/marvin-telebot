package bot

import (
	"log"
	"time"
	"strings"
	"net/url"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)



func StartTeleBot(botToken string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	return bot, updates, err
}

func keyboardToMarkup(keyboard KeyboardLayout) tgbotapi.InlineKeyboardMarkup {
	layout := [][]tgbotapi.InlineKeyboardButton{}
	if len(keyboard) != 0 {
		layout = make([][]tgbotapi.InlineKeyboardButton, len(keyboard))
		for row, rowItem := range(keyboard) {
			layout[row] = make([]tgbotapi.InlineKeyboardButton, len(rowItem))
			for col, colItem := range(rowItem) {
				layout[row][col] = 
					tgbotapi.InlineKeyboardButton{
							Text: colItem.Text,
							CallbackData: &rowItem[col].Id,
					}
			}
		}
	}
	return tgbotapi.InlineKeyboardMarkup{ InlineKeyboard: layout }
}

var globalDialog *Dialog

func sendMessageToBot(bot *tgbotapi.BotAPI, outMsg OutMessage, chatId int64) {
	if outMsg.ReplyToMessageId != 0 && outMsg.IsKeyboardMsg {
		layout := keyboardToMarkup(outMsg.Keyboard)
		msg := tgbotapi.NewEditMessageReplyMarkup(chatId, outMsg.ReplyToMessageId, layout)
		bot.Send(msg)
		return
	}
	if len(outMsg.Audio) != 0 {
		url, errParse := url.Parse(outMsg.Audio)
		if errParse != nil {
			log.Printf("Bad URL to parse: %s, error: %v", outMsg.Audio, errParse)
			return
		}
		audioMsg := tgbotapi.NewAudioShare(chatId, outMsg.Audio)
		if len(outMsg.Keyboard) != 0 {
			layout := keyboardToMarkup(outMsg.Keyboard)
			audioMsg.ReplyMarkup = &layout
		}
		
		z, err := bot.Send(audioMsg)
		log.Printf("Audio upload %s: %v, err: %+v, params: %+v, %v", outMsg.Audio, z, err, audioMsg, url)
		return
	}
	msg := tgbotapi.NewMessage(chatId, outMsg.Text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(outMsg.Keyboard) != 0 {
		layout := keyboardToMarkup(outMsg.Keyboard)
		msg.ReplyMarkup = &layout
	}
	bot.Send(msg)
}

func getDialog(bot *tgbotapi.BotAPI, dataRootDir string, chatId int64, username string) *Dialog {
	if globalDialog != nil {
		return globalDialog
	}
	sender := func (outMsg OutMessage) {
					sendMessageToBot(bot, outMsg, chatId)
				}
	dialogHandler, errDlgHndl := NewNavigationHandler(sender, dataRootDir) //NewInputTestHandler(sender, dataRootDir)
	if errDlgHndl != nil {
		return nil
	}
	globalDialog = NewDialog(sender, time.Second * 30, username, dialogHandler)
	return globalDialog
}

func ProcessTeleBotUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, dataRootDir string) {
	if update.CallbackQuery != nil {
		message := update.CallbackQuery.Message
		if message == nil {
			log.Printf("[%s] No message in callback query, ignored", message.From.UserName)
			return
		}
		callbackData := update.CallbackQuery.Data
		if len(callbackData) != 0 {
			dialog := getDialog(bot, dataRootDir, message.Chat.ID, message.From.UserName)
			if dialog == nil {
				log.Printf("[%s] Failed to create a dialog", message.From.UserName)
				return
			}
			dialog.OnKey(callbackData, message.MessageID)
		}
		return
	}

	if update.Message == nil {
		log.Printf("[%d] No message in update", update.UpdateID)
		return
	}
	
	message := update.Message

	log.Printf("[%s] %s", message.From.UserName, message.Text)
	
	dialog := getDialog(bot, dataRootDir, message.Chat.ID, message.From.UserName)
	if(dialog == nil) {
		log.Printf("[%s] Failed to create a dialog", message.From.UserName)
		return
	}
	
	if message.IsCommand() {
		args := strings.Split(message.CommandArguments(), " ")
		dialog.OnCommand(message.Command(), args)
	} else {
		dialog.OnMessage(InMessage{ Text: message.Text, MessageId: MessageId(message.MessageID)})
	}
}